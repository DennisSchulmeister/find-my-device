// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package msg

import "net"

// Main message structure for passing around network messages in the program.
// All commands of the program, that transmit or receive network datagrams
// use this structure to construct or receive messages in a type-safe way.
type Message struct {
    ClientRequest       *ClientRequestMessage
    DeviceAdvertisement *DeviceAdvertisementMessage
    DeviceInformation   *DeviceInformationMessage
}

// Generic request from client to device
type ClientRequestMessage struct {
    Request    string
    Parameters []string
}

// Local device advertisement multicast
type DeviceAdvertisementMessage struct {
    Group      string
    DeviceName string
    HostName   string
}

// Detailed device information
type DeviceInformationMessage struct {
    Group             string
    DeviceName        string
    HostName          string
    OperatingSystem   string
    NetworkInterfaces []NetworkInterface
}

type NetworkInterface struct {
    net.Interface
    Addresses []NetworkAddress
    Multicast []NetworkAddress
}

type NetworkAddress struct {
    Network string
    Address string
}
