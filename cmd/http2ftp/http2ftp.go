package main

import (
	"github.com/jakobwesthoff/http2ftp"
	"github.com/yob/graval"

	"fmt"
	"log"
	"os"

	"flag"
	"path/filepath"
)

func showUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s [--configurations=] [--port=] [--hostname=]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	os.Exit(42)
}

/**
 * Run the thing :)
 */
func main() {
	flag.Usage = showUsageAndExit
	relativeConfigurationPath := flag.String("config", "config", "Path containing all JSON configuration files")
	serverPort := flag.Int("port", 3333, "Port to run the FTP server on")
	boundHostname := flag.String("hostname", "127.0.0.1", "Hostname/Interface to bind the the server to")
	flag.Parse()

	absoluteConfigurationDirectory, _ := filepath.Abs(*relativeConfigurationPath)

	configurations, err := http2ftp.LoadConfiguration(absoluteConfigurationDirectory)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting ftp server bound to %s:%d", *boundHostname, *serverPort)
	httpFactory := &http2ftp.HTTPDriverFactory{
		Configurations: configurations,
	}
	serverConfiguration := &graval.FTPServerOpts{
		Factory:  httpFactory,
		Port:     *serverPort,
		Hostname: *boundHostname,
	}
	server := graval.NewFTPServer(serverConfiguration)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
