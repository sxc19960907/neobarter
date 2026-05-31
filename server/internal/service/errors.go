package service

import "errors"

// 通用业务错误
var (
	ErrForbidden = errors.New("无权执行此操作")
	ErrNotFound  = errors.New("资源不存在")
)
