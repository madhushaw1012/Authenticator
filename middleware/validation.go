package validation

// Since we are logging in, we'll have to create a token. Now, if the token already exists, we dont create it.
// When we log out, we should delete the token.

import (
	"context"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madhushaw1012/split-bill/database"
	"github.com/madhushaw1012/split-bill/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Claims struct {
	Name     string
	Email    string
	UserType string
	jwt.StandardClaims
}

var jwtKey = []byte("secret")

func Validate(loggingUser *models.User, w http.ResponseWriter) int {
	col := database.OpenCollection(database.Client, "users")
	filter := bson.M{"email": loggingUser.Email}
	result := col.FindOne(context.Background(), filter)
	var user models.User
	err := result.Decode(&user)
	fmt.Println(user)
	// User not found.
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound
	}
	if err != nil {
		fmt.Println("Problem extracting from database", err)
		return http.StatusInternalServerError
	}
	fmt.Println("***************")
	fmt.Println(user.Password + " ---> " + loggingUser.Password)
	// wrong password
	if user.Password != loggingUser.Password {
		return http.StatusUnauthorized
	}

	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Name:     loggingUser.Name,
		Email:    loggingUser.Email,
		UserType: loggingUser.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(http.StatusInternalServerError)
		return http.StatusOK
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		MaxAge:  60 * 60,
		Expires: expirationTime,
	})
	return http.StatusOK
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
	claims := &Claims{}
	tokenString := cookie.Value
	tkn, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Until(time.Unix(claims.ExpiresAt, 0)) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expirationTime := time.Now().Add(10 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh_tokenString, err := refresh_token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   refresh_tokenString,
		Expires: expirationTime,
	})
}
