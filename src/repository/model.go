package repository

import (
	"time"
)


type User struct {
	Id		   		uint64	`gorm:"primaryKey;autoIncrement:true;not nul"`
	Name      	   string
	UserName       string   `gorm:"unique_index:user_name_index"`
	Password       string
	Token          string
	CartId 		   uint64	`gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null;  default: current_timestamp"`
}

type Cart struct {
	Id			uint64		`gorm:"primaryKey;autoIncrement:true;not null"`
	CreatedAt	time.Time	`gorm:"not null;  default: 	current_timestamp"`
	IsPurchased	bool		`gorm:"not null;  default: 	true"`
}

type Item struct {
	Id			uint64		`gorm:"primaryKey;autoIncrement:true;not null"`	
	Name      	string		`gorm:"not null"`
	CreatedAt	time.Time	`gorm:"not null;  default: 	current_timestamp"`
}


type Order struct {
	Id				uint64		`gorm:"primaryKey;autoIncrement:true;not null"`
	CartId			uint64		`gorm:"not null"`
	UserId			uint64		`gorm:"not null"`
	CreatedAt		time.Time	`gorm:"not null;  default: 	current_timestamp"`
}


type CartItemRelation struct {
	CartId			uint64		`gorm:"not null"`
	ItemId			uint64		`gorm:"not null"`
	IsPurchased		bool		`gorm:"not null;  default: 	false"`
}

