package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var imgPath []string
var wg sync.WaitGroup

//LoadFile function
func LoadFile(filepath string) []string {
	path, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	text := string(path)
	imgPath = strings.Split(text, "\n")
	return imgPath
}

//전체주석처리 : Shitf + Alt + A
func main() {
	//"urls.txt"
	filepath := "urls.txt"

	urls := LoadFile(filepath)
	fmt.Println(len(urls))

	start := time.Now()
	file, err := os.OpenFile("error.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Open file is error")
	}
	div := len(urls) / 500
	fmt.Println(div)
	for i := 0; i < div; i++ {

		for _, url := range urls[500*i : 500*(i+1)] {
			tk := strings.Split(url, "/")
			fileName := "imgs/" + tk[len(tk)-1]

			go func(url string) {
				wg.Add(1)
				http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
				res, err := http.Get(url)
				if err != nil {
					fmt.Println("http get error: ", url, "-", err)
					if _, err := file.WriteString(url + "\n"); err != nil {
						fmt.Println(">>>Writing error!", url)
					}
				} else {
					output, err := os.Create(fileName)
					if err != nil {
						fmt.Println("Error while createing", fileName, "-", err)
					}
					_, err = io.Copy(output, res.Body)
					output.Close()
					res.Body.Close()

					if err != nil {
						fmt.Println("Error while Downloading", url, "-", err)
						if _, err := file.WriteString(url + "\n"); err != nil {
							fmt.Println(">>>Writing error!", url)
						}
					}
				}
				wg.Done()
			}(url)
		}
		wg.Wait()
		fmt.Println(i*500, " Done=========================")
	}
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}
