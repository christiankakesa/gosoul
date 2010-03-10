package gosoul

import (
	"net"
	"os"
	"time"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"regexp"
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
// Private methods
func (gs *GoSoul) open(login string, password string) (err os.Error) {
	gs.login = login
	gs.password = password
	gs.connection, err = net.Dial("tcp", "", NSHOST+":"+NSPORT)
	if err != nil {
		return err
	}
	msg, _ := gs.Read()
	gs.salut = strings.Split(msg, " ", 0)
	return err
}

func (gs *GoSoul) md5Auth() string {
	res := fmt.Sprintf("%s-%s/%s%s", gs.salut[2], gs.salut[3], gs.salut[4], gs.password)
	md := md5.New()
	md.Write(strings.Bytes(res))
	authHashString := hex.EncodeToString(md.Sum())
	res = fmt.Sprintf("ext_user_log %s %s %s %s",
		gs.login,
		authHashString,
		GOSOUL_DATA,
		GOSOUL_LOCATION)
	return res
}

// Public methods
func (gs *GoSoul) Authenticate(authType string) (err os.Error) {
	if authType != AUTHTYPE_KRB {
		authType = AUTHTYPE_MD5
	}
	gs.Send("auth_ag ext_user none -")
	gs.Parse()
	switch authType {
	case AUTHTYPE_KRB:
		err = os.NewError("Kerberos authentication not yet implemented")
	case AUTHTYPE_MD5:
		gs.Send(gs.md5Auth())
	}
	msg, _ := gs.Read()
	if msg != "rep 002 -- cmd end" {
		err = os.NewError("Bad login or password")
	} else {
		gs.Send("user_cmd attach")
		myt := time.LocalTime()
		gs.Send(fmt.Sprintf("user_cmd state server:%d", myt.Seconds()))
	}
	return err
}

func (gs *GoSoul) Parse() (err os.Error) {
	res, err := gs.Read()
	// PING CMD
	if state, _ := regexp.MatchString("^ping.*", res); state {
		err = gs.Send(res)
	}
	return err
}

func (gs *GoSoul) Send(s string) (err os.Error) {
	fmt.Fprintf(os.Stdout, "[send:%s] : %s\n", time.LocalTime(), s) //DEBUG
	_, err = gs.connection.Write(strings.Bytes(s + "\n"))
	return err
}

func (gs *GoSoul) Read() (res string, err os.Error) {
	readBuffer := make([]byte, 2048)
	resLen, err := gs.connection.Read(readBuffer)
	if err == nil {
		res = string(readBuffer[0 : resLen-1])
	}
	if len(res) > 0 {
		fmt.Fprintf(os.Stdout, "[read:%s] : %s\n", time.LocalTime(), res)
	} //DEBUG
	return res, err
}

func (gs *GoSoul) Exit() {
	gs.Send("exit")
	gs.connection.Close()
}

func Connect(login string, password string) (gs *GoSoul, err os.Error) {
	gs = new(GoSoul)
	err = gs.open(login, password)
	return gs, err
}
