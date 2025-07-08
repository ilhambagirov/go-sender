package handler

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go-sender/internal/config"
	"net/http"
	"strconv"
	"time"
)

import _ "go-sender/docs"

func StartHttpServing(msgController *MsgController) {
	router := mux.NewRouter()
	port := strconv.Itoa(config.Config.Port)
	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	router.HandleFunc("/message", func(writer http.ResponseWriter, request *http.Request) {
		msgController.GetSentMessages(writer, request)
	})

	router.HandleFunc("/start", func(writer http.ResponseWriter, request *http.Request) {
		msgController.Start(writer, request)
	})

	router.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		msgController.Stop(writer, request)
	})

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
