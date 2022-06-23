package token

type dataPersistence interface {
	// Generate creates a new 6-12 digit authentication token
	Generate() (string, error)

	// Validate checks if a string is registered
	Validate(s string) error
}

type BusinessToken struct {
	dataLayer dataPersistence
}

func (b *BusinessToken) Generate() (string, error) {
	return b.dataLayer.Generate()
}

func (b *BusinessToken) Validate(s string) error {
	return b.dataLayer.Validate(s)
}

func NewBusinessToken(mysqlDataPersistence dataPersistence) *BusinessToken {
	return &BusinessToken{
		dataLayer: mysqlDataPersistence,
	}
}
