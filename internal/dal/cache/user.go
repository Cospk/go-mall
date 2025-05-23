package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Cospk/go-mall/internal/logic/do"
	"github.com/Cospk/go-mall/pkg/enum"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

// SetUserToken 设置用户的AccessToken 和 RefreshToken 缓存
func SetUserToken(ctx context.Context, session *do.SessionInfo) error {
	//log := logger.NewLogger(ctx)
	err := setAccessToken(ctx, session)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	err = setRefreshToken(ctx, session)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	return err
}

func setAccessToken(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, session.AccessToken)
	sessionDataBytes, _ := json.Marshal(session)
	res, err := Redis().Set(ctx, redisKey, sessionDataBytes, time.Hour*2).Result()
	logger.NewLogger(ctx).Debug("redis debug", "res", res, "err", err)
	return err
}

func setRefreshToken(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, session.RefreshToken)
	sessionDataBytes, _ := json.Marshal(session)
	return Redis().Set(ctx, redisKey, sessionDataBytes, 24*time.Hour*7).Err()
}

func DelOldSessionToken(ctx context.Context, session *do.SessionInfo) error {
	oldSession, err := GetUserPlatformSession(ctx, session.UserId, session.Platform)
	if err != nil {
		return err
	}
	if oldSession == nil {
		return nil
	}
	err = DelAccessToken(ctx, oldSession.AccessToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	err = DelayDelRefreshToken(ctx, oldSession.RefreshToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	return nil
}

func GetUserPlatformSession(ctx context.Context, id int64, platform string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, id)
	result, err := Redis().HGet(ctx, redisKey, platform).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	// key 不存在
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	session := new(do.SessionInfo)
	err = json.Unmarshal([]byte(result), &session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func DelAccessToken(ctx context.Context, accessToken string) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, accessToken)
	return Redis().Del(ctx, redisKey).Err()
}

// DelayDelRefreshToken 刷新Token时让旧的RefreshToken 保留一段时间自己过期
func DelayDelRefreshToken(ctx context.Context, refreshToken string) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	// 刷新Token时老的RefreshToken保留的时间(用于发现refresh被窃取)
	return Redis().Expire(ctx, redisKey, 6*time.Hour).Err()
}

func SetUserSession(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, session.UserId)
	sessionDataBytes, _ := json.Marshal(session)
	err := Redis().HSet(ctx, redisKey, session.Platform, sessionDataBytes).Err()
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	return err
}

func LockTokenRefresh(ctx context.Context, refreshToken string) (bool, error) {
	redisLockKey := fmt.Sprintf(enum.REDISKEY_TOKEN_REFRESH_LOCK, refreshToken)
	return Redis().SetNX(ctx, redisLockKey, "locked", 10*time.Second).Result()
}

func UnlockTokenRefresh(ctx context.Context, refreshToken string) error {
	redisLockKey := fmt.Sprintf(enum.REDISKEY_TOKEN_REFRESH_LOCK, refreshToken)
	return Redis().Del(ctx, redisLockKey).Err()
}

func GetRefreshToken(ctx context.Context, refreshToken string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	result, err := Redis().Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	session := new(do.SessionInfo)
	if errors.Is(err, redis.Nil) {
		return session, nil
	}
	json.Unmarshal([]byte(result), &session)

	return session, nil
}

func GetAccessToken(ctx context.Context, accessToken string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, accessToken)
	result, err := Redis().Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	session := new(do.SessionInfo)
	if errors.Is(err, redis.Nil) {
		return session, nil
	}
	json.Unmarshal([]byte(result), &session)

	return session, nil
}

// DelRefreshToken 直接删除RefreshToken缓存  修改密码、退出登录时使用
func DelRefreshToken(ctx context.Context, refreshToken string) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	return Redis().Del(ctx, redisKey).Err()
}

// DelUserSessionOnPlatform Delete user's session on specific platform
func DelUserSessionOnPlatform(ctx context.Context, userId int64, platform string) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, userId)
	return Redis().HDel(ctx, redisKey, platform).Err()
}

func SetPasswordResetToken(ctx context.Context, userId int64, token, code string) error {
	redisKey := fmt.Sprintf(enum.REDISKEY_PASSWORDRESET_TOKEN, token)
	val := fmt.Sprintf("%d:%s", userId, code) // val 以 userId:code 的字符串形式存储
	return Redis().Set(ctx, redisKey, val, 15*time.Minute).Err()
}

func GetPasswordResetToken(ctx context.Context, token string) (userId int64, code string, err error) {
	redisKey := fmt.Sprintf(enum.REDISKEY_PASSWORDRESET_TOKEN, token)
	val, redisErr := Redis().Get(ctx, redisKey).Result()
	if redisErr != nil && redisErr != redis.Nil {
		err = redisErr
		return
	}
	valArr := strings.Split(val, ":")
	if len(valArr) != 2 { // 密码重置Token无对应的缓存, 判定该参数不合法, 此处直接返回
		return
	}
	userId, _ = strconv.ParseInt(valArr[0], 10, 64)
	code = valArr[1]

	return
}

func DelPasswordResetToken(ctx context.Context, token string) error {
	redisKey := fmt.Sprintf(enum.REDISKEY_PASSWORDRESET_TOKEN, token)
	return Redis().Del(ctx, redisKey).Err()
}

// DelUserSessions Delete user's sessions on all platform
func DelUserSessions(ctx context.Context, userId int64) error {
	// 先获取所有平台上的Session信息中
	sessions, err := GetUserAllSessions(ctx, userId)
	if err != nil {
		return err
	}
	// 把所有Session中保存的正在用的Token都过期掉
	for _, sessInfo := range sessions {
		DelOldSessionTokens(ctx, sessInfo)
	}
	// Token过期完成后再删掉Session
	redisKey := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, userId)
	return Redis().Del(ctx, redisKey).Err()
}

// GetUserAllSessions 获取用户在所有platform上的Session
func GetUserAllSessions(ctx context.Context, userId int64) (map[string]*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, userId)
	result, err := Redis().HGetAll(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	// key 不存在
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	sessions := make(map[string]*do.SessionInfo)
	for platform, sessionData := range result {
		session := new(do.SessionInfo)
		err = json.Unmarshal([]byte(sessionData), &session)
		if err != nil {
			return nil, err
		}
		sessions[platform] = session
	}
	//logger.New(ctx).Debug("hgetall user all session", "data", sessions)
	return sessions, nil
}

// DelOldSessionTokens 删除用户旧Session的Token
func DelOldSessionTokens(ctx context.Context, session *do.SessionInfo) error {
	//log := logger.New(ctx)
	oldSession, err := GetUserPlatformSession(ctx, session.UserId, session.Platform)
	if err != nil {
		return err
	}
	if oldSession == nil {
		// 没有旧Session
		return nil
	}
	err = DelAccessToken(ctx, oldSession.AccessToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	err = DelayDelRefreshToken(ctx, oldSession.RefreshToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	return nil
}
