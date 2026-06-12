package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
	"zhanxu-admin/backend/config"
	"zhanxu-admin/backend/internal/dto"
	"zhanxu-admin/backend/internal/model"
	"zhanxu-admin/backend/internal/repository"
	"zhanxu-admin/backend/pkg/cache"
	"zhanxu-admin/backend/pkg/crypto"
	"zhanxu-admin/backend/pkg/jwtutil"
	"zhanxu-admin/backend/pkg/response"
)

var captchaStore = cache.CaptchaStore{}

type AuthService struct {
	cfg      *config.Config
	userRepo *repository.UserRepo
	logRepo  *repository.LogRepo
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepo, logRepo *repository.LogRepo) *AuthService {
	return &AuthService{cfg: cfg, userRepo: userRepo, logRepo: logRepo}
}

func (s *AuthService) GetCaptcha() (*dto.CaptchaResp, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64, _, err := c.Generate()
	if err != nil {
		return nil, err
	}
	return &dto.CaptchaResp{CaptchaID: id, Image: b64}, nil
}

func (s *AuthService) Login(req *dto.LoginReq, ip, userAgent string) (*dto.LoginResp, error) {
	ctx := context.Background()
	captchaAnswer, err := cache.GetDelString(ctx, cache.CaptchaKey(req.CaptchaID))
	if err != nil {
		if cache.IsNotFound(err) {
			return nil, &BizError{Code: response.CodeCaptchaInvalid}
		}
		return nil, err
	}
	if strings.TrimSpace(captchaAnswer) != strings.TrimSpace(req.CaptchaCode) {
		return nil, &BizError{Code: response.CodeCaptchaInvalid}
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.saveLoginLog(0, req.Username, ip, userAgent, 0, "用户不存在")
			return nil, &BizError{Code: response.CodeUserNotFound}
		}
		return nil, err
	}

	if !crypto.CheckPassword(req.Password, user.Password) {
		s.saveLoginLog(user.ID, user.Username, ip, userAgent, 0, "密码错误")
		return nil, &BizError{Code: response.CodePasswordError}
	}

	if user.Status == 0 {
		s.saveLoginLog(user.ID, user.Username, ip, userAgent, 0, "账号已禁用")
		return nil, &BizError{Code: response.CodeUserDisabled}
	}

	sessionID := jwtutil.NewTokenID()
	accessToken, err := jwtutil.GenerateAccessToken(user.ID, user.Username, sessionID,
		s.cfg.Server.JwtSecret, s.cfg.Server.AccessTokenExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwtutil.GenerateRefreshToken(user.ID, user.Username, sessionID,
		s.cfg.Server.JwtSecret, s.cfg.Server.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}

	// 存储 refresh token 到 Redis
	ttl := time.Duration(s.cfg.Server.RefreshTokenExpire) * time.Second
	if err := cache.Set(ctx, cache.RefreshTokenKey(user.ID, sessionID), refreshToken, ttl); err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)
	s.saveLoginLog(user.ID, user.Username, ip, userAgent, 1, "")

	return &dto.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.Server.AccessTokenExpire,
	}, nil
}

func (s *AuthService) Logout(userID uint, token string) error {
	ctx := context.Background()
	ttl := time.Duration(s.cfg.Server.AccessTokenExpire) * time.Second
	claims, err := jwtutil.ParseToken(token, s.cfg.Server.JwtSecret)
	if err != nil || claims.Subject != "access" {
		return jwtutil.ErrTokenInvalid
	}
	if claims.ExpiresAt != nil {
		ttl = time.Until(claims.ExpiresAt.Time)
	}
	if ttl > 0 {
		if err := cache.Set(ctx, cache.BlacklistKey(claims.ID), 1, ttl); err != nil {
			return err
		}
	}
	return cache.Del(ctx, cache.RefreshTokenKey(userID, claims.SessionID))
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenReq) (*dto.LoginResp, error) {
	claims, err := jwtutil.ParseToken(req.RefreshToken, s.cfg.Server.JwtSecret)
	if err != nil {
		return nil, &BizError{Code: response.CodeRefreshTokenInvalid}
	}
	if claims.Subject != "refresh" {
		return nil, &BizError{Code: response.CodeRefreshTokenInvalid}
	}
	if claims.ID == "" || claims.SessionID == "" {
		return nil, &BizError{Code: response.CodeRefreshTokenInvalid}
	}

	// 验证 Redis 中是否存在
	ctx := context.Background()
	key := cache.RefreshTokenKey(claims.UserID, claims.SessionID)
	stored, err := cache.GetString(ctx, key)
	if err != nil {
		if cache.IsNotFound(err) {
			return nil, &BizError{Code: response.CodeRefreshTokenInvalid}
		}
		return nil, err
	}
	if stored != req.RefreshToken {
		return nil, &BizError{Code: response.CodeRefreshTokenInvalid}
	}

	accessToken, err := jwtutil.GenerateAccessToken(claims.UserID, claims.Username, claims.SessionID,
		s.cfg.Server.JwtSecret, s.cfg.Server.AccessTokenExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwtutil.GenerateRefreshToken(claims.UserID, claims.Username, claims.SessionID,
		s.cfg.Server.JwtSecret, s.cfg.Server.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}
	ttl := time.Duration(s.cfg.Server.RefreshTokenExpire) * time.Second
	rotated, err := cache.RotateValue(ctx, key, req.RefreshToken, refreshToken, ttl)
	if err != nil {
		return nil, err
	}
	if !rotated {
		return nil, &BizError{Code: response.CodeRefreshTokenInvalid}
	}

	return &dto.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.Server.AccessTokenExpire,
	}, nil
}

func (s *AuthService) saveLoginLog(userID uint, username, ip, userAgent string, status int8, msg string) {
	_ = s.logRepo.CreateLoginLog(&model.SysLoginLog{
		UserID:   userID,
		Username: username,
		IP:       ip,
		Browser:  parseBrowser(userAgent),
		OS:       parseOS(userAgent),
		Status:   status,
		Message:  msg,
	})
}

func parseBrowser(ua string) string {
	// 简单解析，生产环境可引入 mssola/useragent
	if len(ua) > 64 {
		return ua[:64]
	}
	return ua
}

func parseOS(ua string) string { return "" }
