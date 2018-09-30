package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

type User struct {
	Id         int64  `json:"id,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Email      string `json:"email,omitempty"`
	SummonerId int    `json:"summonerId,omitempty"`
}

type UserRequest struct {
	*User
	ProtectedId  int64  `json:"id"`
	SummonerName string `json:"summonerName"`
	Code         string `json:"code"`
}

type UserResponse struct {
	*User
}

type UserPayload struct {
	*User
}

type RiotResponse struct {
	SummonerId int `json:"id"`
}

func (u *UserRequest) Bind(r *http.Request) error {
	if u.User == nil {
		return errors.New("missing user fields")
	}

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

func checkSummonerId(summonerName, code string) (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://na1.api.riotgames.com/lol/summoner/v3/summoners/by-name/"+summonerName, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("X-Riot-Token", os.Getenv("riotapikey"))
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, errors.New("could not find summoner")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	summonerInfo := &RiotResponse{}
	err = json.Unmarshal(body, summonerInfo)
	if err != nil {
		return 0, err
	}
	req, err = http.NewRequest("GET", "https://na1.api.riotgames.com/lol/platform/v3/third-party-code/by-summoner/"+strconv.Itoa(summonerInfo.SummonerId), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("X-Riot-Token", os.Getenv("riotapikey"))
	resp, err = client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, errors.New("code does not match")
	}
	var userCode string
	// body here is just a string so probably don't need to json.Unmarshal
	body, err = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &userCode)
	if err != nil {
		return 0, err
	}
	if strings.Compare(userCode, code) != 0 {
		return 0, errors.New("code does not match")
	}
	return summonerInfo.SummonerId, nil
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

func dbSearchUsername(searchValue string, offset string) ([]*User, error) {
	var users []*User
	rows, err := db.Query("SELECT id, username, summonerId FROM account WHERE username LIKE CONCAT(?,'%') LIMIT 10 OFFSET ?", searchValue, offset)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.SummonerId)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}
	return users, nil
}
