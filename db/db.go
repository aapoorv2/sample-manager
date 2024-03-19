package database

import (
	"fmt"
	"sample-manager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	username = "abhyudaya.a_ftc"
	password = "password"
	dbName   = "samplemanager"
	sslMode  = "disable"
)

func Connection() *gorm.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbName, sslMode)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	log.Println("Connected to the database")

	err = db.AutoMigrate(&models.Mapping{}, &models.Segment{})
	if err != nil {
		log.Fatalf("Error migrating Mapping table: %v", err)
	}

	return db
}

func GetSampleID(db *gorm.DB, segments []string, itemID string) (string, error) {
	var sample_item_id string

    result := db.Model(&models.Mapping{}).
	Select("sample_item_id").
	Joins("JOIN segments ON segments.mapping_id = mappings.id").
	Where("mappings.item_id = ?", itemID).
	First(&sample_item_id)

	if result.Error != nil {
		return "", result.Error
	}
	return sample_item_id, nil
}
