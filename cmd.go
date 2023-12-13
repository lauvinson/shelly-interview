// Created by vinson on 2023/12/13.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type Stats struct {
	Clients          int64
	MessagesSent     int64
	BytesSent        int64
	MessagesReceived int64
	BytesReceived    int64
}

type Server struct {
	Addr  string
	Stats Stats
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	defer func(ln net.Listener) {
		err = ln.Close()
		if err != nil {
			fmt.Println("Error closing server:", err)
		}
	}(ln)

	fmt.Println("Listening on", s.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		s.Stats.Clients++

		go func(conn net.Conn) {
			defer func(conn net.Conn) {
				err = conn.Close()
				if err != nil {
					fmt.Println("Error closing connection:", err)
				}
			}(conn)

			s.handleConnection(conn)
		}(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		data, err := readLine(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
			} else {
				fmt.Println("Error reading data:", err)
			}
			break
		}

		s.Stats.MessagesReceived++
		s.Stats.BytesReceived += int64(len(data))

		fmt.Println("Received:", string(data))

		n, err := conn.Write([]byte(fmt.Sprintf("%s%s", data, "\n")))
		if err != nil {
			fmt.Println("Error writing data:", err)
			break
		}

		s.Stats.MessagesSent++
		s.Stats.BytesSent += int64(n)
	}
}

func readLine(conn net.Conn) ([]byte, error) {
	var buffer []byte

	for {
		b := make([]byte, 1)
		n, err := conn.Read(b)
		if err != nil {
			return nil, err
		}

		if n == 0 {
			return nil, io.EOF
		}

		buffer = append(buffer, b...)

		if b[0] == '\n' {
			break
		}
	}

	return buffer[:len(buffer)-1], nil
}

func handleStatsRequest(w http.ResponseWriter, s *Server) {
	w.Header().Set("Content-Type", "application/json")

	stats := s.Stats
	b, err := json.Marshal(&stats)
	if err != nil {
		fmt.Println("Error marshalling stats:", err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

func main() {

	// TCP listening on 8080
	tcpAddr := ":8080"

	s := &Server{
		Addr: tcpAddr,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// HTTP listening on 8081
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		handleStatsRequest(w, s)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
