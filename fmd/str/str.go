// Find My Device
// © 2023 Dennis Schulmeister-Zimolong <dennis@wpvs.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package str

import (
    "strings"
    "github.com/lithammer/dedent"
)

// Takes a CamelCasedString and splits it into individual words.
// Words can also be separated with underlines, e.g. "registry_port" or "Device_Id".
// Special care is taken for acronyms (words in all upper case) like "UDP_Port".
func SplitCamelCaseString(original string) (words []string) {
    acronym := 0
    uppers  := []rune(strings.ToUpper(original))
    lowers  := []rune(strings.ToLower(original))
    builder := strings.Builder{}
    result  := []string{}

    for i, runeOriginal := range original {
        runeUpper := uppers[i]
        runeLower := lowers[i]

        if runeOriginal == '_' {
            result = append(result, builder.String())
            builder.Reset()
            acronym = 0
        } else if runeUpper != runeLower && runeOriginal == runeUpper {
            if acronym == 0 {
                result = append(result, builder.String())
                builder.Reset()
            }

            acronym += 1
        } else if runeOriginal == runeLower {
            if acronym > 1 {
                result = append(result, builder.String())
                builder.Reset()
            }

            acronym = 0
        }

        if runeOriginal == '_' && i > 0 {
            continue
        }

        builder.WriteRune(runeOriginal)
    }

    result = append(result, builder.String())

    for _, word := range(result) {
        if word != "" {
            words = append(words, word)
        }
    }

    return words
}

// Dedent multi-line string and remove leading linebreaks, so that multi-line
// string literals can be nicely indented in the source code.
func FixMultiLineString(str string) string {
    return strings.TrimSpace(dedent.Dedent(str))
}
