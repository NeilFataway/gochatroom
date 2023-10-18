package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseError struct {
	s string
}

func (e *BaseError) Error() string {
	return e.s
}

func (e *BaseError) ResponseError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": e.Error(),
		"data":    nil,
	})
}

type ResponsiveError interface {
	error
	ResponseError(ctx *gin.Context)
}

func WrapBaseError(err error) *BaseError {
	return &BaseError{
		s: err.Error(),
	}
}

type ConflictError struct {
	*BaseError
}

func (e *ConflictError) ResponseError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusConflict, gin.H{
		"code":    http.StatusConflict,
		"message": e.Error(),
		"data":    nil,
	})
}

type NotfoundError struct {
	*BaseError
}

func (e *NotfoundError) ResponseError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusConflict, gin.H{
		"code":    http.StatusConflict,
		"message": e.Error(),
		"data":    nil,
	})
}

type BadRequestError struct {
	*BaseError
}

func (e *BadRequestError) ResponseError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": e.Error(),
		"data":    nil,
	})
}

func WrapBadRequestError(err error) *BadRequestError {
	return &BadRequestError{
		BaseError: WrapBaseError(err),
	}
}

type PermissionError struct {
	*BaseError
}

func (e *PermissionError) ResponseError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"code":    http.StatusBadRequest,
		"message": e.Error(),
		"data":    nil,
	})
}

// 基本错误类型
var ClosedMsgReceived = &BaseError{s: "收到客户端关闭信息"}

// 消息解析错误类型
var InValidMessageType = &BadRequestError{&BaseError{s: "消息类型解析失败"}}
var InValidUserID = &BadRequestError{&BaseError{s: "用户ID解析失败"}}
var InValidRoomName = &BadRequestError{&BaseError{s: "房间名解析失败"}}

// 资源冲突错误类型
var roomConflict = &ConflictError{&BaseError{s: "房间已存在，命名冲突"}}
var userAlreadyJoined = &ConflictError{&BaseError{s: "冲突啦，用户已在聊天室中"}}
var UserAlreadyExists = &ConflictError{&BaseError{s: "用户名已存在"}}

// 资源不存在错误类型
var UserAvatarNotExists = &NotfoundError{&BaseError{s: "用户头像地址不存在"}}
var UserNotExists = &NotfoundError{&BaseError{s: "用户ID不存在"}}
var roomNotExist = &NotfoundError{&BaseError{s: "房间不存在"}}
var userNotInRoom = &NotfoundError{&BaseError{s: "用户不在聊天室里"}}

// 权限不足错误类型
var PermissionNotAllowd = &PermissionError{&BaseError{s: "用户权限不足，无法操作"}}
