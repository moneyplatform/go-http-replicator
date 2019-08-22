package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Stat struct {
	Size int64
	Md5 string
}

func main() {
	var source, destination string
	var sourceAuthToken = ""
	var destinationAuthToken = ""

	flag.StringVar(&source, "source", "", "Source auth token")
	flag.StringVar(&sourceAuthToken, "source-auth-token", "", "Source")
	flag.StringVar(&destination, "destination", "", "Destination")
	flag.StringVar(&destinationAuthToken, "destination-auth-token", "", "Destination Auth Token")
	flag.Parse()

	if source == "" || destination == "" {
		log.Fatal("Source, Destination options is required!")
	}
	var stat = *stat(&source, &sourceAuthToken)
	var reader = *download(&source, &sourceAuthToken)
	defer func() {
		e := reader.Close()
		if e != nil {
			log.Fatal("Error close resource " + e.Error())
		}
	}()
	upload(&destination, &destinationAuthToken, &reader, stat.Size)
}

func download(url *string, authToken *string) *io.ReadCloser {
	client := &http.Client{
	}
	var req, err = http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		log.Fatal("Error download file " + *url + " " + err.Error())
	}
	req.Header.Set("X-Auth-Token", *authToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Download error: " + err.Error())
	}
	if resp.StatusCode != 200 {
		log.Fatal("Unexpected response code on download: " + strconv.Itoa(resp.StatusCode))
	}
	return &resp.Body
}

func upload(url *string, authToken *string, reader *io.ReadCloser, contentLength int64) {
	client := &http.Client{}
	var req, err = http.NewRequest(http.MethodPut, *url, *reader)
	if err != nil {
		log.Fatal("Error download file " + *url + " " + err.Error())
	}
	req.Header.Set("X-Auth-Token", *authToken)
	req.ContentLength = contentLength
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	if resp.StatusCode != 201 {
		log.Fatal("Unexpected response code on upload: " + strconv.Itoa(resp.StatusCode))
	}
}

func stat(url *string, authToken *string) *Stat {
	client := &http.Client{}
	var req, err = http.NewRequest(http.MethodHead, *url, nil)
	if err != nil {
		log.Fatal("Error download file " + *url + " " + err.Error())
	}
	req.Header.Set("X-Auth-Token", *authToken)
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Unexpected response code on stat: " + strconv.Itoa(resp.StatusCode))
	}
	return &Stat{
		Size: resp.ContentLength,
		Md5:  resp.Header.Get("etag"),
	}
}