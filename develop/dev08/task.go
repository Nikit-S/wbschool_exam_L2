package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/procfs"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func cd(c []string) {

	myDirr, err := os.Getwd()
	if err != nil {
		return
	}

	switch len(c) {
	case 1:
		home, err := os.UserHomeDir()
		if err != nil {
			return
		}
		os.Setenv("WB_PREVPWD", myDirr)
		os.Chdir(filepath.Join(home))
	case 2:
		switch c[1] {
		default:
			os.Setenv("WB_PREVPWD", myDirr)
			os.Chdir(filepath.Join(myDirr, c[1]))
		case "-":
			prevDir := os.Getenv("WB_PREVPWD")
			os.Setenv("WB_PREVPWD", myDirr)
			os.Chdir(filepath.Join(prevDir))
		}
	default:
		return
	}
	myDirr, _ = os.Getwd()
	os.Setenv("WB_PWD", myDirr)

}

func pwd(c []string) {
	myDirr, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println(myDirr)
}

func echo(c []string) {
	for _, e := range c[1:] {
		fmt.Print(e, " ")
	}
	fmt.Println()
}

func kill(c []string) {
	for _, e := range c[1:] {
		pid, err := strconv.Atoi(e)
		if err != nil {
			continue
		}
		pro, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		err = pro.Kill()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
	}
}

func ps(c []string) {
	processes, err := procfs.AllProcs()
	if err != nil {
		return
	}
	fmt.Printf("%s\t%s\n", "PID", "COMMAND")
	for _, proc := range processes {
		stat, err := proc.Stat()
		if err != nil {
			continue
		}
		fmt.Printf("%d\t%s\n", stat.PID, stat.Comm)
	}
}

var definedCommands map[string]func(c []string) = map[string]func(c []string){
	"cd":   cd,
	"pwd":  pwd,
	"echo": echo,
	"kill": kill,
	"ps":   ps,
}

func mainLoop() {

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		//сделали работу с командами без пайпов
		for _, comm := range strings.Split(reader.Text(), ";") {
			tokens := strings.Fields(comm)
			if len(tokens) == 0 {
				myDirr, err := os.Getwd()
				if err != nil {
					return
				}
				fmt.Print(myDirr, "$ ")
				continue
			}
			if f, ok := definedCommands[tokens[0]]; ok {
				f(tokens)
			} else if tokens[0] == "quit" {
				return
			} else {
				var cmd *exec.Cmd
				if len(tokens) > 1 {
					cmd = exec.Command(tokens[0], tokens[1:]...)
				} else {
					cmd = exec.Command(tokens[0])
				}
				res, _ := cmd.Output()
				fmt.Print(string(res))
			}
		}
		myDirr, err := os.Getwd()
		if err != nil {
			return
		}
		fmt.Print(myDirr, "$ ")
	}
}

func main() {
	myDirr, err := os.Getwd()
	if err != nil {
		return
	}
	os.Setenv("WB_PWD", myDirr)
	fmt.Print(myDirr, "$ ")
	mainLoop()
}
