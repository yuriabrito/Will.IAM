package repositories

// ServiceAccounts repository
type ServiceAccounts interface{}

type serviceAccounts struct{}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts() ServiceAccounts {
	return &serviceAccounts{}
}
