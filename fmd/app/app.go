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
    "os"
    "strings"
)

// Main object of a minimal command line app framework. This embodies some
// generic logic for simple apps that take a command as the first parameter
// followed by arguments and flags for that command. Additionally allows to
// configure the app via env variables and config files as an alternative to
// the command line flags.
//
// Inside the main function of the program simply AddCommand() and Run() must
// be called to register and execute commands. The result of Run() should be
// passed to the ExitOnError() function.
type App interface {
    // Register a named command with the app instance.
    // Must be called one ore more times before Run().
    // The name must be unique and the command not be nil.
    AddCommand(name string, command Command)

    // Run the application.
    // os.Args should usually be used for args.
    // Returned errors should be displayed with ExitOnError().
    Run(args []string) error

    // Get all registered commands by their name
    Commands() map[string]Command

    // Get program name
    Program() string

    // Get name of the current command
    Command() string

    // Get non-flag arguments of the current command
    Arguments() []string
}

type AppStruct struct {
    // Registered commands
    commands map[string]Command

    // Name of the executable
    program string

    // Name of the currently running command
    command string

    // Non-flag arguments of the currently running command
    arguments []string
}

// Create new app instance
func NewApp() App {
    this := &AppStruct{
        commands: make(map[string]Command),
        arguments: make([]string, 0),
    }

    this.AddCommand("help", NewHelpCommand())
    return this
}

func (this *AppStruct) Commands() map[string]Command { return this.commands }
func (this *AppStruct) Program() string { return this.program }
func (this *AppStruct) Command() string { return this.command }
func (this *AppStruct) Arguments() []string { return this.arguments }

// Register command with the app
func (this *AppStruct) AddCommand(name string, command Command) {
    _, found := this.commands[name]

    if found {
        panic(fmt.Sprintf("Command '%v' has already been registered with the app", name))
    }

    if command == nil {
        panic(fmt.Sprintf("Pointer to implementation of command '%v' is nil", name))
    }

    this.commands[name] = command
}

// Run the command requested in the arguments
func (this *AppStruct) Run(args []string) error {
    // Get program name
    if len(args) < 1 {
        panic(fmt.Sprintf("Argument list is empty, cannot determine program name"))
    }

    this.program = args[0]

    // Get name of the requested command
    if len(args) < 2 {
        return fmt.Errorf("Usage: %v command [<flags...>]", this.program)
    }

    this.command = args[1]
    command, found := this.commands[this.command]

    if !found {
        return fmt.Errorf("Unknown command: %v", this.command)
    }

    // Split arguments and flags
    flagIndex := 2

    for i, arg := range args[2:] {
        if strings.HasPrefix(arg, "-") { break }
        this.arguments = append(this.arguments, arg)
        flagIndex = 3 + i   // i always starts at 0
    }

    flags := args[flagIndex:]

    // TODO: Read configuration values
    fmt.Printf("Arguments: %v\n", this.arguments)
    fmt.Printf("Flags: %v\n", flags)
    fmt.Println()

    // Run the requested command
    return command.Run(this)
}

// Does nothing, if err is nil. Otherwise prints the error to Os.Stderr
// (including a trailing newline, if missing) and exits the program with
// return code 1.
func ExitOnError(err error) {
    if err == nil {
        return
    }

    msg := err.Error()
    os.Stderr.WriteString(msg)

    if !strings.HasSuffix(msg, "\n") {
        os.Stderr.WriteString("\n")
    }

    os.Exit(1)
}
