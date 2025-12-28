package auth

import (
	"context"
	"regexp"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/errorx"
	"github.com/Eagle233Fake/omniread/backend/api/token"
	"github.com/Eagle233Fake/omniread/backend/application/assembler"
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/infra/cache"
	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/infra/model"
	"github.com/Eagle233Fake/omniread/backend/infra/repo"
	"github.com/Eagle233Fake/omniread/backend/types/errno"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var AuthServiceSet = wire.NewSet(NewAuthService)

type AuthService struct {
	userRepo  *repo.UserRepo
	authCache *cache.AuthCache
}

func NewAuthService(userRepo *repo.UserRepo, authCache *cache.AuthCache) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		authCache: authCache,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterReq) error {
	// Validate required fields
	if req.Username == "" {
		return errorx.New(errno.ParamErrorCode)
	}
	if req.Email == "" && req.Phone == "" {
		return errorx.New(errno.ParamErrorCode)
	}
	if len(req.Password) < 8 {
		return errorx.New(errno.ParamErrorCode)
	}
	// Password must contain at least one letter and one number
	if matched, _ := regexp.MatchString(`[A-Za-z]`, req.Password); !matched {
		return errorx.New(errno.ParamErrorCode)
	}
	if matched, _ := regexp.MatchString(`[0-9]`, req.Password); !matched {
		return errorx.New(errno.ParamErrorCode)
	}

	// Validate email format
	if req.Email != "" {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, req.Email)
		if !matched {
			return errorx.New(errno.ParamErrorCode)
		}
	}

	// Validate phone format (11 digits)
	if req.Phone != "" {
		matched, _ := regexp.MatchString(`^\d{11}$`, req.Phone)
		if !matched {
			return errorx.New(errno.ParamErrorCode)
		}
	}

	// Check uniqueness
	existingUser, _ := s.userRepo.FindByUsername(ctx, req.Username)
	if existingUser != nil {
		return errorx.New(errno.ErrUserExist)
	}
	if req.Email != "" {
		existingUser, _ = s.userRepo.FindByEmail(ctx, req.Email)
		if existingUser != nil {
			return errorx.New(errno.ErrUserExist)
		}
	}
	if req.Phone != "" {
		existingUser, _ = s.userRepo.FindByPhone(ctx, req.Phone)
		if existingUser != nil {
			return errorx.New(errno.ErrUserExist)
		}
	}

	// Validate gender
	if req.Gender != "" && req.Gender != "male" && req.Gender != "female" && req.Gender != "other" {
		return errorx.New(errno.ParamErrorCode)
	}

	// Validate birthdate
	if req.Birthdate != "" {
		_, err := time.Parse("2006-01-02", req.Birthdate)
		if err != nil {
			return errorx.New(errno.ParamErrorCode)
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errorx.New(errno.ErrUserCreateFailed)
	}

	user := assembler.RegisterReqToUser(req, string(hashedPassword))

	if err := s.userRepo.Create(ctx, user); err != nil {
		return errorx.New(errno.ErrUserCreateFailed)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginReq) (*dto.LoginResp, error) {
	var user *model.User
	var err error

	// Try to find by username, email, or phone
	user, err = s.userRepo.FindByUsername(ctx, req.Username)
	if user == nil {
		user, err = s.userRepo.FindByEmail(ctx, req.Username)
		if user == nil {
			user, err = s.userRepo.FindByPhone(ctx, req.Username)
		}
	}

	if user == nil {
		return nil, errorx.New(errno.ErrUserNotFound)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errorx.New(errno.ErrAuthFailed)
	}

	// Generate token
	t, err := token.GenerateToken(user)
	if err != nil {
		return nil, errorx.New(errno.ErrAuthFailed)
	}

	// Store token in Redis
	cfg := config.GetConfig()
	err = s.authCache.SetSession(ctx, t, user.ID.Hex(), time.Duration(cfg.Auth.AccessExpire)*time.Second)
	if err != nil {
		// Log error but don't fail login? Or fail? Best to fail if cache is critical for security
		return nil, errorx.New(errno.ErrAuthFailed)
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return assembler.UserToLoginResp(user, t), nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileReq) error {
	// Parse user ID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errorx.New(errno.ParamErrorCode)
	}

	user, err := s.userRepo.FindByID(ctx, uid)
	if err != nil || user == nil {
		return errorx.New(errno.ErrUserNotFound)
	}

	// Update fields if provided
	if req.Nickname != "" {
		user.Username = req.Nickname // Assuming username acts as nickname, or add Nickname field to model
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Birthdate != "" {
		if t, err := time.Parse("2006-01-02", req.Birthdate); err == nil {
			user.Birthdate = &t
		}
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errorx.New(errno.ErrUserUpdateFailed)
	}
	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordReq) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errorx.New(errno.ParamErrorCode)
	}

	user, err := s.userRepo.FindByID(ctx, uid)
	if err != nil || user == nil {
		return errorx.New(errno.ErrUserNotFound)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errorx.New(errno.ErrAuthFailed)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errorx.New(errno.ErrUserUpdateFailed)
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errorx.New(errno.ErrUserUpdateFailed)
	}
	return nil
}

func (s *AuthService) UpdatePreferences(ctx context.Context, userID string, req *dto.UpdatePreferencesReq) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errorx.New(errno.ParamErrorCode)
	}

	user, err := s.userRepo.FindByID(ctx, uid)
	if err != nil || user == nil {
		return errorx.New(errno.ErrUserNotFound)
	}

	user.Preferences = model.Preferences{
		FontFamily: req.FontFamily,
		FontSize:   req.FontSize,
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errorx.New(errno.ErrUserUpdateFailed)
	}
	return nil
}

func (s *AuthService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errorx.New(errno.ParamErrorCode)
	}
	return s.userRepo.FindByID(ctx, uid)
}
