package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Ali-iotechsys/gqlgen-example/graph"
	"github.com/Ali-iotechsys/gqlgen-example/graph/generated"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "8080"
const AuthorisationHeaderKey = "Authorization"
const TokenCtxKey = "Token"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := newCustomServer(generated.NewExecutableSchema(graph.New()))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", middleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func newCustomServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			Subprotocols: []string{"graphql-ws", "graphql-transport-ws"},
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			return webSocketInit(ctx, initPayload)
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}

// webSocketInit gets the authorisation token from the websocket payload
func webSocketInit(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
	tokenStr := initPayload.GetString(AuthorisationHeaderKey) // or initPayload.Authorization()
	if len(tokenStr) > 0 {
		ctxNew := context.WithValue(ctx, TokenCtxKey, tokenStr)
		return ctxNew, nil
	}
	return ctx, nil
}

// middleware gets the authorisation token from the http request headers
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get(AuthorisationHeaderKey)
		if len(tokenStr) > 0 {
			ctxNew := context.WithValue(r.Context(), TokenCtxKey, tokenStr)
			r = r.WithContext(ctxNew)
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
