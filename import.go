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
	"fmt"
	"github.com/bboozzoo/q3stats/handlers/api"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type MatchHash string

const (
	emptyHash = ""
)

func readMatchFile(srcfile string) (*os.File, error) {
	in, err := os.Open(srcfile)
	if err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf("failed to open match file %s", srcfile))
	}

	return in, nil
}

func sendMatchData(src io.Reader, addr string) (MatchHash, error) {
	cl := http.Client{}

	url := fmt.Sprintf("http://%s%s%s", addr, uriApi, api.UriAddMatch)

	log.Printf("posting to URL: %s", url)

	resp, err := cl.Post(url, "application/vnd.q3-match-stats", src)
	if err != nil {
		return emptyHash, errors.Wrap(err, "posting match data failed")
	}

	log.Printf("response: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return emptyHash, fmt.Errorf("request failed with status: %d",
			resp.StatusCode)
	}

	matchhash, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return emptyHash, errors.Wrap(err, "failed to receive response")
	}

	return MatchHash(matchhash), nil
}

func runImport(srcfile string, addr string) (MatchHash, error) {

	if srcfile == "" {
		return emptyHash, errors.New("no path to match file")
	}

	indata, err := readMatchFile(srcfile)
	if err != nil {
		return emptyHash, err
	}
	// close input file
	defer indata.Close()

	matchhash, err := sendMatchData(indata, addr)
	if err != nil {
		return emptyHash, errors.Wrap(err, "failed to send match data")
	}

	return matchhash, nil
}
