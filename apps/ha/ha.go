package main

import (
	"bufio"
	"errors"
	"espressif.com/chip/factory/api"
	"espressif.com/chip/factory/db"
	"espressif.com/chip/factory/rpc"
	"io"
	"log"
	"net"
)

func run(addr string) (err error) {
	defer db.Defer()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return err
		}
		go handle(conn)
	}
	return nil
}

func handle(conn net.Conn) error {
	defer func() {
		conn.Close()
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	var err error
	reader := bufio.NewReaderSize(conn, 1024*16)
	for {
		line, isPrefix, lerr := reader.ReadLine()
		if lerr != nil {
			err = lerr
			break
		}
		if isPrefix {
			err = errors.New("line too long")
			break
		}
		if len(line) <= 2 {
			break
		}
		req, resp := rpc.NewRequest(), rpc.NewResponse()
		err = api.Router.DispatchData(line, req, resp)
		if err != nil {
			resp.Status = 500
			resp.Err = err
		}
		conn.Write(append(resp.Json(), '\n'))
	}
	if err != nil && err != io.EOF {
		log.Println(err)
	}
	return err
}

func main() {
	log.Println("ha...ha")
	run(":6666")
}
