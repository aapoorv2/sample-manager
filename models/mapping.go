package models

type Mapping struct {
	ID uint64 `gorm:"primaryKey"`
	Segments []string 
	SampleItemId string
	ItemId string
}