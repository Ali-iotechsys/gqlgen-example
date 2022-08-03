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

// CreateGroup is the resolver for the createGroup field.
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
	for _, cgObserver := range r.groupObservers.CreateGroup {
		cgObserver <- group
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

	// Notify all user observers
	for _, cuObserver := range r.userObservers.CreateUser {
		cuObserver <- user
	}
	return user, nil
}

func (r *mutationResolver) AssociateUserToGroup(ctx context.Context, input model.NewAssociate) (*model.Group, error) {
	// Add new User-Group association
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.ID == input.UserID {
			for _, g := range r.groups {
				if g.ID == input.GroupID {
					g.Users = append(g.Users, u)
					// Notify the update group observer (if existed)
					if ugObserver, ok := r.groupObservers.UpdateGroup[g.ID]; ok {
						ugObserver <- g
					}
					return g, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("associate Error")
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UserUpdate) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == input.UserID {
			if input.NewName != nil {
				u.Name = *input.NewName
			}
			if input.NewAddress != nil {
				u.Address = *input.NewAddress
			}
			// Notify the update user observer (if existed)
			if uuObserver, ok := r.userObservers.UpdateUser[u.ID+u.Name]; ok {
				uuObserver <- u
			}
			if uuObserver, ok := r.userObservers.UpdateUser[u.ID+u.Address]; ok {
				uuObserver <- u
			}
			return u, nil
		}
	}
	return nil, fmt.Errorf("update user Error")
}

// Groups is the resolver for the groups field.
func (r *queryResolver) Groups(ctx context.Context) ([]*model.Group, error) {
	return r.groups, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.users, nil
}

// UserCreated is the resolver for the userCreated field.
func (r *subscriptionResolver) UserCreated(ctx context.Context) (<-chan *model.User, error) {
	eventID := randString(8)
	userEvents := make(chan *model.User, 1)

	go func() {
		// Un-register the user event
		<-ctx.Done()
		r.mu.Lock()
		delete(r.userObservers.CreateUser, eventID)
		fmt.Printf("deleted user events '%s'\n", eventID)
		r.mu.Unlock()
	}()

	// Register new user event
	r.userObservers.CreateUser[eventID] = userEvents

	return userEvents, nil
}

// GroupCreated is the resolver for the groupCreated field.
func (r *subscriptionResolver) GroupCreated(ctx context.Context) (<-chan *model.Group, error) {
	eventID := randString(8)
	groupEvents := make(chan *model.Group, 1)

	go func() {
		// Un-register the group event
		<-ctx.Done()
		r.mu.Lock()
		delete(r.groupObservers.CreateGroup, eventID)
		fmt.Printf("deleted group events '%s'\n", eventID)
		r.mu.Unlock()
	}()

	// Register new group event
	r.groupObservers.CreateGroup[eventID] = groupEvents

	return groupEvents, nil
}

// GroupUpdated is the resolver for the groupUpdated field.
func (r *subscriptionResolver) GroupUpdated(ctx context.Context, groupID string) (<-chan *model.Group, error) {
	eventID := groupID
	groupEvents := make(chan *model.Group, 1)

	go func() {
		// Un-register the group event
		<-ctx.Done()
		r.mu.Lock()
		delete(r.groupObservers.UpdateGroup, eventID)
		fmt.Printf("deleted group events '%s'\n", eventID)
		r.mu.Unlock()
	}()

	// Register new group event
	r.groupObservers.UpdateGroup[eventID] = groupEvents

	return groupEvents, nil
}

// OnUserNameUpdated is the resolver for the onUserNameUpdated field.
func (r *subscriptionResolver) OnUserNameUpdated(ctx context.Context, userID string, name string) (<-chan *model.User, error) {
	eventID := userID + name
	userEvents := make(chan *model.User, 1)

	go func() {
		// Un-register the group event
		<-ctx.Done()
		r.mu.Lock()
		delete(r.groupObservers.UpdateGroup, eventID)
		fmt.Printf("deleted group events '%s'\n", eventID)
		r.mu.Unlock()
	}()

	// Register new group event
	r.userObservers.UpdateUser[eventID] = userEvents

	return userEvents, nil
}

// OnUserAddressUpdated is the resolver for the onUserAddressUpdated field.
func (r *subscriptionResolver) OnUserAddressUpdated(ctx context.Context, userID string, address string) (<-chan *model.User, error) {
	eventID := userID + address
	userEvents := make(chan *model.User, 1)

	go func() {
		// Un-register the group event
		<-ctx.Done()
		r.mu.Lock()
		delete(r.groupObservers.UpdateGroup, eventID)
		fmt.Printf("deleted group events '%s'\n", eventID)
		r.mu.Unlock()
	}()

	// Register new group event
	r.userObservers.UpdateUser[eventID] = userEvents

	return userEvents, nil
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
