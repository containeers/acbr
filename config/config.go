package config

// Config holds the configuration for backup/restore operations
type Config struct {
	Mode       string
	PoolID     string // Single pool ID field
	Region     string
	BackupPath string
	UsersOnly  bool
	MaxResults int32
	DefaultPwd string // Add default password field
}

// GetMaxResults returns the configured MaxResults or a default value
func (c *Config) GetMaxResults() int32 {
	if c.MaxResults <= 0 || c.MaxResults > 50 {
		return 50 // default value, respecting AWS limit
	}
	return c.MaxResults
}
