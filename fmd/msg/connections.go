// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package msg

import (
    "errors"
    "fmt"
    "net"
    "os"
    "strings"
    "time"
    "golang.org/x/exp/slices"
    "github.com/DennisSchulmeister/find-my-device/fmd/conf"
)

// Connection to multiple remote sites, encapsulating a list of net.PacketConn objects.
// Write() will send the data to all receivers. If data needs to be received,
// concurrent goroutines can monitor the connections and put received data
// into a shared channel.
type Connections interface {
    // Get all open connections
    Connections() []net.PacketConn

    // Listen for incoming data
    Start()

    // Stop listening for incoming data
    Stop()

    // Get channel for connections with readable data
    Read() chan ReadResult

    // Write data to all connections
    Write(b []byte) (n int, err error)

    // Close all connections
    Close() error
}

type ConnectionsStruct struct {
    connections []net.PacketConn
    started     bool
    read        chan ReadResult
    notify      map[net.PacketConn]chan string
}

// Concurrent Write(): Written number of bytes and last error
type writeResult struct{
    n int
    err error
}

// Concurrent Read: Channel with available data or error
type ReadResult struct {
    Connection net.PacketConn
    Error error
}

// Check, if the configuration allows dialling at least one address
func ValidateConfig(config *conf.Config) error {
    if config.General.MulticastIP4 == "" && config.General.MulticastIP6 == "" {
        return fmt.Errorf("Neither IPv4 nor IPv6 multicast address has been defined")
    }

    if config.General.Port == 0 {
        return fmt.Errorf("No UDP port number has been defined")
    }

    return nil
}

// Dial all IPv4 and IPv6 multicast addresses from global config
func DialMulticast(config *conf.Config) (Connections, error) {
    this := &ConnectionsStruct{
        connections: make([]net.PacketConn, 0),
        started:     false,
        read:        make(chan ReadResult),
        notify:      make(map[net.PacketConn]chan string),
    }

    if config.General.MulticastIP4 != "" {
        for i := 0; i < 1; i++ {
            addr := fmt.Sprintf("%v:%v", config.General.MulticastIP4, config.General.Port)
            UDPAddr, err := net.ResolveUDPAddr("udp", addr)
            if err != nil { continue }

            conn, err := net.DialUDP("udp", nil, UDPAddr)
            if err != nil { continue }

            this.connections = append(this.connections, conn)
        }
    }

    if config.General.MulticastIP6 != "" {
        allowedInterfaces := make([]string, 0)

        if strings.Contains(config.General.InterfaceIP6, ",") {
            strings.Split(config.General.InterfaceIP6, ",")
        }

        netInterfaces, err := net.Interfaces()
        if err != nil { return nil, err }

        for _, netInterface := range netInterfaces {
            if netInterface.Flags & net.FlagUp == 0 { continue }
            if netInterface.Flags & net.FlagMulticast == 0 { continue }

            if len(allowedInterfaces) > 0 {
                if !slices.Contains(allowedInterfaces, netInterface.Name) {
                    continue
                }
            }

            addr := fmt.Sprintf("[%v%%%v]:%v", config.General.MulticastIP6, netInterface.Name, config.General.Port)
            UDPAddr, err := net.ResolveUDPAddr("udp", addr)
            if err != nil { continue }

            conn, err := net.DialUDP("udp", nil, UDPAddr)
            if err != nil { continue }

            this.connections = append(this.connections, conn)
        }
    }

    if len(this.connections) == 0 {
        return nil, fmt.Errorf("Unable to dial any address")
    }

    return this, nil
}

// Get all open connections
func (this *ConnectionsStruct) Connections() []net.PacketConn {
    return this.connections
}

// Listen for incoming data
func (this *ConnectionsStruct) Start() {
    if this.started { return }
    this.started = true

    // For each connection call a goroutine that repeatedly probes the connection
    // for available data. Send connections with data or error to the this.read
    // channel. Additionally open a dedicated notify channel for each goroutine,
    // this is used by Stop() to break the loops.
    for _, connection := range this.connections {
        notify := make(chan string)
        this.notify[connection] = notify

        dummy := make([]byte, 0)

        go func() {
            for {
                err := connection.SetReadDeadline(time.Now().Add(1 * time.Second))

                if err == nil {
                    var action string

                    select {
                        case action = <- notify:
                        default:
                            action = "read"
                    }

                    if action == "quit" { break }
                    _, _, err = connection.ReadFrom(dummy) // Besser?
                }

                if err == nil {
                    this.read <- ReadResult{Connection: connection, Error: nil}
                } else if !errors.Is(err, os.ErrDeadlineExceeded) {
                    this.read <- ReadResult{ Connection: connection, Error: err}
                    break
                }
            }
        }()
    }
}

// Stop listening for incoming data
func (this *ConnectionsStruct) Stop() {
    if !this.started { return }
    this.started = false

    for _, notify := range this.notify {
        notify <- "quit"
    }
}

// Get channel for connections with readable data
func (this *ConnectionsStruct) Read() chan ReadResult {
    return this.read
}

// Write data to all connections. Blocks until all data is written.
// n will be the maximum number bytes written, which should be the len(b).
// err wraps all errors from all connections.
func (this *ConnectionsStruct) Write(b []byte) (n int, err error) {
    results := make(chan writeResult)

    for _, connection := range this.connections {
        go func() {
            result := writeResult{}

            for result.n < len(b) {
                n1, err1 := connection.Write(b)
                result.n += n1

                if err1 != nil {
                    result.err = err1
                    break
                }
            }

            results <- result
        }()
    }

    for _, connection := range this.connections {
        result := <- results

        if result.n > n {
            result.n = n
        }

        if result.err != nil {
            err = wrapError(connection, err, result.err)
        }
    }

    return
}

// Close all connections.
// err wraps all errors from all connections.
func (this *ConnectionsStruct) Close() error {
    var err error

    for _, connection := range this.connections {
        if err1 := connection.Close(); err1 != nil {
            err = wrapError(connection, err, err1)
        }
    }

    return err
}

// Wrap multiple connection errors by added newErr to oldErr.
// oldErr can be nil, if there is no previous error.
func wrapError(conn net.PacketConn, err error, add error) error {
    if err == nil {
        return fmt.Errorf("%v - %v", conn.LocalAddr(), add.Error())
    } else {
        return fmt.Errorf("%w; %v - %v", err, conn.LocalAddr(), add.Error())
    }
}
