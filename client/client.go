package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hasura/go-graphql-client"
	"time"
)

func main() {
	// Create (Query/Mutation) GraphQL client
	client1 := graphql.NewClient("http://localhost:8080/query", nil)
	if client1 == nil {
		fmt.Println("cannot create client1")
		return
	}
	go func() {
		userCount := 0
		for {
			time.Sleep(3 * time.Second)
			userCount++

			var m struct {
				CreateUser struct {
					Id      graphql.ID
					Name    graphql.String
					Address graphql.String
				} `graphql:"createUser(input: {name: $name, address: $address})"`
			}
			variables := map[string]interface{}{
				"name":    graphql.String(fmt.Sprintf("User_%d", userCount)),
				"address": graphql.String(fmt.Sprintf("This is User_%d address", userCount)),
			}
			mErr := client1.Mutate(context.Background(), &m, variables)
			if mErr != nil {
				fmt.Println(mErr)
				return
			}
			fmt.Printf("client1: created new User (ID: %s, Name: %s, Address: %s)\n",
				m.CreateUser.Id, m.CreateUser.Name, m.CreateUser.Address)
		}
	}()

	// Create (Subscription) GraphQL client
	client2 := graphql.NewSubscriptionClient("wss://localhost:8080/query")
	defer func() {
		_ = client2.Close()
	}()

	var userSub struct {
		UserCreated struct {
			Id      graphql.ID
			Name    graphql.String
			Address graphql.String
		}
	}
	userSubId, sErr := client2.Subscribe(&userSub, nil, func(dataValue *json.RawMessage, errValue error) error {
		if errValue != nil {
			return errValue
		}
		data := userSub.UserCreated
		jsonErr := json.Unmarshal(*dataValue, &data)
		if jsonErr != nil {
			return jsonErr
		}
		fmt.Printf("client2: received User (ID: %s, Name: %s, Address: %s)\n", data.Id, data.Name, data.Address)
		return nil
	})
	if sErr != nil {
		fmt.Println(sErr)
		return
	}
	fmt.Printf("clinet2: user subscribed success, subscribeId: %s\n", userSubId)

	rErr := client2.Run()
	if rErr != nil {
		fmt.Println(rErr)
		return
	}
}
