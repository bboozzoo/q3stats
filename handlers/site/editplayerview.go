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
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

type editViewAlias struct {
	models.Alias
	Selected bool
}

func (s *Site) commonEditPlayerViewHandler(w http.ResponseWriter, req *http.Request, err string) {

	id := mux.Vars(req)["id"]
	pid := strToUint(id)

	log.Printf("edit player ID: %s", id)

	var data = struct {
		Aliases []editViewAlias
		Error   string
	}{
		Aliases: make([]editViewAlias, 0, 10),
		Error:   err,
	}

	// first aliases owned by current player
	for _, a := range s.p.ListPlayerAliases(pid) {
		log.Printf("owned: %s", a.Alias)
		data.Aliases = append(data.Aliases, editViewAlias{
			a,
			true,
		})
	}
	// then list unclaimed aliases
	for _, a := range s.p.ListUnclaimedAliases() {
		log.Printf("unowned: %s", a.Alias)
		data.Aliases = append(data.Aliases, editViewAlias{
			a,
			false,
		})
	}
	s.loadRenderOrError(w, data, "editplayer.tmpl", "base.tmpl")
}

func (s *Site) editPlayerViewHandler(w http.ResponseWriter, req *http.Request) {
	s.commonEditPlayerViewHandler(w, req, "")
}

func (s *Site) editSubmitPlayerViewHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	pid := strToUint(id)

	log.Printf("edit player ID: %s", id)

	defer func() {
		if r := recover(); r != nil {
			s.commonEditPlayerViewHandler(w, req,
				fmt.Sprintf("%v", r))
		}
	}()

	if err := req.ParseForm(); err != nil {
		panic(errors.Wrap(err, "error processing form data"))
	}

	aliases := req.Form["aliases"]
	// expecting name to be 1 element list, aliases a list with
	// many aliases
	if len(aliases) == 0 {
		panic(errors.New("no aliases selected"))
	}

	log.Printf("claimed aliases: %s", aliases)

	err := s.p.UpdatePlayerAliases(pid, aliases)
	if err != nil {
		panic(errors.Wrap(err, "error creating player"))
	}

	// redirect to player view
	http.Redirect(w, req, s.playerViewURL(pid), http.StatusFound)
}

func (s *Site) playerEditURL(id uint) string {
	u, _ := s.r.Get("playeredit").
		URL("id", strconv.FormatUint(uint64(id), 10))
	return u.String()
}
