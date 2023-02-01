// Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package str

import "testing"

func Test_SplitCamelCaseString(t *testing.T) {
    failed := false
    tests  := make(map[string][]string)

    tests["HelloWorld"]        = []string{"Hello", "World"}
    tests["Hello_WorldAgain"]  = []string{"Hello", "World", "Again"}
    tests["hello_world_again"] = []string{"hello", "world", "again"}
    tests["_hello_world"]      = []string{"_hello", "world"}
    tests["HelloWORLD"]        = []string{"Hello", "WORLD"}
    tests["HelloWORLDagain"]   = []string{"Hello", "WORLD", "again"}
    tests["HelloWORLD_Again"]  = []string{"Hello", "WORLD", "Again"}
    tests["HELLO_World"]       = []string{"HELLO", "World"}

    for str, expectedWords := range tests {
        actualWords := SplitCamelCaseString(str)
        failed = false

        for i, actualWord := range actualWords {
            if i < len(expectedWords) {
                expectedWord := expectedWords[i]

                if actualWord != expectedWord {
                    failed = true
                    break
                }
            } else {
                failed = true
                break
            }
        }

        if failed {
            t.Errorf("SplitCamelCaseString(%v) returned %v instead of %v", str, actualWords, expectedWords)
        }
    }
}

func Test_FixMultiLineString(t *testing.T) {
    if FixMultiLineString("Hello") != "Hello" {
        t.Fail()
    }

    str1 := `
             Hello, world!
            `

    if FixMultiLineString(str1) != "Hello, world!" {
        t.Fail()
    }
}
