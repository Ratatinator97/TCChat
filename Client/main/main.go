package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	//On se connecte
	connexion, _ := net.Dial("tcp", "127.0.0.1:8081") //LocalHost
	//On se connecte au serveur
	//Penser a handle les erreurs
	message, _ := bufio.NewReader(connexion).ReadString('\n')

	for {
		TabS := strings.Split(message, "\t")
		if TabS[0] == "TCCHAT_WELCOME" {
			fmt.Println(message)
			break
		}
	}

	nameOfUser := ecritureMsgServeur(1, connexion)
	fmt.Println(nameOfUser)
	fmt.Println("Etape connexion terminee")

	go read(connexion, nameOfUser)
	go write(connexion)

	exit := false
	for exit == false {
		exit = false
	}

}

func read(conn net.Conn, yourName string) {
	fmt.Println(yourName)
	for {
		message1, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
				os.Exit(1)
			}
		}
		if message1 != "" {
			tabS := strings.Split(message1, "\t")
			switch tabS[0] {
			case "TCCHAT_BCAST":
				split_tab := strings.Split(tabS[1], ":")
				fmt.Println(split_tab[0])
				inputName := "[" + yourName + "]"
				fmt.Println(inputName)
				if inputName == split_tab[0] {
					fmt.Println("On a detecte le nom")
				} else {
					fmt.Println("Nom non detecte")
					fmt.Println(tabS[1])
				}
			case "TCCHAT_USERIN":
				fmt.Println(tabS[1])
			case "TCCHAT_USEROUT":
				fmt.Println(tabS[1])
			case "TCCHAT_PERSO":
				fmt.Println(tabS[1])
			default:
				fmt.Println("Unexpected type of msg")
			}
		}
	}
}

func write(conn net.Conn) {
	for {
		ecritureMsgServeur(2, conn)
	}
}

func ecritureMsgServeur(msgType int, conn net.Conn) (name string) {

	reader := bufio.NewReader(os.Stdin)
	name = ""
	switch msgType {

	case 1:
		fmt.Print("Qui etes vous ?: ")
		texte, _ := reader.ReadString('\n')
		for {
			if texte != "\n" {
				break
			}
		}

		fmt.Println("Votre nom est " + texte)
		nvTexte := strings.TrimSuffix(texte, "\n")

		fmt.Fprintf(conn, "TCCHAT_REGISTER"+"\t"+nvTexte+"\n")
		name = nvTexte
		fmt.Println(name)

	case 2:
		name = ""
		texte, _ := reader.ReadString('\n')
		if texte != "\n" && texte != "" {
			texte := strings.TrimSuffix(texte, "\n")
			//fmt.Print("Envoi de message" + texte)
			fmt.Fprintf(conn, "TCCHAT_MESSAGE\t"+texte+"\n") //A le reception du serveur corriger ca
		}
	}
	return

}
