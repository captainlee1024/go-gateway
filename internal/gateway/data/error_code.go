package data

import "errors"

var (
	ErrServiceNotExit = errors.New("服务不存在")
	ErrInvalidID      = errors.New("无效ID")
)
