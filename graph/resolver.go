//go:generate go get github.com/99designs/gqlgen/cmd@v0.14.0
//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"github.com/Ali-iotechsys/gqlgen-example/graph/generated"
	"github.com/Ali-iotechsys/gqlgen-example/graph/model"
	"math/rand"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type EventID = string

type UserObservers struct {
	CreateUser map[EventID]chan *model.User
	UpdateUser map[EventID]chan *model.User
}

type GroupObservers struct {
	CreateGroup map[EventID]chan *model.Group
	UpdateGroup map[EventID]chan *model.Group
}

type Resolver struct {
	groups         []*model.Group
	users          []*model.User
	groupObservers GroupObservers
	userObservers  UserObservers
	mu             sync.Mutex
}

func New() generated.Config {
	return generated.Config{
		Resolvers: &Resolver{
			groupObservers: GroupObservers{
				CreateGroup: map[EventID]chan *model.Group{},
				UpdateGroup: map[EventID]chan *model.Group{},
			},
			userObservers: UserObservers{
				CreateUser: map[EventID]chan *model.User{},
				UpdateUser: map[EventID]chan *model.User{},
			},
		},
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
