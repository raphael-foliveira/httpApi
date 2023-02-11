package main

import (
	_ "github.com/lib/pq"
	"github.com/raphael-foliveira/httpApi/database"
	"github.com/raphael-foliveira/httpApi/routes"
)

func main() {
	routes.Run(database.Start())
}
