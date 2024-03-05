package main

import (
	"context"
	"log"
	"main/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main(){
	l := log.New(os.Stdout, "product-api",log.LstdFlags)
	productsHandler := handlers.NewProducts(l)
	// helloHandler := handlers.NewHello(l)
	// byeHandler := handlers.NewBye(l)
	// sm := http.NewServeMux()

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/",productsHandler.GetProducts)


	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}",productsHandler.UpdateProducts)
	putRouter.Use(productsHandler.MiddlewareValidateProduct)


	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/",productsHandler.AddProduct)
	postRouter.Use(productsHandler.MiddlewareValidateProduct)



	// sm.Handle("/",productsHandler)
	// sm.Handle("/",helloHandler)

	// sm.Handle("/bye",byeHandler)

	s := http.Server{
		Addr : ":9090",
		Handler: sm,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	
	go func(){
		log.Printf("server starting at port :9090")
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