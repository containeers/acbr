# AWS Cognito Backup & Restore (ACBR)

A tool to backup and restore AWS Cognito User Pools, supporting both CLI and AWS Lambda deployments.

## Features

- Backup Cognito User Pools (users, groups, settings)
- Restore to new or existing pools
- Support for SSO and native Cognito users
- Local file system and S3 storage support
- CLI and AWS Lambda deployment options

## Installation


*Install from source*
```bash
go install github.com/yourusername/acbr@latest
```

*Or build locally*
```bash
git clone https://github.com/yourusername/acbr.git
cd acbr
go build
```

## CLI Usage

### Backup

```bash
# Backup to local file system
./acbr -mode backup \
       -pool us-east-1_xxxxx \
       -region us-east-1 \
       -backup-path ./backups/

# Backup to S3
./acbr -mode backup \
       -pool us-east-1_xxxxx \
       -region us-east-1 \
       -backup-path s3://my-bucket/cognito/backups/
```

### Restore

```bash
# Restore from local backup
./acbr -mode restore \
       -pool us-east-1_yyyyy \
       -region us-east-1 \
       -backup-path ./backups/cognito-backup-xxxxx.json \
       -default-pwd "TempPass123!"

# Restore from S3
./acbr -mode restore \
       -pool us-east-1_yyyyy \
       -region us-east-1 \
       -backup-path s3://my-bucket/cognito/backups/cognito-backup-xxxxx.json \
       -default-pwd "TempPass123!"

# Restore only users and groups
./acbr -mode restore \
       -pool us-east-1_yyyyy \
       -region us-east-1 \
       -backup-path ./backups/cognito-backup-xxxxx.json \
       -default-pwd "TempPass123!" \
       -users-only
```

## AWS Lambda Usage

Deploy the Lambda function and invoke with this event structure:

```json
{
  "mode": "backup",
  "poolId": "us-east-1_xxxxx",
  "region": "us-east-1",
  "backupPath": "s3://my-bucket/cognito/backups/",
  "usersOnly": false,
  "defaultPwd": "TempPass123!"
}
```

### Required IAM Permissions

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "cognito-idp:DescribeUserPool",
                "cognito-idp:ListUsers",
                "cognito-idp:ListGroups",
                "cognito-idp:ListResourceServers",
                "cognito-idp:ListUserPoolClients",
                "cognito-idp:ListIdentityProviders",
                "cognito-idp:CreateUserPool",
                "cognito-idp:UpdateUserPool",
                "cognito-idp:CreateGroup",
                "cognito-idp:AdminCreateUser"
            ],
            "Resource": "arn:aws:cognito-idp:*:*:userpool/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject"
            ],
            "Resource": "arn:aws:s3:::my-bucket/*"
        }
    ]
}
```

## Configuration Options

| Flag | Description | Required |
|------|-------------|----------|
| mode | Operation mode: backup or restore | Yes |
| pool | Pool ID (source for backup, target for restore) | Yes |
| region | AWS Region | Yes |
| backup-path | Path to store/read backup files | Yes |
| users-only | Restore only users and groups | No |
| default-pwd | Default password for Cognito-created users | Yes (for restore) |
| max-results | Maximum results per page for AWS API calls (max 50) | No |

## Notes

- SSO users are restored without passwords
- Non-SSO users require a default password during restore
- When restoring to an existing pool, only specified components are updated
- S3 paths must use the format: s3://bucket-name/path/
- Local paths can be relative or absolute

## Development

```bash
# Run tests
go test -v ./...

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o acbr-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o acbr-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o acbr-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o acbr-windows-amd64.exe
```

## License

MIT License - see LICENSE file for details

