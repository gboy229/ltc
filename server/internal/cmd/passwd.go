// Package cmd
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package cmd

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/grand"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
)

const (
	// defaultPasswordLength 自动生成密码的长度，与后台密码长度限制保持一致
	defaultPasswordLength = 16
)

var (
	passwordLowerChars   = []byte("abcdefghijkmnpqrstuvwxyz")
	passwordUpperChars   = []byte("ABCDEFGHJKLMNPQRSTUVWXYZ")
	passwordDigitChars   = []byte("23456789")
	passwordSpecialChars = []byte("!@#%^&*")
	passwordAllChars     = []byte("abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#%^&*")
)

var (
	Passwd = &gcmd.Command{
		Name:  "passwd",
		Brief: "修改后台账号密码",
		Description: `用于在命令行自动修改后台账号的登录密码，支持通过账号名或账号ID指定账号。
不传 -p 时将自动生成一个16位复杂随机密码并输出，也可通过 -p 手动指定（6~16位）。`,
		Examples: `
修改指定账号的密码：  go run main.go passwd -u=admin -p=123456
修改指定账号ID的密码：go run main.go passwd -i=1 -p=123456
自动生成复杂密码：    go run main.go passwd -u=admin
`,
		Arguments: []gcmd.Argument{
			{Name: "username", Short: "u", Brief: "后台账号（登录用户名）"},
			{Name: "id", Short: "i", Brief: "后台账号ID"},
			{Name: "password", Short: "p", Brief: "新密码，长度为6~16位"},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return handlePasswd(ctx, parser)
		},
	}
)

// handlePasswd 修改后台账号密码
func handlePasswd(ctx context.Context, parser *gcmd.Parser) (err error) {
	var (
		username  = gstr.Trim(parser.GetOpt("username").String())
		memberId  = parser.GetOpt("id").Int64()
		password  = parser.GetOpt("password").String()
		generated bool
	)

	if username == "" && memberId <= 0 {
		return gerror.New("请通过 -u=账号 或 -i=账号ID 指定要修改密码的后台账号")
	}

	cols := dao.AdminMember.Columns()
	mod := dao.AdminMember.Ctx(ctx)
	if memberId > 0 {
		mod = mod.WherePri(memberId)
	}
	if username != "" {
		mod = mod.Where(cols.Username, username)
	}

	var member *entity.AdminMember
	if err = mod.Scan(&member); err != nil {
		return gerror.Wrap(err, "查询后台账号失败")
	}
	if member == nil {
		return gerror.New("后台账号不存在，请确认指定的账号是否正确")
	}

	if password == "" {
		if password, err = generatePassword(defaultPasswordLength); err != nil {
			return gerror.Wrap(err, "生成随机密码失败")
		}
		generated = true
	} else if n := len([]rune(password)); n < 6 || n > 16 {
		return gerror.New("密码长度需在6~16位之间")
	}

	// 兼容历史数据中没有密码盐的账号
	salt := member.Salt
	if salt == "" {
		salt = grand.S(6)
	}

	data := g.Map{
		cols.PasswordHash: gmd5.MustEncryptString(password + salt),
	}
	if member.Salt == "" {
		data[cols.Salt] = salt
	}

	result, err := dao.AdminMember.Ctx(ctx).WherePri(member.Id).Data(data).Update()
	if err != nil {
		return gerror.Wrap(err, "修改后台账号密码失败")
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return gerror.New("修改后台账号密码失败，账号可能已被删除")
	}

	if generated {
		fmt.Printf("未指定 -p，已自动生成新密码：%s\n", password)
	}

	g.Log().Infof(ctx, "后台账号[%s](id:%d)密码修改成功", member.Username, member.Id)
	return
}

// generatePassword 生成包含大小写字母、数字和特殊字符的随机密码
func generatePassword(length int) (string, error) {
	// 先各取一个字符，确保密码同时包含四类字符
	groups := [][]byte{passwordLowerChars, passwordUpperChars, passwordDigitChars, passwordSpecialChars}
	chars := make([]byte, 0, length)
	for _, group := range groups {
		c, err := randomPasswordByte(group)
		if err != nil {
			return "", err
		}
		chars = append(chars, c)
	}

	for len(chars) < length {
		c, err := randomPasswordByte(passwordAllChars)
		if err != nil {
			return "", err
		}
		chars = append(chars, c)
	}

	// 打乱顺序，避免固定的字符类型排列
	for i := len(chars) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := n.Int64()
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars), nil
}

// randomPasswordByte 从字符集中随机取出一个字符
func randomPasswordByte(chars []byte) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, err
	}
	return chars[n.Int64()], nil
}
