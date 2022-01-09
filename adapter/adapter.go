package adapter

const (
	MYSQL = "mysql"
	// POSTGRES = "postgres"
)

type Adapter interface {
	Ping() error
}
