package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	listen := flag.String("listen", "", "address to listen")
	upstreams := flag.String("upstreams", "", "upstreams separated by commas")
	timeout := flag.Int("timeout", 5, "timeout in seconds")
	flag.Parse()

	if *listen == "" || *upstreams == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	servers := strings.Split(*upstreams, ",")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(len(servers))

		for _, server := range servers {
			go func(server string) {
				defer wg.Done()

				err := replicate(server, r, body, *timeout)
				if err != nil {
					log.Printf("error replicating to %s: %s\n", server, err)
					return
				}

				log.Printf("successfully replicated to %s\n", server)
			}(server)
		}

		wg.Wait()

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("listening:", *listen)
	log.Fatal(http.ListenAndServe(*listen, mux))
}

func replicate(server string, r *http.Request, body []byte, timeout int) error {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	u, err := url.Parse(server)
	if err != nil {
		return err
	}

	u.Path = r.URL.Path

	req, err := http.NewRequest(r.Method, u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}

	for h, hv := range r.Header {
		for _, v := range hv {
			req.Header.Add(h, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
