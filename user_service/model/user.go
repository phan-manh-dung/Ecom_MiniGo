package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID uint  // foreign key đến User
	User   *User // GORM tự map theo UserID
	RoleID uint  // foreign key đến Role
	Role   Role
}

type User struct {
	gorm.Model
	Name    string   `json:"name" binding:"required,min=2,max=50" validate:"required,min=2,max=50"`
	SDT     string   `json:"sdt" gorm:"uniqueIndex" binding:"required,len=10" validate:"required,len=10"`
	Account *Account // 1-1: mỗi user có 1 account
}

type Role struct {
	gorm.Model
	Name     string     `gorm:"uniqueIndex"` // 'ADMIN', 'USER'
	Accounts []*Account // Một role có thể được dùng bởi nhiều account
}

/* Sử dụng *Account , *User để tránh chuỗi tham chiếu lẫn nhau vô hạn
- User chứa Account
- Account lại chứa User
->  Go không thể xử lý được kiểu dữ liệu đệ quy vòng tròn như vậy trong cùng một struct đầy đủ.
 dùng con trỏ * để tránh tham chiếu đệ quy trực tiếp*/
