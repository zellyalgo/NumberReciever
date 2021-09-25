package main

import (
    "os"
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