// Package auth
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	apiAuth "hotgo/addons/users/api/api/auth"
	usersconsts "hotgo/addons/users/consts"
	"hotgo/addons/users/library/walletcache"
	"hotgo/addons/users/library/walletsign"
	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/addons/users/service"
	"hotgo/addons/users/utility"
	"hotgo/internal/consts"
	"hotgo/internal/library/location"
	"hotgo/internal/library/ratelimit"
	"hotgo/internal/library/token"
	"hotgo/internal/model"
	basesvc "hotgo/internal/service"
)

type sAuth struct{}

func NewAuth() *sAuth {
	return &sAuth{}
}

func init() {
	service.RegisterAuth(NewAuth())
}

// Nonce 获取钱包签名挑战值
func (s *sAuth) Nonce(ctx context.Context, in *apiAuth.NonceReq) (res *apiAuth.NonceRes, err error) {
	address := strings.ToLower(strings.TrimSpace(in.Address))
	if !utility.IsValidWalletAddress(address) {
		err = gerror.New("钱包地址不合法")
		return
	}

	purpose := strings.TrimSpace(in.Purpose)
	if purpose == "" {
		purpose = "login"
	}
	if purpose != "login" && purpose != "action" {
		return nil, gerror.New("签名用途不合法")
	}

	ip := location.GetClientIp(ghttp.RequestFromCtx(ctx))
	if limitErr := s.checkRates(ctx,
		rateLimit{key: fmt.Sprintf("rate:nonce:ip:%s", ip), limit: 30},
		rateLimit{key: fmt.Sprintf("rate:nonce:addr:%s:%s:%s", ip, address, purpose), limit: 30},
	); limitErr != nil {
		return nil, limitErr
	}

	nonce, err := walletsign.Generate(ctx, address)
	if err != nil {
		return nil, err
	}

	res = &apiAuth.NonceRes{Nonce: nonce}
	return
}

// Status 查询账号注册状态
func (s *sAuth) Status(ctx context.Context, in *apiAuth.StatusReq) (res *apiAuth.StatusRes, err error) {
	account, err := normalizeAccount(in.Type, in.Account)
	if err != nil {
		return nil, err
	}

	ip := location.GetClientIp(ghttp.RequestFromCtx(ctx))
	if limitErr := s.checkRates(ctx,
		rateLimit{key: fmt.Sprintf("rate:status:ip:%s", ip), limit: 30},
		rateLimit{key: fmt.Sprintf("rate:status:account:%s:%s:%s", ip, in.Type, account), limit: 30},
	); limitErr != nil {
		return nil, limitErr
	}

	cols := dao.UsersUserAccount.Columns()
	exist, err := dao.UsersUserAccount.Ctx(ctx).
		Where(cols.Type, in.Type).
		Where(cols.Account, account).
		Exist()
	if err != nil {
		return nil, gerror.Wrap(err, "查询账号状态失败")
	}

	isFirstUser, err := s.isFirstUser(ctx)
	if err != nil {
		return nil, err
	}

	forceInvite := false
	if !exist && !isFirstUser {
		value, cfgErr := basesvc.SysAddonsConfig().GetConfigByKey(ctx, "users", "forceInvite")
		if cfgErr != nil {
			return nil, cfgErr
		}
		forceInvite = gconv.Bool(value)
	}

	res = &apiAuth.StatusRes{
		Registered:  exist,
		IsFirstUser: isFirstUser,
		ForceInvite: forceInvite,
	}
	return
}

type rateLimit struct {
	key   string
	limit int
}

func (s *sAuth) checkRates(ctx context.Context, limits ...rateLimit) error {
	for _, item := range limits {
		allowed, err := ratelimit.AllowCount(ctx, item.key, item.limit, time.Minute)
		if err != nil {
			return err
		}
		if !allowed {
			return gerror.New("请求过于频繁，请1分钟后再试")
		}
	}
	return nil
}

func normalizeAccount(accountType, account string) (string, error) {
	account = strings.TrimSpace(account)
	if account == "" {
		return "", gerror.New("账号不能为空")
	}
	switch accountType {
	case usersconsts.UserAuthTypeWallet:
		return strings.ToLower(account), nil
	case usersconsts.UserAuthTypeUsername, usersconsts.UserAuthTypeMobile, usersconsts.UserAuthTypeEmail:
		return account, nil
	default:
		return "", gerror.New("账号类型不合法")
	}
}

// Login 统一登录
func (s *sAuth) Login(ctx context.Context, in *apiAuth.LoginReq) (res *apiAuth.LoginRes, err error) {
	switch in.Type {
	case usersconsts.UserAuthTypeWallet:
		return s.walletLogin(ctx, in)
	default:
		err = gerror.Newf("暂不支持的登录类型：%s", in.Type)
		return
	}
}

