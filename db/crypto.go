package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/config"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents the JWT claims
type Claims struct {
	UserId int    `json:"userID"`
	Name   string `json:"userName"`
	Role   string `json:"role"`
	PWD    string `json:"pwd"`
	jwt.RegisteredClaims
}

// hashPassword hashes a plain text password with bcrypt and returns the hashed password.
func hashPassword(password string) (string, error) {
	// Generate a hash of the password using bcrypt
	log.Println(len(password), len([]byte(password)))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func verifyPassword(password string, hashedPassword string) (bool, error) {
	// Generate a hash of the password using bcrypt

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (database *Database) CreateSession(body gen.PostLoginJSONBody, cxt context.Context) (string, gen.User, error) {
	var user dbModels.User
	query := database.DB.NewSelect()
	query.Model(&user)
	query.Where("\"user\".\"email\" = ?", body.Email)
	query.Relation("DefaultChoice")
	err := query.Scan(cxt) //sql.ErrNoRows
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", gen.User{}, dbModels.ErrUserNotFound
		}
		return "", gen.User{}, err
	}
	verified, err := verifyPassword(*body.Password, user.PwdHash)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", gen.User{}, dbModels.ErrInvalidPassword
	}
	if verified {
		// Generate a new session token
		token, err := generateSessionToken(user, time.Now().AddDate(1, 0, 0))
		if err != nil {
			return "", gen.User{}, err
		}
		return token, user.ToGen(), err
	}
	return "", gen.User{}, err

}
func generateSessionToken(user dbModels.User, expirationTime time.Time) (string, error) {
	// Generate a random token

	// Create the JWT claims, which includes the email and expiration time
	claims := &Claims{
		UserId: user.Id,
		Name:   user.Name,
		Role:   user.Role,
		PWD:    generateSHA256Hash(user.PwdHash)[:8],
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token using the claims and the signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.Config.Crypto.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func generateSHA256Hash(input string) string {
	// Create a new SHA-256 hash instance
	hash := sha256.New()

	// Write the input data to the hash
	hash.Write([]byte(input))

	// Compute the SHA-256 checksum
	checksum := hash.Sum(nil)

	// Convert the truncated checksum to a hexadecimal string
	return base64.StdEncoding.EncodeToString(checksum)
}

// VerifySession verifies the JWT token and returns the user claims if valid.
func unpackToken(tokenString string) (*Claims, error) {
	// Parse the token using the JWT library
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is HMAC and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.Config.Crypto.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Optionally, you can add additional checks here (e.g., user existence in the database)

	return claims, nil
}
func (database *Database) verifySession(tokenString string, cxt context.Context) (dbModels.User, error) {
	claims, err := unpackToken(tokenString)
	if err != nil {
		return dbModels.User{}, err
	}
	var user dbModels.User
	query := database.DB.NewSelect()
	query.Model(&user)
	query.Where("\"user\".\"id\" = ?", claims.UserId)
	query.Relation("DefaultChoice")
	err = query.Scan(cxt) //sql.ErrNoRows
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbModels.User{}, dbModels.ErrUserNotFound
		}
		return dbModels.User{}, err
	}

	if claims.PWD == generateSHA256Hash(user.PwdHash)[:8] {
		return user, nil
	}
	return dbModels.User{}, dbModels.ErrInvalidPassword

}
func (database *Database) VerifySession(tokenString string, cxt context.Context) (gen.User, error) {
	user, err := database.verifySession(tokenString, cxt)
	return user.ToGen(), err
}
