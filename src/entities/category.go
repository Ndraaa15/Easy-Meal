package entities

var (
	Fruits     = "Buah"
	Vegetables = "Sayur"
)

type Category struct {
	ID   uint
	Name string `gorm:"type:VARCHAR(30)"`
}
