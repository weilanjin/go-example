package validator

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// https://github.com/go-playground/validator
// 每一个结构体都可以看成是一棵树
type Nested struct {
	Email string `validate:"email"`
}

type T struct {
	Age    int `validate:"eq=10"`
	Nested Nested
}

func validateEmail(input string) bool {
	if pass, _ := regexp.MatchString(`^([\w._]{2,10})@(\w+).([a-z]{2,4})$`, input); pass {
		return true
	}
	return false
}

func validate(v any) (validateResult bool, errmsg string) {
	validateResult = true
	errmsg = "success"
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	for i := 0; i < vv.NumField(); i++ {
		fieldVal := vv.Field(i)
		tagContent := vt.Field(i).Tag.Get("validate")
		k := fieldVal.Kind()
		switch k {
		case reflect.Int:
			val := fieldVal.Int()
			tagValStr := strings.Split(tagContent, "=")
			tagVal, err := strconv.ParseInt(tagValStr[1], 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			if val != tagVal {
				errmsg = "validate int failed, tag is:" + strconv.FormatInt(
					tagVal, 10,
				)
				validateResult = false
				return
			}
		case reflect.String:
			val := fieldVal.String()
			tagValStr := tagContent
			switch tagValStr {
			case "email":
				isEmail := validateEmail(val)
				if !isEmail {
					errmsg = "validate mail failed, field val is:" + val
					validateResult = false
					return
				}
			}
		case reflect.Struct:
			// 深度优先遍历
			// 一个递归过程
			valInter := fieldVal.Interface()
			validateResult, errmsg = validate(valInter)
			if !validateResult {
				return
			}
		}
	}
	return
}

func Test1(t *testing.T) {
	var a = T{Age: 11, Nested: Nested{Email: "1237@qq.com"}}
	ok, msg := validate(a)
	log.Println(ok, msg)

	var a1 = T{Age: 10, Nested: Nested{Email: "1237@qq.com"}}
	ok, msg = validate(a1)
	log.Println(ok, msg)
}
