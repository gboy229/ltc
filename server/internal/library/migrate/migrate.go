// Package migrate 数据库迁移
// 结构迁移（migrate）与数据填充（seed）分离：
//   - 结构迁移：manifest/migrations/ 下按版本号命名的 SQL 文件，记录在 hg_schema_migrations 表，只执行一次。
//   - 数据填充：manifest/seeders/ 下的 SQL 文件，幂等、可重复执行，不记录版本。
//
// 插件对应目录：addons/<name>/migrations/ 与 addons/<name>/seeders/。
package migrate

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	// DefaultScope 默认作用域：底座
	DefaultScope = "base"
	// DefaultDir 默认结构迁移目录
	DefaultDir = "manifest/migrations"
	// DefaultSeedDir 默认数据填充目录
	DefaultSeedDir = "manifest/seeders"
	// migrationTable 迁移记录表
	migrationTable = "hg_schema_migrations"
)

// Migration 单个结构迁移
type Migration struct {
	Version  int    // 版本号
	Name     string // 名称
	UpFile   string // 升级 SQL 文件
	DownFile string // 降级 SQL 文件（预留，暂未使用）
}

// ensureTable 确保迁移记录表存在
func ensureTable(ctx context.Context) (err error) {
	_, err = g.DB().Exec(ctx, `CREATE TABLE IF NOT EXISTS `+migrationTable+` (
		id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		scope      VARCHAR(64)  NOT NULL DEFAULT 'base',
		version    INT          NOT NULL,
		name       VARCHAR(255) NOT NULL DEFAULT '',
		applied_at DATETIME     DEFAULT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY uk_scope_version (scope, version)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库迁移记录表'`)
	return
}

// load 扫描目录中的迁移文件，按版本号升序返回
func load(dir string) (list []*Migration, err error) {
	if !gfile.Exists(dir) {
		err = gerror.Newf("迁移目录不存在: %s", dir)
		return
	}

	files, err := gfile.ScanDirFile(dir, "*.sql", false)
	if err != nil {
		return
	}

	m := make(map[int]*Migration)
	for _, f := range files {
		base := gfile.Basename(f)
		if strings.HasSuffix(base, ".down.sql") {
			continue
		}
		version, name, ok := parseName(strings.TrimSuffix(base, ".sql"))
		if !ok {
			continue
		}
		m[version] = &Migration{Version: version, Name: name, UpFile: f}
	}

	for _, mig := range m {
		list = append(list, mig)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Version < list[j].Version })
	return
}

// parseName 解析文件名，返回版本号和名称。文件名格式：{版本号}_{名称}
func parseName(name string) (version int, n string, ok bool) {
	idx := strings.Index(name, "_")
	if idx <= 0 {
		return
	}
	v, err := strconv.Atoi(name[:idx])
	if err != nil {
		return
	}
	return v, name[idx+1:], true
}

// execScript 执行一段可能包含多条语句的 SQL 脚本，按分号逐条执行
func execScript(ctx context.Context, script string) (err error) {
	for _, stmt := range splitStatements(script) {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err = g.DB().Exec(ctx, stmt); err != nil {
			return
		}
	}
	return
}

// splitStatements 按分号拆分 SQL 脚本，跳过字符串字面量（单引号）中的分号，
// 并兼容双单引号转义（”）与反斜杠转义（\）。
func splitStatements(script string) (statements []string) {
	var buf strings.Builder
	inQuote := false

	for i := 0; i < len(script); i++ {
		c := script[i]
		switch c {
		case '\\': // 反斜杠转义：连同下一个字符一起保留
			buf.WriteByte(c)
			if i+1 < len(script) {
				buf.WriteByte(script[i+1])
				i++
			}
		case '\'':
			buf.WriteByte(c)
			if inQuote && i+1 < len(script) && script[i+1] == '\'' {
				// 双单引号 = 字符串内的转义单引号
				buf.WriteByte(script[i+1])
				i++
			} else {
				inQuote = !inQuote
			}
		case ';':
			if inQuote {
				buf.WriteByte(c)
			} else {
				statements = append(statements, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteByte(c)
		}
	}
	if buf.Len() > 0 {
		statements = append(statements, buf.String())
	}
	return
}

// appliedVersions 获取指定作用域下已应用的版本号集合
func appliedVersions(ctx context.Context, scope string) (map[int]bool, error) {
	var rows []struct {
		Version int
	}
	if err := g.DB().Model(migrationTable).Ctx(ctx).Fields("version").Where("scope", scope).Scan(&rows); err != nil {
		return nil, err
	}
	m := make(map[int]bool, len(rows))
	for _, r := range rows {
		m[r.Version] = true
	}
	return m, nil
}

// Up 执行未应用的结构迁移（建表/改表），幂等：已应用的版本自动跳过
func Up(ctx context.Context, scope, dir string) (err error) {
	if err = ensureTable(ctx); err != nil {
		return
	}
	list, err := load(dir)
	if err != nil {
		return
	}
	applied, err := appliedVersions(ctx, scope)
	if err != nil {
		return
	}

	appliedCount := 0
	for _, mig := range list {
		if applied[mig.Version] {
			continue
		}
		sql := gfile.GetContents(mig.UpFile)
		if strings.TrimSpace(sql) == "" {
			g.Log().Warningf(ctx, "迁移文件为空，跳过: %s", mig.UpFile)
			continue
		}
		if err = execScript(ctx, sql); err != nil {
			return gerror.Wrapf(err, "迁移失败 scope:%s version:%d name:%s", scope, mig.Version, mig.Name)
		}
		if _, err = g.DB().Model(migrationTable).Ctx(ctx).Data(g.Map{
			"scope":      scope,
			"version":    mig.Version,
			"name":       mig.Name,
			"applied_at": gtime.Now(),
		}).Insert(); err != nil {
			return gerror.Wrapf(err, "记录迁移失败 scope:%s version:%d", scope, mig.Version)
		}
		appliedCount++
		g.Log().Infof(ctx, "迁移成功 scope:%s version:%d name:%s", scope, mig.Version, mig.Name)
	}
	if appliedCount == 0 {
		g.Log().Infof(ctx, "无待执行迁移 scope:%s", scope)
	}
	return
}

// Seed 执行数据填充，幂等、可重复执行，不记录版本
func Seed(ctx context.Context, dir string) (err error) {
	if !gfile.Exists(dir) {
		err = gerror.Newf("种子目录不存在: %s", dir)
		return
	}
	files, err := gfile.ScanDirFile(dir, "*.sql", false)
	if err != nil {
		return
	}
	sort.Strings(files)

	for _, f := range files {
		sql := gfile.GetContents(f)
		if strings.TrimSpace(sql) == "" {
			g.Log().Warningf(ctx, "种子文件为空，跳过: %s", f)
			continue
		}
		if err = execScript(ctx, sql); err != nil {
			return gerror.Wrapf(err, "填充数据失败: %s", gfile.Basename(f))
		}
		g.Log().Infof(ctx, "填充数据成功: %s", gfile.Basename(f))
	}
	return
}
