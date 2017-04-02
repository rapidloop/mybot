/*

mybot - Illustrative Slack bot in Go

Copyright (c) 2015 RapidLoop

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	sTeamInfo = `teamInfo:TBD`
	sBridge   = ` call 400-800-400`
	sHr       = ` call 400-800-401`
)

type WeatherInfoJson struct {
	Weatherinfo WeatherinfoObject
}

type WeatherinfoObject struct {
	City    string
	CityId  string
	Temp    string
	WD      string
	WS      string
	SD      string
	WSE     string
	Time    string
	IsRadar string
	Radar   string
	Rain    string
}

func main() {
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
			} else if (len(parts) == 3 && parts[1] == "weather" && parts[2] == "Guangzhou") || (len(parts) == 2 && parts[1] == "weather") {
				go func(m Message) {
					resp, err := http.Get("http://www.weather.com.cn/data/sk/101280101.html")
					if err != nil {
						log.Fatal(err)
					}

					defer resp.Body.Close()
					input, err := ioutil.ReadAll(resp.Body)

					var jsonWeather WeatherInfoJson
					json.Unmarshal(input, &jsonWeather)

					m.Text = fmt.Sprintf("weatherInfo: \n City:%s \n Wind:%s-%s\n Rain:%s\n Temp:%s\n Time:%s", jsonWeather.Weatherinfo.City, jsonWeather.Weatherinfo.WD, jsonWeather.Weatherinfo.WS, jsonWeather.Weatherinfo.Rain, jsonWeather.Weatherinfo.Temp, jsonWeather.Weatherinfo.Time)

					postMessage(ws, m)
				}(m)
			} else if len(parts) == 3 && parts[1] == "weather" && parts[2] == "Beijing" {
				go func(m Message) {
					resp, err := http.Get("http://www.weather.com.cn/data/sk/101010100.html")
					if err != nil {
						log.Fatal(err)
					}

					defer resp.Body.Close()
					input, err := ioutil.ReadAll(resp.Body)

					var jsonWeather WeatherInfoJson
					json.Unmarshal(input, &jsonWeather)

					m.Text = fmt.Sprintf("weatherInfo: \n City:%s \n Wind:%s-%s\n Rain:%s\n Temp:%s\n Time:%s", jsonWeather.Weatherinfo.City, jsonWeather.Weatherinfo.WD, jsonWeather.Weatherinfo.WS, jsonWeather.Weatherinfo.Rain, jsonWeather.Weatherinfo.Temp, jsonWeather.Weatherinfo.Time)

					postMessage(ws, m)
				}(m)
			} else if len(parts) == 2 && parts[1] == "team" {
				go func(m Message) {
					m.Text = sTeamInfo
					postMessage(ws, m)
				}(m)
			} else if len(parts) == 2 && parts[1] == "bridge" {
				go func(m Message) {
					m.Text = sBridge
					postMessage(ws, m)
				}(m)
			} else if len(parts) == 2 && parts[1] == "hr" {
				go func(m Message) {
					m.Text = sHr
					postMessage(ws, m)
				}(m)
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				postMessage(ws, m)
			}
		}
	}
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
	defer resp.Body.Close()
	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if len(rows) >= 1 && len(rows[0]) == 5 {
		return fmt.Sprintf("%s (%s) is trading at $%s", rows[0][0], rows[0][1], rows[0][2])
	}
	return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}
