package api

import (
	"errors"
	"net/http"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/db"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/golang-jwt/jwt/v4"
)

// optional code omitted

type Server struct {
	DB db.Database
}

func NewServer(DB db.Database) Server {
	return Server{
		DB: DB,
	}
}
func (server Server) isLoggedIn(w http.ResponseWriter, r *http.Request) (gen.User, *db.Claims, error) {
	var token string

	// First, try to read the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		// If the Authorization header is not present or malformed, check the cookie
		cookie, err := r.Cookie("session_token")
		if err == nil {
			token = cookie.Value
		} else if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "No token provided.", http.StatusUnauthorized)
			return gen.User{}, &db.Claims{}, dbModels.ErrInvalidToken
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return gen.User{}, &db.Claims{}, dbModels.ErrInvalidToken
		}
	}

	user, claims, err := server.DB.VerifySession(token, r.Context())
	if err != nil {
		if errors.Is(err, dbModels.ErrInvalidPassword) || errors.Is(err, dbModels.ErrUserNotFound) {
			if w != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			return gen.User{}, claims, dbModels.ErrInvalidToken
		}
		var jwterr *jwt.ValidationError
		if errors.As(err, &jwterr) {
			errCode := jwterr.Errors
			if errCode&jwt.ValidationErrorMalformed != 0 ||
				errCode&jwt.ValidationErrorUnverifiable != 0 ||
				errCode&jwt.ValidationErrorSignatureInvalid != 0 {
				if w != nil {
					http.Error(w, "Token malformed or Signature Invalid.", http.StatusBadRequest)
				}
				return gen.User{}, claims, dbModels.ErrInvalidToken
			} else if errCode&jwt.ValidationErrorExpired != 0 ||
				errCode&jwt.ValidationErrorNotValidYet != 0 {
				http.Error(w, "Token currently not Valid.", http.StatusUnauthorized)
				return gen.User{}, claims, dbModels.ErrInvalidToken
			} else if errCode&jwt.ValidationErrorId != 0 ||
				errCode&jwt.ValidationErrorIssuedAt != 0 ||
				errCode&jwt.ValidationErrorIssuer != 0 ||
				errCode&jwt.ValidationErrorClaimsInvalid != 0 {
				if w != nil {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				}
				return gen.User{}, claims, dbModels.ErrInvalidToken
			} else {
				if w != nil {
					http.Error(w, "Malformed Authorization", http.StatusUnauthorized)
				}
				return gen.User{}, claims, dbModels.ErrInvalidToken
			}
		}
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return gen.User{}, claims, dbModels.ErrInvalidToken
	}
	return user, claims, nil

}
