package main

import (
	"log"
	"main/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	server := controller.NewServer("localhost", ":8080")
	r := mux.NewRouter()

	r.HandleFunc("/", server.GetLogin)
	r.HandleFunc("/post-login", server.PostLogin)
	r.HandleFunc("/data_list", server.GetMainPage)
	r.HandleFunc("/sutartys/update/{id}", server.GetUpdateSutartys).Methods("GET")
	r.HandleFunc("/sutartys/update/{id}", server.PostUpdateSutartys).Methods("POST")
	r.HandleFunc("/sutartys/create", server.GetCreateSutartys).Methods("GET")
	r.HandleFunc("/sutartys/create", server.PostCreateSutartys).Methods("POST")
	r.HandleFunc("/sutartys/delete/{id}", server.PostDeleteSutartys).Methods("POST")
	r.HandleFunc("/breziniai/update/{id}", server.GetUpdateBreziniai).Methods("GET")
	r.HandleFunc("/breziniai/update/{id}", server.PostUpdateBreziniai).Methods("POST")
	r.HandleFunc("/breziniai/create", server.GetCreateBreziniai).Methods("GET")
	r.HandleFunc("/breziniai/create", server.PostCreateBreziniai).Methods("POST")
	r.HandleFunc("/breziniai/delete/{id}", server.PostDeleteBreziniai).Methods("POST")
	r.HandleFunc("/leidimai/create", server.GetCreateLeidimai).Methods("GET")
	r.HandleFunc("/leidimai/create", server.PostCreateLeidimai).Methods("POST")
	r.HandleFunc("/leidimai/update/{id}", server.GetUpdateLeidimai).Methods("GET")
	r.HandleFunc("/leidimai/update/{id}", server.PostUpdateLeidimai).Methods("POST")
	r.HandleFunc("/leidimai/delete/{id}", server.PostDeleteLeidimai).Methods("POST")
	log.Fatal(http.ListenAndServe(server.Host+server.Port, r))
}
