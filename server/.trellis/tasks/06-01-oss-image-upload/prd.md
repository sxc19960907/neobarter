# 图片上传对接阿里云 OSS

## Goal

实现物品图片/用户头像的上传功能。后端提供上传接口，存储抽象为接口层，开发环境用本地文件存储跑通，生产环境切换到阿里云 OSS。前端发布物品/编辑资料时真正能上传图片并预览。

## Requirements

### 存储抽象（参照现有 SMS provider 模式）
- 定义 `Storage` 接口：`Upload(file, filename) (url, error)` + `Delete(url) error`
- `LocalProvider`：开发环境，存到本地 `uploads/` 目录，通过静态文件服务访问
- `OSSProvider`：生产环境，对接阿里云 OSS SDK
- 通过配置 `oss.provider`（local / aliyun）切换

### 上传接口
- `POST /v1/upload/image` — 单图上传（multipart/form-data）
  - 校验文件类型（jpg/png/webp/gif）
  - 校验大小（≤5MB，参照 SRS F-I-003）
  - 返回可访问的 URL
- 需要登录鉴权
- 文件名按 `年月/uuid.ext` 规则生成，避免冲突和路径穿越

### 前端对接
- 物品发布页 Upload 组件接入真实上传（替换当前 `beforeUpload={() => false}` 占位）
- 上传后把返回的 URL 收集进表单的 images 字段
- 个人中心头像上传
- 上传中 loading、失败提示、预览

### 安全要求
- 文件类型白名单（按扩展名 + MIME 双校验）
- 文件大小限制
- 防止路径穿越（不信任客户端文件名）
- 上传接口限流

## Acceptance Criteria

- [ ] 后端上传接口正常返回可访问 URL（本地 provider）
- [ ] 非法文件类型/超大文件被拒绝
- [ ] 本地存储的图片能通过返回的 URL 访问到
- [ ] 前端物品发布页能上传多张图片并预览
- [ ] 前端个人中心能上传头像
- [ ] OSS provider 代码完成（生产可用，开发可不实际连接）
- [ ] 上传相关逻辑有单元测试
- [ ] 编译/测试/CI 通过

## Constraints

- 最多 9 张图片/物品，单张 ≤5MB（SRS F-I-003）
- 本地存储目录 `uploads/` 加入 .gitignore
- 阿里云 OSS SDK：`github.com/aliyun/aliyun-oss-go-sdk`
