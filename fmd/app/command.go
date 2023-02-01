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
    "os/signal"
    "syscall"
    "golang.org/x/sync/errgroup"
    "github.com/DennisSchulmeister/find-my-device/fmd/str"
)

//------------------------------------------------------------------------------
// Command
//------------------------------------------------------------------------------

// Object implementing an app command called by "./program command".
// This is where the interesting things happen.
type Command interface {
    // Provide help information
    Help() *CommandHelp

    // Run the command
    Run(app App) error
}

type CommandStruct struct {
    Notify chan string
    Steps  CommandSteps
}

// Help texts for a command
type CommandHelp struct {
    // Short description
    Description string

    // Non-flag command arguments
    Arguments string

    // Detailed help text, with the following placeholders:
    //  * $program$ for the executable name.
    //  * $command$ for the command name.
    Help string
}

// Template methods for the default implementation of the Command Run() method.
// The default implementation contains shared logic to validate configuration
// and start one or more goroutines to perform the actual logic.
//
// The "Notify" channel in the CommandStruct is used for communication with the
// goroutines. An interrupt signal (usually triggered by Ctrl+C) will cause the
// string "quit" to be sent to the channel. Other than that, the channel has no
// pre-defined use and can be freely used at will.
type CommandSteps interface {
    // Set app instance
    App(app App)

    // Get header string with command description and configuration values
    // TODO: Refactor so that the command must only tell, which configuration values to print
    Header() string

    // Perform sanity checks on the configuration values
    Validate() error

    // Return one or more go-routines for actual command execution
    Go() []CommandFunc
}

// Go-routine with the actual command logic
type CommandFunc func() error

// Default implementation of the Help() method
func (this *CommandStruct) Help() *CommandHelp {
    return &CommandHelp{}
}

// Default implementation of the Run() method. Assumes, that any non-trivial
// implementation will set the steps attribute of its parent to provide the
// actual command logic.
func (this *CommandStruct) Run(app App) error {
    if this.Steps == nil { panic("The default Run() method needs this.Steps") }

    // Print configuration values
    header := str.FixMultiLineString(this.Steps.Header()) + "\n"
    fmt.Printf(header)

    // Check configuration values
    if err := this.Steps.Validate(); err != nil {
        return err
    }

    // Handle interrupt signals caused by Ctrl+C
    this.Notify = make(chan string)

    signals := make(chan os.Signal)
    signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

    go func() {
        <- signals
        this.Notify <- "quit"
    }()

    // Execute go-routines
    waitgroup := new(errgroup.Group)

    for _, goroutine := range this.Steps.Go() {
        waitgroup.Go(goroutine)
    }

    return waitgroup.Wait()
}
