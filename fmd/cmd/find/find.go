// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package find

import (
    "fmt"
    "github.com/DennisSchulmeister/find-my-device/fmd/app"
    "github.com/DennisSchulmeister/find-my-device/fmd/conf"
)

// Command "find": Try to find a device on the local network
type FindCommandStruct struct {
    app.Command
    config *conf.Config
}

// Create new command instance
func New(config *conf.Config) app.Command {
    return &FindCommandStruct{
        config: config,
    }
}

// Provide help information
func (this *FindCommandStruct) Help() *app.CommandHelp {
    return &app.CommandHelp{
        Description: "Find devices on the local network or remote registry server",
        Help: `
            TODO: Help for find command
        `,
    }
}

// Run the command
func (this *FindCommandStruct) Run(app app.App) error {
    fmt.Println("TODO: Find")
    return nil
}
