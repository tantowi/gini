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
	"errors"
	"fmt"
	"strings"
	"testing"
)

var ini01 = `
; This is a comment
# Another comment

[kamus]
makan=eat
Minum =  drink
LIHAT   =   see = watch    # double = should OK

  [  STATUS  ]
  web = active`

//
// TestLoadReader
//
func TestLoadReader(t *testing.T) {
	reader := strings.NewReader(ini01)
	ini, err := LoadReader(reader)
	if err != nil {
		t.Fatal(err)
	}

	chkkey(t, ini, "kamus", "makan", "eat")
	chkkey(t, ini, "kamus", "minum", "drink")
	chkkey(t, ini, "kamus", "lihat", "see = watch")
	chkkey(t, ini, "status", "web", "active")

}

//
// TestLoadFile
//
func TestLoadFile(t *testing.T) {
	ini, err := LoadFile("test_a.ini")
	if err != nil {
		t.Fatal(err)
	}

	chkkey(t, ini, "setting", "color", "red")
	chkkey(t, ini, "setting", "width", "700")
	chkkey(t, ini, "setting", "height", "450")

	chkkey(t, ini, "server", "host", "10.10.20.20")
	chkkey(t, ini, "server", "port", "3344")
}

//
// chkkey
//
func chkkey(t *testing.T, ini *Ini, section, key, expect string) {
	value := ini.Read(section, key)
	if value != expect {
		t.Errorf("Expect %s, got %s (section: %s, key: %s)", expect, value, section, key)
	}

}

//
// TestSectionList
//
func TestSectionList(t *testing.T) {
	reader := strings.NewReader(ini01)
	ini, err := LoadReader(reader)
	if err != nil {
		t.Fatal(err)
	}

	var expect = []string{"kamus", "status"}
	secs := ini.SectionList()

	err = sliceCompare(secs, expect)
	if err != nil {
		t.Fatal(err)
	}
}

//
// TestKeyList
//
func TestKeyList(t *testing.T) {
	reader := strings.NewReader(ini01)
	ini, err := LoadReader(reader)
	if err != nil {
		t.Fatal(err)
	}

	var expect = []string{"makan", "minum", "lihat"}
	lst := ini.KeyList("kamus")

	err = sliceCompare(lst, expect)
	if err != nil {
		t.Fatal(err)
	}
}

//
// SliceCompare
//
func sliceCompare(value, expect []string) error {
	if value == nil {
		if expect != nil {
			return errors.New("Expect NOT NIL, got NIL")
		}
		return nil
	}

	if expect == nil {
		if value != nil {
			return errors.New("Expect NIL, got NOT NIL")
		}
		return nil
	}

	if len(value) != len(expect) {
		return fmt.Errorf("Expect: %d items, got: %d items", len(expect), len(value))
	}

	for i, v := range value {
		if v != expect[i] {
			return fmt.Errorf("Expect: %s, got: %s at item: %d", expect[i], v, i)
		}
	}

	return nil
}

//
// TestLostKey
//
func TestLostKey(t *testing.T) {
	ini := "# This is a comment\n" +
		"satu = Is One    # English translation\n" +
		"[kamus]\n" +
		"makan=eat\n"

	reader := strings.NewReader(ini)
	_, err := LoadReader(reader)
	if err == nil {
		t.Errorf("Undetected Lost Key")
	}
}

//
// TestEmptySection
//
func TestEmptySection(t *testing.T) {
	ini := "# This is a comment" +
		"[]\n" +
		"makan=eat\n"

	reader := strings.NewReader(ini)
	_, err := LoadReader(reader)
	if err == nil {
		t.Errorf("Undetected Empty Section")
	}
}

//
// TestMalformedSection
//
func TestMalformedSection(t *testing.T) {
	ini := "# This is a comment" +
		"[   \n" +
		"makan=eat\n"

	reader := strings.NewReader(ini)
	_, err := LoadReader(reader)
	if err == nil {
		t.Errorf("Undetected Malformed Section")
	}
}

//
// TestMalformedKey
//
func TestMalformedKey(t *testing.T) {
	ini := "# This is a comment" +
		"[STATUS]\n" +
		"makan makan\n"

	reader := strings.NewReader(ini)
	_, err := LoadReader(reader)
	if err == nil {
		t.Errorf("Undetected Malformed Key")
	}
}

//
// TestComment
//
func TestComment(t *testing.T) {
	data := "# Test comment\n" +
		"[ENGLISH]\n" +
		"#satu = one\n" +
		"dua = two\n" +
		";tiga = three\n"

	reader := strings.NewReader(data)
	ini, err := LoadReader(reader)
	if err != nil {
		t.Fatal(err)
	}

	if ini.Read("english", "satu") == "one" {
		t.Errorf("Error comment detection '#one'")
	}

	if ini.Read("english", "#satu") == "one" {
		t.Errorf("Error comment detection '#one'")
	}

	if ini.Read("english", "dua") != "two" {
		t.Errorf("Error comment detection 'two'")
	}

	if ini.Read("english", "tiga") == "three" {
		t.Errorf("Error comment detection ';tiga'")
	}

	if ini.Read("english", "tiga") == ";three" {
		t.Errorf("Error comment detection ';tiga'")
	}

}
