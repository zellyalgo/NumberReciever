package main

import (
    "bufio"
    "errors"
    "fmt"
    "net"
    "os"
    "time"
    "sync"
    "strconv"
    "log"
)

const (
    portConn = "4000"
    connType = "tcp"
    maxConn = 5
    numbersLen = 1024
)

func main() {
  f := StartLogger()
  defer f.Close()
  StartServer()
}

func StartLogger () *os.File {
  f, err := os.OpenFile("numbers.log", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
  if err != nil {
    log.Fatalf("error opening file: %v", err)
  }
  log.SetFlags(0)
  log.SetOutput(f)
  return f
}

type Numbers interface {
  ProcessNumbers()
  Add(number int)
}
type Number struct {
  duplicate int
  unique int
  numbers chan int
  totalNumbers []int
}
func NewNumbers () *Number {
  return &Number{0, 0, make (chan int, numbersLen), make([]int, 0)}
}
func (n *Number) ProcessNumbers () {
  for {
    number := <- n.numbers
    if n.exist(number) {
      n.duplicate++
    } else {
      n.unique++
      log.Println(number)
      fmt.Println(number)
      n.totalNumbers = append(n.totalNumbers, number)
    }
  }
}
func (n *Number) exist(number int) bool {
    for _, item := range n.totalNumbers {
        if item == number {
            return true
        }
    }
    return false
}
func (n *Number) Add(number int) {
  n.numbers <- number
}

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
    c.sem <- true
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