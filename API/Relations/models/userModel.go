package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID         int    `gorm:"primaryKey;autoIncrement"`
	Username       string `gorm:"unique"`
	Password       string
	FollowersCount int
	FollowingCount int
}

type Follower struct {
	UserID     int `gorm:"foreignKey:UserID"`
	FollowerID int `gorm:"foreignKey:UserID"`
}

type Following struct {
	UserID      int `gorm:"foreignKey:UserID"`
	FollowingID int `gorm:"foreignKey:UserID"`
}
