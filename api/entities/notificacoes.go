package entities
import "time"
type Notificacao struct {
	ID            int64
	DispositivoID int64
	Tipo          string
	Titulo        string
	Body          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
