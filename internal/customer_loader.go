package internal

// LoaderCustomer is the interface that wraps the basic methods that a customer loader should implement.
type LoaderCustomer interface {
	// ReadAll returns all customers saved in the database.
	ReadAll() (err error)
}
