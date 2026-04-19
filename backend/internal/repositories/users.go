package repositories

import (
	"backend/internal/db"
	sqlcdb "backend/internal/db/sqlc"
	"backend/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type UserRepository struct {
	store  *db.Store
	logger *zap.Logger
}

func NewUserRepository(store *db.Store, logger *zap.Logger) *UserRepository {
	return &UserRepository{store: store, logger: logger}
}

// GetUserByEmail fetches a user by their email, checking the linked accounts for provider and provider_id.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Fetch the user from the users table
	user, err := r.store.GetUserByEmail(ctx, sqlcdb.GetUserByEmailParams{Email: email})
	if err != nil {
		r.logger.Debug("Failed to fetch user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	// Fetch the linked account(s) for the user (provider and provider_id will be in the linked accounts table)
	linkedAccounts, err := r.store.GetLinkedAccountsByEmail(ctx, sqlcdb.GetLinkedAccountsByEmailParams{Email: email})
	if err != nil {
		r.logger.Debug("Failed to fetch linked accounts", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	// Construct the user object by combining data from both tables
	var accounts []models.LinkedAccount
	for _, linkedAccount := range linkedAccounts {
		accounts = append(accounts, models.LinkedAccount{
			ID:         int(linkedAccount.ID),
			UserID:     int(linkedAccount.UserID),
			Provider:   linkedAccount.Provider,
			ProviderID: linkedAccount.ProviderID,
			CreatedAt:  linkedAccount.CreatedAt.Time,
		})
	}

	return &models.User{
		ID:             int(user.ID),
		Email:          user.Email,
		Name:           user.Name,
		Avatar:         user.Avatar.String,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt.Time,
		LastLogin:      user.LastLogin.Time,
		LinkedAccounts: accounts, // Link the accounts to the user
	}, nil
}

// CreateOrUpdateUser either creates a new user or updates an existing user's last login timestamp.
func (r *UserRepository) CreateOrUpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	// Check if the user exists by email
	existingUser, err := r.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		// Update last login timestamp if user exists
		err = r.store.UpdateUserLogin(ctx, sqlcdb.UpdateUserLoginParams{
			LastLogin: pgtype.Timestamp{Time: time.Now(), Valid: true},
			Email:     user.Email,
		})
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	// Create new user in the users table
	createdUser, err := r.store.CreateUser(ctx, sqlcdb.CreateUserParams{
		Email:  user.Email,
		Name:   user.Name,
		Avatar: pgtype.Text{String: user.Avatar, Valid: user.Avatar != ""},
		Role:   user.Role,
	})
	if err != nil {
		return nil, err
	}

	// Link the user to the provider in the linked_accounts table
	for _, linkedAccount := range user.LinkedAccounts {
		_, err = r.store.CreateLinkedAccount(ctx, sqlcdb.CreateLinkedAccountParams{
			UserID:     createdUser.ID,
			Provider:   linkedAccount.Provider,
			ProviderID: linkedAccount.ProviderID,
		})
		if err != nil {
			return nil, err
		}
	}

	// Return the newly created user
	return &models.User{
		ID:             int(createdUser.ID),
		Email:          createdUser.Email,
		Name:           createdUser.Name,
		Avatar:         createdUser.Avatar.String,
		Role:           createdUser.Role,
		CreatedAt:      createdUser.CreatedAt.Time,
		LastLogin:      createdUser.LastLogin.Time,
		LinkedAccounts: user.LinkedAccounts, // Include linked accounts
	}, nil
}

// UpdateUserLogin updates the last login timestamp for the given email.
func (r *UserRepository) UpdateUserLogin(ctx context.Context, email string) error {
	// Update the last login timestamp for the user in the users table
	err := r.store.UpdateUserLogin(ctx, sqlcdb.UpdateUserLoginParams{
		LastLogin: pgtype.Timestamp{Time: time.Now(), Valid: true},
		Email:     email,
	})
	if err != nil {
		r.logger.Debug("Failed to update last login", zap.String("email", email), zap.Error(err))
		return err
	}

	return nil
}

func (r *UserRepository) GetUserGroups(ctx context.Context, userID int) ([]models.Group, error) {
	groups, err := r.store.GetUserGroups(ctx, sqlcdb.GetUserGroupsParams{UserID: int32(userID)})
	if err != nil {
		r.logger.Debug("Failed to fetch user groups", zap.Int("user_id", userID), zap.Error(err))
		return nil, err
	}

	var result []models.Group
	for _, g := range groups {
		result = append(result, models.Group{
			ID:   int(g.ID),
			Name: g.Name,
		})
	}
	return result, nil
}

func (r *UserRepository) GetGroupPermissions(ctx context.Context, groupID int) ([]models.Permission, error) {
	permissions, err := r.store.GetGroupPermissions(ctx, sqlcdb.GetGroupPermissionsParams{GroupID: int32(groupID)})
	if err != nil {
		r.logger.Debug("Failed to fetch group permissions", zap.Int("group_id", groupID), zap.Error(err))
		return nil, err
	}

	var result []models.Permission
	for _, p := range permissions {
		result = append(result, models.Permission{
			ID:       int(p.ID),
			Codename: p.Codename,
			Name:     p.Name,
		})
	}
	return result, nil
}

func (r *UserRepository) GetUserPermissions(ctx context.Context, userID int) ([]models.Permission, error) {
	permissions, err := r.store.GetUserPermissions(ctx, sqlcdb.GetUserPermissionsParams{UserID: int32(userID)})
	if err != nil {
		r.logger.Debug("Failed to fetch user permissions", zap.Int("user_id", userID), zap.Error(err))
		return nil, err
	}

	var result []models.Permission
	for _, p := range permissions {
		result = append(result, models.Permission{
			ID:       int(p.ID),
			Codename: p.Codename,
			Name:     p.Name,
		})
	}
	return result, nil
}

func (r *UserRepository) AddUserToGroup(ctx context.Context, userID int, groupID int) error {
	err := r.store.AddUserToGroup(ctx, sqlcdb.AddUserToGroupParams{
		UserID:  int32(userID),
		GroupID: int32(groupID),
	})
	if err != nil {
		r.logger.Debug("Failed to add user to group", zap.Int("user_id", userID), zap.Int("group_id", groupID), zap.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) CreateGroup(ctx context.Context, name string) (*models.Group, error) {
	group, err := r.store.CreateGroup(ctx, sqlcdb.CreateGroupParams{Name: name})
	if err != nil {
		r.logger.Debug("Failed to create group", zap.String("name", name), zap.Error(err))
		return nil, err
	}
	return &models.Group{
		ID:   int(group.ID),
		Name: group.Name,
	}, nil
}

func (r *UserRepository) CreatePermission(ctx context.Context, codename, name string) (*models.Permission, error) {
	permission, err := r.store.CreatePermission(ctx, sqlcdb.CreatePermissionParams{
		Codename: codename,
		Name:     name,
	})
	if err != nil {
		r.logger.Debug("Failed to create permission", zap.String("codename", codename), zap.String("name", name), zap.Error(err))
		return nil, err
	}
	return &models.Permission{
		ID:       int(permission.ID),
		Codename: permission.Codename,
		Name:     permission.Name,
	}, nil
}

func (r *UserRepository) AssignPermissionToGroup(ctx context.Context, groupID int, permissionID int) error {
	err := r.store.AssignPermissionToGroup(ctx, sqlcdb.AssignPermissionToGroupParams{
		GroupID:      int32(groupID),
		PermissionID: int32(permissionID),
	})
	if err != nil {
		r.logger.Debug("Failed to assign permission to group", zap.Int("group_id", groupID), zap.Int("permission_id", permissionID), zap.Error(err))
		return err
	}
	return nil
}
