package dao

import (
	"crypto/tls"
	"log"
	"net"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/wilsonth122/money-tracker-api/pkg/model"
)

// DAO struct
type DAO struct {
	Addresses         []string
	Username          string
	Password          string
	AdminDatabase     string
	AppDatabase  string
	UserCollection    string
	ExpenseCollection string
}

var DBConn = DAO{}

var db *mgo.Database

// Connect MongoDB session
func (dao *DAO) Connect() {
	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs:    dao.Addresses,
		Database: dao.AdminDatabase,
		Username: dao.Username,
		Password: dao.Password,
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	log.Printf("Dialing MongoDB Server...")
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Fatal(err)
	}

	db = session.DB(dao.AppDatabase)
	log.Println("Successfully connected to " + dao.AppDatabase)
}

// InsertUser - Inserts a user into the users collection
func (dao *DAO) InsertUser(user model.User) error {
	err := db.C(dao.UserCollection).Insert(&user)

	return err
}

// FindUserByEmail - Runs a find on the users collection and returns the first user with the email
func (dao *DAO) FindUserByEmail(email string) (model.User, error) {
	var user model.User

	err := db.C(dao.UserCollection).Find(bson.M{"email": email}).One(&user)

	return user, err
}

// RemoveUserByEmail - Rmeoves a user from the users collection
func (dao *DAO) RemoveUserByEmail(email string) error {
	err := db.C(dao.UserCollection).Remove(bson.M{"email": email})

	return err
}

// UserExists - Checks whether a user is already using the provided email
func (dao *DAO) UserExists(email string) (bool, error) {
	n, err := db.C(dao.UserCollection).Find(bson.M{"email": email}).Limit(1).Count()

	if err != nil && err != mgo.ErrNotFound {
		return true, err
	}

	if n > 0 {
		return true, nil
	}

	return false, nil
}

// FindAllExpenses - Runs a find on the expenses collection
// and returns all expense records relating to the user specified
func (dao *DAO) FindAllExpenses(user string) ([]model.Expense, error) {
	var expenses []model.Expense

	err := db.C(dao.ExpenseCollection).Find(bson.M{"userID": user}).All(&expenses)

	return expenses, err
}

// FindExpenseByID - Runs a find on the expenses collection and returns the first expense with the id
func (dao *DAO) FindExpenseByID(id string) (model.Expense, error) {
	var expense model.Expense

	err := db.C(dao.ExpenseCollection).FindId(bson.ObjectIdHex(id)).One(&expense)

	return expense, err
}

// InsertExpense - Inserts an expense record into the expenses collection
func (dao *DAO) InsertExpense(expense model.Expense) error {
	err := db.C(dao.ExpenseCollection).Insert(&expense)

	return err
}

// RemoveExpenseByID - Removes an expense by id
func (dao *DAO) RemoveExpenseByID(id string) error {
	err := db.C(dao.ExpenseCollection).RemoveId(bson.ObjectIdHex(id))

	return err
}

// UpdateExpense - Updates an expense record in the expenses collection
func (dao *DAO) UpdateExpense(expense model.Expense) error {
	err := db.C(dao.ExpenseCollection).UpdateId(expense.ID, &expense)

	return err
}

// RemoveUserExpenses - Removes all expenses relating to a user
func (dao *DAO) RemoveUserExpenses(email string) error {
	_, err := db.C(dao.ExpenseCollection).RemoveAll(bson.M{"userID": email})

	return err
}
