package repositories

// Healthcheck must implement Do()
type Healthcheck interface {
	Do() error
}

type healthcheck struct {
	storage *Storage
}

// NewHealthcheck returns a Healthcheck implementer
func NewHealthcheck(storage *Storage) Healthcheck {
	return &healthcheck{storage: storage}
}

func (h *healthcheck) Do() error {
	var result int
	_, err := h.storage.PG.DB.Query(&result, "SELECT 1")
	return err
}
