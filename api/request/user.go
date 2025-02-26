package request

// UserLogin 用户登录请求，同时绑定请求头和请求体 gin的issue也提到过：https://github.com/gin-gonic/gin/issues/2309#issuecomment-2020168668
type UserLogin struct {
	Body struct {
		LoginName string `json:"login_name" binding:"required,e164|email"`
		Password  string `json:"password" binding:"required,min=8"`
	}
	Header struct {
		Platform string `json:"platform" header:"platform" binding:"required,oneof=H5 APP"`
	}
}
