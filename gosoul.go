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

// Deafult values for the GoSoul client
const (
	GOS_DATA     = "GoSoul, by Christian KAKESA"
	GOS_LOCATION = "GoSoul@HOME"
)

// Athentication type : Kerberos or MD5
const (
	AUTHTYPE_KRB = "kerberos"
	AUTHTYPE_MD5 = "md5"
)

type UserData struct {
	login    string
	password string
	data     string
	State    string
	Location string
}

const (
	UserStateActif  string = "actif"   // User is connected and interaction is possible
	UserStateIdle   string = "idle"    // User is connected but no interaction is possible
	UserStateServer string = "service" // Service connection to coonect a device or hardware
)

type GoSoul struct {
	User  UserData
	conn  net.Conn
	salut []string
}

func (gs *GoSoul) md5Auth() string {
	res := fmt.Sprintf("%s-%s/%s%s", gs.salut[2], gs.salut[3], gs.salut[4], gs.User.password)
	md := md5.New()
	md.Write([]byte(res))
	authHashString := hex.EncodeToString(md.Sum(nil))
	res = fmt.Sprintf("ext_user_log %s %s %s %s",
		gs.User.login,
		authHashString,
		gs.User.data,
		gs.User.Location)
	return res
}

func (gs *GoSoul) Authenticate(authType string) error {
	gs.send("auth_ag ext_user none -")
	err := gs.Parse()
	if err != nil {
		return err
	}
	switch authType {
	case AUTHTYPE_KRB:
		return errors.New("Kerberos authentication not yet implemented")
	case AUTHTYPE_MD5:
		gs.send(gs.md5Auth())
	}
	msg, _ := gs.read()
	if msg != "rep 002 -- cmd end" {
		return errors.New("Bad login or password")
	} else {
		gs.send("user_cmd attach")
		gs.SetState(UserStateServer)
	}
	return nil
}

func (gs *GoSoul) Parse() error {
	res, err := gs.read()
	if err != nil {
		return err
	}
	if state, _ := regexp.MatchString("^ping.*", res); state {
		err = gs.send(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *GoSoul) send(s string) error {
	_, err := gs.conn.Write([]byte(s + "\r\n"))
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("[gosoul-send] : %s", s))
	return nil
}

func (gs *GoSoul) read() (string, error) {
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

func (gs *GoSoul) SetState(state string) {
	gs.User.State = state
	gs.send(fmt.Sprintf("user_cmd state %s:%d", state, time.Now().Unix()))
}

//TODO: Others netsoul send command here...

func (gs *GoSoul) Exit() {
	gs.send("exit")
	gs.conn.Close()
}

// Provides a GoSoul instance for netsoul server interaction.
func New(login, password, addr string) (gs *GoSoul, err error) {
	gs = &GoSoul{
		User: UserData{login: login,
			password: password,
			data:     url.QueryEscape(GOS_DATA),
			State:    UserStateServer,
			Location: url.QueryEscape(GOS_LOCATION)}}
	gs.conn, err = net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	msg, err := gs.read()
	if err != nil {
		return nil, err
	}
	gs.salut = strings.Split(msg, " ")
	return gs, err
}
