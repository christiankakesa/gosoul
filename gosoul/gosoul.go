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
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	gosHost     = "ns-server.epita.fr"
	gosPort     = "4242"
	gosData     = "GoSoul, by Christian KAKESA"
	gosLocation = "GoSoul@HOME"
)

// Athentication type : Kerberos or MD5
const (
	AUTHTYPE_KRB = "kerberos"
	AUTHTYPE_MD5 = "md5"
)

type GoSoul struct {
	login      string
	password   string
	conn       net.Conn
	salut      []string
	NsData     string
	NsLocation string
}

func (gs *GoSoul) md5Auth() string {
	res := fmt.Sprintf("%s-%s/%s%s", gs.salut[2], gs.salut[3], gs.salut[4], gs.password)
	md := md5.New()
	md.Write([]byte(res))
	authHashString := hex.EncodeToString(md.Sum(nil))
	res = fmt.Sprintf("ext_user_log %s %s %s %s",
		gs.login,
		authHashString,
		gs.NsData,
		gs.NsLocation)
	return res
}

func (gs *GoSoul) Authenticate(authType string) error {
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
		gs.Send(fmt.Sprintf("user_cmd state actif:%d", time.Now().Unix()))
	}
	return nil
}

func (gs *GoSoul) Parse() error {
	res, err := gs.Read()
	if err != nil {
		return err
	}
	if state, _ := regexp.MatchString("^ping.*", res); state {
		err = gs.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *GoSoul) Send(s string) error {
	_, err := gs.conn.Write([]byte(s + "\r\n"))
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("[gosoul-send] : %s", s))
	return nil
}

func (gs *GoSoul) Read() (string, error) {
	readBuffer := make([]byte, 2048)
	resLen, err := gs.conn.Read(readBuffer)
	if err != nil {
		return "", err
	}
	res := string(readBuffer[0 : resLen-1])
	if err != nil {
		return "", err
	}
	log.Println(fmt.Sprintf("[gosoul-read] : %s", res))
	return res, nil
}

func (gs *GoSoul) Exit() {
	gs.Send("exit")
	gs.conn.Close()
}

func NewGoSoul(login, password, host, port string) (gs *GoSoul, err error) {
	gs = &GoSoul{
		login:      login,
		password:   password,
		NsData:     url.QueryEscape(gosData),
		NsLocation: url.QueryEscape(gosLocation)}
	if host == "" {
		host = gosHost
	}
	if port == "" {
		port = gosPort
	}
	gs.conn, err = net.Dial("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}
	msg, err := gs.Read()
	if err != nil {
		return nil, err
	}
	gs.salut = strings.Split(msg, " ")
	return gs, err
}

func New(login, password string) (*GoSoul, error) {
	return NewGoSoul(login, password, "", "")
}
