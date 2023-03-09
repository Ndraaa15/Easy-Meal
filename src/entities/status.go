package entities

var (
	Process = "Sedang Diproses"
	Paid    = "Sudah Dibayar"
	Done    = "Sudah Diambil"
)

type Status struct {
	ID   uint
	Name string `gorm:"type:VARCHAR(30)"`
}
