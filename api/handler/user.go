package handler

import (
	"encoding/json"
	"go_rest_api_with_mysql/api/presenter"
	"go_rest_api_with_mysql/entity"
	logger "go_rest_api_with_mysql/pkg/log"
	"go_rest_api_with_mysql/usecase/user"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *user.Service
}

var log *zap.SugaredLogger = logger.GetLogger().Sugar()

func RegisterUserHandler(r *mux.Router, userService *user.Service) {
	handler := &UserHandler{
		userService: userService,
	}
	r.HandleFunc("/user", handler.getUserList).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user", handler.createUser).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/{id}", handler.UpdateUser).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc("/user/{id}", handler.getUser).Methods(http.MethodGet, http.MethodOptions)
}

// get user list
func (handler *UserHandler) getUserList(w http.ResponseWriter, r *http.Request) {
	ConfigureCorsHeader(w, r, "*", "*")
	if r.Method == http.MethodOptions {
		return
	}

	userList, err := handler.userService.GetUsers()
	if err != nil {
		log.Errorf("Error retrieving user list. %v", err)
		return
	}

	responsePayload := make([]*presenter.UserResponse, len(userList))

	for i, userData := range userList {
		responsePayload[i] = &presenter.UserResponse{
			ID:        userData.ID,
			Email:     userData.Email,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			CreatedAt: userData.CreatedAt,
		}
	}

	response, err := json.Marshal(responsePayload)
	if err != nil {
		log.Errorf("Error parsing user data to JSON. %v", err)
		processResponseErrorStatus(w, err, http.StatusExpectationFailed)
		_, _ = w.Write([]byte("Error parsing app configuration data to JSON"))
		return
	}

	addDefaultHeaders(w, r)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// create a user
func (handler *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ConfigureCorsHeader(w, r, "*", "*")
	if r.Method == http.MethodOptions {
		return
	}

	createUserReq := &presenter.UserCreateRequest{}
	_ = json.NewDecoder(r.Body).Decode(&createUserReq)

	userData := &entity.User{
		Email:     createUserReq.Email,
		FirstName: createUserReq.FirstName,
		LastName:  createUserReq.LastName,
		Password:  createUserReq.Password,
	}

	id, err := handler.userService.CreateUser(userData)
	if err != nil {
		log.Errorf("Error creating user details. %v", err)
		processResponseErrorStatus(w, err, http.StatusExpectationFailed)
		_, _ = w.Write([]byte("User creation failed."))
		return
	}
	responsePayload := &presenter.UserResponse{
		ID: id,
	}

	response, err := json.Marshal(responsePayload)
	if err != nil {
		log.Errorf("Error parsing user creation response data to JSON. %v", err)
		processResponseErrorStatus(w, err, http.StatusExpectationFailed)
		_, _ = w.Write([]byte("Error processing user creation"))
		return
	}

	addDefaultHeaders(w, r)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// update a user
func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ConfigureCorsHeader(w, r, "*", "*")
	if r.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Errorf("Error parsing uuid. %v", err)
		return
	}

	updateUserReq := &presenter.UserUpdateRequest{}
	_ = json.NewDecoder(r.Body).Decode(&updateUserReq)

	userData := &entity.User{
		Email:     updateUserReq.Email,
		FirstName: updateUserReq.FirstName,
		LastName:  updateUserReq.LastName,
	}

	err = handler.userService.UpdateUser(userData, userId)
	if err != nil {
		log.Errorf("Error updating user details. %v", err)
		processResponseErrorStatus(w, err, http.StatusExpectationFailed)
		_, _ = w.Write([]byte("User updating failed."))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("User details update successful."))
}

// get a user by id
func (handler *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	ConfigureCorsHeader(w, r, "*", "*")
	if r.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(r)
	userId, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Errorf("Error parsing uuid. %v", err)
		return
	}

	user, err := handler.userService.GetUser(userId)
	if err != nil {
		log.Errorf("Error retrieving user data. %v", err)
		return
	}

	responsePayload := &presenter.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}

	response, err := json.Marshal(responsePayload)
	if err != nil {
		log.Errorf("Error parsing user data to JSON. %v", err)
		processResponseErrorStatus(w, err, http.StatusExpectationFailed)
		_, _ = w.Write([]byte("Error parsing app configuration data to JSON"))
		return
	}

	addDefaultHeaders(w, r)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
