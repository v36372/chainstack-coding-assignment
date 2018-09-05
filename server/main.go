package main

import (
	"chainstack/cmd"
	"chainstack/config"
	"chainstack/infra"
	"fmt"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
)

func main() {
	cmd.Execute()
	conf := config.Get()

	setupInfra(conf)
	defer infra.ClosePostgreSql()
	address := fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port)
	server := http.Server{
		Addr: address,
	}

	if err := gracehttp.Serve(&server); err != nil {
		panic(err)
	}
}

func setupInfra(conf config.Config) {
	// Postgresql
	infra.InitPostgreSQL()
}
