package main

import (
  "testing"
  "log"
  "os"

  "io/ioutil"
)

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
  f.WriteString("Just Something to Override\n")
  f.Close()
  f = StartLogger()
  f.Close()
  bodyByte, _ := ioutil.ReadFile("numbers.log")
  body := string(bodyByte)
  if body != "" {
    t.Error("Expected = [], actual =", body)
  }
}