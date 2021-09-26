package main

import (
  "fmt"
	"log"
	"time"
)

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
    select {
    case number := <- n.numbers:
      n.addNumber(number)
    case <- time.After(10 * time.Second):
      n.summary()
    }
  }
}
func (n *Number) addNumber (number int) {
  if n.exist(number) {
    n.duplicate++
  } else {
    n.unique++
    log.Println(number)
    n.totalNumbers = append(n.totalNumbers, number)
  }
}
func (n *Number) summary () {
  fmt.Printf("Received %v unique numbers, %v duplicates. Unique total: %v\n", n.unique, n.duplicate, len(n.totalNumbers))
  n.unique = 0
  n.duplicate = 0
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