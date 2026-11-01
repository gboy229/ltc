// Package modules 插件隐式注册
// 每新增一个插件，在此目录下新建 <addonName>.go 文件，内容仅一行匿名导入，例如：
//
//	//go:build !no_<addonName>
//
//	package modules
//
//	import _ "hotgo/addons/<addonName>"
//
// 顶部 //go:build !no_<addonName> 用于编译期隔离：当某插件编译不过时，
// 可用 `go build -tags "no_<addonName>"` 临时将其排除出构建，不影响主站发布。
package modules
