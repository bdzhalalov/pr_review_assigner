package database

import (
	"fmt"
	"github.com/bdzhalalov/pr-review-assigner/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func MakeMigrations(db *gorm.DB, logger *logrus.Logger) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.PullRequest{},
	)
	if err != nil {
		return err
	}

	err = SeedTeamsAndUsers(db, logger)
	if err != nil {
		return err
	}
	return nil
}

func SeedTeamsAndUsers(db *gorm.DB, logger *logrus.Logger) error {
	err := gofakeit.Seed(time.Now().UnixNano())
	if err != nil {
		return err
	}

	var teams []models.Team
	var users []models.User

	userCounter := 1

	for i := 0; i < 20; i++ {
		team := models.Team{
			TeamName: gofakeit.Company() + fmt.Sprintf("-%d", i),
		}
		teams = append(teams, team)
	}

	if err := db.Create(&teams).Error; err != nil {
		return fmt.Errorf("failed to insert teams: %w", err)
	}

	// теперь создаём пользователей
	for _, team := range teams {
		for j := 0; j < 10; j++ {
			user := models.User{
				UserID:   fmt.Sprintf("u%d", userCounter),
				Username: gofakeit.Name(),
				IsActive: rand.Intn(2) == 1,
				TeamID:   team.ID,
			}

			users = append(users, user)
			userCounter++
		}
	}

	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}

	logger.Infof("Seeded %d teams and %d users", len(teams), len(users))
	return nil
}
