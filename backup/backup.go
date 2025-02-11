package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"acbr/aws"
	"acbr/config"
	"acbr/storage"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoBackup struct {
	UserPoolConfig    *cognitoidentityprovider.DescribeUserPoolOutput
	Users             []types.UserType
	Groups            []types.GroupType
	ResourceServers   []types.ResourceServerType
	Clients           []types.UserPoolClientDescription
	IdentityProviders []types.ProviderDescription
}

type Backup struct {
	client aws.CognitoClient
	config *config.Config
}

func NewBackup(client aws.CognitoClient, config *config.Config) *Backup {
	return &Backup{
		client: client,
		config: config,
	}
}

func (b *Backup) Execute() error {
	backup := &CognitoBackup{}

	// Get User Pool configuration
	userPool, err := b.getUserPool()
	if err != nil {
		return fmt.Errorf("failed to get user pool: %w", err)
	}
	backup.UserPoolConfig = userPool

	// Get Users
	users, err := b.getUsers()
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	backup.Users = users

	// Get Groups
	groups, err := b.getGroups()
	if err != nil {
		return fmt.Errorf("failed to get groups: %w", err)
	}
	backup.Groups = groups

	// Get Resource Servers
	servers, err := b.getResourceServers()
	if err != nil {
		return fmt.Errorf("failed to get resource servers: %w", err)
	}
	backup.ResourceServers = servers

	// Get Clients
	clients, err := b.getClients()
	if err != nil {
		return fmt.Errorf("failed to get clients: %w", err)
	}
	backup.Clients = clients

	// Get Identity Providers
	providers, err := b.getIdentityProviders()
	if err != nil {
		return fmt.Errorf("failed to get identity providers: %w", err)
	}
	backup.IdentityProviders = providers

	// Save backup to file
	return b.saveBackup(backup)
}

func (b *Backup) saveBackup(backup *CognitoBackup) error {
	// Create storage based on backup path
	storage, err := storage.NewStorage(b.config.BackupPath)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// Marshal backup data
	data, err := json.Marshal(backup)
	if err != nil {
		return fmt.Errorf("failed to marshal backup: %w", err)
	}

	// Generate backup filename
	filename := fmt.Sprintf("cognito-backup-%s-%s.json",
		b.config.PoolID,
		time.Now().Format("20060102-150405"))

	// For local storage, join path with filename
	// For S3, the path handling is already correct in S3Storage
	path := filename
	if !strings.HasPrefix(b.config.BackupPath, "s3://") {
		path = filepath.Join(b.config.BackupPath, filename)
	}

	// Save backup
	if err := storage.Save(context.Background(), data, path); err != nil {
		return fmt.Errorf("failed to save backup: %w", err)
	}

	return nil
}

func (b *Backup) getUserPool() (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	return b.client.DescribeUserPool(context.Background(), &cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: &b.config.PoolID,
	})
}

func (b *Backup) getUsers() ([]types.UserType, error) {
	var users []types.UserType
	paginator := cognitoidentityprovider.NewListUsersPaginator(b.client, &cognitoidentityprovider.ListUsersInput{
		UserPoolId: &b.config.PoolID,
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		users = append(users, output.Users...)
	}
	return users, nil
}

func (b *Backup) getGroups() ([]types.GroupType, error) {
	var groups []types.GroupType
	paginator := cognitoidentityprovider.NewListGroupsPaginator(b.client, &cognitoidentityprovider.ListGroupsInput{
		UserPoolId: &b.config.PoolID,
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		groups = append(groups, output.Groups...)
	}
	return groups, nil
}

func (b *Backup) getResourceServers() ([]types.ResourceServerType, error) {
	output, err := b.client.ListResourceServers(context.Background(), &cognitoidentityprovider.ListResourceServersInput{
		UserPoolId: &b.config.PoolID,
		MaxResults: awssdk.Int32(b.config.GetMaxResults()),
	})
	if err != nil {
		return nil, err
	}
	return output.ResourceServers, nil
}

func (b *Backup) getClients() ([]types.UserPoolClientDescription, error) {
	var clients []types.UserPoolClientDescription
	paginator := cognitoidentityprovider.NewListUserPoolClientsPaginator(b.client, &cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId: &b.config.PoolID,
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		clients = append(clients, output.UserPoolClients...)
	}
	return clients, nil
}

func (b *Backup) getIdentityProviders() ([]types.ProviderDescription, error) {
	output, err := b.client.ListIdentityProviders(context.Background(), &cognitoidentityprovider.ListIdentityProvidersInput{
		UserPoolId: &b.config.PoolID,
	})
	if err != nil {
		return nil, err
	}
	return output.Providers, nil
}
