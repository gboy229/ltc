# LTC H5

LTC 黄金AI量化移动端 H5，基于 **Vue 3 + Vite + TypeScript + Tailwind CSS**，独立于 `web`（后台管理端）与 `server`。

首页按照 750 × 1754 设计稿实现，共用 `bg.png` 作为全站背景；PC 端以 400px 固定宽度的手机画布等比放大展示（页面可纵向滚动），移动端宽度 100%。

## 技术栈

| 能力 | 方案 |
| --- | --- |
| 框架 | Vue 3（Composition API + `<script setup>`） |
| 构建 | Vite 7 |
| 语言 | TypeScript（strict） |
| 样式 | Tailwind CSS 3.4 + PostCSS + Autoprefixer |
| 路由 | Vue Router 4（History 模式，页面懒加载） |
| 状态 | Pinia 3 |
| 请求 | Axios（统一封装，自动解包与错误提示） |

> Tailwind 选择 3.4 是为了获得更宽的移动端浏览器兼容性（Tailwind 4 要求 Safari 16.4+ / Chrome 111+）。

## 环境要求

- Node.js `>= 20.19.0`
- pnpm `>= 10`

## 快速开始

```bash
cd h5
pnpm install
pnpm dev
```

默认访问 <http://localhost:5173>。`server.host` 已开启，终端会输出局域网地址，手机与电脑在同一网络下可直接访问调试。

```bash
pnpm build      # 类型检查 + 生产构建，产物在 dist/
pnpm preview    # 预览生产构建
pnpm typecheck  # 仅类型检查
```

## 页面

| 路径 | 页面 | 说明 |
| --- | --- | --- |
| `/` | 首页 | LTC 钱包、公告、入金报单、报单明细 |
| `/notice` | 平台公告 | 公告分类列表、顶部 `gg_bg1.png` Banner |
| `/notice/:id` | 公告详情 | 标题背景使用 `gg_bg2.png`，正文与签名 |
| `/orders` | 报单明细 | 入金总额概览、状态筛选、报单列表、提前退本入口 |
| `/earnings` | 收益明细 | 收益总额、入金报单、直推/静态收益、收支流水 |
| `/community` | 我的社区 | UID/等级、直推/静态收益、邀请好友、我的直推 |
| `/my-assets` | 我的资产 | 总资产、在线提币、支出/收入明细 |
| `/withdraw` | 在线提币 | 选择链、提币数量、到账/手续费、提币记录 |
| `/example` | 能力示例 | Pinia、轻提示、请求封装演示 |
| `*` | 404 | 页面不存在 |

首页的钱包地址、挡位、周期、订单明细目前是本地示例数据，集中在 `src/views/HomeView.vue` 和 `src/config/demo.ts` 中。切换挡位/周期会实时更新选中态，点击钱包地址可复制；「立即报单」会打开购买确认弹窗，确认支付目前为 Web3 钱包占位提示；「查看更多」跳转报单明细页，公告入口跳转公告列表；接入后端时替换为对应接口即可。

公告列表与详情共用 `src/data/notices.ts` 中的演示数据；列表分类图标使用内联 SVG。点击列表条目进入 `/notice/:id`，详情页标题背景使用 `gg_bg2.png`，按设计稿去掉了浏览次数。
首页顶部公告条和侧边导航的「平台公告」均已指向 `/notice`。

首页「立即报单」会打开独立组件 `BuyModal.vue`，弹窗顶部使用 `buy_bg.png`，支付方式图标使用 `icon6.png`；挡位、周期、金额和预计总收益随首页当前选择实时变化，确认支付当前只弹出 Web3 钱包占位提示。

报单明细页的数据在 `src/data/orders.ts`，包含入金总额、今日/累计释放和 4 条演示报单；筛选「全部 / 进行中 / 已结束」为前端本地筛选。点击「申请提前退本」会打开独立组件 `EarlyRefundModal.vue`，按 20% 违约金计算实际到账，确认退本当前以 toast 占位。

