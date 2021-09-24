package main

import (
    "net"
    "testing"
    "bytes"
)

var srv Server

func init() {
    srv := &Server{make(chan bool, max_conn), make (chan bool), nil}

    go func() {
        srv.Run()
    }()
}

// Below init function
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

/*
func TestNETServer_Request(t *testing.T) {
    tt := []struct {
        test    string
        payload []byte
        want    []byte
    }{
        {
            "Sending a simple request returns result",
            []byte("hello world\n"),
            []byte("Request received: hello world"),
        },
        {
            "Sending another simple request works",
            []byte("goodbye world\n"),
            []byte("Request received: goodbye world"),
        },
    }

    for _, tc := range tt {
        t.Run(tc.test, func(t *testing.T) {
            conn, err := net.Dial("tcp", ":4000")
            if err != nil {
                t.Error("could not connect to TCP server: ", err)
            }
            defer conn.Close()

            if _, err := conn.Write(tc.payload); err != nil {
                t.Error("could not write payload to TCP server:", err)
            }

            out := make([]byte, 1024)
            if _, err := conn.Read(out); err == nil {
                if bytes.Compare(out, tc.want) == 0 {
                    t.Error("response did match expected output")
                }
            } else {
                t.Error("could not read from connection")
            }
        })
    }
}*/
