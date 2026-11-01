// Package addons
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package addons

import "embed"

// SQLFS 内嵌所有插件的 migrations/seeders SQL。
// 这样部署时只需要一个二进制，不需要额外同步 addons 目录。
//
//go:embed */migrations/*.sql */seeders/*.sql
var SQLFS embed.FS
