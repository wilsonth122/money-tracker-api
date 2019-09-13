package app

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
	"github.com/wilsonth122/money-tracker-api/pkg/stream"
)

// Setup - Should be called by the init() function upon service start up.
// Moved logic here to support multiple environments which may need custom code in the init() function to work.
func Setup() {
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
	dao.DBConn.ExpenseCollection = conf.Database.ExpenseCollection
	dao.DBConn.Connect()
}

// Start - Should be called by the main() function upon service start up.
// Moved logic here to support multiple environments which may need custom code in the main() function to work.
func Start() {
	conf := config.New()

	// Start the websocket used for streaming expenses
	stream.Init()

	r := mux.NewRouter()

	// Attach JWT auth middleware
	r.Use(auth.JwtAuthentication)

	c := cors.New(cors.Options{
		AllowedOrigins: conf.API.AllowedOrigins,
		AllowedMethods: conf.API.AllowedMethods,
		AllowedHeaders: conf.API.AllowedHeaders,
	})

	r.HandleFunc("/api/user/new", api.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/login", api.LoginUser).Methods("POST")
	r.HandleFunc("/api/user/delete", api.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/stream/expenses", api.StreamAllExpenses).Methods("GET")
	r.HandleFunc("/api/expenses", api.AllExpenses).Methods("GET")
	r.HandleFunc("/api/expenses/{id}", api.GetExpense).Methods("GET")
	r.HandleFunc("/api/expenses", api.CreateExpense).Methods("POST")
	r.HandleFunc("/api/expenses", api.UpdateExpense).Methods("PUT")
	r.HandleFunc("/api/expenses/{id}", api.DeleteExpense).Methods("DELETE")

	handler := c.Handler(r)

	port := conf.API.Port
	log.Printf("Origins: %s", conf.API.AllowedOrigins)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
