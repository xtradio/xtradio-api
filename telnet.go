package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func telnet(command string) []string {
	// connect to this socket

	liquidsoapHost := getEnv("LIQUIDSOAP_HOST")
	liquidsoapPort := getEnv("LIQUIDSOAP_PORT")

	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%s", liquidsoapHost, liquidsoapPort))
	defer conn.Close()
	// read in input from stdin

	// send to socket
	fmt.Fprintf(conn, command+"\n")
	// listen for reply
	data := bufio.NewScanner(conn)

	var response []string

	x := 0
	for data.Scan() {
		if x > 3 {
			break
		} else if strings.HasPrefix(data.Text(), "[playing]") {
			fmt.Println("Skipping: " + data.Text())
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

	return response
}
