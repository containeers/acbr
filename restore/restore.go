package restore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"acbr/aws"
	"acbr/backup"
	"acbr/config"
	"acbr/storage"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type Restore struct {
	client aws.CognitoClient
	config *config.Config
}

func NewRestore(client aws.CognitoClient, config *config.Config) *Restore {
	return &Restore{
		client: client,
		config: config,
	}
}

func (r *Restore) Execute() error {
	// Load backup
	backup, err := r.loadBackup()
	if err != nil {
		return fmt.Errorf("failed to load backup: %w", err)
	}

	// If users-only mode, only restore users and groups
	if r.config.UsersOnly {
		return r.restoreUsersAndGroups(backup)
	}

	// Check if target pool exists
	_, err = r.client.DescribeUserPool(context.Background(), &cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: &r.config.PoolID,
	})
	if err != nil {
		// If pool doesn't exist, create new one with provided name
		backup.UserPoolConfig.UserPool.Name = awssdk.String(r.config.PoolID)
		poolID, err := r.createUserPool(backup.UserPoolConfig)
		if err != nil {
			return fmt.Errorf("failed to create user pool: %w", err)
		}
		r.config.PoolID = poolID
		fmt.Printf("Created new pool: %s\n", poolID)
	} else {
		fmt.Printf("Using existing pool: %s\n", r.config.PoolID)
	}

	// Restore full configuration
	if err := r.restoreUserPool(backup); err != nil {
		return fmt.Errorf("failed to restore user pool: %w", err)
	}

	// Always restore users and groups after pool configuration
	if err := r.restoreUsersAndGroups(backup); err != nil {
		return fmt.Errorf("failed to restore users and groups: %w", err)
	}

	fmt.Printf("Successfully restored to pool: %s\n", r.config.PoolID)
	return nil
}

func (r *Restore) restoreUsersAndGroups(backup *backup.CognitoBackup) error {

	fmt.Printf("Restoring groups: %v\n", backup.Groups)
	// Restore groups first
	for _, group := range backup.Groups {
		if err := r.createGroup(&group); err != nil {
			return fmt.Errorf("failed to create group %s: %w", *group.GroupName, err)
		}
	}

	// Restore users
	for _, user := range backup.Users {
		if err := r.createUser(&user); err != nil {
			return fmt.Errorf("failed to create user %s: %w", *user.Username, err)
		}
	}

	return nil
}

func (r *Restore) loadBackup() (*backup.CognitoBackup, error) {
	// Create storage based on backup path
	storage, err := storage.NewStorage(r.config.BackupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	// For S3, use just the filename
	path := r.config.BackupPath
	if strings.HasPrefix(path, "s3://") {
		parts := strings.Split(path, "/")
		// Get the last part (filename)
		path = parts[len(parts)-1]
	}

	// Load backup data
	data, err := storage.Load(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to load backup: %w", err)
	}

	// Unmarshal backup data
	var backup backup.CognitoBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup: %w", err)
	}

	return &backup, nil
}

func (r *Restore) createUserPool(config *cognitoidentityprovider.DescribeUserPoolOutput) (string, error) {
	input := &cognitoidentityprovider.CreateUserPoolInput{
		PoolName: config.UserPool.Name,
		// Copy relevant settings from config.UserPool
		Policies:               config.UserPool.Policies,
		AutoVerifiedAttributes: config.UserPool.AutoVerifiedAttributes,
		MfaConfiguration:       config.UserPool.MfaConfiguration,
		EmailConfiguration:     config.UserPool.EmailConfiguration,
		SmsConfiguration:       config.UserPool.SmsConfiguration,
		// Add other configurations as needed
	}

	result, err := r.client.CreateUserPool(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("failed to create user pool: %w", err)
	}

	return *result.UserPool.Id, nil
}

func (r *Restore) createGroup(group *types.GroupType) error {
	input := &cognitoidentityprovider.CreateGroupInput{
		GroupName:   group.GroupName,
		UserPoolId:  &r.config.PoolID,
		Description: group.Description,
		Precedence:  group.Precedence,
		RoleArn:     group.RoleArn,
	}

	_, err := r.client.CreateGroup(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	return nil
}

func (r *Restore) createUser(user *types.UserType) error {
	// Check if user is from SSO (has identities attribute)
	isSSO := false
	for _, attr := range user.Attributes {
		if *attr.Name == "identities" {
			isSSO = true
			break
		}
	}

	if isSSO {
		fmt.Printf("Skipping password for SSO user: %s\n", *user.Username)
	} else if r.config.DefaultPwd == "" {
		return fmt.Errorf("default-pwd is required for non-SSO user: %s", *user.Username)
	}

	// Filter out non-mutable attributes
	var filteredAttrs []types.AttributeType
	for _, attr := range user.Attributes {
		if *attr.Name != "sub" && *attr.Name != "identities" {
			filteredAttrs = append(filteredAttrs, attr)
		}
	}

	input := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:     &r.config.PoolID,
		Username:       user.Username,
		UserAttributes: filteredAttrs,
		MessageAction:  types.MessageActionTypeSuppress,
	}

	// Set temporary password only for non-SSO users
	if !isSSO {
		input.TemporaryPassword = awssdk.String(r.config.DefaultPwd)
	}

	_, err := r.client.AdminCreateUser(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *Restore) restoreUserPool(backup *backup.CognitoBackup) error {
	// Update user pool settings
	updateInput := &cognitoidentityprovider.UpdateUserPoolInput{
		UserPoolId: &r.config.PoolID,
		// Copy settings from backup.UserPoolConfig.UserPool
		Policies:               backup.UserPoolConfig.UserPool.Policies,
		AutoVerifiedAttributes: backup.UserPoolConfig.UserPool.AutoVerifiedAttributes,
		MfaConfiguration:       backup.UserPoolConfig.UserPool.MfaConfiguration,
		EmailConfiguration:     backup.UserPoolConfig.UserPool.EmailConfiguration,
		SmsConfiguration:       backup.UserPoolConfig.UserPool.SmsConfiguration,
		// Add other configurations as needed
	}

	_, err := r.client.UpdateUserPool(context.Background(), updateInput)
	if err != nil {
		return fmt.Errorf("failed to update user pool: %w", err)
	}

	// Restore resource servers
	for _, server := range backup.ResourceServers {
		_, err := r.client.CreateResourceServer(context.Background(), &cognitoidentityprovider.CreateResourceServerInput{
			UserPoolId: &r.config.PoolID,
			Identifier: server.Identifier,
			Name:       server.Name,
			Scopes:     server.Scopes,
		})
		if err != nil {
			return fmt.Errorf("failed to create resource server %s: %w", *server.Identifier, err)
		}
	}

	// Restore app clients
	for _, client := range backup.Clients {
		_, err := r.client.CreateUserPoolClient(context.Background(), &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: &r.config.PoolID,
			ClientName: client.ClientName,
		})
		if err != nil {
			return fmt.Errorf("failed to create client %s: %w", *client.ClientName, err)
		}
	}

	// Restore identity providers
	for _, provider := range backup.IdentityProviders {
		_, err := r.client.CreateIdentityProvider(context.Background(), &cognitoidentityprovider.CreateIdentityProviderInput{
			UserPoolId:   &r.config.PoolID,
			ProviderName: provider.ProviderName,
			ProviderType: provider.ProviderType,
			// Remove fields that don't exist in ProviderDescription
			// These fields need to be fetched separately if needed
		})
		if err != nil {
			return fmt.Errorf("failed to create identity provider %s: %w", *provider.ProviderName, err)
		}
	}

	return nil
}
