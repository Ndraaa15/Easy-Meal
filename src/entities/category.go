package entities

var (
	Fruits     = "Fruit"
	Vegetables = "Vegetables"
)

type Category struct {
	ID   uint
	Name string `gorm:"type:VARCHAR(30)"`
}
