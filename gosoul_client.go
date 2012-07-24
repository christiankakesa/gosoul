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
	confDir  = os.Getenv("HOME") + "/.config/gosoul"
	confPath = confDir + "/config.txt"
)

func checkConfig() (login, password string, err error) {
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		os.Mkdir(confDir, 0755)
	}
	f, err := os.Open(confPath)
	if err != nil {
		return "", "", err
	}
	content := make([]byte, 2048)
	l, err := f.Read(content)
	if err != nil {
		return "", "", err
	}
	res := strings.SplitN(string(content[0:l-1]), ":", 2)
	login = res[0]
	password = res[1]
	return login, password, nil
}

func main() {
	login, password, err := checkConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "The file %s don't exist.\n", confPath)
		fmt.Fprintf(os.Stderr, "Create it and put one line with \"login_l:socks_pass\"\n")
		fmt.Fprintf(os.Stderr, "Example : echo \"login_l:socks_pass\" > %s\n", confPath)
		os.Exit(1)
	}
	gos, err := gosoul.New(login, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}
	err = gos.Authenticate(gosoul.AUTHTYPE_MD5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		select {
		case <-sigChan:
			gos.Exit()
			fmt.Fprintf(os.Stdout, "Thanks for using GO-Soul, the NetSoul ident service writen in GO language !!!\n")
			os.Exit(0)
		}
	}()
	for {
		err = gos.Parse()
		if err != nil {
			gos.Exit()
			fmt.Fprintf(os.Stderr, "[ERROR caught] : %v\n", err.Error())
			tryReconnect := 0
			for {
				if tryReconnect > 100 {
					fmt.Fprintf(os.Stderr, "Try to reconnect 100 times without success !!!\n")
					os.Exit(0)
				}
				time.Sleep(time.Second)
				gos, err = gosoul.New(login, password)
				if err != nil {
					tryReconnect += 1
					continue
				}
				gos.Authenticate(gosoul.AUTHTYPE_MD5)
				break
			}
		}
	}
}
