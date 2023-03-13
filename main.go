package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/madhushaw1012/split-bill/database"
	validation "github.com/madhushaw1012/split-bill/middleware"
	"github.com/madhushaw1012/split-bill/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getAll", ShowAll).Methods("GET")
	r.HandleFunc("/", IndexRegister).Methods("GET")
	r.HandleFunc("/login", IndexLogin).Methods("GET")
	r.HandleFunc("/user/login", Login).Methods("POST")
	r.HandleFunc("/user/register", Register).Methods("POST")
	r.HandleFunc("/user/logout", Logout).Methods("GET")
	r.HandleFunc("/home", Home).Methods("GET")

	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", r)
}
func Home(w http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("token")
	// if err != nil {
	// 	fmt.Println("Bhagg yaha se ", err)
	// 	return
	// }
	// if time.Since(cookie.Expires) > 0 {
	// 	fmt.Println("Session expired")
	// 	return
	// }
	currentUser := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(currentUser); err != nil {
		fmt.Println("User theek se nhi aaya ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "home.html", currentUser)

}
func IndexRegister(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "register.html", nil)
}
func IndexLogin(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		MaxAge:  0,
		Expires: time.Now(),
	})
	fmt.Fprintln(w, "User successfully logged out")
}
func Login(w http.ResponseWriter, r *http.Request) {
	newLoginUser := new(models.User)
	json.NewDecoder(r.Body).Decode(newLoginUser)
	statusCode := validation.Validate(newLoginUser, w)
	w.WriteHeader(statusCode)
}
func ShowAll(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}
	cursor, err := col.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Database problem")
	}
	defer cursor.Close(context.Background())
	var allUsers []models.User
	cursor.All(context.Background(), &allUsers)
	fmt.Fprintln(w, allUsers)
}
func Register(w http.ResponseWriter, r *http.Request) {

	newRegisterUser := new(models.User)
	// Try to decode the request body into the struct. If there is an error, send response to client
	err := (json.NewDecoder(r.Body)).Decode(newRegisterUser)
	if err != nil {
		// Log activity
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("New user register: ", newRegisterUser)
	w.Header().Set("Content-type", "application/json")
	//check if user exists or not
	if !UserExists(newRegisterUser) {
		fmt.Println(w, "User already exists. Please login...")
		w.WriteHeader(http.StatusPermanentRedirect)
		return
	}
	//Added new user
	insertNewUser(newRegisterUser)
	fmt.Println("Congratulations, You are now registered.")
	json.NewEncoder(w).Encode(newRegisterUser)
	w.Header().Add("headerName", newRegisterUser.Name)
	w.WriteHeader(http.StatusOK)
}

var col *mongo.Collection = database.OpenCollection(database.Client, "users")

func UserExists(user *models.User) bool {
	filter := bson.M{"email": user.Email}
	result := col.FindOne(context.Background(), filter)
	var u models.User
	err := result.Decode(u)
	fmt.Println("Check user exists error:", u)
	return err == mongo.ErrNoDocuments
}

func insertNewUser(user *models.User) {
	inserted, err := col.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal("Error inserting", err)
	}
	fmt.Println("Successfully inserted", inserted.InsertedID)
}
