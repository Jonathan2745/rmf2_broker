package config

type Config struct {
	Addr        string
	DataDir     string
	DefaultOrg  string
	VDA5050Base string // e.g. http://localhost:8000 — empty disables live positions
}
