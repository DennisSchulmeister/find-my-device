// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package app

import (
    "fmt"
    "strings"
    "github.com/DennisSchulmeister/find-my-device/fmd/str"
    //"golang.org/x/exp/maps"
)

// Build in help command. Either shows a list of all available commands or the
// help text and configuration parameters of a given command.
type HelpCommandStruct struct {
    CommandStruct
    app App
}

// Create new help command instance
func NewHelpCommand() Command {
    return &HelpCommandStruct{}
}

// Provide help information
func (this *HelpCommandStruct) Help() *CommandHelp {
    return &CommandHelp{
        Description: "Show available commands or help for a single command",
        Arguments: "[command]",
        Help: `
            See Dilbert strip from August 15, 2000: Our disaster recovery plan is to
            panic and cry 'Help! Help!'. Obviously some disaster just happened or you
            are simply curious how to use this progrm. This page is here to help.

            This is the help page for the help command. So you already know, how to
            get help for a specific command. Each command comes with its own help
            page that can be retrieved with '$program$ help <command>'.
        `,
    }
}

// Run the command
func (this *HelpCommandStruct) Run(app App) error {
    this.app = app
    arguments := app.Arguments()

    if len(arguments) == 0 {
        this.showListOfCommands()
        return nil
    } else {
        return this.showCommandHelp(arguments[0])
    }
}

func (this *HelpCommandStruct) showListOfCommands() {
    // TODO: Print list of available commands
}

func (this *HelpCommandStruct) showCommandHelp(name string) error {
    // Print command name and long help-text
    command, found := this.app.Commands()[name]

    if !found {
        return fmt.Errorf("No help for command '%v' found", name)
    }

    help := command.Help()

    text := str.FixMultiLineString(help.Help)
    text = strings.ReplaceAll(text, "$program$", this.app.Program())
    text = strings.ReplaceAll(text, "$command$", name);


    fmt.Printf("%v %v %v\n", this.app.Program(), name, help.Arguments)
    fmt.Println()
    fmt.Println(text)

    // TODO: Print flag list (with empty line between general and command flags)
    // TODO: Print mapping table that maps flags to env variables and config file entries

    return nil
}
