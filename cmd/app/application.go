package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/wilsonth122/money-tracker-api/pkg/api"
	"github.com/wilsonth122/money-tracker-api/pkg/auth"
	"github.com/wilsonth122/money-tracker-api/pkg/config"
	"github.com/wilsonth122/money-tracker-api/pkg/dao"
)

func init() {
	// Load config
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}

	conf := config.New()

	// Connect to server
	dao.DBConn.Addresses = conf.Database.Addresses
	dao.DBConn.Username = conf.Database.Username
	dao.DBConn.Password = conf.Database.Password
	dao.DBConn.AdminDatabase = conf.Database.AdminDatabase
	dao.DBConn.AppDatabase = conf.Database.AppDatabase
	dao.DBConn.UserCollection = conf.Database.UserCollection
	dao.DBConn.Connect()
}

func main() {
	conf := config.New()

	r := mux.NewRouter()

	// Attach JWT auth middleware
	r.Use(auth.JwtAuthentication)

	c := cors.New(cors.Options{
		AllowedOrigins: conf.API.AllowedOrigins,
		AllowedMethods: conf.API.AllowedMethods,
		AllowedHeaders: conf.API.AllowedHeaders,
		Debug:          true,
	})

	r.HandleFunc("/api/user/new", api.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/login", api.LoginUser).Methods("POST")
	r.HandleFunc("/api/user/delete", api.DeleteUser).Methods("DELETE")

	handler := c.Handler(r)

	port := conf.API.Port
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