// walletLogin 钱包签名登录
func (s *sAuth) walletLogin(ctx context.Context, in *apiAuth.LoginReq) (res *apiAuth.LoginRes, err error) {
	address := strings.ToLower(strings.TrimSpace(in.Account))
	if !utility.IsValidWalletAddress(address) {
		err = gerror.New("钱包地址不合法")
		return
	}

	// 验证钱包签名（开发阶段可通过配置 users.wallet.skipVerify 关闭，方便测试）
	if !g.Cfg().MustGet(ctx, "users.wallet.skipVerify", false).Bool() {
		if err = s.verifyWalletSignature(ctx, address, in.Nonce, in.Signature); err != nil {
			return
		}
	}

	// 查找或创建用户
	user, err := s.findOrCreateByWallet(ctx, address, in.Uid)
	if err != nil {
		return
	}
	if user.Status != consts.StatusEnabled {
		err = gerror.New("账号已被禁用")
		return
	}

	// 保证钱包地址与用户ID映射存在，扫块时需要通过付款地址反查用户
	if err = walletcache.Set(ctx, address, user.Id); err != nil {
		return
	}

	// 签发 token
	identity := &model.Identity{
		Id:       user.Id,
		Username: user.Uid,
		Avatar:   user.Avatar,
		App:      consts.AppApi,
		LoginAt:  gtime.Now(),
	}
	tokenStr, expires, err := token.Login(ctx, identity)
	if err != nil {
		return
	}

	// 更新最后登录信息
	s.updateLastLogin(ctx, user.Id)

	res = &apiAuth.LoginRes{
		Token:   tokenStr,
		Expires: expires,
	}
	return
}

// updateLastLogin 更新最后登录信息（所有登录方式通用）
func (s *sAuth) updateLastLogin(ctx context.Context, userId int64) {
	_, _ = dao.UsersUser.Ctx(ctx).WherePri(userId).Data(g.Map{
		dao.UsersUser.Columns().LastLoginIp:   location.GetClientIp(ghttp.RequestFromCtx(ctx)),
		dao.UsersUser.Columns().LastLoginTime: gtime.Now(),
	}).Update()
}

// verifyWalletSignature 校验钱包签名（nonce 校验 + 签名验证 + 防重放）
func (s *sAuth) verifyWalletSignature(ctx context.Context, address, nonce, signature string) (err error) {
	return walletsign.VerifyAndConsume(ctx, address, nonce, signature)
}

// Logout 退出登录
func (s *sAuth) Logout(ctx context.Context, in *apiAuth.LogoutReq) (res *apiAuth.LogoutRes, err error) {
	err = token.Logout(ghttp.RequestFromCtx(ctx))
	return
}

// findOrCreateByWallet 根据钱包地址查找或创建用户（并发安全）
func (s *sAuth) findOrCreateByWallet(ctx context.Context, address, referrerUid string) (user *entity.UsersUser, err error) {
	// 先查
	user, err = s.getUserByWallet(ctx, address)
	if err != nil || user != nil {
		return
	}

	// 第一个账号不校验推荐关系；后续账号按强制邀请配置处理
	var pid int64
	isFirstUser, err := s.isFirstUser(ctx)
	if err != nil {
		return
	}
	if !isFirstUser {
		if referrerUid == "" {
			forceInvite, cfgErr := basesvc.SysAddonsConfig().GetConfigByKey(ctx, "users", "forceInvite")
			if cfgErr != nil {
				return nil, cfgErr
			}
			if gconv.Bool(forceInvite) {
				return nil, gerror.New("请填写邀请码")
			}
		} else if pid, err = s.resolveReferrer(ctx, referrerUid); err != nil {
			return
		}
	}

	// 首次登录自动注册
	nickname := address
	if len(address) > 12 {
		nickname = address[:6] + "..." + address[len(address)-4:]
	}

	cacheWritten := false
	err = dao.UsersUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		uid := strconv.Itoa(grand.N(10000000, 99999999))
		userData := g.Map{
			dao.UsersUser.Columns().Uid:      uid,
			dao.UsersUser.Columns().LevelId:  1,
			dao.UsersUser.Columns().Nickname: nickname,
			dao.UsersUser.Columns().Status:   consts.StatusEnabled,
		}
		if pid > 0 {
			userData[dao.UsersUser.Columns().Pid] = pid
		}
		id, err := dao.UsersUser.Ctx(ctx).Data(userData).InsertAndGetId()
		if err != nil {
			return gerror.Wrap(err, "创建用户失败")
		}

		if _, err = dao.UsersUserAccount.Ctx(ctx).Data(g.Map{
			dao.UsersUserAccount.Columns().UserId:   id,
			dao.UsersUserAccount.Columns().Type:     usersconsts.UserAuthTypeWallet,
			dao.UsersUserAccount.Columns().Account:  address,
			dao.UsersUserAccount.Columns().Verified: 1,
			dao.UsersUserAccount.Columns().Status:   consts.StatusEnabled,
		}).Insert(); err != nil {
			return gerror.Wrap(err, "创建用户凭证失败")
		}

		// 创建用户关系树
		if err = s.createTree(ctx, id, pid); err != nil {
			return
		}

		// 同步创建资产账户，失败则回滚注册
		if err = basesvc.Asset().EnsureAccount(ctx, uint64(id), consts.AssetAccountUSDT); err != nil {
			return gerror.Wrap(err, "创建资产账户失败")
		}

		// 注册成功前必须写入钱包地址映射，失败则回滚用户创建
		if err = walletcache.Set(ctx, address, id); err != nil {
			return gerror.Wrap(err, "写入钱包用户缓存失败")
		}
		cacheWritten = true
		return
	})
	if err != nil {
		if cacheWritten {
			_ = walletcache.Delete(ctx, address)
		}
		// 并发下其他请求可能已创建同地址凭证（唯一索引冲突），重新查询确认
		if existing, qerr := s.getUserByWallet(ctx, address); qerr == nil && existing != nil {
			// 重新写入缓存，确保后续扫块可以解析
			if cerr := walletcache.Set(ctx, address, existing.Id); cerr != nil {
				err = gerror.Wrap(cerr, "写入钱包用户缓存失败")
				return nil, err
			}
			return existing, nil
		}
		return
	}

	// 创建成功后重新查询，返回完整用户信息
	user, err = s.getUserByWallet(ctx, address)
	return
}

