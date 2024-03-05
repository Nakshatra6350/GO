package main

import (
	"context"
	"log"
	"main/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main(){
	l := log.New(os.Stdout, "product-api",log.LstdFlags)
	helloHandler := handlers.NewHello(l)
	byeHandler := handlers.NewBye(l)
	sm := http.NewServeMux()
	sm.Handle("/",helloHandler)

	sm.Handle("/bye",byeHandler)

	s := http.Server{
		Addr : ":9090",
		Handler: sm,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	
	go func(){
		err := s.ListenAndServe()
		if err!= nil{
			log.Fatal(err)
		}
	}()

	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)
	signal.Notify(sign,os.Kill)

	sig := <- sign
	l.Println("Received terminate, graceful shutdown ",sig)

	tc,_ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}