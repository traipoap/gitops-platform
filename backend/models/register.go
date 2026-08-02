package models

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"` // Stores email or username
	PasswordHash string `gorm:"not null"`             // Stores bcrypt hash
	Role         string `gorm:"default:'user'"`       // e.g., "admin", "user"
}

// CreateUserInput defines the expected JSON payload
type CreateUserInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"` // Optional, default to "user"
}
