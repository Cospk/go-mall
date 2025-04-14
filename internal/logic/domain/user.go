package domain

import (
	"github.com/Cospk/go-mall/internal/dal/cache"
	"github.com/Cospk/go-mall/internal/dal/dao"
	"github.com/Cospk/go-mall/internal/logic/do"
	"github.com/Cospk/go-mall/pkg/auth"
	"github.com/Cospk/go-mall/pkg/enum"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type UserDomain struct {
	ctx     *gin.Context
	userDao *dao.UserDao
}

func NewUserDomain(ctx *gin.Context) *UserDomain {
	return &UserDomain{
		ctx:     ctx,
		userDao: dao.NewUserDao(ctx),
	}
}

func (domain *UserDomain) LoginUser(Name, password, platform string) (*do.TokenInfo, error) {
	existedUser, err := domain.userDao.FindUserByName(Name)
	if err != nil {
		return nil, errcode.Wrap("UserDomainSvcLoginUserError", err)
	}
	if existedUser.ID == 0 {
		return nil, errcode.ErrUserPasswordError
	}
	if !utils.BcryptCompare(existedUser.Password, password) {
		return nil, errcode.ErrUserPasswordError
	}
	token, err := domain.GenAuthToken(existedUser.ID, platform, "")
	return token, err
}

// GenAuthToken 生成AccessToken和RefreshToken
func (domain *UserDomain) GenAuthToken(userId int64, platform string, sessionId string) (*do.TokenInfo, error) {
	user := domain.GetUserBaseInfo(userId)
	// 处理参数异常情况, 用户不存在、被删除、被禁用
	if user.ID == 0 || user.IsBlocked == enum.UserBlockStateBlocked {
		err := errcode.ErrUserInvalid
		return nil, err
	}

	userSession := new(do.SessionInfo)
	userSession.UserId = userId
	userSession.Platform = platform
	if sessionId == "" {
		sessionId = auth.GenSessionId(userId)
	}
	userSession.SessionId = sessionId
	accessToken, RefreshToken, err2 := auth.GenUserAuthToken(userSession.UserId, userSession.Platform, userSession.SessionId)
	if err2 != nil {
		return nil, errcode.Wrap("UserDomainSvcGenAuthTokenError", err2)
	}
	// 设置userSession 缓存
	userSession.AccessToken = accessToken
	userSession.RefreshToken = RefreshToken

	//向缓存中写入session
	err := cache.SetUserToken(domain.ctx, userSession)
	if err != nil {
		return nil, errcode.Wrap("设置Token缓存时发生错误", err)
	}
	err = cache.DelOldSessionToken(domain.ctx, userSession)
	if err != nil {
		return nil, errcode.Wrap("删除旧Token缓存时发生错误", err)
	}
	err = cache.SetUserSession(domain.ctx, userSession)
	if err != nil {
		return nil, errcode.Wrap("设置Session缓存时发生错误", err)
	}

	srvCreateTime := time.Now()

	return &do.TokenInfo{
		AccessToken:   userSession.AccessToken,
		RefreshToken:  userSession.RefreshToken,
		Duration:      int64((time.Hour * 2).Seconds()),
		SrvCreateTime: srvCreateTime,
	}, nil
}

func (domain *UserDomain) RefreshToken(refreshToken string) (*do.TokenInfo, error) {
	log := logger.NewLogger(domain.ctx)
	ok, err := cache.LockTokenRefresh(domain.ctx, refreshToken)
	defer cache.UnlockTokenRefresh(domain.ctx, refreshToken)
	if err != nil {
		err = errcode.Wrap("刷新Token时设置Redis锁发生错误", err)
		return nil, err
	}
	if !ok {
		err = errcode.ErrTooManyRequests
		return nil, err
	}
	tokenSession, err := cache.GetRefreshToken(domain.ctx, refreshToken)
	if err != nil {
		log.Error("GetRefreshTokenCacheErr", "err", err)
		// 服务断发生错误一律提示客户端Token有问题
		// 生产环境可以做好监控日志中这个错误的监控
		err = errcode.ErrToken
		return nil, err
	}
	// refreshToken没有对应的缓存
	if tokenSession == nil || tokenSession.UserId == 0 {
		err = errcode.ErrToken
		return nil, err
	}
	userSession, err := cache.GetUserPlatformSession(domain.ctx, tokenSession.UserId, tokenSession.Platform)
	if err != nil {
		log.Error("GetUserPlatformSessionErr", "err", err)
		err = errcode.ErrToken
		return nil, err
	}
	// 请求刷新的RefreshToken与UserSession中的不一致, 证明这个RefreshToken已经过时
	// RefreshToken被窃取或者前端页面刷Token不是串行的互斥操作都有可能造成这种情况
	if userSession.RefreshToken != refreshToken {
		// 记一条警告日志
		log.Warn("ExpiredRefreshToken", "requestToken", refreshToken, "newToken", userSession.RefreshToken, "userId", userSession.UserId)
		// 错误返回Token不正确, 或者更精细化的错误提示已在xxx登录如不是您本人操作请xxx
		err = errcode.ErrToken
		return nil, err
	}

	// 重新生成Token  因为不是用户主动登录所以sessionID与之前的保持一致
	tokenInfo, err := domain.GenAuthToken(tokenSession.UserId, tokenSession.Platform, tokenSession.SessionId)
	if err != nil {
		err = errcode.Wrap("GenAuthTokenErr", err)
		return nil, err
	}
	return tokenInfo, nil
}

func (domain *UserDomain) GetUserBaseInfo(userId int64) *do.UserBaseInfo {
	return &do.UserBaseInfo{
		ID:        12345678,
		Nickname:  "Kevin",
		LoginName: "kev@gomall.com",
		Verified:  1,
		Avatar:    "",
		Slogan:    "",
		IsBlocked: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (domain *UserDomain) RegisterUser(info *do.UserBaseInfo, password string) (*do.UserBaseInfo, error) {
	existedUser, err := domain.userDao.FindUserByLoginName(info.LoginName)
	if err != nil {
		return nil, errcode.Wrap("UserDomainSvcRegisterUserError", err)
	}
	if existedUser.LoginName != "" { // 用户名已经被占用
		return nil, errcode.ErrUserNameOccupied
	}
	passwordHash, err := utils.BcryptPassword(password)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}
	userModel, err := domain.userDao.CreateUser(info, passwordHash)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}
	err = utils.CopyStruct(info, userModel)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}

	return info, nil
}

func (domain *UserDomain) LogoutUser(userId int64, platform string) error {
	log := logger.NewLogger(domain.ctx)
	userSession, err := cache.GetUserPlatformSession(domain.ctx, userId, platform)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}
	// 删掉用户当前会话中的AccessToken和RefreshToken
	err = cache.DelAccessToken(domain.ctx, userSession.AccessToken)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}
	err = cache.DelRefreshToken(domain.ctx, userSession.RefreshToken)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}
	// 删掉用户在对应平台上的Session
	err = cache.DelUserSessionOnPlatform(domain.ctx, userId, platform)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}

	return nil
}
