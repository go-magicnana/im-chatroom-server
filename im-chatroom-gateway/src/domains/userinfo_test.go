package domains

import (
	"fmt"
	"testing"
)

func TestUserInfo_Validate(t *testing.T) {
	user := UserInfo{
		UserId: "123123123",
	}
	fmt.Println(user)

	//user.Validate()
}
