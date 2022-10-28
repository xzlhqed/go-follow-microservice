package data

import (
	"fmt"
	"encoding/json"
	"io"
)

var usersList = []*User {}

type User struct {
	ID        int   `json:"id"`
	Followers []int `json:"followers"`
	Following []int `json:"following"`
}

func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

func getBool(flag bool, _ int) bool {
	return flag
}

func (u1 *User) isFollowing(u2 *User) (bool, int) {
	for idx, b := range u1.Following {
		if b == u2.ID {
			return true, idx
		}
	}
	return false, -1
}

func (u1 *User) isFollowed(u2 *User) (bool, int) {
	for idx, b := range u1.Followers {
		if b == u2.ID {
			return true, idx
		}
	}
	return false, -1
}

func GetUsers() []*User {
	return usersList
}

func AddUser(u *User) {
	usersList = append(usersList, u)
}

func FollowUser(id_1 int, id_2 int) {
	_, pos_1, _ := findUser(id_1)
	_, pos_2, _ := findUser(id_2)

	user_1 := usersList[pos_1]
	user_2 := usersList[pos_2]

	if getBool(user_1.isFollowing(user_2)) {
		return
	}
	
	user_1.Following = append(user_1.Following, user_2.ID)
	user_2.Followers = append(user_2.Followers, user_1.ID)
}

func UnfollowUser(id_1 int, id_2 int) {
	_, pos_1, _ := findUser(id_1)
	_, pos_2, _ := findUser(id_2)

	user_1 := usersList[pos_1]
	user_2 := usersList[pos_2]

	flag, idx_1 := user_1.isFollowing(user_2)
	_,    idx_2 := user_2.isFollowed(user_1)

	if flag {
		user_1.Following = append(user_1.Following[:idx_1], user_1.Following[idx_1+1:]...)
		user_2.Followers = append(user_2.Followers[:idx_2], user_2.Followers[idx_2+1:]...)
	}
}

func UpdateUser(id int, u *User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	usersList[pos] = u

	return nil
}

var ErrUserNotFound = fmt.Errorf("User not found")

func findUser(id int) (*User, int, error) {
	for i, u := range usersList {
		if u.ID == id {
			return u, i, nil
		}
	}
	return nil, -1, ErrUserNotFound
} 