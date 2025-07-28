package main

import (
	"gin/order_service/db"
)

func main() {
	db.ConnectDatabase()
}
