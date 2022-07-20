package domains

type UserInfo struct {
	UserId string `validate:"required" form:"userId" query:"userId" json:"userId"`
	Token  string `form:"token" query:"token" json:"token"`
	Name   string `validate:"required" form:"name" query:"name" json:"name"`
	Avatar string `validate:"required" form:"avatar" query:"avatar" json:"avatar"`
	Gender string `validate:"required" form:"gender" query:"gender" json:"gender"	`
	Role   string `form:"role" query:"role" json:"role"`
}
