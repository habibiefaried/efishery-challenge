package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

var envlist = map[string]string{}
var mutex = &sync.RWMutex{}

func DownloadFile(url string) error {
	// Get the data
	s := strings.Split(url, "/")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(s[len(s)-1])
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	_, err = io.Copy(out, resp.Body)
	return err
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		fmt.Printf("Input diterima: %v\n", s)
		switch command := strings.ToLower(s[0]); command {
		case "set":
			if len(s) == 3 {
				mutex.Lock()
				envlist[s[1]] = s[2]
				mutex.Unlock()
				conn.Write([]byte("Penulisan key berhasil\n"))
			} else {
				conn.Write([]byte("Penggunaan: set \"<key>\" \"<value>\"\n"))
			}
		case "get":
			if len(s) == 2 {
				mutex.RLock()
				val, ok := envlist[s[1]]
				if ok {
					conn.Write([]byte(val + "\n"))
				} else {
					conn.Write([]byte("key " + s[1] + " tidak ditemukan\n"))
				}
				mutex.RUnlock()
			} else {
				conn.Write([]byte("Penggunaan: get \"<key>\"\n"))
			}
		case "list":
			mutex.RLock()
			for k, _ := range envlist {
				conn.Write([]byte(k + "\n"))
			}
			mutex.RUnlock()
		case "unset":
			if len(s) == 2 {
				mutex.Lock()
				_, ok := envlist[s[1]]
				if ok {
					delete(envlist, s[1])
					conn.Write([]byte("key " + s[1] + " berhasil dihapus\n"))
				} else {
					conn.Write([]byte("key " + s[1] + " tidak ditemukan\n"))
				}
				mutex.Unlock()
			} else {
				conn.Write([]byte("Penggunaan: unset \"<key>\"\n"))
			}
		case "import":
			fmt.Println("Download all env")
			if len(s) == 3 {
				err := DownloadFile(s[2])
				if err != nil {
					fmt.Println(err)
				}
			} else {
				conn.Write([]byte("Penggunaan: import \"<tipe:json/.env/yaml>\" \"<url>\" \n"))
			}
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
