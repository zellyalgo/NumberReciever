# NumberReciever
New Relic Code Challenge

This is a TPC Server wich allows at least 5 connetions to 4000 port.

### Rules
As user, you can enter:

* A 9 digit number: add this number to a list if that number is unique.
* A _terminated_ sequence, this will terminate the program for everyone.
* Otherwhise, the connection with that user will finish and nothing happens.

### How to Build
You can build the application just typing the following:

```
$ go build
```

this will create a NumberReciever bin that you can execute yo start up the server.

### How to Test
To check the test just write:

```
$ go test -v -cover
```

if you wanna get more information about coverage:

```
$ go test -v -coverprofile cover.out
$ go tool cover -func cover.out
```

These commands give you more detail information.

### Choises and Assumptions

First of all, I choose Golang just for simplicity and velocity, I have doubts about using Java with Vert.x, this kind of framework is really useful for event-driven applications. Finally, use Golang for the size of the application and because I want to learn more about the use of chans and goroutines.

I focus on Event-driven, in golang with chans, just because I think that is the most efficient way if you want to share information, in this case, between threads or users.

The schema of the application is really simple:

* __main.go__: Just run a Server and configure log file.
* __net.go__: The net.Listener implementation just to control the number of users.
* __numbers.go__: This is the numbers processor, the logic under unique numbers.
* __server.go__: Control the channels, handle the request and controls the flow.