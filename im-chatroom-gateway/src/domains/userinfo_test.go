package domains

import "testing"

func TestUserInfo_Validate(t *testing.T) {
	user := UserInfo{
		UserId: "123123123",
	}

	user.Validate()
}
