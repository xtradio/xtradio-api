package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func telnet(command string) ([]string, error) {
	// connect to this socket

	var response []string

	liquidsoapHost := getEnv("LIQUIDSOAP_HOST")
	liquidsoapPort := getEnv("LIQUIDSOAP_PORT")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", liquidsoapHost, liquidsoapPort))

	if err != nil {
		log.Printf("Failed to connect to %s:%s, reason: %s", liquidsoapHost, liquidsoapPort, err)
		liqConnectionFailure.With(prometheus.Labels{"host": liquidsoapHost, "port": liquidsoapPort}).Inc()
		return response, err
	}
	defer conn.Close()
	// read in input from stdin

	// send to socket
	fmt.Fprintf(conn, command+"\n")
	// listen for reply
	data := bufio.NewScanner(conn)

	x := 0
	for data.Scan() {
		if x > 3 {
			break
		} else if strings.HasPrefix(data.Text(), "[playing]") {
			log.Println("Skipping: " + data.Text())
		} else if strings.HasPrefix(data.Text(), "[ready]") {
			song := strings.Replace(data.Text(), "[ready] ", "", -1)
			response = append(response, song)
		} else if data.Text() == "END" {
			break
		} else if data.Text() == "" {
			break
		} else {
			response = append(response, data.Text())
		}

		x = x + 1
	}

	return response, nil
}
