package main

import (
	"github.com/kien14502/ecommerce-be/internal/routers"
)

func main() {
	r := routers.AppRouter()

	r.Run(":8000") 
}