// isFirstUser 是否还没有注册用户
func (s *sAuth) isFirstUser(ctx context.Context) (bool, error) {
	exist, err := dao.UsersUser.Ctx(ctx).Limit(1).Exist()
	if err != nil {
		return false, gerror.Wrap(err, "查询用户是否存在失败")
	}
	return !exist, nil
}

// resolveReferrer 根据推荐人用户号解析其用户ID
func (s *sAuth) resolveReferrer(ctx context.Context, uid string) (pid int64, err error) {
	if uid == "" {
		return 0, nil
	}

	var referrer *entity.UsersUser
	if err = dao.UsersUser.Ctx(ctx).Where(dao.UsersUser.Columns().Uid, uid).Scan(&referrer); err != nil {
		return 0, gerror.Wrap(err, "查询推荐人失败")
	}
	if referrer == nil {
		return 0, gerror.New("推荐人不存在")
	}
	return referrer.Id, nil
}

// createTree 创建用户关系树（直接上级 level=1，其余复制推荐人关系并层级递增）
func (s *sAuth) createTree(ctx context.Context, userId, pid int64) (err error) {
	treeCols := dao.UsersTree.Columns()
	records := g.List{
		g.Map{
			treeCols.Uid:   int(userId),
			treeCols.Pid:   int(userId),
			treeCols.Level: 0,
		},
	}

	// 有推荐人时，复制推荐人的全部上级关系，层级 +1
	if pid > 0 {
		var ancestors []*entity.UsersTree
		if err = dao.UsersTree.Ctx(ctx).
			Where(treeCols.Uid, int(pid)).
			OrderAsc(treeCols.Level).
			Scan(&ancestors); err != nil {
			return gerror.Wrap(err, "查询推荐人关系树失败")
		}
		for _, a := range ancestors {
			records = append(records, g.Map{
				treeCols.Uid:   int(userId),
				treeCols.Pid:   a.Pid,
				treeCols.Level: a.Level + 1,
			})
		}
	}

	if _, err = dao.UsersTree.Ctx(ctx).Data(records).Insert(); err != nil {
		return gerror.Wrap(err, "创建用户关系树失败")
	}
	return
}

// getUserByWallet 根据钱包地址查询用户（校验凭证状态）
func (s *sAuth) getUserByWallet(ctx context.Context, address string) (user *entity.UsersUser, err error) {
	accountCols := dao.UsersUserAccount.Columns()

	var account *entity.UsersUserAccount
	if err = dao.UsersUserAccount.Ctx(ctx).
		Where(accountCols.Type, usersconsts.UserAuthTypeWallet).
		Where(accountCols.Account, address).
		Scan(&account); err != nil {
		err = gerror.Wrap(err, "查询用户凭证失败")
		return
	}
	if account == nil {
		return nil, nil
	}

	// 校验凭证状态
	if account.Status != consts.StatusEnabled {
		return nil, gerror.New("登录凭证已被禁用")
	}

	if err = dao.UsersUser.Ctx(ctx).WherePri(account.UserId).Scan(&user); err != nil {
		err = gerror.Wrap(err, "查询用户信息失败")
	}
	return
}
