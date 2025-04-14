package auth

import (
	"errors"
	"fmt"
	"github.com/Cospk/go-mall/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWT相关配置
const (
	accessTokenSecret  = "your_access_token_secret_key"  // 访问token密钥
	refreshTokenSecret = "your_refresh_token_secret_key" // 刷新token密钥
	accessTokenExpiry  = time.Hour * 2                   // 访问token有效期2小时
	refreshTokenExpiry = time.Hour * 24 * 7              // 刷新token有效期7天
)

// CustomClaims 自定义JWT Claims
type CustomClaims struct {
	UserId    int64  `json:"user_id"`
	Platform  string `json:"platform"`   // 平台信息
	SessionId string `json:"session_id"` // 会话ID
	jwt.RegisteredClaims
}

// genAccessToken 生成JWT格式的访问token
func genAccessToken(uid int64, platform string, sessionId string) (string, error) {
	claims := CustomClaims{
		UserId:    uid,
		Platform:  platform,
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-mall",
			Subject:   fmt.Sprintf("%d", uid),
			ID:        utils.RandNumStr(10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessTokenSecret))
}

// genRefreshToken 生成JWT格式的刷新token
func genRefreshToken(uid int64, platform string, sessionId string) (string, error) {
	claims := CustomClaims{
		UserId:    uid,
		Platform:  platform,
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-mall",
			Subject:   fmt.Sprintf("%d", uid),
			ID:        utils.RandNumStr(10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshTokenSecret))
}

// GenUserAuthToken 生成用户认证token对
func GenUserAuthToken(uid int64, platform string, sessionId string) (accessToken, refreshToken string, err error) {
	accessToken, err = genAccessToken(uid, platform, sessionId)
	if err != nil {
		return
	}
	refreshToken, err = genRefreshToken(uid, platform, sessionId)
	if err != nil {
		return
	}

	return
}

// GenPasswordResetToken 生成密码重置token
func GenPasswordResetToken(userId int64) (string, error) {
	// 使用特殊的过期时间和用途
	claims := CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24小时有效期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-mall",
			Subject:   fmt.Sprintf("%d", userId),
			ID:        utils.RandNumStr(10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessTokenSecret))
}

// GenSessionId 生成会话ID
func GenSessionId(userId int64) string {
	return fmt.Sprintf("%d-%d-%s", userId, time.Now().Unix(), utils.RandNumStr(6))
}

// ParseUserIdFromToken 从Token中解析出userId
func ParseUserIdFromToken(accessToken string) (userId int64, platform string, sessionId string, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessTokenSecret), nil
	})

	if err != nil {
		return 0, "", "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserId, claims.Platform, claims.SessionId, nil
	}

	return 0, "", "", errors.New("invalid token")
}

// ParseRefreshToken 从刷新Token中解析出userId
func ParseRefreshToken(refreshToken string) (userId int64, platform string, sessionId string, err error) {
	token, err := jwt.ParseWithClaims(refreshToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshTokenSecret), nil
	})

	if err != nil {
		return 0, "", "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserId, claims.Platform, claims.SessionId, nil
	}

	return 0, "", "", errors.New("invalid refresh token")
}

// ValidateAccessToken 验证访问Token是否有效
func ValidateAccessToken(tokenString string) (bool, error) {
	_, _, _, err := ParseUserIdFromToken(tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ValidateRefreshToken 验证刷新Token是否有效
func ValidateRefreshToken(tokenString string) (bool, error) {
	_, _, _, err := ParseRefreshToken(tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}

// RefreshAccessToken 使用刷新Token获取新的访问Token
func RefreshAccessToken(refreshToken string) (newAccessToken string, err error) {
	userId, platform, sessionId, err := ParseRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// 生成新的访问Token
	newAccessToken, err = genAccessToken(userId, platform, sessionId)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
