package entities

var (
	Process = "Berlangsung"
	Done    = "Berhasil"
	Failed  = "Tidak Berhasil"
)

type Status struct {
	ID     uint
	Status string `gorm:"type:VARCHAR(30)"`
}
