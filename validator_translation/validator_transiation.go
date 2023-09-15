package validator_translation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// Trans 定义一个全局翻译器T
var Trans ut.Translator

// InitTrans 初始化表单参数验证器的翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		//注册验证器
		_ = v.RegisterValidation("mobile", ValidateMobile)

		//初始化翻译器
		zhT := zh.New()
		enT := en.New()
		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)
		// locale 通常取决于 http 请求头的 'Accept-Language'
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		//注册翻译器
		//默认注册英文，en 注册英文 zh 注册中文
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		}
		//自定义错误内容
		_ = v.RegisterTranslation("mobile", Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}非法的手机号码!", false)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
		return
	}
	return
}

// RemoveTopStruct 将返回的结构体名去除掉，只留下需要的字段名
func RemoveTopStruct(fields map[string]string) string {
	//res := map[string]string{}
	var res string
	for _, err := range fields {
		//res[field[strings.LastIndex(field, ".")+1:]] = err
		res += err + " "
	}
	return res
}

// 自定义校验函数
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则表达式判断mobile是否合法
	//pattern := "^((13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(18[0,5-9]))\\d{8}$"
	pattern := `^((13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(18[0,5-9]))\d{8}$`
	ok, _ := regexp.MatchString(pattern, mobile)
	if !ok {
		return false
	}
	return true
}
