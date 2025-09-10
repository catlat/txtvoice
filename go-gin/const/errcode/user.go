package errcode

import (
	"go-gin/internal/errorx"
)

var (
	// 以下定义业务上的错误,注意1开头的是系统错误
	ErrUserNotFound       = errorx.New(20001, "用户不存在")
	ErrUserNameOrPwdFaild = errorx.New(20002, "用户名或者密码错误")
	ErrUserMustLogin      = errorx.New(20003, "请先登录")
	ErrUserNeedLoginAgain = errorx.New(20004, "token已过期,请重新登录")

	// 第三方服务错误
	ErrDLYTUpstream = errorx.New(20020, "视频服务错误")
	ErrASRUpstream  = errorx.New(20021, "语音识别服务错误")
	ErrTranslateUp  = errorx.New(20022, "翻译服务错误")
	ErrTTSUpstream  = errorx.New(20023, "语音合成服务错误")

	// 业务扩展错误
	ErrUserVoiceNotConfigured = errorx.New(20030, "未配置我的声音")
)
