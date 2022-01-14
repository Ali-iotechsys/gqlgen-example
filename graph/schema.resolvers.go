package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Ali-iotechsys/gqlgen-example/graph/generated"
	"github.com/Ali-iotechsys/gqlgen-example/graph/model"
)

func (r *mutationResolver) CreateGroup(ctx context.Context, input model.NewGroup) (*model.Group, error) {
	group := &model.Group{
		Text: input.Text,
		ID:   fmt.Sprintf("G%d", rand.Int()),
	}
	r.groups = append(r.groups, group)
	return group, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user := &model.User{
		Name:    input.Name,
		Address: input.Address,
		ID:      fmt.Sprintf("U%d", rand.Int()),
	}
	r.users = append(r.users, user)
	return user, nil
}

func (r *mutationResolver) AssociateUserToGroup(ctx context.Context, input *model.NewAssociate) (*model.Group, error) {
	for _, u := range r.users {
		if u.ID == input.UserID {
			for _, g := range r.groups {
				if g.ID == input.GroupID {
					g.Users = append(g.Users, u)
					return g, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("associate Error")
}

func (r *queryResolver) Groups(ctx context.Context) ([]*model.Group, error) {
	return r.groups, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
