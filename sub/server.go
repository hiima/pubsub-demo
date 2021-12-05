package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sub/graph"
	"sub/graph/generated"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

func main() {
	redisURL := mustGetEnv("REDIS_URL")

	// WebSocket を使う場合はトランスポートの追加を自分で行う必要があるため、
	// `handler.NewDefaultServer` ではなく `handler.New` で初期化すること
	// https://vallettaio.hatenablog.com/entry/2020/05/10/024551
	srv := handler.New(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: graph.NewResolver(
					context.Background(),
					redisURL,
				),
			},
		),
	)
	// Subscription を使うには `transport.POST` と `transport.Websocket` をトランスポートとして追加すればよい
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	// CORS の設定
	c := cors.New(cors.Options{
		// プロダクトコードであれば正しくオリジンを設定すること
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", c.Handler(srv))

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
