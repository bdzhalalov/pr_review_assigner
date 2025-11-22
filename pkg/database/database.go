package database

import (
	"fmt"
	"github.com/bdzhalalov/pr-review-assigner/config"
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
)

func ConnectToDB(config *config.Config, logger *logrus.Logger) (*gorm.DB, error) {
	databaseURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)

	db, err := gorm.Open(mysql.Open(databaseURI), &gorm.Config{
		Logger: dbLogger.Default.LogMode(dbLogger.Silent),
	})

	if err != nil {
		logger.Errorf("Failed to connect to the database: %v", err)
		return nil, err
	}

	logger.Info("Successfully connected to the database")

	err = db.AutoMigrate(models.User{}, models.Team{})
	if err != nil {
		logger.Errorf("Failed to migrate tables: %v", err)
		return nil, err
	}

	return db, nil
}
