## 内容管理

### 简介




### 使用说明




### 迁移或安装

1、安装 HotGo (2.19.0及以上)

项目介绍：https://github.com/bufanyun/hotgo

2、将当前插件项目拷贝进 HotGo 根目录的 server/addons 目录下

3、在 HotGo 根目录的 server/addons/modules 目录下创建go文件:cms.go，内容如下：
```go
package modules

import _ "hotgo/addons/cms"
```

4、HotGo 后台进入 开发工具->插件管理->找到 内容管理 (cms) 进行安装

5、重启服务即可生效


### 升级说明（已安装旧版本）

v1.0.1 起公告表新增前台公开 ID 与首页推荐字段。已有数据库执行：

```sql
ALTER TABLE `hg_addon_cms_article`
  ADD COLUMN `article_no` varchar(32) NOT NULL DEFAULT '' COMMENT '前台公开ID' AFTER `id`,
  ADD COLUMN `is_recommend` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否推荐：1=推荐，0=不推荐' AFTER `status`;

UPDATE `hg_addon_cms_article`
SET `article_no` = LOWER(SUBSTRING(MD5(CONCAT(`id`, RAND(), NOW())), 1, 20))
WHERE `article_no` = '';

ALTER TABLE `hg_addon_cms_article`
  ADD UNIQUE KEY `uk_article_no` (`article_no`),
  ADD KEY `idx_recommend` (`is_recommend`, `status`);
```

然后重新生成 DAO：

```shell
gf gen dao -l "mysql:root:root@tcp(127.0.0.1:3306)/zztest?loc=Local&parseTime=true&charset=utf8mb4" \
  -p ./addons/cms -t hg_addon_cms_article -r hg_addon_ -j CamelLower -c -m -g default -d model/dao
```


### 常用命令行

```shell
# 接口维护-gen service
gf gen service -s=addons/cms/logic -d=addons/cms/service

```
