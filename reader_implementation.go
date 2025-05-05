package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
)

type Reader struct {
	data io.Reader
}

var redis_client = CreateRedisClient()

const URL_Matcher = `(?i)^(http|https):\/\/([a-z0-9-]+\.)+[a-z]{2,}(:\d+)?(\/.*)?$`

func main() {

	osArgs := os.Args

	if len(os.Args) == 1 {
		log.Println("Please provide source as an argument")
		return
	}

	source := osArgs[1]

	if source == "" {
		log.Println("Please provide source as an argument")
		return
	}

	NewReader(source)

}

func NewReader(source string) (*Reader, error) {
	match, _ := regexp.Match(URL_Matcher, []byte(source))
	if match {
		log.Println("Fetching files from browser")
		file, err := getFileFromBrowser(source)

		if err != nil {
			log.Println("Error while fetching file from URL: ", err)
			return nil, err
		}

		if file.Header.Get("Content-Encoding") == "gzip" {
			bodyReader, err := gzip.NewReader(file.Body)
			if err != nil {
				log.Println("Error while creating gzip reader: ", err)
				return nil, err
			}
			return &Reader{data: bodyReader}, err
		}
		return &Reader{data: file.Body}, err
	} else {
		log.Println("Fetching files from local")
		file, err := getFileFromLocal(source)

		if err != nil {
			log.Println("Error file does not exist: ", err)
			return nil, err
		}

		return &Reader{data: file}, err
	}

}

func getFileFromBrowser(url string) (*http.Response, error) {

	redisData, err := redis_client.Get(ctx, url).Result()

	if err != nil {
		log.Println("Error while getting file from redis: ", err)
	} else {
		log.Println("File found in redis")

		readResp := bufio.NewReader(strings.NewReader(redisData))
		getReq, err := http.NewRequest("GET", url, nil)
		if err == nil {

			reconstructedResp, err := http.ReadResponse(readResp, getReq)

			if err != nil {
				log.Println("Error while reading response from redis: ", err)
			} else {
				return reconstructedResp, nil
			}
		}
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error while seting new HTTP request ", err)
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip")

	file, err := client.Do(req)

	if err != nil {
		log.Println("Error while fetching file from URL: ", err)
		return nil, err
	}

	dumpedFileResp, err := httputil.DumpResponse(file, true)

	if err != nil {
		log.Println("Error while fetching file from URL: ", err)
	} else {
		err = (redis_client.Set(ctx, url, string(dumpedFileResp), 0)).Err()
		if err != nil {
			log.Println("Error while setting file in redis: ", err)
		}
	}

	log.Println("Response Status Code: ", file.StatusCode)

	return file, nil
}

func getFileFromLocal(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Print("Error file does not exist: ", err)
		return nil, err
	}

	log.Println("File name: ", file.Name())
	return file, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.data == nil {
		log.Println("Invalid file")
		return 0, io.EOF
	}
	return r.data.Read(p)
}
