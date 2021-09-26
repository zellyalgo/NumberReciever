package main

import (
  "bytes"
  "net"
  "testing"
)

var srv Server
var mockNumber MockNumber

func init() {
  mockNumber = MockNumber{0}
  srv := &Server{make(chan bool, maxConn), make (chan bool), nil, &mockNumber}
  go srv.Run()
}

type MockNumber struct {
  numberAdded int
}
func (n *MockNumber) ProcessNumbers() {}
func (n *MockNumber) Add(number int)() {
  n.numberAdded = number
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

func TestCreate_listener (t *testing.T) {
  server := &Server{make(chan bool, maxConn), make (chan bool), nil, nil}
  server.createListener("tcp", "3500")
  if server.listen == nil {
    t.Error("server listener should not be nil", server.listen)
  }
  defer server.listen.Close()
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
  srvIngest := &Server{make(chan bool, maxConn), make (chan bool), nil, &mockNumber}
  if err := srvIngest.processMessage("109458278"); err != nil {
    t.Error("should not throw any error")
  }
  if mockNumber.numberAdded != 109458278 {
    t.Error("Excepted to add 109458278 but was", mockNumber.numberAdded)
  }
}