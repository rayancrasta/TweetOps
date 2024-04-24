package models

type User struct {
	UserID         int    `gorm:"primaryKey;autoIncrement"`
	Username       string `gorm:"unique"`
	Password       string
	FollowersCount int
	FollowingCount int
	Verified       bool
}
