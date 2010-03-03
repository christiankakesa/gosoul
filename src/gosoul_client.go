package main

import (
	"exp/exception"
	"fmt"
	"os"
	//	"os/signal"
	"strings"
	//	"syscall"
	"./gosoul"
)

var (
	confDIR  = os.Getenv("HOME") + "/.config/gosoul"
	confPATH = confDIR + "/config.txt"
)

func checkConfig() (login, password string, err os.Error) {
	f, err := os.Open(confPATH, os.O_RDONLY, 0600)
	if err != nil {
		defer func() { os.Mkdir(confDIR, 0755) }()
		return
	}
	content := make([]byte, 2048)
	l, err := f.Read(content)
	if err != nil {
		return
	}
	res := strings.Split(string(content[0:l-1]), ":", 0)
	login = res[0]
	password = res[1]

	return login, password, err
}

func main() {
	login, password, err := checkConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "The file %s don't exist.\n", confPATH)
		fmt.Fprintf(os.Stderr, "Create it and put one line with \"login_l:socks_pass\"\n")
		fmt.Fprintf(os.Stderr, "Example : echo \"login_l:socks_pass\" > %s\n", confPATH)
		os.Exit(1)
	}
	gc := gosoul.Connect(login, password)
	gc.Authenticate(gosoul.AUTHTYPE_MD5)

	//sig := (<-signal.Incoming).(signal.UnixSignal)
	//switch sig {
	//case syscall.SIGINT:
	//	gc.Exit()
	//	println("SIGHUP SIGINT")
	//	os.Exit(1)
	//}

	for true {
		exception.Try(func(throw exception.Handler) { gc.Parse() }).Catch(func(x interface{}) {
			go gc.Exit()
			fmt.Fprintf(os.Stderr, "%v\n", x)
			fmt.Fprintf(os.Stdout, "Thanks for using GO-Soul, the NetSoul ident service writen in GO language !!!\n")
		})
	}
	os.Exit(0)
}
