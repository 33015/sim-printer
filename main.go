package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

type Printertype struct {
	Printertype string `json:"printertype"`
	Port        string `json:"port"`
	Name        string `json:"name"`
}
type Printers struct {
	Details []Printertype `json:"printers"`
}

func start_server(port, printertype string) string {
	log.Printf("Starting Server on port %s", port)
	var name string
	service := ":" + port
	jcount := 0
	name1 := printertype + port + "pdf"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// run as a goroutine
		jcount = jcount + 1
		name = name1 + strconv.Itoa(jcount) + ".pdf"
		switch printertype {
		case "raw":
			go handle_raw(conn, name)
		case "lpr":
			go handle_lpr(conn, name)

		}
		fmt.Println("after go routine starts")
	}
}

func handle_raw(conn net.Conn, filename string) {
	// close connection on exit
	defer conn.Close()

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(os.Stderr, "Client Fatal error: %s", err.Error())
		return
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(os.Stderr, "Client Fatal error: %s", err.Error())
		return
	}
	log.Printf("Saved file %s \n", filename)
}

func handle_lpr(conn net.Conn, filename string) {}

func load_config(filename string) Printers {
	config, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer config.Close()

	configdata, _ := ioutil.ReadAll(config)
	var printers Printers
	json.Unmarshal(configdata, &printers)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return printers
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	var printers = load_config("config.json")
	fmt.Println(len(printers.Details))
	for p := range printers.Details {
		fmt.Printf("p: %v\n", p)
	}
	for i := 0; i < len(printers.Details); i++ {
		go start_server(printers.Details[i].Port, printers.Details[i].Printertype)
	}
	for {
	}
}
