package gosoul

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	NSHOST          = "ns-server.epita.fr"
	NSPORT          = "4242"
	GOSOUL_DATA     = "GO-Soul,%20by%20Christian%20KAKESA"
	GOSOUL_LOCATION = "Gosoul@HOME"
	AUTHTYPE_KRB    = "kerberos"
	AUTHTYPE_MD5    = "md5"
)

type GoSoul struct {
	connection net.Conn
	login      string
	password   string
	salut      []string
}

func (gs *GoSoul) open(login string, password string) (err error) {
	gs.login = login
	gs.password = password
	gs.connection, err = net.Dial("tcp", NSHOST+":"+NSPORT)
	if err != nil {
		return err
	}
	msg, _ := gs.Read()
	gs.salut = strings.Split(msg, " ")
	return err
}

func (gs *GoSoul) md5Auth() string {
	res := fmt.Sprintf("%s-%s/%s%s", gs.salut[2], gs.salut[3], gs.salut[4], gs.password)
	md := md5.New()
	md.Write([]byte(res))
	authHashString := hex.EncodeToString(md.Sum(nil))
	res = fmt.Sprintf("ext_user_log %s %s %s %s",
		gs.login,
		authHashString,
		GOSOUL_DATA,
		GOSOUL_LOCATION)
	return res
}

func (gs *GoSoul) Authenticate(authType string) (err error) {
	if authType != AUTHTYPE_KRB {
		authType = AUTHTYPE_MD5
	}
	gs.Send("auth_ag ext_user none -")
	gs.Parse()
	switch authType {
	case AUTHTYPE_KRB:
		err = errors.New("Kerberos authentication not yet implemented")
	case AUTHTYPE_MD5:
		gs.Send(gs.md5Auth())
	}
	msg, _ := gs.Read()
	if msg != "rep 002 -- cmd end" {
		err = errors.New("Bad login or password")
	} else {
		gs.Send("user_cmd attach")
		myt := time.Now()
		gs.Send(fmt.Sprintf("user_cmd state server:%d", myt.Unix()))
	}
	return err
}

func (gs *GoSoul) Parse() (err error) {
	res, err := gs.Read()
	// PING CMD
	if state, _ := regexp.MatchString("^ping.*", res); state {
		err = gs.Send(res)
	}
	return err
}

func (gs *GoSoul) Send(s string) (err error) {
	fmt.Fprintf(os.Stdout, "[send:%s] : %s\n", time.Now(), s) //DEBUG
	_, err = gs.connection.Write([]byte(s + "\n"))
	return err
}

func (gs *GoSoul) Read() (res string, err error) {
	readBuffer := make([]byte, 2048)
	resLen, err := gs.connection.Read(readBuffer)
	if err == nil {
		res = string(readBuffer[0 : resLen-1])
	}
	if len(res) > 0 {
		fmt.Fprintf(os.Stdout, "[read:%s] : %s\n", time.Now(), res)
	}
	return res, err
}

func (gs *GoSoul) Exit() {
	gs.Send("exit")
	gs.connection.Close()
}

func Connect(login string, password string) (gs *GoSoul, err error) {
	gs = new(GoSoul)
	err = gs.open(login, password)
	return gs, err
}
