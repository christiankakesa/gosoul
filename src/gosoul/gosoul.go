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
func (gs *GoSoul) open(login string, password string) {
	var err os.Error
	gs.login = login
	gs.password = password
	gs.connection, err = net.Dial("tcp", "", NSHOST+":"+NSPORT)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection server: %s\n", err.String())
		os.Exit(1)
	}
	gs.salut = strings.Split(gs.Read(), " ", 0)
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
func (gs *GoSoul) Authenticate(authType string) {
	if authType != AUTHTYPE_KRB {
		authType = AUTHTYPE_MD5
	}
	gs.Send("auth_ag ext_user none -")
	gs.Parse()
	switch authType {
	case AUTHTYPE_KRB:
		fmt.Fprintf(os.Stderr, "Kerberos authentication not yet implemented\n")
		os.Exit(1)
	case AUTHTYPE_MD5:
		gs.Send(gs.md5Auth())
	}
	if gs.Read() != "rep 002 -- cmd end" {
		fmt.Fprintf(os.Stderr, "Bad login or password\n")
		os.Exit(1)
	} else {
		gs.Send("user_cmd attach")
		myt := time.LocalTime()
		gs.Send(fmt.Sprintf("user_cmd state server:%d", myt.Seconds()))
	}
}

func (gs *GoSoul) Parse() {
	res := gs.Read()
	// PING CMD
	if state, _ := regexp.MatchString("^ping.*", res); state {
		gs.Send(res)
	}
}

func (gs *GoSoul) Send(s string) {
	fmt.Fprintf(os.Stdout, "%s\n", s) //DEBUG
	gs.connection.Write(strings.Bytes(s + "\n"))
}

func (gs *GoSoul) Read() string {
	res := ""
	readBuffer := make([]byte, 2048)
	resLen, err := gs.connection.Read(readBuffer)
	if err == nil {
		res = string(readBuffer[0 : resLen-1])
	}
	if len(res) > 0 {
		fmt.Fprintf(os.Stdout, "[%s] : %s\n", time.LocalTime(), res)
	} //DEBUG
	return res
}

func (gs *GoSoul) Exit() {
	gs.Send("exit")
	gs.connection.Close()
}

func Connect(login string, password string) (gs *GoSoul) {
	gs = new(GoSoul)
	gs.open(login, password)
	return gs
}
