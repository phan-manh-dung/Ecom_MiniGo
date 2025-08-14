package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID uint  `gorm:"not null"` // foreign key đến User
	User   *User `gorm:"foreignKey:UserID"` // GORM tự map theo UserID
	RoleID uint  `gorm:"not null"` // foreign key đến Role
	Role   *Role `gorm:"foreignKey:RoleID"` // GORM tự map theo RoleID
}

type User struct {
	gorm.Model
	Name     string   `json:"name" binding:"required,min=2,max=50" validate:"required,min=2,max=50"`
	SDT      string   `json:"sdt" gorm:"uniqueIndex" binding:"required,len=10" validate:"required,len=10"`
	Password string   `json:"password" binding:"required,min=8" validate:"required,min=8"`
	Account  *Account `gorm:"foreignKey:UserID"` // 1-1: mỗi user có 1 account
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
