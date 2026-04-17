package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignInRequest struct {
	Password string `json:"password"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, errors.New("invalid json"))
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	if req.Password != pass {
		writeError(w, errors.New("invalid password"))
		return
	}

	h := sha256.Sum256([]byte(pass))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": fmt.Sprintf("%x", h),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJson(w, map[string]string{"token": tokenString})
}
