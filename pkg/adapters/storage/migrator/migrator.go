package migrator

import (
	"github.com/RezaMokaram/ExchangeService/pkg/adapters/storage/types"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&types.Outbox{},
		&types.User{},
		&types.Notification{},
		&types.Crypto{},
	)
}
