package model

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100;not null"`
	BobotPertama int    `gorm:"not null"`
	BobotKedua   int    `gorm:"not null"`
	BobotKetiga  int    `gorm:"not null"`
	LabelPertama string `gorm:"size:100;not null"`
	LabelKedua   string `gorm:"size:100;not null"`
	LabelKetiga  string `gorm:"size:100;not null"`
	Nilai        int    `gorm:"not null"`
}
