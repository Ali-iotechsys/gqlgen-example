# gqlgen-example
GraphQL Example

- Run `$ go generate ./...`
- Run `$ go run server.go`
- Open http://localhost:8080 in a browser.
- Try the following queries:

```
query viewUsers {
  users {
    id
    name
    address
  }
}

query viewGroups {
  groups {
    id
    text
    users {
      id
      name
    }
  }
}

mutation createUserA {
  createUser(input: {name: "UserA", address: "This is UserA address"}) {
    id
    name
    address
  }
}

mutation createUserB {
  createUser(input: {name: "UserB", address: "This is UserB address"}) {
    id
    name
    address
  }
}

mutation createGroup1 {
  createGroup(input: {text: "Group1", }) {
    id
    text
  }
}

mutation createGroup2 {
  createGroup(input: {text: "Group2", }) {
    id
    text
  }
}

mutation associateUserToGroup {
  associateUserToGroup(input: {userId: "U5577006791947779410", groupId: "G6129484611666145821"}) {
    id
    text
    users {
      id
      name
    }
  }
}
```
