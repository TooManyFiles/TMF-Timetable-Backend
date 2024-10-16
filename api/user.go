package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	untisApiStructs "github.com/Mr-Comand/goUntisAPI/structs"
	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/config"
	"github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/untisDataCollectors"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/bcrypt"
)

// Get all users
// (GET /users)
func (server Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		http.Error(w, "Insufficient permission.", http.StatusForbidden)
		return
	}
	var resp []gen.User
	resp, err = server.DB.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Create a new user
// (POST /users)
func (server Server) PostUsers(w http.ResponseWriter, r *http.Request) {
	if !config.Config.CanSignUp {
		user, _, err := server.isLoggedIn(w, r)
		if err != nil {
			return
		}
		if user.Role == nil || *user.Role != gen.UserRoleAdmin {
			http.Error(w, "SignUp is currently disabled on this server.", http.StatusForbidden)
			return
		}
	}
	var userWithPW gen.PostUsersJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&userWithPW)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	_ = json.NewEncoder(log.Writer()).Encode(userWithPW)
	if userWithPW.UserData == nil || userWithPW.Password == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if (userWithPW.UserData.Role != nil) &&
		(*userWithPW.UserData.Role !=
			gen.UserRole("student")) &&
		(*userWithPW.UserData.Role !=
			gen.UserRole("teacher")) {

		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}
	resp, err := server.DB.CreateUser(*userWithPW.UserData, *userWithPW.Password, r.Context())

	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Field('C') {
			case "23505":
				http.Error(w, "User already exists.", http.StatusConflict)
			default:
				http.Error(w, "Internal server error."+pgErr.Field('C'), http.StatusInternalServerError)
				log.Printf("Error type: %T, Details: %s", err, err.Error())
			}
		} else if errors.Is(err, dbModels.ErrUsernameNotMachRequirements) {
			http.Error(w, " Username dose not match the requirements.", http.StatusUnprocessableEntity)
		} else if errors.Is(err, dbModels.ErrPasswordNotMachRequirements) {
			http.Error(w, " Password dose not match the requirements.", http.StatusUnprocessableEntity)
		} else if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			http.Error(w, "Password to long.", http.StatusUnprocessableEntity)
		} else {
			log.Printf("Error type: %T, Details: %s", err, err.Error())
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// Delete a user by ID
// (DELETE /users/{userId})
func (server Server) DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		http.Error(w, "Insufficient permission.", http.StatusForbidden)
		return
	}
	err = server.DB.DeleteUserByID(userId, r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Get a user by ID
// (GET /users/{userId})
func (server Server) GetUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		http.Error(w, "Insufficient permission.", http.StatusForbidden)
		return
	}
	resp, err := server.DB.GetUserByID(userId, r.Context())
	if err != nil {
		if errors.Is(err, dbModels.ErrUserNotFound) {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Update a user by ID
// (PUT /users/{userId})
func (server Server) PutUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if userId == -1 {
		userId = *user.Id
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		if userId != *user.Id {
			http.Error(w, "Insufficient permission.", http.StatusForbidden)
			return
		}
	}
	var JSONRequestBody gen.PutUsersUserIdJSONRequestBody
	err = json.NewDecoder(r.Body).Decode(&JSONRequestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	update_user := false
	if JSONRequestBody.Name != nil {
		update_user = true
		user.Name = *JSONRequestBody.Name
	}
	if JSONRequestBody.Email != nil {
		update_user = true
		user.Email = JSONRequestBody.Email
	}
	if JSONRequestBody.DefaultChoiceID != nil {
		update_user = true
		user.DefaultChoice.Id = JSONRequestBody.DefaultChoiceID
	}
	if update_user {
		err = server.DB.UpdateUser(user, r.Context())
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(JSONRequestBody)
}

// Returns currently logged in user.
// (GET /currentUser)
func (server Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

// Update the untisAcc of the active user
// (PUT /user/untisAcc)
func (server Server) PutUserUntisAcc(w http.ResponseWriter, r *http.Request) {
	user, claims, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	var JSONRequestBody gen.PutUserUntisAccJSONBody
	err = json.NewDecoder(r.Body).Decode(&JSONRequestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	key, err := base64.StdEncoding.DecodeString(claims.CryptoKey)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	err = server.DB.UpdateUntisLogin(user, *JSONRequestBody.UserName, *JSONRequestBody.Forename, *JSONRequestBody.Surname, *JSONRequestBody.UntisPWD, key, r.Context())
	if err != nil {
		if errors.Is(err, dbModels.ErrUserNotFound) {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}
		if errors.Is(err, untisDataCollectors.ErrStudentNotFound) {
			http.Error(w, "Student not found!", http.StatusNotFound)
			return
		}
		var rpcError *untisApiStructs.RPCError
		if errors.As(err, &rpcError) {
			if rpcError.Code == -8504 {
				http.Error(w, "Bad Untis credentials.", http.StatusUnprocessableEntity)
			} else {
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
			}
			return
		}
		print(err.Error())
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
