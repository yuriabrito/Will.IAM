package repositories

import "github.com/ghostec/Will.IAM/models"

// Services repository
type Services interface {
	Create(*models.Service) error
}

type services struct {
	storage *Storage
}

func (ss services) Create(s *models.Service) error {
	_, err := ss.storage.PG.DB.Query(
		s, `INSERT INTO services (name, permission_name,
		creator_service_account_id) VALUES (?name, ?permission_name,
		?creator_service_account_id) RETURNING id`, s,
	)
	return err
}

// NewServices services ctor
func NewServices(s *Storage) Services {
	return &services{storage: s}
}
