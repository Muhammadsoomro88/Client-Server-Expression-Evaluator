package storage

import (
	"log"

	"bitbucket.org/getcake/capservice/stores"
	"github.com/cakemarketing/go-common/v5/settings"
)

var client *stores.Redis

func Redisconnection() {
	settings.AddConfigPath("config")
	settings.SetConfigName("local")
	err := settings.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	client = stores.ConnectRedis(settings.GetString("REDIS_HOST"), settings.GetString("REDIS_PORT"))
}

func Clientinsert(name string) {
	client.LPush("Clients", name)
}

func Expressioninsert(name string, exp []byte) {
	client.SAdd(name, string(exp))
}

/*func Expressionget(x string) []string {
	res, err := client.SMembers(x).Result()
	if err != nil {
		log.Fatal(err)
	}
	return res
}*/

/*func Clientget(totalClients int) []string {
	res, err := client.LRange("Clients", int64(totalClients), int64(totalClients)).Result()
	if err != nil {
		log.Fatal(err)
	}
	return res
}*/

func Clientget1() []string {
	res, err := client.LRange("Clients", 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func Poprecord(x string) string {
	res, err := client.SPop(x).Result()
	if err != nil {
		log.Fatal(err)
	}
	return res
}
