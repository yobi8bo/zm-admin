package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// 业务错误码
const (
	CodeSuccess = 200

	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeTooManyReqs  = 429
	CodeServerError  = 500

	// 用户模块 1000x
	CodeUserNotFound       = 10001
	CodeUserDisabled       = 10002
	CodePasswordError      = 10003
	CodeUsernameExists     = 10004
	CodeAdminUserProtected = 10005
	CodeAdminRoleProtected = 10006

	// 角色模块 1001x
	CodeRoleNotFound             = 10011
	CodeRoleCodeExists           = 10012
	CodeRoleInUse                = 10013
	CodeAdminRoleDeleteProtected = 10014

	// 菜单模块 1002x
	CodeMenuNotFound    = 10021
	CodeMenuHasChildren = 10022

	// 部门模块 1003x
	CodeDeptNotFound    = 10031
	CodeDeptHasChildren = 10032
	CodeDeptHasUsers    = 10033

	// 认证模块 1004x
	CodeTokenExpired        = 10041
	CodeTokenInvalid        = 10042
	CodeRefreshTokenInvalid = 10043
	CodeCaptchaInvalid      = 10044
)

var msgMap = map[int]string{
	CodeSuccess:                  "success",
	CodeBadRequest:               "请求参数错误",
	CodeUnauthorized:             "未登录或token已失效",
	CodeForbidden:                "无权限访问",
	CodeNotFound:                 "资源不存在",
	CodeTooManyReqs:              "请求过于频繁，请稍后再试",
	CodeServerError:              "服务器内部错误",
	CodeUserNotFound:             "用户不存在",
	CodeUserDisabled:             "用户已被禁用",
	CodePasswordError:            "密码错误",
	CodeUsernameExists:           "用户名已存在",
	CodeAdminUserProtected:       "管理员账户不能删除",
	CodeAdminRoleProtected:       "管理员账户不能重新分配角色",
	CodeRoleNotFound:             "角色不存在",
	CodeRoleCodeExists:           "角色标识已存在",
	CodeRoleInUse:                "角色已分配用户，无法删除",
	CodeAdminRoleDeleteProtected: "管理员角色不能删除",
	CodeMenuNotFound:             "菜单不存在",
	CodeMenuHasChildren:          "存在子菜单，无法删除",
	CodeDeptNotFound:             "部门不存在",
	CodeDeptHasChildren:          "存在子部门，无法删除",
	CodeDeptHasUsers:             "部门下存在用户，无法删除",
	CodeTokenExpired:             "token已过期",
	CodeTokenInvalid:             "token无效",
	CodeRefreshTokenInvalid:      "refresh token无效",
	CodeCaptchaInvalid:           "验证码错误或已过期",
}

func GetMsg(code int) string {
	if msg, ok := msgMap[code]; ok {
		return msg
	}
	return "未知错误"
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: CodeSuccess, Msg: "success", Data: data})
}

func SuccessPage(c *gin.Context, list any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: PageData{List: list, Total: total, Page: page, PageSize: pageSize},
	})
}

func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{Code: code, Msg: GetMsg(code), Data: nil})
}

func FailWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Msg: msg, Data: nil})
}

func ServerError(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Code: CodeServerError, Msg: GetMsg(CodeServerError), Data: nil})
}
