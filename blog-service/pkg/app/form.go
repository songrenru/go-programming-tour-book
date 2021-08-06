package app

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ValidError struct {
	Key string
	Message string
}

type ValidErrors []*ValidError

func (this ValidError) Error() string {
	return this.Message
}

func (this ValidErrors) Error() string {
	return strings.Join(this.Errors(), ",")
}

func (this ValidErrors) Errors() []string {
	var errs []string
	for _, validError := range this {
		errs = append(errs, validError.Error())
	}

	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, val := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key: key,
				Message: val,
			})
		}

		return false, errs
	}

	return true, nil
}
