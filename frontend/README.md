# Our Records 前端

生活记录应用的前端项目。

## 技术栈

- Vue 3 (Composition API)
- Vue Router 4
- Element Plus
- Axios
- Vite

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:3000

### 3. 构建生产版本

```bash
npm run build
```

## 项目结构

```
src/
├── api/           # API 接口
├── assets/        # 静态资源
├── components/    # 通用组件
├── router/        # 路由配置
├── utils/         # 工具函数
├── views/         # 页面组件
├── App.vue        # 根组件
└── main.js        # 入口文件
```

## 功能模块

- 首页：时光轴
- 记录列表：查看、筛选、删除记录
- 上传记录：上传文字、图片、音频、视频

## 代理配置

开发环境下，API 请求会代理到后端服务 **http://localhost:8088**

配置见 `vite.config.js`:

```javascript
proxy: {
  '/api': {
    target: 'http://localhost:8088',
    changeOrigin: true
  },
  '/uploads': {
    target: 'http://localhost:8088',
    changeOrigin: true
  }
}
```
