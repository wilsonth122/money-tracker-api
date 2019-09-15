package api

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"

	"github.com/wilsonth122/money-tracker-api/pkg/dao"
	"github.com/wilsonth122/money-tracker-api/pkg/model"
	"github.com/wilsonth122/money-tracker-api/pkg/stream"
	u "github.com/wilsonth122/money-tracker-api/pkg/utils"
)

// StreamAllExpenses - Endpoint to stream all expenses instead of just get them once
func StreamAllExpenses(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket connection")
	stream.WsHandler(w, r)
}

// AllExpenses - Endpoint to retrieve all expenses
func AllExpenses(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	expenses, err := dao.DBConn.FindAllExpenses(user)

	if err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RespondWithJSON(w, http.StatusOK, expenses)
}

// GetExpense - Endpoint to get a specific expense by id
func GetExpense(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	params := mux.Vars(r)
	expense, err := dao.DBConn.FindExpenseByID(params["id"])

	if err != nil || expense.UserID != user {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid Expense ID")
		return
	}

	u.RespondWithJSON(w, http.StatusOK, expense)
}

// CreateExpense - Endpoint to create an expense
func CreateExpense(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := r.Context().Value("user").(string)

	var expense model.Expense
	expense.UserID = user

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	expense.ID = bson.NewObjectId()

	if err := dao.DBConn.InsertExpense(expense); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send new expense to expense stream
	go stream.Writer(&expense)

	u.RespondWithJSON(w, http.StatusCreated, expense)
}

// UpdateExpense - Endpoint to update an expense
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := r.Context().Value("user").(string)

	var expense model.Expense
	expense.UserID = user

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dao.DBConn.UpdateExpense(expense); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send updated expense to expense stream
	go stream.Writer(&expense)

	u.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// DeleteExpense - Endpoint to delete an expense
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := dao.DBConn.RemoveExpenseByID(params["id"])

	if err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid Expense ID")
		return
	}

	u.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
