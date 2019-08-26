package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Stat struct {
	Size int64
}

func main() {
	var source, destination string
	var sourceAuthToken, destinationAuthToken string
	var fromByte, toByte string

	flag.StringVar(&source, "source", "", "Source")
	flag.StringVar(&sourceAuthToken, "source-auth-token", "", "Source auth token")
	flag.StringVar(&destination, "destination", "", "Destination")
	flag.StringVar(&destinationAuthToken, "destination-auth-token", "", "Destination Auth Token")
	flag.StringVar(&fromByte, "from-byte", "", "Starting byte")
	flag.StringVar(&toByte, "to-byte", "", "Ending byte")
	flag.Parse()

	if source == "" || destination == "" {
		log.Fatal("Source, Destination options is required!")
	}
	var stat = *stat(&source, &sourceAuthToken)
	var from, to *int64

	if fromByte != "" {
		i, e := strconv.ParseInt(fromByte, 10, 64)
		if e != nil {
			log.Fatal(e)
		}
		from = &i
	}
	if toByte != "" {
		i, e := strconv.ParseInt(toByte, 10, 64)
		if e != nil {
			log.Fatal(e)
		}
		to = &i
	}

	var reader = *download(&source, &sourceAuthToken, from, to)
	defer func() {
		e := reader.Close()
		if e != nil {
			log.Fatal("Error close resource " + e.Error())
		}
	}()
	upload(&destination, &destinationAuthToken, &reader, stat.Size)
	log.Printf("success, replicated %d bytes", stat.Size)
}

func download(url *string, authToken *string, fromByte *int64, toByte *int64) *io.ReadCloser {
	client := &http.Client{}
	var req, err = http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if *authToken != "" {
		req.Header.Add("X-Auth-Token", *authToken)
	}
	if fromByte != nil && toByte != nil {
		var header = fmt.Sprintf("bytes=%d-%d", *fromByte, *toByte)
		log.Println("Apply range: " + header)
		req.Header.Add("Range", header)
	}
	log.Printf("starting downloadinig from %s", *url)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error download " + *url + " " + err.Error())
	}
	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		log.Fatal("Download: unexpected response code: " + strconv.Itoa(resp.StatusCode))
	}
	return &resp.Body
}

func upload(url *string, authToken *string, reader *io.ReadCloser, contentLength int64) {
	client := &http.Client{}
	var req, err = http.NewRequest(http.MethodPut, *url, *reader)
	if err != nil {
		log.Fatal(err)
	}
	if *authToken != "" {
		req.Header.Add("X-Auth-Token", *authToken)
	}
	req.ContentLength = contentLength
	var resp *http.Response
	log.Printf("starting uploading to %s", *url)
	resp, err = client.Do(req)
	if err != nil {
		// handle error
		log.Fatal("Error upload " + *url + " " + err.Error())
	}
	if resp.StatusCode != 201 {
		log.Fatal("Upload: unexpected response code: " + strconv.Itoa(resp.StatusCode))
	}
}

func stat(url *string, authToken *string) *Stat {
	client := &http.Client{}
	var req, err = http.NewRequest(http.MethodHead, *url, nil)
	if err != nil {
		log.Fatal("Error get stat " + *url + " " + err.Error())
	}
	if *authToken != "" {
		req.Header.Add("X-Auth-Token", *authToken)
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Stat: unexpected response code: " + strconv.Itoa(resp.StatusCode))
	}
	return &Stat{
		Size: resp.ContentLength,
	}
}
