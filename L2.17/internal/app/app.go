package app

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"telnet/internal/config"
)

func Run(cfg config.Config) error {
	address := net.JoinHostPort(cfg.Host, cfg.Port)
	var resstr string
	conn, err := net.DialTimeout("tcp", address, cfg.Timeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("connected: address: %s\n", address)
	defer func() {
		conn.Close()
		fmt.Printf("Connection closed by %s\n", resstr)
	}()

	doneWrite := make(chan error, 1)
	doneRead := make(chan error, 1)

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if tcpWtrite, ok := conn.(*net.TCPConn); ok {
			tcpWtrite.CloseWrite()
		}

		doneWrite <- err
	}()

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		doneRead <- err
	}()

	select {
	case err := <-doneRead:
		resstr = "server"
		// Сервер закрыл соединение
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err

	case err := <-doneWrite:
		resstr = "client"

		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		// Ввод закончен, но сервер ещё может прислать ответ
		err = <-doneRead
		resstr += " и дослушали сервер"

		if err == nil || errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
}
