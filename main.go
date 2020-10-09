package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		fmt.Println(s)
		switch command := strings.ToLower(s[0]); command {
		case "set":
			fmt.Println("SET Key")
		case "get":
			fmt.Println("GET Key")
		case "list":
			fmt.Println("List keys")
		case "unset":
			fmt.Println("unset")
		case "download":
			fmt.Println("Download all env")
		default:
			conn.Write([]byte("Perintah tidak ditemukan\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Accept connection....")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Calling handleConnection")
		go handleConnection(conn)
	}
}
