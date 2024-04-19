package initializers

import "Users/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Follower{}, &models.Following{})

}
