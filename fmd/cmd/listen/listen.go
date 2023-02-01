// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package listen

import (
    "fmt"
    "github.com/DennisSchulmeister/find-my-device/fmd/app"
    "github.com/DennisSchulmeister/find-my-device/fmd/conf"
)

// Command "listen": Listen for device announcements on the local network and log
// them on the console.
type ListenCommandStruct struct {
    app.Command
    config *conf.Config
}

// Create new command instance
func New(config *conf.Config) app.Command {
    return &ListenCommandStruct{
        config: config,
    }
}

// Provide help information
func (this *ListenCommandStruct) Help() *app.CommandHelp {
    return &app.CommandHelp{
        Description: "Listen for device announcements on the local network",
        Help: `
            TODO: Help for listen command
        `,
    }
}

// Run the command
func (this *ListenCommandStruct) Run(app app.App) error {
    fmt.Println("TODO: Listen")
    return nil
}
