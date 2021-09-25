package main

import (
	"testing"
	"net"
	"fmt"
	"os"
	"time"
)

func TestGet_max_connections (t *testing.T) {
	limitedListenerLocal := NewLimitedListener(5, nil)
	if limitedListenerLocal.GetNumberConections() != 5 {
		t.Error("Expected 5 but found:", limitedListenerLocal.GetNumberConections())
	}
}

func TestNew_client_reduce_connections (t *testing.T) {
	listen, err := net.Listen(connType, ":3000")
	if err != nil {
	  fmt.Println("Error listening:", err.Error())
	  os.Exit(1)
	}
	limitedListenerLocal := NewLimitedListener(5, listen)
	go limitedListenerLocal.Accept()
	defer limitedListenerLocal.Close()
	time.Sleep(time.Millisecond)
	if limitedListenerLocal.GetNumberConections() != 4 {
		t.Error("Expected 4 but found:", limitedListenerLocal.GetNumberConections())
	}
}
func TestWhen_client_goes_release_connections (t *testing.T) {
	listen, err := net.Listen(connType, ":3000")
	if err != nil {
	  fmt.Println("Error listening:", err.Error())
	  os.Exit(1)
	}
	limitedListener := NewLimitedListener(5, listen)
	go func () {
		defer limitedListener.Close()
		for {
			conn, err := limitedListener.Accept()
			if err != nil {
				fmt.Println("problem listening", err)
			}
			conn.Close()
		}
	} ()

   	conn, err := net.Dial("tcp", "localhost:3000")
   	if err != nil {
   		t.Error("Can get connection to local Server", err)
   	}
	conn.Close()
	time.Sleep(time.Millisecond)
	if limitedListener.GetNumberConections() != 4 {
		t.Error("Expected 4 but found:", limitedListener.GetNumberConections())
	}
}