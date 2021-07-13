package domain

type Menu struct {
	ID      uint32 `gorm:"primaryKey;autoIncrement;index"`
	Name    string
	Protein float64
	Fat     float64
	Carbs   float64
}
