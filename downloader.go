package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
)

func main() {
	file, err := os.Open("./output/meta")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	line, err := lineCounter(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("line number", line)

	var wg sync.WaitGroup
	wg.Add(line)

	file, err = os.Open("./output/meta")
	scanner := bufio.NewScanner(file)
	// scan on comma? https://golang.org/src/bufio/example_test.go
	scanner.Split(bufio.ScanLines)

	// Scan.
	for scanner.Scan() {
		go func(url string) {
			defer wg.Done()

			fmt.Println("url is ", url)

			resp, err := http.Get(url)
			if err != nil {
				log.Println("err", err)
			} else if resp.StatusCode != 200 {
				log.Println("invalid", url)
			} else {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("ioutil.ReadAll -> %v", err)
				}

				resp.Body.Close()
				// You can now save it to disk or whatever...
				ioutil.WriteFile("./output/"+randStringBytes(20)+".jpg", bodyBytes, 0666)
			}
			// defer resp.Body.Close()
		}(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	fmt.Println("go")
	wg.Wait()
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
