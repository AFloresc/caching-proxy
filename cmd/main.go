package main

import (
	"caching-proxy/proxy"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	port := flag.Int("port", 0, "Port to run the proxy server on")
	origin := flag.String("origin", "", "Origin server URL")
	clear := flag.Bool("clear-cache", false, "Clear the cache")

	flag.Parse()

	if *clear {
		proxy.ClearCache()
		fmt.Println("Cache cleared.")
		return
	}

	if *port == 0 || *origin == "" {
		fmt.Println("Usage: caching-proxy --port <number> --origin <url>")
		os.Exit(1)
	}

	log.Printf("Starting proxy on port %d forwarding to %s\n", *port, *origin)
	proxy.StartServer(*port, *origin)
}
