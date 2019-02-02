package repositories

import "github.com/ghostec/Will.IAM/models"

// Services repository
type Services interface {
	List() ([]models.Service, error)
	Create(*models.Service) error
	Clone() Services
	setStorage(*Storage)
}

type services struct {
	*withStorage
}

func (ss *services) Clone() Services {
	return NewServices(ss.storage.Clone())
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

// List returns all services in storage
func (ss services) List() ([]models.Service, error) {
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
	return &services{&withStorage{storage: s}}
}
