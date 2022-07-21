package domains

import (
	"github.com/go-playground/validator/v10"
	"im-chatroom-gateway/apierror"
	"sync"
)

type CustomValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	c.lazyInit()
	e := c.validate.Struct(i)

	if e != nil {
		return apierror.InvalidParameter.Replace(e.Error())
	} else {
		return nil
	}
}

func (c *CustomValidator) lazyInit() {
	c.once.Do(func() {
		c.validate = validator.New()
	})
}
