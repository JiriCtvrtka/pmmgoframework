package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	parseArgs()
}

func parseArgs() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	res := map[string]string{}
	for _, v := range args {
		v := strings.ReplaceAll(v, "--", "")
		// TODO: equal char can be in URL (query params)
		if strings.Contains(v, "=") {
			p := strings.Split(v, "=")
			res[p[0]] = p[1]
			continue
		}

		res[v] = ""
	}

	checkClient(res)
}

func checkClient(res map[string]string) {
	valid := []string{"client", "pmm2", "link-client"}
	tarballURL := "https://downloads.percona.com/downloads/TESTING/pmm/pmm2-client-2.31.0.tar.gz"
	for _, v := range valid {
		if _, ok := res[v]; !ok {
			continue
		}
		if res[v] == "" {
			continue
		}

		tarballURL = res[v]
	}

	resp, err := http.Get(tarballURL)
	if err != nil {
		fmt.Printf("invalid tarball URL: %s\n", tarballURL)
		os.Exit(1)
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(tarballURL)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	split := strings.Split(parsedURL.Path, "/")
	out, err := os.Create(split[len(split)-1])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
