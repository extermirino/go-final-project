package main

import (
	"fmt"
	"go1f/pkg/db"
	"go1f/pkg/server"
	"os"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = server.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
