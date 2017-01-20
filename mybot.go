/*

mybot - Illustrative Slack bot in Go

Copyright (c) 2015 RapidLoop
Copyright (c) 2017 sndnvaps

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"encoding/csv"
	"fmt"
	//"log"
	"net/http"
	"os"
	"strings"
	"github.com/urfave/cli"
	"time"
)


var (
	Token string
)

func actionStartSlack(c *cli.Context) error {
	if c.String("token") != "" || c.String("t") != "" {
		slackRun(Token)
	} else {
		cli.ShowCommandHelp(c, "start")
	}
	return nil
}


func main() {
	app := cli.NewApp()
	app.Name = "aiicySlackBot"
	app.Usage = "Slack bot to get the stock infomation"
	app.Version = "0.5.0"
	app.Compiled = time.Now()
	app.Copyright = "Copyright (c) 2015 RapidLoop\n\t Copyright (c) 2017 sndnvaps<admin@sndnvaps.com>"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "RapidLoop",
			Email: "mdevan@gaia.local",
		},
	}
	
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start slack bot",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "token,t",
					Usage:       "myslack bot token",
					Destination: &Token,
				},
			},
			Action: actionStartSlack,
		},
	}
		
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
	
/*
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("mybot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 3 && parts[1] == "stock" {
				// looks good, get the quote and reply with the result
				go func(m Message) {
					m.Text = getQuote(parts[2])
					postMessage(ws, m)
				}(m)
				// NOTE: the Message object is copied, this is intentional
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				postMessage(ws, m)
			}
		}
	}
*/
}

// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getQuote(sym string) string {
	sym = strings.ToUpper(sym)
	url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if len(rows) >= 1 && len(rows[0]) == 5 {
		return fmt.Sprintf("%s (%s) is trading at $%s", rows[0][0], rows[0][1], rows[0][2])
	}
	return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}
