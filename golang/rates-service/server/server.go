package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jeremyseow/rates-service/external/rate"
	"github.com/jeremyseow/rates-service/server/handler"
	"github.com/jeremyseow/rates-service/server/route"
	"github.com/jeremyseow/rates-service/storage"
)

type Server struct {
	router *gin.Engine
}

func NewServer(rateClient *rate.RatesClient, ratesStorage *storage.RatesStorage) *Server {
	router := gin.Default()
	handlers := handler.NewHandlers(rateClient, ratesStorage)
	route.SetupRoutes(router, handlers)
	return &Server{
		router: router,
	}
}

func (s *Server) StartServer() {
	srv := &http.Server{
		Addr:    ":4001",
		Handler: s.router.Handler(),
	}

	go srv.ListenAndServe()

	// graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT)

	<-shutdownChan

	s.Shutdown()
}

func (s *Server) Shutdown() {
	fmt.Println("shutting down")
}
