package main

import (
	"flag"
	"go.uber.org/zap"
	"log"
	"time"

	"github.com/AtCliffUnderline/golang-diploma/internal/config"
	"github.com/AtCliffUnderline/golang-diploma/internal/database"
	"github.com/AtCliffUnderline/golang-diploma/internal/entities"
	router "github.com/AtCliffUnderline/golang-diploma/internal/http"
	"github.com/AtCliffUnderline/golang-diploma/internal/integrations"
	"github.com/caarlos0/env/v6"
)

func main() {
	c := config.Config{}
	err := env.Parse(&c)
	if err != nil {
		panic(err)
	}

	flag.StringVar(&c.RunAddress, "a", c.RunAddress, "a 127.0.0.1:8080")
	flag.StringVar(&c.DatabaseURI, "d", c.DatabaseURI, "d postgres://username:password@host:port/database_name")
	flag.StringVar(&c.AccrualSystemAddress, "r", c.AccrualSystemAddress, "r http://127.0.0.1:8081")
	flag.Parse()

	s := database.InitStorage(c)
	ur := entities.UserStorage{Storage: *s}
	or := database.OrderStorage{Storage: *s}
	wr := entities.WithdrawnStorage{Storage: *s}
	as := integrations.AccrualService{Address: c.AccrualSystemAddress, OrderRepository: or}
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	h := router.Handler{Config: c, UserRepository: ur, OrderRepository: or, WithdrawnRepository: wr}
	r := router.SetupRouter(h, ur, l)

	initOrdersChecker(as)

	err = r.Run(c.RunAddress)
	if err != nil {
		log.Fatal(err)
	}
}

func initOrdersChecker(as integrations.AccrualService) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			as.SyncAllOrders()
		}
	}
}
