package models

type Mapping struct {
	ID uint64 `gorm:"primaryKey"`
	segments []string 
	sample_item_id string
	item_id string
}