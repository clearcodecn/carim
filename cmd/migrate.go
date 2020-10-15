package main

import (
	"context"
	"github.com/clearcodecn/carim/ent"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	client, err := ent.Open("postgres", "postgres://postgres:postgres@localhost:5432/carim?sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

}