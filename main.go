package main

import (
	"github.com/aler9/gortsplib"
	"github.com/aler9/gortsplib/pkg/base"
	"github.com/aler9/gortsplib/pkg/url"
	"log"
	"os"
)

func main() {
	c := gortsplib.Client{
		OnRequest: func(request *base.Request) {
			if request.Method == base.Describe {
				request.Header["Require"] = base.HeaderValue{
					"www.onvif.org/ver20/backchannel",
				}
			}
		},
		OnResponse: func(response *base.Response) {
			log.Printf(response.String())
		},
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	tracks, _, _, err := c.Describe(u)
	if err != nil {
		panic(err)
	}

	log.Printf("available tracks: %v\n", tracks)
}
