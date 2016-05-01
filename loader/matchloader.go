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
package loader

import (
	"bytes"
	"encoding/xml"
	"github.com/bboozzoo/q3stats/util"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type RawStat struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
}

type RawItem struct {
	Name    string `xml:"name,attr"`
	Pickups uint   `xml:"pickups,attr"`
	// item held time in ms
	Time uint `xml:"time,attr"`
}

type RawWeapon struct {
	Name  string `xml:"name,attr"`
	Hits  uint   `xml:"hits,attr"`
	Shots uint   `xml:"shots,attr"`
	Kills uint   `xml:"kills,attr"`
}

type RawPlayer struct {
	Name     string      `xml:"name,attr"`
	Stats    []RawStat   `xml:"stat"`
	Items    []RawItem   `xml:"items>item"`
	Weapons  []RawWeapon `xml:"weapons>weapon"`
	Powerups []RawItem   `xml:"powerups>item"`
}

type RawMatch struct {
	Datetime string      `xml:"datetime,attr"`
	Map      string      `xml:"map,attr"`
	Type     string      `xml:"type,attr"`
	Duration uint        `xml:"duration,attr"`
	Players  []RawPlayer `xml:"player"`
	DataHash string
}

func LoadMatch(src io.Reader) (*RawMatch, error) {
	raw, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load match data")
	}

	var match RawMatch
	err = xml.Unmarshal(raw, &match)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	match.DataHash = util.DataHash(raw)

	log.Printf("match: %+v", match)
	return &match, nil
}

func LoadMatchFile(path string) (*RawMatch, error) {

	in, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open match file")
	}
	defer in.Close()

	return LoadMatch(in)
}

func LoadMatchData(data []byte) (*RawMatch, error) {
	buf := bytes.NewBuffer(data)
	return LoadMatch(buf)
}
