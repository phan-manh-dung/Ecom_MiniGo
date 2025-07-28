package main

import (
	"gin/user_service/db"
)

func main() {
	db.ConnectDatabase()
}
