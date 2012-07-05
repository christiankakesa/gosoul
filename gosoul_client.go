package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fenicks/gosoul/gosoul"
)

var (
	confDIR  = os.Getenv("HOME") + "/.config/gosoul"
	confPATH = confDIR + "/config.txt"
)

func checkConfig() (login, password string, err error) {
	f, err := os.Open(confPATH)
	if err != nil {
		defer func() { os.Mkdir(confDIR, 0755) }()
		return
	}
	content := make([]byte, 2048)
	l, err := f.Read(content)
	if err != nil {
		return
	}
	res := strings.SplitN(string(content[0:l-1]), ":", 2)
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
	gc, err := gosoul.Connect(login, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	err = gc.Authenticate(gosoul.AUTHTYPE_MD5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		select {
		case <-sigChan:
			gc.Exit()
			fmt.Fprintf(os.Stdout, "Thanks for using GO-Soul, the NetSoul ident service writen in GO language !!!\n")
			os.Exit(0)
		}
	}()
	for {
		err = gc.Parse()
		if err != nil {
			gc.Exit()
			fmt.Fprintf(os.Stderr, "[ERROR caught] : %v\n", err.Error())
			time.Sleep(time.Second)
			gc, _ = gosoul.Connect(login, password)
			gc.Authenticate(gosoul.AUTHTYPE_MD5)
		}
	}
}
