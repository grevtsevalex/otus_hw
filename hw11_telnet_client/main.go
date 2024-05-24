package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	timeout   time.Duration
	ErrNoPort = errors.New("you must write a port")
	ErrNoHost = errors.New("you must write host")
)

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout for connection")
}

func main() {
	flag.Parse()
	switch len(flag.Args()) {
	case 0:
		fmt.Println(ErrNoHost)
		return
	case 1:
		fmt.Println(ErrNoPort)
		return
	default:
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)
	in := os.Stdin
	out := os.Stdout

	client := NewTelnetClient(address, timeout, in, out)

	err := client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect error: %v", err)
		return
	}
	defer client.Close()

	interruptCh := make(chan os.Signal, 1)

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "send error: %v", err)
			interruptCh <- os.Interrupt
		}

		os.Stderr.WriteString("\nEOF...")
		interruptCh <- os.Interrupt
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error from server: %v", err)
			interruptCh <- os.Interrupt
		}

		os.Stderr.WriteString("\nConnection was closed by peer...")
		interruptCh <- os.Interrupt
	}()

	signal.Notify(interruptCh, os.Interrupt)
	<-interruptCh
	os.Stderr.WriteString("\nBye-bye, client!")
}
