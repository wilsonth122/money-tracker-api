package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/wilsonth122/money-tracker-api/pkg/dao"
	"github.com/wilsonth122/money-tracker-api/pkg/model"
	u "github.com/wilsonth122/money-tracker-api/pkg/utils"
)

// CreateUser - Endpoint ofr creating a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if resp, ok := validate(&user); !ok {
		log.Println(resp)
		u.RespondWithError(w, http.StatusBadRequest, resp)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := dao.DBConn.InsertUser(user); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Token = model.GenerateToken(user.Email)

	// Delete password before response
	user.Password = ""

	u.RespondWithJSON(w, http.StatusOK, user)
}

// LoginUser - Endpoint for logging a user in and returning a signed token
func LoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var unauthUser model.User

	if err := json.NewDecoder(r.Body).Decode(&unauthUser); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := dao.DBConn.FindUserByEmail(unauthUser.Email)
	if err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid login credentials. Please try again")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(unauthUser.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid login credentials. Please try again")
		return
	}

	user.Token = model.GenerateToken(user.Email)

	// Delete password before response
	user.Password = ""

	u.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser - Endpoint for deleting a user based on auth token
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := r.Context().Value("user").(string)

	if err := dao.DBConn.RemoveUserByEmail(user); err != nil {
		log.Println(err)
		u.RespondWithError(w, http.StatusBadRequest, "User doesn't exist or has already been deleted")
		return
	}

	u.RespondWithJSON(w, http.StatusOK, "User deleted")
}

// Validate User details
func validate(user *model.User) (string, bool) {

	if !strings.Contains(user.Email, "@") {
		return "Email address is required", false
	}

	if len(user.Password) < 6 {
		return "Password is required", false
	}

	// Check for errors and duplicate emails
	exists, err := dao.DBConn.UserExists(user.Email)
	if err != nil {
		return err.Error(), false
	}
	if exists {
		return "Email address already in use by another user", false
	}

	return "Requirement passed", true
}
