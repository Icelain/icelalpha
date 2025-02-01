package cmd

import "flag"

type Flags struct {
	HttpPort uint
}

func getFlags() Flags {

	port := flag.Uint("Port", 8080, "Port to run http server on")
	flag.Parse()

	return Flags{HttpPort: *port}

}
