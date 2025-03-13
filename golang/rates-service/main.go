package main

import (
	"github.com/jeremyseow/rates-service/external/rate"
	"github.com/jeremyseow/rates-service/server"
	"github.com/jeremyseow/rates-service/storage"
)

func main() {
	ratesClient := rate.NewRatesClient()
	ratesStorage := storage.NewRatesStorage()
	s := server.NewServer(ratesClient, ratesStorage)
	s.StartServer()
}
