package main

import (
	"gin/product_service/db"
)

func main() {
	db.ConnectDatabase()
}
