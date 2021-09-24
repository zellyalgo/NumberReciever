package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "time"
    "sync"
)

const (
    port_conn = "4000"
    conn_type = "tcp"
    max_conn = 5
)

func main() {
  StartServer()
}

type Server struct{ 
  quit chan bool
  exit chan bool
  listen net.Listener
}

func StartServer() {
    s := &Server{make(chan bool, max_conn), make (chan bool), nil}
    s.Run()
}

func (s *Server) Run() {
    s.createListener(conn_type, port_conn)
    defer s.listen.Close()
    go s.acceptConnetions()
    <- s.exit 
    time.Sleep(time.Second)
}
func (s *Server) createListener (conn_type string, port string) {
  listen, err := net.Listen(conn_type, ":"+port)
  if err != nil {
      fmt.Println("Error listening:", err.Error())
      os.Exit(1)
  }
  s.listen = NewLimitedListener(max_conn, listen)
  fmt.Println("Listening on", port_conn)
}
func (s *Server) acceptConnetions () {
  for {
    conn, err := s.listen.Accept()
    if err != nil {
        fmt.Println("Error accepting: ", err.Error())
        os.Exit(1)
    }
    // Handle connections in a new goroutine.
    go s.handleConnection(conn)
  }
}
func (s *Server) handleConnection (conn net.Conn) {
    defer conn.Close()
    input := make(chan string) // input stream
    errc := make (chan error)
    go s.readAsync(conn, input, errc) // go and read
    for {
        select { // read or close
        case message := <-input:
          s.processMessage(message)
        case <-s.quit:
          conn.Write([]byte("The Server is stopping...\n"))
          return
        case <- errc:
          return
        }
    }
}
func (s *Server) processMessage (message string) {
  fmt.Println("Message Received:", message)
  if message == "terminated" {
    s.Stop()
  }
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
    for i:= 0; i < max_conn; i++ {
      s.quit <- true
    }
    s.exit <- true
}


type LimitedListener struct {
    sync.Mutex
    net.Listener
    sem chan bool
}
func NewLimitedListener(count int, l net.Listener) *LimitedListener {
    sem := make(chan bool, count)
    for i := 0; i < count; i++ {
        sem <- true
    }
    return &LimitedListener{
        Listener: l,
        sem:      sem,
    }
}
func (l *LimitedListener) Addr() net.Addr { 
  return l.Listener.Addr()
}
func (l *LimitedListener) Close() error {
  return l.Listener.Close()
}
func (l *LimitedListener) Accept() (net.Conn, error) {
    <-l.sem
    c, err := l.Listener.Accept()
    if err!=nil {
      return nil, err
    }
    return NewLimitedConn(l.sem, c), nil
}
func (l *LimitedListener) GetNumberConections () int {
  return len(l.sem)
}

type LimitedConn struct { 
  sem chan bool
  net.Conn
}
func NewLimitedConn (sem chan bool, c net.Conn) *LimitedConn {
  return &LimitedConn{sem, c}
}
func (c *LimitedConn) LocalAddr() net.Addr {
  return c.Conn.LocalAddr()
}
func (c *LimitedConn) RemoteAddr() net.Addr {
  return c.Conn.RemoteAddr()
}
func (c *LimitedConn) Read(b []byte) (n int, err error) {
  return c.Conn.Read(b)
}
func (c *LimitedConn) Close() error {
    c.sem <- true // release
    return c.Conn.Close()
}
func (c *LimitedConn) Write(b []byte) (int, error) {
  return c.Conn.Write(b)
}
func (c *LimitedConn) SetDeadline(t time.Time) error {
  return c.Conn.SetDeadline(t)
}
func (c *LimitedConn) SetReadDeadline(t time.Time) error {
  return c.Conn.SetReadDeadline(t)
}
func (c *LimitedConn) SetWriteDeadline(t time.Time) error {
  return c.Conn.SetWriteDeadline(t)
}