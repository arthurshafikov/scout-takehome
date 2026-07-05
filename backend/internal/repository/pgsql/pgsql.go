package pgsql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arthurshafikov/scout-takehome/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultConnectToDatabaseAttempts = 15

func ConnectToDatabase(ctx context.Context, databaseConfig *config.DBConfig, isDebug bool) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
		databaseConfig.Host,
		databaseConfig.User,
		databaseConfig.Password,
		databaseConfig.DBName,
		databaseConfig.Port,
		databaseConfig.SSLMode,
		databaseConfig.TimeZone,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	if isDebug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var db *gorm.DB
	var err error
	for i := 1; i <= defaultConnectToDatabaseAttempts; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			if i == defaultConnectToDatabaseAttempts {
				log.Fatalln(fmt.Errorf("connect to db error: %w", err))
			}

			time.Sleep(time.Second)
			continue
		}

		break
	}

	return db
}