提前退本弹窗顶部使用 `end_bg.png`，显示入金金额、违约金、累计已释放和实际到账；支持遮罩、ESC、右上角 X 关闭，打开时锁定页面滚动。

收益明细页的数据在 `src/data/earnings.ts`，收益/支出切换为前端本地筛选；收益总额卡片使用 `sy_bg1.png`，直推/静态收益卡片使用 `sy_bg2.png`，图标分别使用 `icon8.png` 和 `icon7.png`。查看更多当前以 toast 占位。

我的社区页的数据在 `src/data/community.ts`，顶部卡片使用 `my_bg1.png`，邀请好友区使用 `my_bg2.png`；顶部头像使用 `icon9.png`，直推成员头像使用 `icon10.png`，直推/静态收益图标使用 `icon8.png`/`icon5.png`，等级徽章复用 `side/level_bg.png`。邀请地址、邀请码复制已接入剪贴板；「查看全部」当前以 toast 占位。

我的资产页的数据在 `src/data/assets.ts`；总资产卡片使用 `zc_bg1.png` 和 `icon11.png`，在线提币使用 `icon12.png`，支出/收入明细分别使用 `icon_out.png` / `icon_in.png`。支出、收入切换为前端本地筛选；点击「在线提币」进入 `/withdraw`。

在线提币页使用 `w_bg1.png` 作为可用余额卡片背景，按需求去掉了提币地址输入栏；提币数量最少 10 USDT，支持 TRC20 / ERC20 / BEP20 切换，到账数量和手续费随输入/链路实时变化，确认提币当前以 toast 占位。提币记录数据在 `src/data/withdraw.ts`。

## 响应式与 PC 预览

| 视口 | 行为 |
| --- | --- |
| `< 768px`（移动端） | 宽度 100%，内容自然撑高，页面可纵向滚动；适配 iPhone 安全区 |
| `>= 768px`（PC / 平板） | 固定 400px 画布（375 设计宽 × 1.0667），内容等比放大，水平居中；画布高于视口时页面纵向滚动 |

PC 缩放逻辑在 `src/components/H5Shell.vue`：通过 `ResizeObserver` 获取 375px 基准宽度下的内容高度，再统一乘以 `DESKTOP_SCALE = 1.0667`，得到 400px 画布。这样字号和间距同步放大，不会出现“宽度变大了但字还是很小”的情况；画布高于视口时页面正常纵向滚动。所有页面都通过 `H5Shell` 使用 `src/assets/images/bg.png` 作为背景。

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `VITE_APP_TITLE` | `LTC` | 应用标题，注入 `index.html` 与页面标题 |
| `VITE_API_BASE_URL` | `/api` | Axios `baseURL` |
| `VITE_PROXY_TARGET` | `http://localhost:8000` | 开发环境代理目标（GoFrame 默认端口） |
| `VITE_PORT` | `5173` | 开发服务器端口 |

- `.env`：通用配置
- `.env.development`：开发环境，配置后端代理地址
- `.env.production`：生产环境，可改为独立接口域名
- 个人本地覆盖使用 `.env.local` / `.env.development.local`（已加入 `.gitignore`）

开发服务器会代理 `/api` 与 `/home` 到 `VITE_PROXY_TARGET`，不重写路径。生产环境由 Nginx 或网关转发。

## 目录结构

