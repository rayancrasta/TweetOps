package models

type User struct {
	UserID         int    `gorm:"primaryKey;autoIncrement"`
	Username       string `gorm:"unique"`
	Password       string
	FollowersCount int
	FollowingCount int
	Verified       bool
	Accountlang    string
}

type Follower struct {
	UserID     int `gorm:"foreignKey:UserID"`
	FollowerID int `gorm:"foreignKey:UserID"`
}

type Following struct {
	UserID      int `gorm:"foreignKey:UserID"`
	FollowingID int `gorm:"foreignKey:UserID"`
}
