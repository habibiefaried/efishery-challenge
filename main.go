package main

import (
	"bufio"
	"encoding/json"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
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
var dirname = "./fsdir"

func sanitizePath(v string) string {
	ret := strings.Replace(v, "..", "", -1)
	ret = strings.Replace(ret, "/", "", -1)
	return ret
}

func putKey(key string, value string) {
	mutex.Lock()
	skey := sanitizePath(key)
	if key == skey {
		err := ioutil.WriteFile(dirname+"/"+skey, []byte(value), 0644)
		if err != nil {
			log.Println(err)
		} else {
			envlist[key] = value
		}
	} else {
		log.Printf("Key %v contains invalid character\n", key)
	}
	mutex.Unlock()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		log.Printf("Input diterima: %v\n", s)
		switch command := strings.ToLower(s[0]); command {
		case "set":
			if len(s) == 3 {
				putKey(s[1], s[2])
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
					_ = os.Remove(dirname + "/" + sanitizePath(s[1]))
					conn.Write([]byte("key " + s[1] + " berhasil dihapus\n"))
				} else {
					conn.Write([]byte("key " + s[1] + " tidak ditemukan\n"))
				}
				mutex.Unlock()
			} else {
				conn.Write([]byte("Penggunaan: unset \"<key>\"\n"))
			}
		case "import":
			if len(s) == 3 {
				resp, err := http.Get(s[2])
				if err != nil {
					log.Println(err)
				} else {
					// Write the body to file
					body, _ := ioutil.ReadAll(resp.Body)
					tipe := strings.ToLower(s[1])

					if tipe == ".env" {
						data, err := godotenv.Unmarshal(string(body))
						if err != nil {
							log.Println(err)
							conn.Write([]byte("file .env tidak valid\n"))
						} else {
							for k, v := range data {
								putKey(k, v)
							}
							conn.Write([]byte("import berhasil\n"))
						}
					} else if tipe == "json" {
						jsonMap := make(map[string]string)
						err := json.Unmarshal(body, &jsonMap)
						if err != nil {
							log.Println(err)
							conn.Write([]byte("file json tidak valid\n"))
						} else {
							for k, v := range jsonMap {
								putKey(k, v)
							}
							conn.Write([]byte("import berhasil\n"))
						}
					} else if tipe == "yaml" {
						yamlMap := make(map[string]string)
						err := yaml.Unmarshal(body, &yamlMap)
						if err != nil {
							log.Println(err)
							conn.Write([]byte("file yaml tidak valid\n"))
						} else {
							for k, v := range yamlMap {
								putKey(k, v)
							}
							conn.Write([]byte("import berhasil\n"))
						}
					} else {
						conn.Write([]byte("Penggunaan: valid tipe -> .env/json/yaml \n"))
					}
				}
				defer resp.Body.Close()

			} else {
				conn.Write([]byte("Penggunaan: import \"<tipe:json/.env/yaml>\" \"<url>\" \n"))
			}
		default:
			conn.Write([]byte("Perintah tidak ditemukan\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		log.Println("Creating dir")
		errDir := os.MkdirAll(dirname, 0755)
		if errDir != nil {
			log.Println(err)
		}
	} else {
		log.Println("Directory exists, loading...")
		files, err := ioutil.ReadDir(dirname)
		if err != nil {
			log.Println(err)
		}

		for _, file := range files {
			content, err := ioutil.ReadFile(dirname + "/" + file.Name())
			if err != nil {
				log.Println(err)
			}

			envlist[file.Name()] = string(content)
		}
	}

	log.Println("Accepting connections....")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("Calling handleConnection")
		go handleConnection(conn)
	}
}
