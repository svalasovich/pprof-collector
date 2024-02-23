package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	_ "unsafe"

	"github.com/google/pprof/profile"
)

const output = "result.pprof"

func main() {
	// http.DefaultTransport.(*http.Transport).DisableCompression = true
	seconds, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	addresses := os.Args[2:]
	fmt.Printf("ADDRESSES: %s\n", addresses)
	profiles := make([]*profile.Profile, 0, len(addresses))

	fmt.Println("START collect data")
	var wg sync.WaitGroup
	for _, address := range addresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()

			get, err := http.DefaultClient.Get(fmt.Sprintf("%s?seconds=%d", strings.TrimSpace(address), seconds))
			if err != nil {
				panic(err)
			}
			defer get.Body.Close()

			p, err := profile.Parse(get.Body)
			if err != nil {
				panic(err)
			}

			if err := p.CheckValid(); err != nil {
				panic(err)
			}

			profiles = append(profiles, p)
		}(address)
	}

	wg.Wait()
	fmt.Println("DONE collect data")

	p, err := profile.Merge(profiles)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Cannot open output to write: %v", err)
	}

	if err := out.Truncate(0); err != nil {
		panic(err)
	}
	if _, err := out.Seek(0, 0); err != nil {
		panic(err)
	}

	if err := p.Write(out); err != nil {
		log.Fatalf("Cannot write merged profile to file: %v", err)
	}

	if err := out.Close(); err != nil {
		log.Printf("Error when closing the output file: %v", err)
	}
}
