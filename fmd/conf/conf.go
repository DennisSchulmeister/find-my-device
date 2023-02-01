// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package conf

import "time"

// Definition of the central program configuration.
// Used to read flags, env variables and config files.
type Config struct {
    General   GeneralConfig   `prefix: "FMD_"           command: ""`
    Advertise AdvertiseConfig `prefix: "FMD_ADVERTISE_" command: "advertise"`
    Find      FindConfig      `prefix: "FMD_FIND_"      command: "find"`
    Listen    ListenConfig    `prefix: "FMD_LISTEN_"    command: "listen"`
    Remote    RemoteConfig    `prefix: "FMD_REMOTE_"    command: "remote"`
    Registry  RegistryConfig  `prefix: "FMD_REGISTRY_"  command: "registry"`
}

// General configuration values for all commands.
// NOTE: Field names must not conflict with fields in the other structures!
type GeneralConfig struct {
    MulticastIP4 string         `default: "224.0.0.1"   hide: "false"   help: "IPv4 multicast address for local network communication"`
    MulticastIP6 string         `default: "ff02::1"     hide: "false"   help: "IPv6 multicast address for local network communication"`
    InterfaceIP6 string         `default: ""            hide: ""        help: "Comma-separated list of network devices for IPv6 multicast"`
    Port         uint32         `default: "54321"       hide: "false"   help: "UDP port for local network communication"`
    URL          string         `default: "https://find-my-device.iot-embedded.de"  hide: "false"   help: "URL of remote registry server"`
    Username     string         `default: ""            hide: "false"   help: "Username to authenticate at the remote registry server"`
    Password     string         `default: ""            hide: "true"    help: "Password to authenticate at the remote registry server"`
    Interactive  bool           `default: "true"        hide: "false"   help: "Ask user to enter missing values interactively"`
}

type AdvertiseConfig struct {
    Respond      bool           `default: "true"        hide: "false"   help: "Respond to find requests on the local network"`
    Multicast    bool           `default: "true"        hide: "false"   help: "Send device announcements on the local network"`
    Registry     bool           `default: "true"        hide: "false"   help: "Advertise device information on remote registry server"`
    Interval     time.Duration  `default: "15"          hide: "false"   help: "Seconds between advertisements"`
    Group        string         `default: ""            hide: "false"   help: "Optional name to group related devices"`
    DeviceName   string         `default: ""            hide: "false"   help: "Name of the device if not the system hostname"`
    SecretKey    string         `default: ""            hide: "true"    help: "Secret key to encrypt and restrict access to device information"`
    AuthKey      string         `default: ""            hide: "true"    help: "Owner authorization key in the remote registry"`
}

type FindConfig struct {
    Local        bool           `default: "true"        hide: "false"   help: "Find devices on the local network"`
    Registry     bool           `default: "true"        hide: "false"   help: "Find devices on remote registry server"`
    DeviceName   bool           `default: ""            hide: "false"   help: "Comma-separated list of searched devices"`
    SecretKey    bool           `default: ""            hide: "true"    help: "Secret key to access the device information"`
}

type ListenConfig struct {
    Timeout      time.Duration  `default: "0"           hide: "false"   help: "Maximum number of seconds to listen"`
}

type RemoteConfig struct {
    Request      string         `default: ""            hide: "false"   help: "Remote request. See help text for allowed values."`
    Value        string         `default: ""            hide: "false"   help: "Parameter value for a remote request. See help text for details."`
}

type RegistryConfig struct {
    UI           bool           `default: "true"        hide: "false"   help: "Serve WEB UI for human users"`
    REST         bool           `default: "true"        hide: "false"   help: "Serve REST webservice for remote devices"`
    Listen       bool           `default: "true"        hide: "false"   help: "Listen to device advertisements on the local network"`
    Scan         time.Duration  `default: "0"           hide: "false"   help: "Scan for devices on the local network every X seconds"`
    Anonymous    bool           `default: "false"       hide: "false"   help: "Allow anonymous access without authentication"`
    SelfSignup   bool           `default: "false"       hide: "false"   help: "Allow users and devices to signup themselves"`
}
