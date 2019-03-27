package repositories

// ListOptions is used by repo.ListOptions functions
type ListOptions struct {
	PageSize int
	Page     int
}

// Limit as used in queries
func (lo ListOptions) Limit() int {
	return lo.PageSize
}

// Offset as used in queries
func (lo ListOptions) Offset() int {
	return lo.PageSize * lo.Page
}
