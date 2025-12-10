package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

var validUser = struct {
	Username string
	Password string
}{
	Username: "user",
	Password: "12345",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if creds.Username != validUser.Username || creds.Password != validUser.Password {
		log.Println("Failed login attempt:", creds.Username)
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(TokenResponse{Error: "Invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(2 * time.Minute)
	claims := jwt.MapClaims{
		"username": creds.Username,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing token"))
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Println("Invalid token attempt")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte("Invalid or expired token"))
		return
	}

	w.Write([]byte("Access granted! Hello, " + claims["username"].(string)))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/protected", protectedHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
