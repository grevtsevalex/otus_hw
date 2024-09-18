package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"time"
)

var ErrConnectionFailed = errors.New("connection failed")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetCl struct {
	in      io.ReadCloser
	out     io.Writer
	timeout time.Duration
	address string
	conn    net.Conn
}

func (t *TelnetCl) Connect() error {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", t.address)
	if err != nil {
		return ErrConnectionFailed
	}

	t.conn = conn
	return nil
}

func readWrite(source io.Reader, dest io.Writer) error {
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		dest.Write([]byte(scanner.Text() + "\n"))
	}
	return scanner.Err()
}

func (t *TelnetCl) Send() error {
	return readWrite(t.in, t.conn)
}

func (t *TelnetCl) Receive() error {
	return readWrite(t.conn, t.out)
}

func (t *TelnetCl) Close() error {
	return t.conn.Close()
}

func NewTelnetClient(addr string, t time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetCl{address: addr, timeout: t, in: in, out: out}
}
