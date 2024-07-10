package main

import "product_service/internal/server"

func main() {
	s := server.New()
	s.Init()
	s.Run()
}
