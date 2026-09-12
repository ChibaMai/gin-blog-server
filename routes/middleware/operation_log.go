package middleware

import (
	"bytes"
	"gin-blog/dao"
	"gin-blog/model"
	"gin-blog/model/dto"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"io"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var optMap = map[string]string{
	"Article":      "文章",
	"BlogInfo":     "博客信息",
	"Category":     "分类",
	"Comment":      "评论",
	"FriendLink":   "友链",
	"Menu":         "菜单",
	"Message":      "留言",
	"OperationLog": "操作日志",
	"Resource":     "资源权限",
	"Role":         "角色",
	"Tag":          "标签",
	"User":         "用户",
	"Page":         "页面",
	// "Talk":         "说说",
	// "Login":        "登录",
	"POST":   "新增或修改",
	"PUT":    "修改",
	"DELETE": "删除",
}

func GetOptString(key string) string {
	return optMap[key]
}

// 在 gin 中获取 Response Body 内容: 对 gin 的 ResponseWriter 进行包装, 每次往请求方响应数据时, 将响应数据返回出去
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer // 响应体缓存
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // 将响应数据存到缓存中
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s) // 将响应数据存到缓存中
	return w.ResponseWriter.WriteString(s)
}

// 记录操作日志中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不记录 GET 请求操作记录 (太多了) 和 文件上传操作记录 (请求体太长)
		if c.Request.Method == "GET" || strings.Contains(c.Request.RequestURI, "upload") {
			c.Next()
			return
		}

		// 包装响应 writer
		blw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		uuid := utils.GetFromContext[string](c, "uuid")

		// 从 session 中取用户的 SessionInfo
		val := sessions.Default(c).Get(KEY_USER + uuid)
		if val == nil {
			utils.Logger.Warn("操作日志中间件: 从session中获取用户信息失败",
				zap.String("uuid", uuid),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.RequestURI),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    r.ERROR_TOKEN_RUNTIME,
				"message": "用户会话已过期",
			})
			return
		}

		var sessionInfo dto.SessionInfo
		err := utils.Json.Unmarshal(val.(string), &sessionInfo) // 反序列化
		if err != nil {
			utils.Logger.Error("操作日志中间件: 反序列化session失败",
				zap.Error(err),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    r.ERROR_TOKEN_RUNTIME,
				"message": "用户信息解析失败",
			})
			return
		}

		reqBody, _ := io.ReadAll(c.Request.Body) // 记录请求信息
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		optResource := getOptResource(c.HandlerName())
		operationLog := model.OperationLog{
			OptModule:     GetOptString(optResource),
			OptType:       GetOptString(c.Request.Method),
			OptUrl:        c.Request.RequestURI,
			OptMethod:     c.HandlerName(),
			OptDesc:       GetOptString(c.Request.Method) + GetOptString(optResource),
			RequestParam:  string(reqBody),
			RequestMethod: c.Request.Method,
			UserId:        sessionInfo.UserInfoId,
			Nickname:      sessionInfo.Nickname,
			IpAddress:     sessionInfo.IpAddress,
			IpSource:      sessionInfo.IpSource,
		}

		c.Next()

		operationLog.ResponseData = blw.body.String() // 从缓存中获取响应体内容

		// 异步写入数据库，避免阻塞请求处理
		go func() {
			if err := dao.Create(&operationLog); err != nil {
				utils.Logger.Error("操作日志中间件: 写入数据库失败",
					zap.Error(err),
					zap.Any("log", operationLog),
				)
			}
		}()
	}
}

// example: "gin-blog/api/v1.(*Resource).Delete-fm" => "Resource"
// 安全的字符串解析，避免 index out of range
func getOptResource(handlerName string) string {
	parts := strings.Split(handlerName, ".")
	if len(parts) < 2 {
		return "Unknown"
	}

	s := parts[1]
	if len(s) < 3 {
		return "Unknown"
	}

	// 去除开头的 "(*" 和结尾的 ")"
	return s[2 : len(s)-1]
}
