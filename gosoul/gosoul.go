// Copyright 2012 Christian Kakesa. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gosoul provides functions to connect to the NetSoul socket.
// For now only authentication is supported and the PING server command.
package gosoul

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
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
	conn     net.Conn
	login    string
	password string
	salut    []string
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

func (gs *GoSoul) Authenticate(authType string) error {
	if authType != AUTHTYPE_KRB {
		authType = AUTHTYPE_MD5
	}
	gs.Send("auth_ag ext_user none -")
	gs.Parse()
	switch authType {
	case AUTHTYPE_KRB:
		return errors.New("Kerberos authentication not yet implemented")
	case AUTHTYPE_MD5:
		gs.Send(gs.md5Auth())
	}
	msg, _ := gs.Read()
	if msg != "rep 002 -- cmd end" {
		return errors.New("Bad login or password")
	} else {
		gs.Send("user_cmd attach")
		myt := time.Now()
		gs.Send(fmt.Sprintf("user_cmd state server:%d", myt.Unix()))
	}
	return nil
}

func (gs *GoSoul) Parse() error {
	res, err := gs.Read()
	if state, _ := regexp.MatchString("^ping.*", res); state {
		err = gs.Send(res)
	}
	return err
}

func (gs *GoSoul) Send(s string) error {
	_, err := gs.conn.Write([]byte(s + "\n"))
	log.Printf("[gosoul-send] : %s\n", s)
	return err
}

func (gs *GoSoul) Read() (string, error) {
	readBuffer := make([]byte, 2048)
	resLen, err := gs.conn.Read(readBuffer)
	if err != nil {
		return "", err
	}
	res := string(readBuffer[0 : resLen-1])
	log.Printf("[gosoul-read] : %s\n", res)
	return res, err
}

func (gs *GoSoul) Exit() {
	gs.Send("exit")
	gs.conn.Close()
}

func New(login string, password string) (*GoSoul, error) {
	conn, err := net.Dial("tcp", NSHOST+":"+NSPORT)
	if err != nil {
		return nil, err
	}
	gs := &GoSoul{login: login, password: password, conn: conn}
	msg, err := gs.Read()
	if err != nil {
		return nil, err
	}
	gs.salut = strings.Split(msg, " ")
	return gs, err
}
