package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"context"
	"github.com/gorilla/mux"
	"github.com/xzlhqed/golang-follow-microservice/data"
)

type Users struct {
	l *log.Logger
}

func ReturnUsers(l *log.Logger) *Users {
	return &Users{l}
}

func stringToInt(id string, rw http.ResponseWriter, r *http.Request) int {
	return_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}
	return return_id
}

func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	userList := data.GetUsers()
	d, err := json.Marshal(userList)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	rw.Write(d)
}

func (u *Users) GetUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	id := stringToInt(vars["id"], rw, r)

	userList := data.GetUsers()
	for _, us := range userList {
		if id == us.ID {
			output, err := json.Marshal(us)
			if err != nil {
				http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
			}
			rw.Write(output)
			return
		}
	}

	u.l.Printf("User does not exist")
}

func (u *Users) AddUser(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(KeyUser{}).(data.User)
	
	userList := data.GetUsers()
	for _, us := range userList {
		if user.ID == us.ID {
			u.l.Printf("User already exists")
			return
		}
	}

	data.AddUser(&user)
}

func (u *Users) FollowUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id_1 := stringToInt(vars["id1"], rw, r)
	id_2 := stringToInt(vars["id2"], rw, r)

	data.FollowUser(id_1, id_2)

}

func (u *Users) UnfollowUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id_1 := stringToInt(vars["id1"], rw, r)
	id_2 := stringToInt(vars["id2"], rw, r)

	data.UnfollowUser(id_1, id_2)

}

type KeyUser struct{}

func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := data.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "unable to read user", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}