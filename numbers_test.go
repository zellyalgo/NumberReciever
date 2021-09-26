package main

import (
  "testing"
  "time"
)

func TestAdd_one_number (t *testing.T) {
  number := *NewNumbers()
  number.addNumber(123)
  if number.totalNumbers[0] != 123 {
    t.Error("Expected 123 but found:", number.totalNumbers[0])
  }
  if number.unique != 1 {
    t.Error("Expected 1 but found:", number.unique)	
  }
  if number.duplicate != 0 {
    t.Error("Expected 0 but found:", number.duplicate)	
  }
}
func TestAdd_two_unique_numbers (t *testing.T) {
  number := *NewNumbers()
  number.addNumber(123)
  number.addNumber(1234)
  if len(number.totalNumbers) != 2 {
    t.Error("Expected 2 but found:", len(number.totalNumbers))
  }
  if number.unique != 2 {
    t.Error("Expected 2 but found:", number.unique)	
  }
  if number.duplicate != 0 {
    t.Error("Expected 0 but found:", number.duplicate)	
  }
}
func TestAdd_two_unique_and_one_duplcated_numbers (t *testing.T) {
  number := *NewNumbers()
  number.addNumber(123)
  number.addNumber(1234)
  number.addNumber(123)
  if len(number.totalNumbers) != 2 {
    t.Error("Expected 2 but found:", len(number.totalNumbers))
  }
  if number.unique != 2 {
    t.Error("Expected 2 but found:", number.unique)	
  }
  if number.duplicate != 1 {
    t.Error("Expected 1 but found:", number.duplicate)	
  }
}

func TestSumary_should_reset_counters (t *testing.T) {
  number := *NewNumbers()
  number.addNumber(123)
  number.addNumber(1234)
  number.summary()
  if number.unique != 0 {
    t.Error("Expected 0 but found:", number.unique)	
  }
  if number.duplicate != 0 {
    t.Error("Expected 0 but found:", number.duplicate)	
  }
}

func TestAdd_using_chans (t *testing.T) {
  number := *NewNumbers()
  number.Add(123)
  num := <- number.numbers
  if num != 123 {
    t.Error("Expected 1 but found:", num)	
  }
}

func TestProcess_Numbers (t *testing.T) {
  number := *NewNumbers()
  number.numbers <- 123
  go number.ProcessNumbers()
  time.Sleep(time.Millisecond)
  if number.unique != 1 {
    t.Error("Expected 1 but found:", number.unique)	
  }
}