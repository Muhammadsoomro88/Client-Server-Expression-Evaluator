package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/apex/log"
	"github.com/cakemarketing/clientcalculator/calculator"
	"github.com/cakemarketing/clientcalculator/handler"
	"github.com/cakemarketing/clientcalculator/server"
	"github.com/cakemarketing/clientcalculator/storage"
	_ "github.com/go-sql-driver/mysql"
)

//total Exp counter
var totalExp int = 0

//user exp counter
var userExp int = 0

//client counter
var totalClients int = 0

//all connected clients
var allClients []string

type Data1 struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	log := log.WithFields(log.Fields{})
	go func() {
		<-c
		handler.CloseHandler(totalExp, totalClients, log, c)
		os.Exit(0)
	}()
}

func main() {
	SetupCloseHandler()
	req := server.ConnectServer()
	defer req.Close()

	for {
		con, err := req.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go getClient(con)
	}
}

func getClient(con net.Conn) {
	defer con.Close()
	log := log.WithFields(log.Fields{})

	log.Info("Client Connected")
	io.WriteString(con, "Enter Username: ")
	name, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}

	newName := calculator.GetUsername(name)

	newName = strings.Replace(newName, "\r\n", "", -1)

	storage.Redisconnection()

	storage.Clientinsert(newName)
	allClients = append(allClients, newName)

	io.WriteString(con, "\n"+"Welcome "+name+"\n")
	scanner := bufio.NewScanner(con)

	totalClients++
	userExp = 0

	//MySql Database setup
	db := storage.Connectionsql()
	defer db.Close()

	for scanner.Scan() {
		var res int = 0
		var flag int = 0

		//testing
		var num int = 0
		var opr int = 0

		io.WriteString(con, "Enter the Expression: ")
		x, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}

		for i := 0; i < len(x); i++ {
			if x[i] == '0' || x[i] == '1' || x[i] == '2' || x[i] == '3' || x[i] == '4' || x[i] == '5' || x[i] == '6' || x[i] == '7' || x[i] == '8' || x[i] == '9' {
				num++
			}
			if x[i] == '+' || x[i] == '-' || x[i] == '*' || x[i] == '/' {
				opr++
			}
		}

		if len(x) == 4 {
			io.WriteString(con, "\nExpression must contain atleast 2 digits\n\n ")
			continue
		}

		for i := 0; i < len(x); i++ {
			if x[i] == 'a' || x[i] == 'b' || x[i] == 'c' || x[i] == 'd' || x[i] == 'e' || x[i] == 'f' || x[i] == 'g' || x[i] == 'h' || x[i] == 'i' || x[i] == 'j' || x[i] == 'k' || x[i] == 'l' || x[i] == 'm' || x[i] == 'n' || x[i] == 'o' || x[i] == 'p' || x[i] == 'q' || x[i] == 'r' || x[i] == 's' || x[i] == 't' || x[i] == 'u' || x[i] == 'v' || x[i] == 'w' || x[i] == 'x' || x[i] == 'y' || x[i] == 'z' {
				flag = 1
			}
		}
		if num > opr {
			res, totalExp, userExp = calculator.Expresison(x, res, totalExp, userExp, name, con)
		} else {
			flag = 1
		}

		if res == 0 || flag == 1 {
			io.WriteString(con, "\nExpression must contain Numbers and Operators +,-,*,/\n\n")
		}

	}
	log.Warn("Client Disconnected")
	resExp := fmt.Sprint(userExp)
	clientName := calculator.GetUsername(name)
	clientName = strings.Replace(clientName, "\r\n", "", -1)
	log.Info("Connected Client: " + clientName + ", Expresison Executed: " + resExp + "\n")
}
