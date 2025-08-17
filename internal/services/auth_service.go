package services

import (
	"errors"
	"fmt"
	"gin-simple-app/internal/models"
	"gin-simple-app/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(req models.RegisterRequest) (*models.LoginResponse, error)
	Login(req models.LoginRequest) (*models.LoginResponse, error)
	ForgotPassword(req models.ForgotPasswordRequest) error
	ResetPassword(req models.ResetPasswordRequest) error
	ChangePassword(userID uint, req models.ChangePasswordRequest) error
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
	RefreshToken(userID uint) (*models.LoginResponse, error)
	Logout(userID uint) error
}

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// JWT Claims
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &AuthServiceImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register creates a new user account
func (s *AuthServiceImpl) Register(req models.RegisterRequest) (*models.LoginResponse, error) {
	// Validate password confirmation
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate email verification token
	emailVerificationToken := uuid.New().String()

	// Create user
	user := &models.User{
		Name:                   req.Name,
		Email:                  req.Email,
		Password:               string(hashedPassword),
		Phone:                  &req.Phone,
		Address:                req.Address,
		IsActive:               true,
		IsEmailVerified:        false,
		EmailVerificationToken: &emailVerificationToken,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, expiresIn, err := s.generateJWTToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	return &models.LoginResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

// Login authenticates a user
func (s *AuthServiceImpl) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, expiresIn, err := s.generateJWTToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	return &models.LoginResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

// ForgotPassword initiates password reset process
func (s *AuthServiceImpl) ForgotPassword(req models.ForgotPasswordRequest) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists or not
			return nil
		}
		return err
	}

	// Generate password reset token
	resetToken := uuid.New().String()
	expiry := time.Now().Add(time.Hour) // Token expires in 1 hour

	// Update user with reset token
	user.PasswordResetToken = &resetToken
	user.PasswordResetExpiry = &expiry

	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	// TODO: Send email with reset token
	// For now, we'll just log it (in production, you'd send an email)
	fmt.Printf("Password reset token for %s: %s\n", user.Email, resetToken)

	return nil
}

// ResetPassword resets user password using reset token
func (s *AuthServiceImpl) ResetPassword(req models.ResetPasswordRequest) error {
	// Validate password confirmation
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	// Find user by reset token
	user, err := s.userRepo.GetByPasswordResetToken(req.Token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Check if token is expired
	if user.PasswordResetExpiry == nil || time.Now().After(*user.PasswordResetExpiry) {
		return errors.New("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update user password and clear reset token
	user.Password = string(hashedPassword)
	user.PasswordResetToken = nil
	user.PasswordResetExpiry = nil

	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

// ChangePassword changes user password
func (s *AuthServiceImpl) ChangePassword(userID uint, req models.ChangePasswordRequest) error {
	// Validate password confirmation
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password
	user.Password = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

// ValidateToken validates JWT token and returns claims
func (s *AuthServiceImpl) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token for the user
func (s *AuthServiceImpl) RefreshToken(userID uint) (*models.LoginResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate new JWT token
	token, expiresIn, err := s.generateJWTToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.LoginResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

// Logout handles user logout (token invalidation would be handled client-side or with a token blacklist)
func (s *AuthServiceImpl) Logout(userID uint) error {
	// In a more sophisticated implementation, you might want to blacklist the token
	// For now, we'll just confirm the user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return nil
}

// generateJWTToken generates a new JWT token for the user
func (s *AuthServiceImpl) generateJWTToken(userID uint, email string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiry)

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(s.jwtExpiry.Seconds()), nil
}
