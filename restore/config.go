package restore

// Config holds the configuration for restore operations
type Config struct {
	Mode         string // backup or restore
	SourcePoolID string
	TargetPoolID string // optional for restore
	Region       string
	BackupPath   string
	UsersOnly    bool
}
