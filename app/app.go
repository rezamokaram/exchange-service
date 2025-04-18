package app

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/rezamokaram/exchange-service/config"
	"github.com/rezamokaram/exchange-service/internal/common"
	"github.com/rezamokaram/exchange-service/internal/crypto"
	cryptoPort "github.com/rezamokaram/exchange-service/internal/crypto/port"
	"github.com/rezamokaram/exchange-service/internal/notification"
	notifPort "github.com/rezamokaram/exchange-service/internal/notification/port"
	"github.com/rezamokaram/exchange-service/internal/user"
	userPort "github.com/rezamokaram/exchange-service/internal/user/port"
	"github.com/rezamokaram/exchange-service/pkg/adapters/storage"
	"github.com/rezamokaram/exchange-service/pkg/adapters/storage/migrator"
	"github.com/rezamokaram/exchange-service/pkg/cache"
	"github.com/rezamokaram/exchange-service/pkg/postgres"

	redisAdapter "github.com/rezamokaram/exchange-service/pkg/adapters/cache"

	"gorm.io/gorm"

	appCtx "github.com/rezamokaram/exchange-service/pkg/context"
)

type app struct {
	db                  *gorm.DB
	cfg                 config.ExchangeConfig
	cryptoService       cryptoPort.Service
	userService         userPort.Service
	notificationService notifPort.Service
	redisProvider       cache.Provider
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) CryptoService(ctx context.Context) cryptoPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.cryptoService == nil {
			a.cryptoService = a.cryptoServiceWithDB(a.db)
		}
		return a.cryptoService
	}

	return a.cryptoServiceWithDB(db)
}

func (a *app) cryptoServiceWithDB(db *gorm.DB) cryptoPort.Service {
	return crypto.NewService(storage.NewCryptoRepo(db, true, a.redisProvider))
}

func (a *app) UserService(ctx context.Context) userPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.userService == nil {
			a.userService = a.userServiceWithDB(a.db)
		}
		return a.userService
	}

	return a.userServiceWithDB(db)
}

func (a *app) userServiceWithDB(db *gorm.DB) userPort.Service {
	return user.NewService(storage.NewUserRepo(db, true, a.redisProvider))
}

func (a *app) notifServiceWithDB(db *gorm.DB) notifPort.Service {
	return notification.NewService(storage.NewNotificationRepo(db, a.redisProvider),
		user.NewService(storage.NewUserRepo(db, true, a.redisProvider)), storage.NewOutboxRepo(db))
}

func (a *app) NotificationService(ctx context.Context) notifPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.notificationService == nil {
			a.notificationService = a.notifServiceWithDB(a.db)
		}
		return a.notificationService
	}

	return a.notifServiceWithDB(db)
}

func (a *app) Config() config.ExchangeConfig {
	return a.cfg
}

func (a *app) setDB() error {
	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		User:   a.cfg.DB.User,
		Pass:   a.cfg.DB.Password,
		Host:   a.cfg.DB.Host,
		Port:   a.cfg.DB.Port,
		DBName: a.cfg.DB.Database,
		Schema: a.cfg.DB.Schema,
	})

	if err != nil {
		return err
	}

	if err := migrator.Migrate(db); err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *app) setRedis() {
	a.redisProvider = redisAdapter.NewRedisProvider(fmt.Sprintf("%s:%d", a.cfg.Redis.Host, a.cfg.Redis.Port))
}

func NewApp(cfg config.ExchangeConfig) (App, error) {
	a := &app{
		cfg: cfg,
	}

	if err := a.setDB(); err != nil {
		return nil, err
	}

	a.setRedis()

	return a, a.registerOutboxHandlers()
}

func NewMustApp(cfg config.ExchangeConfig) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}

func (a *app) registerOutboxHandlers() error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	common.RegisterOutboxRunner(a.notifServiceWithDB(a.db), scheduler)

	scheduler.Start()

	return nil
}
