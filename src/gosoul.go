/**
 * Christian KAKESA
 * filename : gosoul.go
 */
package main
import ("net"
		"http"
		"os"
		"time"
		"fmt"
		"flag"
		"bufio"
		"crypto/md5"
		"encoding/hex"
		"strings"
		"regexp"	)

/* Construction de la chaîne d'authentification MD5 */
func md5NSAuthString(hash string, ip string, port string, pass string) string {
	authStringFormated := fmt.Sprintf("%s-%s/%s%s", hash, ip, port, pass)
	md := md5.New()
	md.Write(strings.Bytes(authStringFormated))
	authString := hex.EncodeToString(md.Sum())
	return authString
}

func main() {
	/* Dialogue utilisateur pour la récupération des paramètres de connexion */
	var login string
	flag.StringVar(&login, "l", "", "IONIS socks login")
	flag.Parse()
	if login == "" {
		fmt.Printf("Usage :\n")
		flag.PrintDefaults()
		os.Exit(-20001)
	}
	fmt.Printf("Enter your IONIS socks password :\n")
	buf := bufio.NewReader(os.Stdin)
	password, err := buf.ReadString('\n');
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20002)
	}
	password = password[0:len(password)-1]
	
	/* Connection au serveur NetSoul */
	conn, err := net.Dial("tcp", "", "ns-server.epita.fr:4242")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20003)
	}
	var msgReader = make([]byte, 2048)
	msgLen, err := conn.Read(msgReader);
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20004)
	}
	res := string(msgReader[0:msgLen])
	salutNSString := strings.Split(res, " ", 0)
	
	/* Authentication au serveur NetSoul */
	msgLen, err = conn.Write( strings.Bytes("auth_ag ext_user none -\n") )
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20005)
	}
	msgLen, err = conn.Read(msgReader);
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20006)
	}
	res = string(msgReader[0:msgLen])
	if res != "rep 002 -- cmd end\n" {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", "Impossible d'initialiser l'authentification au server !!!");
		os.Exit(-20007)
	}
	res = fmt.Sprintf(	"ext_user_log %s %s %s %s\n",
						login,
						md5NSAuthString(salutNSString[2], salutNSString[3], salutNSString[4], password),
						http.URLEscape("GO-Soul, by Christian KAKESA"),
						http.URLEscape("HOME with GO-Soul NetSoul ident client writen in GO language")	)
	msgLen, err = conn.Write( strings.Bytes(res) )
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20006)
	}
	msgLen, err = conn.Read(msgReader);
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20007)
	}
	res = string(msgReader[0:msgLen])
	if res != "rep 002 -- cmd end\n" {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", "Impossible de s'authentifier au server, vérifier vos login et mot de passe !!!");
		os.Exit(-20008)
	}
	
	/* Activation des services du PIE IONIS */
	msgLen, err = conn.Write( strings.Bytes("user_cmd attach\n") )
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20009)
	}
	/* Changement du statut de connexion en tant que client de type service/serveur */
	myt := time.LocalTime()
	msgLen, err = conn.Write( strings.Bytes( fmt.Sprintf("user_cmd state server:%d\n", myt.Seconds()) ) )
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20010)
	}
	
	/* Boucle de traitement des messages reçus */
	fmt.Printf("GO-Soul service started...\n")
	for {
		msgLen, err = conn.Read(msgReader);
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
			break
		}
		res = string(msgReader[0:msgLen-1])
		myt := time.LocalTime()
		fmt.Fprintf(os.Stdout, "[%s] : %s\n", myt.Format(time.ISO8601), res)
		state, err := regexp.MatchString("^ping.*", res)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err);
			break
		} else if state {
			msgLen, err = conn.Write( strings.Bytes(res) )
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
				break
			}
		}
	}
	
	msgLen, err = conn.Write( strings.Bytes("exit\n") )
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(-20011)
	}
	conn.Close()
	fmt.Printf("Thanks for using GO-Soul, the NetSoul ident service writen in GO language !!!\n")
	os.Exit(0)
}
