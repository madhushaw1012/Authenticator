package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	controller "github.com/madhushaw1012/split-bill/controllers"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css", http.FileServer(http.Dir("./static/css"))))
	r.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images", http.FileServer(http.Dir("./static/images"))))
	r.HandleFunc("/getAll", controller.ShowAll).Methods("GET")
	r.HandleFunc("/", controller.IndexRegister).Methods("GET")
	r.HandleFunc("/login", controller.IndexLogin).Methods("GET")
	r.HandleFunc("/user/login", controller.Login).Methods("POST")
	r.HandleFunc("/user/register", controller.Register).Methods("POST")
	r.HandleFunc("/user/logout", controller.Logout).Methods("GET")
	r.HandleFunc("/home", controller.Home).Methods("GET")

	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", r)
}
