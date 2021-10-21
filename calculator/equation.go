package calculator

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/apex/log"
	"github.com/cakemarketing/clientcalculator/storage"
)

type Data struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

var obj Data

func storeData(newName string) {
	//redis connect
	storage.Redisconnection()

	//connect SQL
	db := storage.Connectionsql()
	defer db.Close()

	var index int = 0

	var obj1 Data

	exp := storage.Poprecord(newName)
	err := json.Unmarshal([]byte(exp), &obj1)
	if err != nil {
		log.Fatal(err.Error())
	}
	index++
	time.Sleep(2 * time.Second)
	t1 := time.Now()
	tstamp := t1.Format(time.RFC3339)
	storage.Sqlquery(newName, obj1.Expression, obj1.Result, tstamp)
}

func Expresison(x string, res int, totalExp int, userExp int, name string, con net.Conn) (int, int, int) {
	for i := 0; i < len(x); i++ {
		if x[i] == '*' || x[i] == '+' || x[i] == '-' || x[i] == '/' {
			res = 1

			expression, err := govaluate.NewEvaluableExpression(x)
			if err != nil {
				log.Fatal(err.Error())
			}

			parameters := make(map[string]interface{})

			result, err := expression.Evaluate(parameters)
			if err != nil {
				log.Fatal(err.Error())
			}
			totalExp++
			userExp++

			valStr := fmt.Sprint(result)
			io.WriteString(con, "\n")
			io.WriteString(con, "Result: "+valStr+"\n\n")

			newName := GetUsername(name)
			newName = strings.Replace(newName, "\r\n", "", -1)

			x = strings.Replace(x, "\r\n", "", -1)
			log.Info(x + " = " + valStr)

			newRes, _ := strconv.ParseFloat(valStr, 64) //converting string into float
			obj = Data{Expression: x, Result: newRes}

			//trying marshal
			byteArray, err := json.Marshal(obj)
			if err != nil {
				log.Fatal(err.Error())
			}

			storage.Expressioninsert(newName, byteArray)

			go storeData(newName)
			break
		}
	}
	return res, totalExp, userExp
}

func GetUsername(name string) string {
	x := name
	var arr []byte

	for i := 21; i < len(x); i++ {
		arr = append(arr, x[i])
	}
	var newName string = ""
	newName = string(arr)
	return newName
}
