package models

import (
    "github.com/pkg/errors"
    uuid "github.com/satori/go.uuid"
    "gorm.io/gorm"
)

type User struct {
    ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
    Username     string
    PasswordHash []byte
}

func GetUserByUsername(db *gorm.DB, username string) (User, error) {
	var user User

    if err := db.Where("username = ?", username).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return user, errors.Wrap(err, "user not found")
        }
        return user, errors.Wrap(err, "database query error")
    }

    return user, nil
}
