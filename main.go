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
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func doImport(c *cli.Context) {
	log.Printf("import")

	srcfile := c.Args().First()
	log.Fatal(runImport(srcfile))
}

func doDaemon(c *cli.Context) {
	log.Printf("starting daemon")
	log.Fatal(runDaemon())
}

func main() {
	app := cli.NewApp()
	app.Name = "q3stats"
	app.HideHelp = true
	app.HideVersion = true
	app.Commands = []cli.Command{
		{
			Name:   "import",
			Usage:  "import match logs",
			Action: doImport,
		},
		{
			Name:   "daemon",
			Usage:  "run daemon",
			Action: doDaemon,
		},
	}

	app.Run(os.Args)
}
