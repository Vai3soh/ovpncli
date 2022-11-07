/*
* ovpncli -- Library for wrapping openvpn3 (https://github.com/OpenVPN/openvpn3) functionality in go way.
* Copyright (C) 2022 Vai3soh
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU Affero General Public License Version 3
* as published by the Free Software Foundation.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU Affero General Public License for more details.

* You should have received a copy of the GNU Affero General Public License
* along with this program in the COPYING file.
* If not, see <http://www.gnu.org/licenses/>.
 */
package ovpncli

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func getResponse(client http.Client, url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func getConfigUrl(n *html.Node) string {

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {

				if strings.Contains(a.Val, ".ovpn") {
					link := "https://ipspeed.info/" + a.Val
					return link
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r := getConfigUrl(c)
		if r != "" {
			return r
		}
	}
	return ""
}

type overwriteClient struct {
	ClientAPI_OpenVPNClient
}

func (ocl *overwriteClient) Log(arg2 ClientAPI_LogInfo) {
	log.Printf("log: %s", arg2.GetText())

}

func (ocl *overwriteClient) Event(arg2 ClientAPI_Event) {
	log.Printf("event name: %s", arg2.GetName())
	log.Printf("event info: %s", arg2.GetInfo())
}

func (ocl *overwriteClient) Remote_override_enabled() {

}

func (ocl *overwriteClient) Socket_protect() {

}

func getProfile() string {
	client := http.Client{
		Timeout: 6 * time.Second,
	}
	resp := getResponse(client, `https://ipspeed.info/freevpn_openvpn.php?language=en`)

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(strings.NewReader(string(bytes)))
	if err != nil {
		log.Fatal(err)
	}

	url := getConfigUrl(doc)
	resp = getResponse(client, url)
	defer resp.Body.Close()
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)

}

func TestConnection(t *testing.T) {
	profile := getProfile()
	cfg := NewClientConfig(
		WithConfig(profile),
		WithSslDebugLevel(0),
		WithCompressionMode("yes"),
		WithDisableClientCert(true),
		WithLegacyAlgorithms(true),
		WithNonPreferredDCAlgorithms(true),
		WithTunPersist(true),
	)
	creds := NewClientCreds(
		WithPassword(""),
		WithUsername(""),
	)

	ocl := &overwriteClient{}
	cliObj := NewClient(ocl)
	ocl.ClientAPI_OpenVPNClient = cliObj
	ev := cliObj.Eval_config(cfg)
	if ev.GetError() {
		log.Fatalf("err: config eval failed [%s]\n", ev.GetMessage())
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func(sec int) {
		time.Sleep(time.Duration(sec) * time.Millisecond)
		cancel()
	}(20000)

	status := cliObj.Provide_creds(creds)
	if status.GetError() {
		log.Fatal(status.GetMessage())
	}
	cliObj.StartConnection(ctx)
	err := cliObj.CallbackError()
	if err != nil {
		log.Println(err)
	}
	defer DeleteClient(cliObj)
}
