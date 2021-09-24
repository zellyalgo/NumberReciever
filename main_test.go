package main

import (
    "net"
    "testing"
    "bytes"
    "log"
    "os"

    "io/ioutil"
)

var srv Server

func init() {
    srv := &Server{make(chan bool, max_conn), make (chan bool), make (chan int, numbersLen), nil}

    go func() {
        srv.Run()
    }()
}

func TestServer_Run(t *testing.T) {
    conn, err := net.Dial("tcp", "localhost:4000")
    if err != nil {
        t.Error("could not connect to server: ", err)
    }
    defer conn.Close()
}

func TestSend_terminated_message(t *testing.T) {
    conn, err := net.Dial("tcp", "localhost:4000")
    if err != nil {
        t.Error("could not connect to server: ", err)
    }
    conn.Write([]byte("terminated\n"))
    out := make([]byte, 1024)
    expected := []byte("The Server is stopping...\n")
    if _, err := conn.Read(out); err == nil {
        if bytes.Compare(out, expected) == 0 {
            t.Error("response did match expected output")
        }
    } else {
        t.Error("could not read from connection")
    }

}

func TestLess_than_9_Digist (t *testing.T) {
    if err := srv.processMessage("1094"); err == nil {
        t.Error("should throw error when less than 9 digist")
    }
}
func TestMore_than_9_Digist (t *testing.T) {
    if err := srv.processMessage("1094582768"); err == nil {
        t.Error("should throw error when more than 9 digist")
    }
}
func Test9_Digist (t *testing.T) {
    if err := srv.processMessage("109458278"); err != nil {
        t.Error("should not throw any error")
    }
}

func TestStart_logger (t *testing.T) {
    f := StartLogger()
    log.Println("Just write something")
    f.Close()
    bodyByte, _ := ioutil.ReadFile("numbers.log")
    body := string(bodyByte)
    if body != "Just write something\n" {
        t.Error("Expected = [Just write something], actual =", body)
    }
}
func TestOverride_file (t *testing.T) {
    f, err := os.Create("numbers.log")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    _, _ = f.WriteString("Just Something to Override\n")
    f.Close()
    f = StartLogger()
    f.Close()
    bodyByte, _ := ioutil.ReadFile("numbers.log")
    body := string(bodyByte)
    if body != "" {
        t.Error("Expected = [], actual =", body)
    }
}
/*
func TestUpTo5Connections(t *testing.T) {
    clients := make([]net.Conn, 0, 5)
    for i := 0; i <7; i++ {
        conn, err := net.Dial("tcp", "localhost:4000")
        if err != nil {
            t.Error("could not connect to server: ", err)
        } else {
            clients = append(clients, conn)
        }
    }
    for _, conn := range clients {
        conn.Close()
    }
}*/