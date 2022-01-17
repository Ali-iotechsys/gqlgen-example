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

type Resolver struct {
	groups         []*model.Group
	users          []*model.User
	groupObservers map[string]struct {
		Group chan *model.Group
	}
	userObservers map[string]struct {
		User chan *model.User
	}
	mu sync.Mutex
}

func New() generated.Config {
	return generated.Config{
		Resolvers: &Resolver{
			groupObservers: map[string]struct{ Group chan *model.Group }{},
			userObservers:  map[string]struct{ User chan *model.User }{},
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
