// fmd: Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package msg

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "io"
    "reflect"
    "strings"
)

// Encoder/Decoder to read and write messages from a byte-stream
type MessageCoder interface {
    // Encode message and write it to the byte-stream
    Write(message Message) error

    // Decode next message inside the byte-stream
    Read() (Message, error)

    // Close internal readers and writers
    Close() error
}

type MessageCoderStruct struct {
    gzipReader  *gzip.Reader
    jsonDecoder *json.Decoder

    jsonBuffer  bytes.Buffer
    jsonEncoder *json.Encoder
    gzipWriter  *gzip.Writer
}

// Create new message encoder/decoder instance
func NewMessageCoder(reader io.Reader, writer io.Writer) (MessageCoder, error) {
    var err error
    this := &MessageCoderStruct{}

    if reader != nil {
        this.gzipReader, err = gzip.NewReader(reader)
        if err != nil { return nil, err }
        this.jsonDecoder = json.NewDecoder(this.gzipReader)
    }

    if writer != nil {
        this.jsonEncoder = json.NewEncoder(&this.jsonBuffer)
        this.gzipWriter  = gzip.NewWriter(writer)
    }

    return this, nil
}

// Write a new message into the output stream.
// Message format is gzip compressed json.
func (this *MessageCoderStruct) Write(message Message) error {
    // NOTE: Each field of the message struct is encoded individually to safe
    // some bytes. Because the JsonEncoder doesn't skip nil pointers but rather
    // encodes them as attributes with value "null".
    builder := strings.Builder{}
    builder.WriteString("{")

    messageType  := reflect.TypeOf(message)
    messageValue := reflect.ValueOf(message)
    emptyObject  := true

    for i := 0; i < messageValue.NumField(); i++ {
        fieldValue := messageValue.Field(i)
        if fieldValue.IsZero() { continue }

        if !emptyObject {
            builder.WriteString(",")
        }

        emptyObject = false

        fieldType  := messageType.Field(i)
        builder.WriteString("\"")
        builder.WriteString(fieldType.Name)
        builder.WriteString("\":")

        this.jsonBuffer.Reset()
        err := this.jsonEncoder.Encode(fieldValue.Interface())
        if err != nil { return err }

        builder.WriteString(strings.TrimSpace(this.jsonBuffer.String()))
    }

    builder.WriteString("}")

    _, err := io.WriteString(this.gzipWriter, builder.String())
    return err
}

// Read a message from the input stream
// Message format is gzip compressed json.
func (this *MessageCoderStruct) Read() (message Message, err error) {
    err = this.jsonDecoder.Decode(&message)
    if err != nil { return }
    return
}

// Close internal readers and writers
func (this *MessageCoderStruct) Close() error {
    if this.gzipReader != nil {
        if err := this.gzipReader.Close(); err != nil {
            return err
        }
    }

    if this.gzipWriter != nil {
        if err := this.gzipWriter.Close(); err != nil {
            return err
        }
    }

    return nil
}
