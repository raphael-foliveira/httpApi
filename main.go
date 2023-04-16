package main

import (
	_ "github.com/lib/pq"
	"github.com/raphael-foliveira/httpApi/database"
	"github.com/raphael-foliveira/httpApi/routes"
)

func main() {
	database.Get()
	routes.Run()
}
