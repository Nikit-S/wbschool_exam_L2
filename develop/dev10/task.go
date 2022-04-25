package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func main() {
	timeout := flag.Int("timeout", 10, "Timeout")
	flag.Parse()
	args := flag.Args()
	conn, err := net.DialTimeout("tcp", args[0]+":"+args[1], time.Duration(*timeout)*time.Second)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	reader := bufio.NewScanner(os.Stdin)

	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Println("Connection is closed")
				os.Exit(0)
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Print(message)
		}
	}()

	for reader.Scan() {
		_, err := fmt.Fprintf(conn, reader.Text()+" / HTTP/1.0\r\n\r\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
