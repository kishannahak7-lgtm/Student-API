package main

import (
	"GO-API/internal/config"
	"context"
	//"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){
	//load setup
	cfg:=config.MustLoad()

	//database setup
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to projAPI"))
	})
	//setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}


	slog.Info("server started ",slog.String("address ",cfg.Address))
	

	done:=make(chan os.Signal,1)

	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGALRM)


	go func(){
		err := server.ListenAndServe()
	if err!=nil{
		log.Fatal("failed to start server")
	}
	}()
	
	<-done

	//gracefull shutdown
	slog.Info("shutting down the server ")

	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()


	if err:=server.Shutdown(ctx); err != nil {
	
		slog.Error("failed to shutdown server", slog.String("error",err.Error()))
	}
	slog.Info("shutdown successfully")

}
