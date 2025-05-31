package middleware

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SetCurrentUserDB(db *gorm.DB, userID uuid.UUID) *gorm.DB {
	db.Exec("SET app.current_user_id = ?", userID.String())
	return db
}
