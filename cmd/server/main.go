package main

import (
	"fmt"
	"github.com/zeina1i/ethpay/passwords"
	server2 "github.com/zeina1i/ethpay/server"
	"github.com/zeina1i/ethpay/store/mysql"
)

func main() {
	store, err := mysql.NewStore(&mysql.Config{
		MysqlUsername: "ethpayserver",
		MysqlPassword: "123456",
		MysqlHost:     "localhost",
		MysqlPort:     6033,
		MysqlDb:       "ethpayserver",
	})
	if err != nil {
		fmt.Println(err)
	}

	err = store.InitializeDB()
	if err != nil {
		fmt.Println(err)
	}

	pm := passwords.NewScryptPasswords(nil)
	server := server2.NewHTTPServer(store, pm)
	server.InitRoutes()
	server.ListenAndServe()
	/*server := server.NewHTTPServer()
	server.InitRoutes()
	err := server.ListenAndServe()
	if err!=nil {
		fmt.Println(err)
	}*/
}