```
h5/
├── src/
│   ├── api/                  # 接口封装（按业务模块拆分）
│   ├── assets/
│   │   └── images/           # 设计稿切图：bg / banner / bd_bg1 / buy_bg / my_bg1 / my_bg2 / gg_bg1 / gg_bg2 / sy_bg1 / sy_bg2 / zc_bg1 / logo / icon1-12 / icon_in / icon_out / side
│   ├── components/
│   │   ├── AppHeader.vue     # 首页/公告页/公告详情页共用 Header
│   │   ├── BackButton.vue    # 统一返回箭头（宽度自适应，不使用 w-9）
│   │   ├── EarlyRefundModal.vue # 提前退本确认弹窗
│   │   ├── BuyModal.vue      # 确认入金购买弹窗
│   │   ├── H5Shell.vue       # 全站背景、移动端满宽、PC 400px 画布、toast 根节点
│   │   ├── NoticeCategoryIcon.vue # 公告分类内联 SVG 图标
│   │   ├── PageHeading.vue   # 返回箭头 + 页面标题 + 副标题
│   │   └── SideNav.vue       # 全局侧边导航
│   ├── config/               # demo 数据（UID / 等级 / 钱包地址）
│   ├── data/                 # 公告、报单、收益、社区、资产等演示数据
│   ├── router/               # 路由与页面标题
│   ├── stores/               # Pinia Store（含 sideNav 开关）
│   ├── styles/               # Tailwind 入口、主题与 H5 画布样式
│   ├── types/                # 全局类型声明
│   ├── utils/                # auth / clipboard / http / toast
│   ├── views/                # HomeView / NoticeView / NoticeDetailView / OrderDetailView / EarningsView / CommunityView / AssetsView / WithdrawView / ExampleView / NotFoundView
│   ├── App.vue
│   ├── env.d.ts
│   └── main.ts
├── index.html
├── postcss.config.js
├── tailwind.config.js
├── tsconfig.json
└── vite.config.ts
```

## 开发约定

### 样式

- 优先使用 Tailwind 原子类，避免内联 `style`（`H5Shell` 的缩放值除外，需运行时计算）。
- 主题色在 `tailwind.config.js` 中扩展：`gold-*` 与 `ink-*`，页面背景统一由 `H5Shell` 承担。
- 仅在必要时于 `src/styles/index.css` 的 `@layer` 中补充全局样式。
- 动态拼接的 class 必须保证完整字符串出现在源码中，否则 Tailwind 不会生成对应样式。
- 页面内边距使用 `px-4`，顶部/底部使用 `env(safe-area-inset-*)` 适配刘海屏。
- 页面返回统一使用 `src/components/BackButton.vue`（`h-9`、宽度自适应，不写 `w-9`），列表页、详情页、示例页保持一致。

已提供的移动端工具类：

| 类名 | 作用 |
| --- | --- |
| `pt-safe` / `pb-safe` | 兼容 iPhone 刘海与底部安全区（需配合 `viewport-fit=cover`） |
| `min-h-app` | 全屏高度，带 `100dvh` 回退 |

### 移动端字号与触控

移动端没有强制统一的字号法规，但通常参考以下标准：

| 来源 | 关键建议 |
| --- | --- |
| Apple HIG | 正文约 17pt，辅助文字 11–13pt，点击区域不小于 44×44pt |
| Material 3 | 正文 14–16sp，次级说明 12–13sp，标签 11–14sp，点击区域 48×48dp |
| 国内 H5 实践 | 正文 14–16px，次级 12–13px，辅助文字不低于 11.5px，尽量避免 11px 以下 |

本项目按 375px 基准的约定：

| 用途 | 字号 | 页面示例 |
| --- | --- | --- |
| 品牌标题 | 24px | LTC |
| 卡片标题 | 18px | 入金报单、报单明细 |
| 正文 / 金额 | 14–17px | 入金金额、$1000.00 U |
| 次级说明 | 12–13.5px | 选择挡位、投资周期、公告正文 |
| 辅助信息 | 12px（320px 极窄屏 11px） | 时间、公告日期、钱包地址 |
| 紧凑标签 | 11–12px | 最低 1,000 USDT、挡位金额 |

- 正文建议不低于 14px；常规辅助文字不低于 12px；11px 只用于 320px 极窄屏的密集标签，避免 10px 及以下。
- 主按钮高度 44px（立即报单），分段控件 36px（挡位、周期、查看更多）。
- 750px 设计稿按 2 倍换算到 375px 基准。本项目在换算值基础上整体上调约 1–1.5px，以改善中文阅读，同时保留原有设计的紧凑感。
- 如果后续需要支持更极端的字号无障碍设置，建议引入 rem 或 CSS 变量作为字号 token，而不是继续散落写固定 px。

