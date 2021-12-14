package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net"
	"net/http"
	"playground/app/container"
	"strconv"
	"time"
	"github.com/afex/hystrix-go/hystrix"
)

type Server struct {
	router *chi.Mux
}

func main(){
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8081"), hystrixStreamHandler)

	limitApiClient := container.Init()

	router := chi.NewRouter()

	server := &Server{router: router}
	server.withGlobalMiddleware(router)

	server.router.Get("/test-hystrix", func(w http.ResponseWriter, r *http.Request) {
		id,_ := strconv.Atoi(r.URL.Query().Get("id"))
		limitApiClient.Check(id)

	})
	server.start()

}
func (s *Server) withGlobalMiddleware(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
}


func (s *Server) start() {
	server := &http.Server{Addr: ":" + "8080", Handler: s.router}
	fmt.Println("server started on port: %s",8080)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("cannot start server.", err)
		panic("cannot start server")
	}
}