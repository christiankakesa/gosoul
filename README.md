# GoSoul [![Build Status](https://travis-ci.org/fenicks/gosoul.svg?branch=master)](https://travis-ci.org/fenicks/gosoul)

**GoSoul** is a [netsoul](https://doc.ubuntu-fr.org/netsoul) authentication client written in Go.

## Add your IONIS school credentials in config file

    mkdir -p $HOME/.config/gosoul
    echo "user_name:my_socks_password" > $HOME/.config/gosoul/config.txt

## Installation

### From binaries

Choose your binary platform: https://github.com/fenicks/gosoul/releases/latest.

Uncompress the archive and run the binary like the example bellow:

    cd /tmp
    wget https://github.com/fenicks/gosoul/releases/download/v1.1.0/gosoul-v1.1.0-linux-amd64.tar.gz
    tar -zxvf gosoul-v1.1.0-linux-amd64.tar.gz -C $HOME/bin
    $HOME/bin/gosoul

If $HOME/bin is in your $PATH run: `gosoul`

### From source

You need a GO develomeent environment ready: https://golang.org/doc/install.

    go get github.com/fenicks/gosoul
    go install github.com/fenicks/gosoul
    $GOPATH/bin/gosoul

If $GOPATH/bin is in your $PATH run: `gosoul`

## Output example

    christian@y510p:~/tmp$ $HOME/bin/gosoul
    2016/01/02 03:34:16 [gosoul-read] : salut 63 7cedc8aae60adfe21ca26c7628d0044d 54.256.3.254 46214 1451702055
    2016/01/02 03:34:16 [gosoul-send] : auth_ag ext_user none -
    2016/01/02 03:34:16 [gosoul-read] : rep 002 -- cmd end
    2016/01/02 03:34:16 [gosoul-send] : ext_user_log kakesa_c f0fabc090cbeaee2cbe5eb004124b6e2 GoSoul%2C+by+Christian+KAKESA %40y510p
    2016/01/02 03:34:16 [gosoul-read] : rep 002 -- cmd end
    2016/01/02 03:34:16 [gosoul-send] : user_cmd attach
    2016/01/02 03:34:16 [gosoul-send] : user_cmd state server:1451702056

## Links
 *   [Ubuntu netsoul](http://doc.ubuntu-fr.org/netsoul)
 *   [RubySoul-NG](https://github.com/fenicks/rubysoul-ng/)
 *   etc...
