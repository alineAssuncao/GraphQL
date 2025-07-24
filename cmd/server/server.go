package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/13-GraphQL/graph"
	"github.com/devfullcycle/13-GraphQL/internal/database"
	"github.com/vektah/gqlparser/v2/ast"
	_ "modernc.org/sqlite" // Import SQLite driver
)

const defaultPort = "8080"

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT
		)
	`)
	if err != nil {
		log.Fatalf("erro ao criar tabela de categorias: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS courses (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			category_id TEXT
		)
	`)
	if err != nil {
		log.Fatalf("erro ao criar tabela de courses: %v", err)
	}

	categoryDB := database.NewCategory(db)
	courseDB := database.NewCourse(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CategoryDB: categoryDB,
		CourseDB:   courseDB,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
