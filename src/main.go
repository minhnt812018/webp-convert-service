package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func main() {
	const bindAddress = ":80"
	http.HandleFunc("/", requestHandler)
	fmt.Println("Http server listening on", bindAddress)
	_ = http.ListenAndServe(bindAddress, nil)
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	var imgUrl string
	query := request.URL.Query()
	for q := range query {
		key := q
		element := request.FormValue(key)
		if key == "url" {
			imgUrl = element
		} 
	}
	if imgUrl=="" {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	u, err := url.QueryUnescape(imgUrl)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		log.Println("imageUrl is error")
		return
	}
	resp, err := http.Get(strings.TrimRight(u, "/"))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	cmd := exec.Command("sh", "-c", "rm -rf /src/RGB.jpeg && convert -colorspace RGB -quality 100 - /src/RGB.jpeg && cwebp -mt -exact -quiet /src/RGB.jpeg -o -")
	cmd.Stdin = io.Reader(resp.Body)
	cmd.Stdout = response
	_ = cmd.Start()
	defer cmd.Wait()

	response.Header().Set("Content-Type", "image/webp")
	response.WriteHeader(http.StatusOK)
}
