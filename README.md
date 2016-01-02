# GoSoul [![Build Status](https://travis-ci.org/fenicks/gosoul.svg?branch=master)](https://travis-ci.org/fenicks/gosoul)

**GoSoul** is a [netsoul](https://doc.ubuntu-fr.org/netsoul) authentication client written in Go.

    go get -u github.com/fenicks/gosoul
    go install github.com/fenicks/gosoul

## Add your credentials in config file

    mkdir -p $HOME/.config/gosoul
    echo "user_name:my_socks_password" > $HOME/.config/gosoul/config.txt

## Run the GoSoul client

    gosoul

**Output example**

    christian@christian-GA-MA78GM-UD2H:~/workspace/gosoul$ go run gosoul.go 
    2014/10/11 15:10:04 [gosoul-read] : salut 2368 1d5673511d3b37f9f335a47440c5a698 52.25.45.65 57979 1413032947
    2014/10/11 15:10:04 [gosoul-send] : auth_ag ext_user none -
    2014/10/11 15:10:04 [gosoul-read] : rep 002 -- cmd end
    2014/10/11 15:10:04 [gosoul-send] : ext_user_log kakesa_c 2c49cd2513b6a773d284126c50cbd4c1 GoSoul%2C+by+Christian+KAKESA %40christian-GA-MA78GM-UD2H
    2014/10/11 15:10:04 [gosoul-read] : rep 002 -- cmd end
    2014/10/11 15:10:04 [gosoul-send] : user_cmd attach
    2014/10/11 15:10:04 [gosoul-send] : user_cmd state server:1413033004

# Links
 *   [Ubuntu netsoul](http://doc.ubuntu-fr.org/netsoul)
 *   [RubySoul-NG](https://github.com/fenicks/rubysoul-ng/)
 *   etc...
