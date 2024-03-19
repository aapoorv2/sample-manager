package models

type Segment struct {
	ID        uint64 `gorm:"primaryKey"`
	Segment   string
	MappingID uint64
}

type Mapping struct {
	ID          uint64 `gorm:"primaryKey"`
	SampleItemID string
	ItemID       string
	Segments     []Segment `gorm:"foreignKey:MappingID"`
}