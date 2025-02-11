package restore

import (
	"context"
	"testing"

	"acbr/config"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type mockCognitoClient struct {
	describeUserPoolOutput *cognitoidentityprovider.DescribeUserPoolOutput
	describeUserPoolError  error
}

func (m *mockCognitoClient) DescribeUserPool(ctx context.Context, params *cognitoidentityprovider.DescribeUserPoolInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	return m.describeUserPoolOutput, m.describeUserPoolError
}

func (m *mockCognitoClient) ListUsers(ctx context.Context, params *cognitoidentityprovider.ListUsersInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
	return &cognitoidentityprovider.ListUsersOutput{}, nil
}

func (m *mockCognitoClient) ListGroups(ctx context.Context, params *cognitoidentityprovider.ListGroupsInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListGroupsOutput, error) {
	return &cognitoidentityprovider.ListGroupsOutput{}, nil
}

func (m *mockCognitoClient) ListResourceServers(ctx context.Context, params *cognitoidentityprovider.ListResourceServersInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListResourceServersOutput, error) {
	return &cognitoidentityprovider.ListResourceServersOutput{}, nil
}

func (m *mockCognitoClient) ListUserPoolClients(ctx context.Context, params *cognitoidentityprovider.ListUserPoolClientsInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUserPoolClientsOutput, error) {
	return &cognitoidentityprovider.ListUserPoolClientsOutput{}, nil
}

func (m *mockCognitoClient) ListIdentityProviders(ctx context.Context, params *cognitoidentityprovider.ListIdentityProvidersInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListIdentityProvidersOutput, error) {
	return &cognitoidentityprovider.ListIdentityProvidersOutput{}, nil
}

func (m *mockCognitoClient) CreateUserPool(ctx context.Context, params *cognitoidentityprovider.CreateUserPoolInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateUserPoolOutput, error) {
	return &cognitoidentityprovider.CreateUserPoolOutput{}, nil
}

func (m *mockCognitoClient) UpdateUserPool(ctx context.Context, params *cognitoidentityprovider.UpdateUserPoolInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.UpdateUserPoolOutput, error) {
	return &cognitoidentityprovider.UpdateUserPoolOutput{}, nil
}

func (m *mockCognitoClient) CreateGroup(ctx context.Context, params *cognitoidentityprovider.CreateGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
	return &cognitoidentityprovider.CreateGroupOutput{}, nil
}

func (m *mockCognitoClient) AdminCreateUser(ctx context.Context, params *cognitoidentityprovider.AdminCreateUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	return &cognitoidentityprovider.AdminCreateUserOutput{}, nil
}

func (m *mockCognitoClient) CreateResourceServer(ctx context.Context, params *cognitoidentityprovider.CreateResourceServerInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateResourceServerOutput, error) {
	return &cognitoidentityprovider.CreateResourceServerOutput{}, nil
}

func (m *mockCognitoClient) CreateUserPoolClient(ctx context.Context, params *cognitoidentityprovider.CreateUserPoolClientInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateUserPoolClientOutput, error) {
	return &cognitoidentityprovider.CreateUserPoolClientOutput{}, nil
}

func (m *mockCognitoClient) CreateIdentityProvider(ctx context.Context, params *cognitoidentityprovider.CreateIdentityProviderInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateIdentityProviderOutput, error) {
	return &cognitoidentityprovider.CreateIdentityProviderOutput{}, nil
}

func TestNewRestore(t *testing.T) {
	client := &mockCognitoClient{}
	cfg := &config.Config{
		Mode:   "restore",
		PoolID: "test-pool",
		Region: "us-east-1",
	}

	restore := NewRestore(client, cfg)
	if restore == nil {
		t.Error("NewRestore() returned nil")
	}
}
