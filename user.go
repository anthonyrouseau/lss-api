package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type User struct {
	Id         int64  `json:"id,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Email      string `json:"email,omitempty"`
	SummonerId string `json:"summonerId,omitempty"`
}

type UserRequest struct {
	*User
	ProtectedId int `json:"id"`
}

type UserResponse struct {
	*User
}

type UserPayload struct {
	*User
}

func (u *UserRequest) Bind(r *http.Request) error {
	if u.User == nil {
		return errors.New("missing user fields")
	}
	// this is where you could clear out a protected id field
	u.ProtectedId = -1 //reset protected Id
	return nil
}

func NewUserPayload(user *User) *UserPayload {
	return &UserPayload{User: user}
}

func NewUserResponse(user *User) *UserResponse {
	resp := &UserResponse{User: user}
	if resp.User == nil {
		// try getting user from database
		// if you get a user
		// set the resp.User = NewUserPayload(user)
	}
	return resp
}

func (ur *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// pre-processing before resposne is changed to json and sent
	ur.Password = "" //make sure password isn't sent
	return nil
}

type UserListResponse []*UserResponse

func NewUserListResponse(users []*User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		list = append(list, NewUserResponse(user))
	}
	return list
}

func dbNewUser(user *User) (int64, error) {
	//create new user in database
	stmt, err := db.Prepare("INSERT INTO account(username,password,email,summonerId) VALUES(?,?,?,?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(user.Username, user.Password, user.Email, user.SummonerId)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Println("created new user")
	return id, nil
}
