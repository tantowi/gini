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
// Ini : parsed INI config file
//
type Ini struct {
	sections map[string]keys
}

type keys map[string]string

//
// Read : Read a value. Return "" if section or key not found
//
func (f *Ini) Read(sectionName, keyName string) string {
	keys := f.sections[sectionName]
	if keys == nil {
		return ""
	}

	return keys[keyName]
}

//
// SectionExists : check whether a section exists
//
func (f *Ini) SectionExists(sectionName string) bool {
	_, fnd := f.sections[sectionName]
	return fnd
}

//
// KeyExists : check whether a key exists
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
// SectionList : get list of sections
//
func (f *Ini) SectionList() []string {
	var lst []string

	for k := range f.sections {
		lst = append(lst, k)
	}

	return lst
}

//
// KeyList : get list of keys on a section
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
// parseIni : parse ini file
//
func parseIni(in *bufio.Reader) (map[string]keys, error) {
	var data = make(map[string]keys)
	var sectionName string = ""
	var sectionItems map[string]string
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
		n := strings.IndexRune(line, '#')
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

			name := strings.TrimSpace(strings.ToLower(line[1 : ln-1]))
			//log.Println(">> section: \"" + name + "\"")

			if sectionItems != nil {
				// save section to data
				data[sectionName] = sectionItems
				sectionName = ""
				sectionItems = nil
			}

			// create new section
			sectionName = name
			sectionItems = make(map[string]string)
			continue
		}

		n = strings.IndexRune(line, '=')
		if n < 0 {
			return nil, errors.New("Invalid format at line " + strconv.Itoa(lineNumber))
		}

		name := strings.ToLower(strings.TrimSpace(line[:n]))
		value := strings.TrimSpace(line[n+1:])
		//log.Println(">> key: \"" + name + "\"  value: \"" + value + "\"")

		if sectionItems == nil {
			return nil, errors.New("Key without section at line " + strconv.Itoa(lineNumber))
		}

		if name == "" {
			return nil, errors.New("Empty key at line " + strconv.Itoa(lineNumber))
		}

		sectionItems[name] = value
	}

	// save last section
	if sectionItems != nil {
		data[sectionName] = sectionItems
	}

	return data, nil
}

//
// LoadReader : load INI file as io.Reader
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
//
func LoadFile(path string) (*Ini, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	return LoadReader(in)
}
