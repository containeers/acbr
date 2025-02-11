package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"acbr/aws"
	"acbr/backup"
	"acbr/config"
	"acbr/restore"
)

// Add this type at the top of main.go
type LambdaEvent struct {
	Mode       string `json:"mode"`
	PoolID     string `json:"poolId"`
	Region     string `json:"region"`
	BackupPath string `json:"backupPath"`
	UsersOnly  bool   `json:"usersOnly,omitempty"`
	MaxResults int32  `json:"maxResults,omitempty"`
}

var Version = "dev" // This will be set during build

/*
# Local storage
./acbr -mode backup -source us-east-1_xxxxx -region us-east-1 -backup-path ./backups/
./acbr -mode backup -source us-east-1_xxxxx -region us-east-1 -backup-path /absolute/path/backups/

# S3 storage
./acbr -mode backup -source us-east-1_xxxxx -region us-east-1 -backup-path s3://my-bucket/cognito/backups/
*/
func main() {
	cfg := &config.Config{}

	// Check if running as Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(handleLambda)
		return
	}

	// CLI flags
	flag.StringVar(&cfg.Mode, "mode", "", "Operation mode: backup or restore")
	flag.StringVar(&cfg.PoolID, "pool", "", "Pool ID (source for backup, target for restore)")
	flag.StringVar(&cfg.Region, "region", "", "AWS Region")
	flag.StringVar(&cfg.BackupPath, "backup-path", "", "Path to store/read backup files")
	flag.BoolVar(&cfg.UsersOnly, "users-only", false, "Restore only users and groups")
	var maxResults int
	flag.IntVar(&maxResults, "max-results", 50, "Maximum results per page for AWS API calls (max 50)")
	flag.StringVar(&cfg.DefaultPwd, "default-pwd", "", "Default password for Cognito-created users (required for non-SSO users)")
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Parse()

	if *showVersion {
		fmt.Printf("acbr version %s\n", Version)
		os.Exit(0)
	}

	if cfg.Mode == "" || cfg.PoolID == "" || cfg.Region == "" || cfg.BackupPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg.MaxResults = int32(maxResults)

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

// Add this function to main.go
func handleLambda(event LambdaEvent) error {
	cfg := &config.Config{
		Mode:       event.Mode,
		PoolID:     event.PoolID,
		Region:     event.Region,
		BackupPath: event.BackupPath,
		UsersOnly:  event.UsersOnly,
		MaxResults: event.MaxResults,
	}

	if cfg.MaxResults == 0 || cfg.MaxResults > 50 {
		cfg.MaxResults = 50
	}

	return run(cfg)
}

func run(config *config.Config) error {
	client, err := aws.NewCognitoClient(config.Region)
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	switch config.Mode {
	case "backup":
		b := backup.NewBackup(client, config)
		return b.Execute()
	case "restore":
		r := restore.NewRestore(client, config)
		return r.Execute()
	default:
		return fmt.Errorf("invalid mode: %s", config.Mode)
	}
}
