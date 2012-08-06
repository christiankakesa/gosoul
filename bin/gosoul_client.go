package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fenicks/gosoul"
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

var server string

func init() {
	flag.StringVar(&server, "server", "ns-server.epita.fr:4242", "NetSoul server host in this form : host:port, ip:port, [ipv6]:port")
}

func main() {
	flag.Parse()
	login, password, err := checkConfig()
	if err != nil {
		log.Println(fmt.Sprintf("The file %s does not exist", confPath))
		log.Fatalln(fmt.Sprintf(`Example: echo "%s:my_socks_pass" > %s`, login, confPath))
	}
	gos, err := gosoul.New(login, password, server)
	if err != nil {
		log.Fatalln("[ERROR caught]: ", err)
	}
	err = gos.Authenticate(gosoul.AUTHTYPE_MD5)
	if err != nil {
		log.Fatalln("[ERROR caught]: ", err)
	}
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		_ = <-sigChan
		gos.Exit()
		log.Println("Thanks for using GoSoul, the NetSoul ident service writen in Go language !!!")
		os.Exit(0)
	}()
	for {
		err = gos.Parse()
		if err != nil {
			log.Println("[ERROR caught]: ", err)
			gos.Exit() // Ensure socket close
			tryReconnect := 1
			for {
				if tryReconnect == 100 {
					log.Fatalln("[ERROR caught]: Try to reconnect 100 times without success !!!")
				}
				time.Sleep(1 * time.Second)
				log.Println("Try to reconnect to the NetSoul server: retry =>", strconv.Itoa(tryReconnect))
				gos, err = gosoul.New(login, password, server)
				if err != nil {
					log.Println("[ERROR caught]: ", err)
					tryReconnect += 1
					continue
				}
				err = gos.Authenticate(gosoul.AUTHTYPE_MD5)
				if err != nil {
					log.Println("[ERROR caught]: ", err)
					tryReconnect += 1
					continue
				}
				break
			}
		}
	}
}
