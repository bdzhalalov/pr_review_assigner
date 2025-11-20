package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"

	"github.com/bdzhalalov/pr-review-assigner/config"
)

type Database struct {
	Connection *gorm.DB
	logger     *logrus.Logger
}

func ConnectToDB(config *config.Config, logger *logrus.Logger) (*Database, error) {
	databaseURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)

	db, err := gorm.Open("mysql", databaseURI)

	if err != nil {
		logger.Errorf("Failed to connect to the database: %v", err)
		return nil, err
	}

	logger.Info("Successfully connected to the database")

	return &Database{
		Connection: db,
		logger:     logger,
	}, nil
}

func (db *Database) Close() error {
	if db.Connection != nil {
		return db.Connection.Close()
	}

	db.logger.Info("Connection to database was closed")
	return nil
}
