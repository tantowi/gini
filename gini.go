/*
MIT License

Copyright (c) 2020 Tantowi Mustofa, ttw@tantowi.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

//
// Package gini read INI format configuration file.
// Many use INI format for configuration file because of its simplicity
//
//    1. Ini file grouped by section. Section name appears on a line by itself, in square bracket ( [] )
//    2. Properties appears below its section. Property has name and value, separated by equal sign ( = )
//    3. Characters after # or ; is comment and ignored
//
// Example :
//    # my configuration file
//
//    [dbserver]
//    host = 192.168.0.10
//    port = 5432
//    user = postgres
//    pass = postgres
//
//    [apiserver]
//    host = 192.168.0.20
//    port = 8080
//
package gini

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

//
// Ini configuration type
//
type Ini struct {
	sections map[string]keys
}

type keys map[string]string

//
// Read a value from configuration with specified sectionName and keyName.
// Return "" if section or key not found
//
func (f *Ini) Read(sectionName, keyName string) string {
	keys := f.sections[sectionName]
	if keys == nil {
		return ""
	}

	return keys[keyName]
}

//
// SectionExists check whether a section exists.
// Return `true` if the section exists
//
func (f *Ini) SectionExists(sectionName string) bool {
	_, fnd := f.sections[sectionName]
	return fnd
}

//
// KeyExists check whether a key in a section exists.
// Return `true` if the key exists
//
func (f *Ini) KeyExists(sectionName, keyName string) bool {
	sect, fnd := f.sections[sectionName]
	if !fnd {
		return false
	}

	_, fnd = sect[keyName]
	return fnd
}

//
// SectionList list of available sections.
// Return `array of string` contains available sections
//
func (f *Ini) SectionList() []string {
	var lst []string

	for k := range f.sections {
		lst = append(lst, k)
	}

	return lst
}

//
// KeyList list of keys on a specified section.
// Return `array of string` contains available keys in that section
//
func (f *Ini) KeyList(sectionName string) []string {
	var lst []string

	sect, fnd := f.sections[sectionName]
	if !fnd {
		return lst
	}

	for k := range sect {
		lst = append(lst, k)
	}

	return lst
}

//
// parseIni parse ini string in `Reader`.
// Return map of sections
// Also return error if occured while reading and parsing the INI. On successful, error is nil
//
func parseIni(in *bufio.Reader) (map[string]keys, error) {
	var data = make(map[string]keys)
	var sectionName string = ""
	var done = false
	var lineNumber = 0

	for !done {
		line, err := in.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				done = true
			} else {
				return nil, err
			}
		}
		lineNumber++

		// remove comment
		n := strings.IndexAny(line, "#;")
		if n >= 0 {
			line = line[0:n]
		}

		// trim the line
		line = strings.TrimSpace(line)
		//log.Println(">" + line)

		// skip blank line
		if len(line) == 0 {
			continue
		}

		ln := len(line)

		// section
		if line[0] == '[' {
			if ln <= 2 || line[ln-1] != ']' {
				return nil, errors.New("Invalid section at line " + strconv.Itoa(lineNumber))
			}

			sectionName = strings.TrimSpace(strings.ToLower(line[1 : ln-1]))
			//log.Println(">> section: " + sectionName)
			continue
		}

		// key
		n = strings.IndexRune(line, '=')
		if n < 0 {
			return nil, errors.New("Invalid format at line " + strconv.Itoa(lineNumber))
		}

		if sectionName == "" {
			return nil, errors.New("Key without section at line " + strconv.Itoa(lineNumber))
		}

		keyName := strings.ToLower(strings.TrimSpace(line[:n]))
		keyValue := strings.TrimSpace(line[n+1:])
		//log.Println(">> key: " + name + ", value: " + value )

		if keyName == "" {
			return nil, errors.New("Empty key at line " + strconv.Itoa(lineNumber))
		}

		section, fnd := data[sectionName]
		if !fnd {
			section = make(map[string]string)
			data[sectionName] = section
		}

		section[keyName] = keyValue
	}

	return data, nil
}

//
// LoadReader load INI file as io.Reader
// Return *Ini ready to read
// Return error if occured while reading and parsing the INI. On successful, error is nil
//
func LoadReader(in io.Reader) (*Ini, error) {
	bufin, ok := in.(*bufio.Reader)
	if !ok {
		bufin = bufio.NewReader(in)
	}
	data, err := parseIni(bufin)
	if err != nil {
		return nil, err
	}

	f := new(Ini)
	f.sections = data
	return f, nil
}

//
// LoadFile : load INI file as path
// Return *Ini ready to read
// Return error if occured while reading and parsing the INI. On successful, error is nil
//
func LoadFile(path string) (*Ini, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	return LoadReader(in)
}
