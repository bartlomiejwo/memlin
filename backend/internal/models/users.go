package models

import "time"

type User struct {
	ID             int             `json:"id"`
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	Avatar         string          `json:"avatar"`
	Role           string          `json:"role"`
	CreatedAt      time.Time       `json:"created_at"`
	LastLogin      time.Time       `json:"last_login"`
	LinkedAccounts []LinkedAccount `json:"linked_accounts"`
	Groups         []Group         `json:"groups"`
}

type LinkedAccount struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Group struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	ID       int    `json:"id"`
	Codename string `json:"codename"`
	Name     string `json:"name"`
}
