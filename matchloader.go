// The MIT License (MIT)

// Copyright (c) 2016 Maciej Borzecki

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"bytes"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type rawStat struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
}

type rawItem struct {
	Name    string `xml:"name,attr"`
	Pickups int    `xml:"pickups,attr"`
	Time    int    `xml:"time,attr"`
}

type rawWeapon struct {
	Name  string `xml:"name,attr"`
	Hits  int    `xml:"hits,attr"`
	Shots int    `xml:"shots,attr"`
	Kills int    `xml:"kills,attr"`
}

type rawPlayer struct {
	Name     string      `xml:"name,attr"`
	Stats    []rawStat   `xml:"stat"`
	Items    []rawItem   `xml:"items>item"`
	Weapons  []rawWeapon `xml:"weapons>weapon"`
	Powerups []rawItem   `xml:"powerups>item"`
}

type rawMatch struct {
	Datetime string      `xml:"datetime,attr"`
	Map      string      `xml:"map,attr"`
	Type     string      `xml:"type,attr"`
	Players  []rawPlayer `xml:"player"`
}

func LoadMatch(src io.Reader) (*Match, error) {
	raw, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load match data")
	}

	var match rawMatch
	xml.Unmarshal(raw, &match)

	log.Printf("match: %+v", match)
	return nil, nil
}

func LoadMatchFile(path string) (*Match, error) {

	in, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open match file")
	}
	defer in.Close()

	return LoadMatch(in)
}

func LoadMatchData(data []byte) (*Match, error) {
	buf := bytes.NewBuffer(data)
	return LoadMatch(buf)
}
