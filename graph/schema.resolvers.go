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
	// Add new Group
	group := &model.Group{
		Text: input.Text,
		ID:   fmt.Sprintf("G%d", rand.Int()),
	}
	r.mu.Lock()
	r.groups = append(r.groups, group)
	r.mu.Unlock()

	// Notify all group observers
	for _, gObserver := range r.groupObservers {
		gObserver.Group <- group
	}
	return group, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	// Add new User
	user := &model.User{
		Name:    input.Name,
		Address: input.Address,
		ID:      fmt.Sprintf("U%d", rand.Int()),
	}
	r.mu.Lock()
	r.users = append(r.users, user)
	r.mu.Unlock()

	// Notify all group observers
	for _, uObserver := range r.userObservers {
		uObserver.User <- user
	}
	return user, nil
}

func (r *mutationResolver) AssociateUserToGroup(ctx context.Context, input *model.NewAssociate) (*model.Group, error) {
	// Add new User-Group association
	r.mu.Lock()
	defer r.mu.Unlock()
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

func (r *subscriptionResolver) UserCreated(ctx context.Context) (<-chan *model.User, error) {
	id := randString(8)
	userEvents := make(chan *model.User, 1)

	go func() {
		// Un-register the user event
		<-ctx.Done()
		r.mu.Lock()
		delete(r.userObservers, id)
		fmt.Printf("deleted user events '%s'\n", id)
		r.mu.Unlock()
	}()
	// Register new user event
	r.userObservers[id] = struct{ User chan *model.User }{User: userEvents}

	return userEvents, nil
}

func (r *subscriptionResolver) GroupCreated(ctx context.Context) (<-chan *model.Group, error) {
	id := randString(8)
	groupEvents := make(chan *model.Group, 1)

	go func() {
		// Un-register the group event
		<-ctx.Done()
		r.mu.Lock()
		delete(r.groupObservers, id)
		fmt.Printf("deleted group events '%s'\n", id)
		r.mu.Unlock()
	}()

	// Register new group event
	r.groupObservers[id] = struct{ Group chan *model.Group }{Group: groupEvents}

	return groupEvents, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
