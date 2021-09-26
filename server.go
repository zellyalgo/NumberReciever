package main

import (
  "bufio"
  "errors"
  "fmt"
  "net"
  "os"
  "strconv"
  "time"
)

type Server struct{ 
  quit chan bool
  exit chan bool
  listen net.Listener
  number Numbers
}

func StartServer() {
  s := &Server{make(chan bool, maxConn), make (chan bool), nil, nil}
  s.Run()
}

func (s *Server) Run() {
  s.number = NewNumbers()
  go s.number.ProcessNumbers()
  s.createListener(connType, portConn)
  defer s.listen.Close()
  go s.acceptConnetions()
  <- s.exit 
  time.Sleep(time.Second)
}
func (s *Server) createListener (connType string, port string) {
  listen, err := net.Listen(connType, ":"+port)
  if err != nil {
    fmt.Println("Error listening:", err.Error())
    os.Exit(1)
  }
  s.listen = NewLimitedListener(maxConn, listen)
  fmt.Println("Listening on", port)
}
func (s *Server) acceptConnetions () {
  for {
    conn, err := s.listen.Accept()
    if err != nil {
      fmt.Println("Error accepting: ", err.Error())
      os.Exit(1)
    }
    go s.handleConnection(conn)
  }
}
func (s *Server) handleConnection (conn net.Conn) {
  defer conn.Close()
  input := make(chan string)
  errc := make (chan error)
  go s.readAsync(conn, input, errc)
  for {
    select {
    case message := <-input:
      err := s.processMessage(message)
      if err != nil {
        return
      }
    case <-s.quit:
      conn.Write([]byte("The Server is stopping...\n"))
      return
    case <- errc:
      return
    }
  }
}
func (s *Server) processMessage (message string) error {
  fmt.Println("Message Received:", message)
  if message == "terminated" {
    s.Stop()
    return nil
  }
  if len(message) != 9 {
    return errors.New("Message should have more than 9 digist")
  }
  number, err := strconv.Atoi(message)
  if err != nil {
    return err
  }
  s.number.Add(number)
  return nil
}
func (s *Server) readAsync(conn net.Conn, lines chan string, errc chan error) {
  defer conn.Close()
  scann := bufio.NewScanner(conn)
  for scann.Scan() {
    lines <- scann.Text()
  }
  errc <- scann.Err()
}
func (s *Server) Stop() {
  for i:= 0; i < maxConn; i++ {
    s.quit <- true
  }
  s.exit <- true
}