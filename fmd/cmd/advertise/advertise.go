// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package advertise

import (
    "fmt"
    "log"
    "net"
    "os"
    "runtime"
    "strings"
    "time"
    "github.com/DennisSchulmeister/find-my-device/fmd/app"
    "github.com/DennisSchulmeister/find-my-device/fmd/conf"
    "github.com/DennisSchulmeister/find-my-device/fmd/msg"
)

// Command "advertise": Advertise device information on the local network
// and/or on a remote registry server.
type AdvertiseCommandStruct struct {
    app.CommandStruct

    app    app.App
    config *conf.Config
}

// Create new command instance
func New(config *conf.Config) app.Command {
    this := &AdvertiseCommandStruct{
        config: config,
    }

    this.CommandStruct.Steps = this
    return this
}

// Provide help information
func (this *AdvertiseCommandStruct) Help() *app.CommandHelp {
    return &app.CommandHelp{
        Description: "Send device announcements on the local network or remote registry server",
        Help: `
            TODO: Help for advertise command
        `,
    }
}

// Set app instance
func (this *AdvertiseCommandStruct) App(app app.App) {
    this.app = app
}

// Return header string with name and configuration
func (this *AdvertiseCommandStruct) Header() string {
    builder := strings.Builder{}

    builder.WriteString("Advertise device information\n")
    builder.WriteString("============================\n")
    builder.WriteString("\n")
    builder.WriteString(fmt.Sprintf(" - IPv4 multicast address for local network communication: %v\n", this.config.General.MulticastIP4))
    builder.WriteString(fmt.Sprintf(" - IPv6 multicast address for local network communication: %v\n", this.config.General.MulticastIP6))
    builder.WriteString(fmt.Sprintf(" - UDP port for local network communication: %v\n", this.config.General.Port))
    builder.WriteString(fmt.Sprintf(" - Respond to queries on the local network: %v\n", this.config.Advertise.Respond))
    builder.WriteString(fmt.Sprintf(" - Send advertisement multicasts on the local network: %v\n", this.config.Advertise.Multicast))
    builder.WriteString(fmt.Sprintf(" - Seconds between advertisements: %v\n", this.config.Advertise.Interval * time.Second))

    return builder.String()
}

// Check configuration values
func (this *AdvertiseCommandStruct) Validate() error {
    if this.config.Advertise.Respond || this.config.Advertise.Multicast {
        return msg.ValidateConfig(this.config)
    }

    return nil
}

// Return go-routines to be started
func (this *AdvertiseCommandStruct) Go() []app.CommandFunc {
    functions := make([]app.CommandFunc, 0)

    if this.config.Advertise.Multicast {
        functions = append(functions, this.sendLocalAnnouncements)
    }

    if this.config.Advertise.Respond {
        functions = append(functions, this.respondToLocalRequests)
    }

    return functions
}

// Send periodic device announcements on the local network
func (this *AdvertiseCommandStruct) sendLocalAnnouncements() error {
    // Create network connections
    conns, err := msg.DialMulticast(this.config)
    if err != nil { return err }
    defer conns.Close()

    fmt.Println()
    for _, conn := range conns.Connections() {
        fmt.Printf("Advertisement multicasts will be sent to %v\n", conn.RemoteAddr())
    }
    fmt.Println()

    // Periodically send advertisement datagrams
    messageCoder, err := msg.NewMessageCoder(nil, conns)
    if err != nil { return err }

    for {
        // Send advertisements
        log.Println("Sending advertisement multicast")

        message := this.newDeviceAdvertisementMessage()

        if err := messageCoder.Write(message); err != nil {
            log.Printf("%v", err)
        }

        // Wait for notification or timeout
        var action string

        select {
            case action = <- this.CommandStruct.Notify:
            case <- time.After(this.config.Advertise.Interval * time.Second):
                action = "send"
        }

        if action == "quit" {
            break
        }
    }

    return messageCoder.Close()
}

// Answer find requests on the local network
func (this *AdvertiseCommandStruct) respondToLocalRequests() error {
    return nil
}

// Create new device advertisement message
func (this *AdvertiseCommandStruct) newDeviceAdvertisementMessage() msg.Message {
    var err error
    message := msg.Message{}
    message.DeviceAdvertisement = &msg.DeviceAdvertisementMessage{}

    message.DeviceAdvertisement.Group      = this.config.Advertise.Group
    message.DeviceAdvertisement.DeviceName = this.config.Advertise.DeviceName

    message.DeviceAdvertisement.HostName, err = os.Hostname()
    if err != nil { log.Printf("%v", err) }

    if message.DeviceAdvertisement.DeviceName == "" {
        message.DeviceAdvertisement.DeviceName = message.DeviceAdvertisement.HostName
    }

    return message
}

// Create new device information message
func (this *AdvertiseCommandStruct) newDeviceInformationMessage() msg.Message {
    var err error
    message := msg.Message{}
    message.DeviceInformation = &msg.DeviceInformationMessage{}

    // General system information
    message.DeviceInformation.Group           = this.config.Advertise.Group
    message.DeviceInformation.DeviceName      = this.config.Advertise.DeviceName
    message.DeviceInformation.OperatingSystem = runtime.GOOS

    message.DeviceInformation.HostName, err = os.Hostname()
    if err != nil { log.Printf("%v", err) }

    if message.DeviceInformation.DeviceName == "" {
        message.DeviceInformation.DeviceName = message.DeviceAdvertisement.HostName
    }

    // Network information
    message.DeviceInformation.NetworkInterfaces = make([]msg.NetworkInterface, 0)

    netInterfaces, err := net.Interfaces()
    if err != nil { log.Printf("%v", err) }

    for _, netInterface := range netInterfaces {
        networkInterface := msg.NetworkInterface{}
        networkInterface.Interface = netInterface
        networkInterface.Addresses = make([]msg.NetworkAddress, 0)
        networkInterface.Multicast = make([]msg.NetworkAddress, 0)

        netAddress, err := netInterface.Addrs()
        if err != nil { log.Printf("%v", err) }

        for _, netAddress := range netAddress {
            networkAddress := msg.NetworkAddress{}
            networkAddress.Network = netAddress.Network()
            networkAddress.Address = netAddress.String()

            networkInterface.Addresses = append(networkInterface.Addresses, networkAddress)
        }

        netMulticasts, err := netInterface.MulticastAddrs()
        if err != nil { log.Printf("%v", err) }

        for _, netMulticast := range netMulticasts {
            networkMulticast := msg.NetworkAddress{}
            networkMulticast.Network = netMulticast.Network()
            networkMulticast.Address = netMulticast.String()

            networkInterface.Multicast = append(networkInterface.Multicast, networkMulticast)
        }
    }

    return message
}
