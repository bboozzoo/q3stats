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
package site

import (
	"fmt"
	"github.com/bboozzoo/q3stats/models"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

func (s *Site) commonNewPlayerViewHandler(w http.ResponseWriter, req *http.Request, err string) {

	var data = struct {
		Aliases []models.Alias
		Error   string
	}{
		Error: err,
	}

	data.Aliases = s.p.ListUnclaimedAliases()
	s.loadRenderOrError(w, data, "newplayer.tmpl", "base.tmpl")
}

func (s *Site) newPlayerViewHandler(w http.ResponseWriter, req *http.Request) {
	s.commonNewPlayerViewHandler(w, req, "")
}

func (s *Site) createNewPlayerViewHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.commonNewPlayerViewHandler(w, req,
				fmt.Sprintf("%v", r))
		}
	}()

	if err := req.ParseForm(); err != nil {
		panic(errors.Wrap(err, "error processing form data"))
	}

	names := req.Form["name"]
	aliases := req.Form["aliases"]
	// expecting name to be 1 element list, aliases a list with
	// many aliases
	if len(names) != 1 || len(aliases) == 0 {
		panic(errors.New("empty name or no aliases selected"))
	}

	name := names[0]
	if len(name) == 0 {
		panic(errors.New("empty name"))
	}

	log.Printf("new player name: %s", name)
	log.Printf("claimed aliases: %s", aliases)

	id, err := s.p.CreateNewPlayer(name, aliases)
	if err != nil {
		panic(errors.Wrap(err, "error creating player"))
	}

	log.Printf("created player: %d", id)

	// redirect to player view
	http.Redirect(w, req, s.playerViewURL(id), http.StatusFound)
}
