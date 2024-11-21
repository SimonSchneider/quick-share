package main

import (
	"context"
	"log"
	"os"
)

func main() {
	if err := Run(context.Background(), os.Args, os.Stdin, os.Stdout, os.Stderr, os.Getenv, os.Getwd); err != nil {
		log.Fatalf("Critical: %v", err)
	}
}
