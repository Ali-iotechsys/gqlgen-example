// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Group struct {
	ID    string  `json:"id"`
	Text  string  `json:"text"`
	Users []*User `json:"users"`
}

type NewAssociate struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
}

type NewGroup struct {
	Text string `json:"text"`
}

type NewUser struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type UserUpdate struct {
	UserID     string  `json:"userID"`
	NewName    *string `json:"newName"`
	NewAddress *string `json:"newAddress"`
}