### 图片与静态资源

| 目录 | 适用场景 | 引用方式 |
| --- | --- | --- |
| `src/assets/images/` | 页面图片、背景、插画、图标 | `import` 后绑定到 `src`，或 CSS `url()` |
| `public/` | 必须保持固定路径的文件，如第三方校验文件 | 以 `/` 开头的绝对路径 |
| 后端返回的图片 | 用户头像、上传图片、CDN 资源 | 直接使用接口返回的 URL，不要放进仓库 |

当前设计切图均位于 `src/assets/images`：

- `bg.png`：全站背景（在 `src/styles/index.css` 中引用）
- `banner.jpg`：首页主视觉
- `gg_bg1.png`：平台公告页顶部 Banner 背景
- `gg_bg2.png`：公告详情页标题背景
- `bd_bg1.png`：报单明细页「入金总额」卡片背景
- `buy_bg.png`：确认入金弹窗顶部钱包插画
- `icon6.png`：DAPP 钱包支付图标
- `icon7.png`：静态收益图标
- `icon8.png`：直推收益图标
- `sy_bg1.png`：收益明细页「收益总额」卡片背景
- `sy_bg2.png`：直推/静态收益卡片背景
- `icon9.png`：社区顶部带桂冠头像
- `icon10.png`：直推成员头像
- `my_bg1.png`：我的社区顶部身份卡片背景
- `my_bg2.png`：邀请好友区域背景
- `icon11.png`：USDT 总资产图标
- `icon12.png`：在线提币图标
- `icon_in.png` / `icon_out.png`：收入/支出流水图标
- `zc_bg1.png`：我的资产总资产卡片背景
- `w_bg1.png`：在线提币可用余额卡片背景
- `end_bg.png`：提前退本弹窗顶部徽章插画
- `logo.png`：头部品牌图与 favicon
- `icon1.png` ~ `icon5.png`：公告、铃铛、菜单、钱包、USDT 金币
- `side/bg.png`：侧边导航整体背景
- `side/center_bg.png`：侧边导航 UID / 等级卡片背景
- `side/logo.png`：侧边导航品牌徽章
- `side/level_bg.png`：等级徽章底图（叠加 T2 文本）
- `side/nav1.png` ~ `side/nav5.png`：报单明细、我的社区、我的资产、平台公告、收益明细图标

组件中使用：

```vue
<script setup lang="ts">
import banner from '@/assets/images/banner.jpg'
</script>

<template>
  <img :src="banner" alt="黄金AI量化交易模型" />
</template>
```

- 小于 4KB 的图片会被 Vite 自动内联为 base64，其余输出到 `dist/assets` 并带哈希。
- 不要用 `import` 引用 `public/` 下的文件，也不要把图片放在 `views/`、`components/` 里。
- 缺图的小图标使用内联 SVG（Heroicons 风格），通过 Tailwind 的尺寸/颜色类控制。

### 路由

路由集中在 `src/router/index.ts`，页面按需加载。新增页面时：

```ts
{
  path: '/member',
  name: 'member',
  component: () => import('@/views/MemberView.vue'),
  meta: { title: '我的' },
}
```

`meta.title` 会自动追加到 `document.title`。使用 History 模式，部署时需要服务端回退到 `index.html`。

### 状态

Store 放在 `src/stores`，使用组合式写法，可参考 `stores/counter.ts`。

### 接口

接口统一放在 `src/api` 下按模块封装，页面中不要直接使用 axios：

```ts
// src/api/member.ts
import { http } from '@/utils/http'

export interface MemberInfo {
  id: number
  nickname: string
}

export function getMemberInfo() {
  return http.get<MemberInfo>('/member/info')
}

export function saveMember(data: { nickname: string }) {
  return http.post<void>('/member/save', data)
}
```

`http` 的行为：

