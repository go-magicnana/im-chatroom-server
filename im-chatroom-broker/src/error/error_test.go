package error

import (
	"fmt"
	"testing"
)

func TestError_Format(t *testing.T) {
	fmt.Println(InvalidRequest.Format("haha"))
}
