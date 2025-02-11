package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoClient interface {
	DescribeUserPool(ctx context.Context, params *cognitoidentityprovider.DescribeUserPoolInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DescribeUserPoolOutput, error)
	ListUsers(ctx context.Context, params *cognitoidentityprovider.ListUsersInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error)
	ListGroups(ctx context.Context, params *cognitoidentityprovider.ListGroupsInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListGroupsOutput, error)
	ListResourceServers(ctx context.Context, params *cognitoidentityprovider.ListResourceServersInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListResourceServersOutput, error)
	ListUserPoolClients(ctx context.Context, params *cognitoidentityprovider.ListUserPoolClientsInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUserPoolClientsOutput, error)
	ListIdentityProviders(ctx context.Context, params *cognitoidentityprovider.ListIdentityProvidersInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListIdentityProvidersOutput, error)
	CreateUserPool(ctx context.Context, params *cognitoidentityprovider.CreateUserPoolInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateUserPoolOutput, error)
	UpdateUserPool(ctx context.Context, params *cognitoidentityprovider.UpdateUserPoolInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.UpdateUserPoolOutput, error)
	CreateGroup(ctx context.Context, params *cognitoidentityprovider.CreateGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error)
	AdminCreateUser(ctx context.Context, params *cognitoidentityprovider.AdminCreateUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	CreateResourceServer(ctx context.Context, params *cognitoidentityprovider.CreateResourceServerInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateResourceServerOutput, error)
	CreateUserPoolClient(ctx context.Context, params *cognitoidentityprovider.CreateUserPoolClientInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateUserPoolClientOutput, error)
	CreateIdentityProvider(ctx context.Context, params *cognitoidentityprovider.CreateIdentityProviderInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateIdentityProviderOutput, error)
}
