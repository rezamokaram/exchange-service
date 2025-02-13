package postgres

import (
	"fmt"
	"os"

	"github.com/RezaMokaram/ExchangeService/config"
	"github.com/RezaMokaram/ExchangeService/models"
	bankModels "github.com/RezaMokaram/ExchangeService/models/bank"
	cryptoModels "github.com/RezaMokaram/ExchangeService/models/crypto"
	tradeModels "github.com/RezaMokaram/ExchangeService/models/trade"
	userModels "github.com/RezaMokaram/ExchangeService/models/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDatabase(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Port,
		cfg.SSLMode,
		cfg.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		fmt.Printf("migrations failed: %v\n", err.Error())
	}

	if !hasTestData(db) {
		sqlFile, err := os.ReadFile("./main-data.sql")
		if err != nil {
			fmt.Printf("reading sql dump file failed: %v\n", err.Error())
		}

		result := db.Exec(string(sqlFile))
		if result.Error != nil {
			fmt.Printf("executing sql dump file failed: %v\n", result.Error)
		}

		fmt.Println("Fake Data Inserted.")
	}

	return db, nil
}

type DBConnOptions struct {
	User   string
	Pass   string
	Host   string
	Port   uint
	DBName string
	Schema string
}

func (o DBConnOptions) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		o.Host, o.Port, o.User, o.Pass, o.DBName, o.Schema)
}

func NewPsqlGormConnection(opt DBConnOptions) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(opt.PostgresDSN()), &gorm.Config{
		Logger: logger.Discard,
	})
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&userModels.User{},
		&userModels.Profile{},
		&cryptoModels.Crypto{},
		&bankModels.PaymentInfo{},
		&models.Transaction{},
		&bankModels.BankingInfo{},
		&models.SupportTicket{},
		&models.TicketMessage{},
		&tradeModels.OpenTrade{},
		&tradeModels.ClosedTrade{},
		&tradeModels.FutureOrder{},
	)
}

func hasTestData(db *gorm.DB) bool {
	var count int64
	db.Model(&userModels.User{}).Count(&count)
	return count > 0
}
