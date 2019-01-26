package repositories

import "github.com/ghostec/Will.IAM/models"

// Services repository
type Services interface {
	All() ([]models.Service, error)
	Create(*models.Service) error
}

type services struct {
	storage *Storage
}

func (ss services) Create(s *models.Service) error {
	_, err := ss.storage.PG.DB.Query(
		s, `INSERT INTO services (name, permission_name, service_account_id,
		creator_service_account_id, am_url) VALUES (?name, ?permission_name,
		?service_account_id, ?creator_service_account_id, ?am_url) RETURNING id`,
		s,
	)
	return err
}

// All returns all services in storage
func (ss services) All() ([]models.Service, error) {
	var allServices []models.Service
	if _, err := ss.storage.PG.DB.Query(
		&allServices, `SELECT * FROM services`,
	); err != nil {
		return nil, err
	}
	return allServices, nil
}

// NewServices services ctor
func NewServices(s *Storage) Services {
	return &services{storage: s}
}
