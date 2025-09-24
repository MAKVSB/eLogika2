package initializers

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDB(debug bool) {
	var err error
	dsn := GlobalAppConfig.DB_URL

	var gormLogger logger.Interface
	if debug {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel:      logger.Info,
				SlowThreshold: time.Second,
			},
		)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger,
	})

	if err != nil {
		log.Fatal("Failed to connect to database server")
	}
}
