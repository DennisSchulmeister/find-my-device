// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package main

import (
    "os"
    "github.com/DennisSchulmeister/find-my-device/fmd/app"
    "github.com/DennisSchulmeister/find-my-device/fmd/conf"
    "github.com/DennisSchulmeister/find-my-device/fmd/cmd/advertise"
    "github.com/DennisSchulmeister/find-my-device/fmd/cmd/find"
    "github.com/DennisSchulmeister/find-my-device/fmd/cmd/listen"
)

// Main function :-)
func main() {
    myApp  := app.NewApp()
    config := &conf.Config{}

    // TODO: Read configuration
    config.General.Port         = 54321
    config.General.MulticastIP4 = "224.0.0.1"
    config.General.MulticastIP6 = "ff02::1"
    config.Advertise.Respond    = true
    config.Advertise.Multicast  = true
    config.Advertise.Interval   = 15
    config.Find.Local           = true
    config.Listen.Timeout       = 0
    ////

    myApp.AddCommand("advertise", advertise.New(config))
    myApp.AddCommand("find", find.New(config))
    myApp.AddCommand("listen", listen.New(config))

    // TODO: Add other commands

    err := myApp.Run(os.Args)
    app.ExitOnError(err)
}