- `baseURL` 取自 `VITE_API_BASE_URL`
- 自动携带 `Authorization: Bearer <token>`
- 响应 `{ code, message, data }` 中 `code === 0` 时直接返回 `data`
- 业务失败与网络异常统一 toast，并 reject `ApiError`
- `code === 61`（GoFrame 未授权）时清除本地 token

Token 读写使用 `src/utils/auth.ts`（存储键为 `h5_token`）。登录页与路由守卫尚未内置，接入登录流程时可在接口 61 分支中跳转登录页。

### 侧边导航

`src/components/SideNav.vue` 全局挂载在 `H5Shell` 内，任意页面都可以通过 Pinia store 打开：

```ts
import { useSideNavStore } from '@/stores/sideNav'

const sideNav = useSideNavStore()
sideNav.openNav()
```

- 从右侧滑出，宽度为画布的 76.5%，带遮罩、关闭按钮，支持 ESC 和点击遮罩关闭。
- 支持手势：左滑打开、右滑关闭（水平位移 ≥60px 且大于垂直位移，避免与页面纵向滚动冲突）。
- 打开时锁定 body 滚动；PC 端会先回到画布顶部，避免抽屉顶部被滚动裁切。
- 抽屉位于 `H5Shell` 内，会随手机画布一起等比缩放；PC 端遮罩只覆盖画布，不覆盖整个浏览器。
- 支持 `uid` / `level` / `address` 三个 props，默认取 `src/config/demo.ts` 中的演示数据，首页与抽屉共用同一个钱包地址。
- 菜单行为集中在 `SideNav.vue` 的 `handleItem`：
  - 报单明细：跳转 `/orders`
  - 平台公告：跳转 `/notice`
  - 收益明细：跳转 `/earnings`
  - 我的社区：跳转 `/community`
  - 我的资产：跳转 `/my-assets`
  - 侧边导航五个入口均已接入页面
- 地址复制复用 `src/utils/clipboard.ts`，与首页头部保持一致。

### 轻提示

`src/utils/toast.ts` 会挂载到 `H5Shell` 的 `#h5-toast-root`：

- 移动端：固定在视口顶部，滚动时始终可见
- PC 端：固定在视口顶部，与 400px 画布居中对齐；页面滚动时提示仍然可见
- 未找到根节点时自动回退到 `document.body`

## 部署

```bash
pnpm build
```

将 `dist/` 部署到静态服务器，并配置 SPA 回退与接口转发，例如 Nginx：

```nginx
server {
    listen 80;
    server_name h5.example.com;
    root /path/to/h5/dist;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

如使用 `/home` 等其它前缀，按需追加 `location`。接口独立域名部署时，把 `VITE_API_BASE_URL` 改为完整地址并处理跨域。

## 常见问题

**手机无法访问开发服务器？**  
确认手机与电脑在同一局域网，使用终端输出的 Network 地址，并检查系统防火墙；`vite.config.ts` 已设置 `server.host = true`。

**PC 上字体还是偏小 / 想调整画布大小？**  
PC 画布由 `H5Shell.vue` 中的 `PC_CANVAS_WIDTH` 控制，当前为 `400`（375 × 1.0667）。调大即可整体放大字号与间距，页面会相应变长并可滚动；调小则更接近手机原尺寸。也可以直接调整 `DESKTOP_SCALE`。

**接口 404 或跨域？**  
开发环境检查后端是否运行在 `VITE_PROXY_TARGET`，以及请求路径是否带 `/api` 前缀；生产环境检查 Nginx 转发与 `VITE_API_BASE_URL`。

**Tailwind 类名不生效？**  
确认文件位于 `tailwind.config.js` 的 `content` 范围内，且动态类名为完整字符串。

**如何升级 Tailwind 4？**  
安装 `tailwindcss@4` 与 `@tailwindcss/vite`，将 `src/styles/index.css` 改为 `@import "tailwindcss";`，移除 PostCSS 中的 Tailwind 配置，并检查旧工具类与浏览器兼容性。
