package models

import "github.com/lib/pq"


type Mapping struct {
	ID          uint64 `gorm:"primaryKey"`
	SampleItemID string
	ItemID       string
	Segments     pq.StringArray `gorm:"type:text[]"`
}