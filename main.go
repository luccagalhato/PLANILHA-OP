package main

import (
	"flag"
	"log"
	"op/api"
	c "op/config"
	database "op/database"
	s "op/server"
)

var createConfig bool
var connectionLinx *database.SQLStr

var err error

func init() {
	
	flag.BoolVar(&createConfig, "c", false, "create config.yaml file")
	flag.Parse()

	if createConfig {
		c.CreateConfigFile()
		return
	}

	log.Print("loading config file")
	if err := c.LoadConfig(); err != nil {
		log.Fatal(err)
	}
	log.Print("connecting sql ...")
	connectionLinx, err = database.MakeSQL(c.Yml.SQLLinx.Host, c.Yml.SQLLinx.Port, c.Yml.SQLLinx.User, c.Yml.SQLLinx.Password)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	api.SetSQLConn(connectionLinx)
	s.Controllers()
}
