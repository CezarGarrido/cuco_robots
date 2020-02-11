package entities
import "time"
type Dispositivo struct {
	ID           int64
	AlunoID      int64
	Model        string
	Platform     string
	UUID         string
	Version      string
	Manufacturer string
	IsVirtual    string
	Serial       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
