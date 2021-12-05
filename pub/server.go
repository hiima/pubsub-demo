package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"pub/graph"
	"pub/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(
		context.Background(),
		mustGetEnv("REDIS_URL"),
	)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	port := mustGetEnv("PORT")
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustGetEnv(k string) (v string) {
	v, ok := os.LookupEnv(k)
	if !ok {
		panic(fmt.Sprintf("the environment variable `$%s` must be set", k))
	}
	return
}
