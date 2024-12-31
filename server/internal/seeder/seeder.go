package seeder

import (
	"fmt"
	"funcedup/internal/schema"
)

func (d *Domain) seedUsers() error {
	users := []schema.User{
		{
			Username:     "michael",
			Email:        "michael.chen@elmntri.com",
			PasswordHash: "testtesttest",
			Points:       0,
		},
		{
			Username:     "alan",
			Email:        "vimalan.renganattan@elmntri.com",
			PasswordHash: "testtesttest",
			Points:       0,
		},
		{
			Username:     "jeff",
			Email:        "jeff.hsu@elmntri.com",
			PasswordHash: "testtesttest",
			Points:       0,
		},
	}

	for _, user := range users {
		var existingUser schema.User

		if err := d.params.DB.GetDB().Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			continue
		}

		if err := d.params.DB.GetDB().Create(&user).Error; err != nil {
			fmt.Printf("failed to seed user: %v", err)
		}
	}

	return nil
}
