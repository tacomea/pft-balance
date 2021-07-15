package domain

type Menu struct {
	ID      string `gorm:"primaryKey;index"`
	Name    string
	Protein float64
	Fat     float64
	Carbs   float64
}
