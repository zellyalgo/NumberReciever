package main

import (
  "net"
	"sync"
	"time"
)

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