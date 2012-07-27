package main

import (
	"fmt"
	"log"
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
		log.Println(fmt.Sprintf("The file %s does not exist.", confPath))
		log.Fatalln(fmt.Sprintf("Example : echo \"%s:my_socks_pass\" > %s", login, confPath))
	}
	gos, err := gosoul.New(login, password)
	if err != nil {
		log.Fatalln("[ERROR caught] : ", err)
	}
	err = gos.Authenticate(gosoul.AUTHTYPE_MD5)
	if err != nil {
		log.Fatalln("[ERROR caught] : ", err)
	}
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		select {
		case <-sigChan:
			gos.Exit()
			log.Println("Thanks for using GoSoul, the NetSoul ident service writen in Go language !!!")
			os.Exit(0)
		}
	}()
	for {
		err = gos.Parse()
		if err != nil {
			gos.Exit() // Ensure socket close
			log.Println("[ERROR caught] : ", err)
			tryReconnect := 1
			for {
				if tryReconnect >= 100 {
					log.Fatalln("[ERROR caught] : Try to reconnect 100 times without success !!!")
				}
				time.Sleep(time.Second)
				log.Println("Try to reconnect : ", string(tryReconnect))
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
