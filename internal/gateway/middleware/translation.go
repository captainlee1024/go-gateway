package middleware

import (
	"github.com/captainlee1024/go-gateway/internal/pkg/public"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// TranslationMiddleware 设置 Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置支持语言
		enT := en.New() // 英文翻译器
		zhT := zh.New() // 中文翻译器

		// 设置国际化翻译器
		// 第一个参数是备用 (fallback) 的语言环境，后面的参数是支持的语言环境，可以是多个
		uni := ut.New(zhT, zhT, enT)
		val := validator.New()

		// 根据参数取翻译器实例
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		local := c.DefaultQuery("local", "zh")
		trans, _ := uni.GetTranslator(local)

		// 把翻译器注册到验证器（validator）
		// 校验器是真正做校验的，这里注册到在 gin 中拿到的校验其中
		switch local {
		case "en":
			enTranslations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			zhTranslations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			// 自定义校验和翻译的方法
			// 1. 自定义验证方法
			// 2. 自定义翻译器
			usernameValid(val, trans)
			passwordValid(val, trans)

			break
		}
		c.Set(public.CtxTranslatorKey, trans)
		c.Set(public.CtxValidatorKey, val)
		c.Next()
	}
}

// 管理员 username 校验和翻译
func usernameValid(val *validator.Validate, trans ut.Translator) {
	// 自定义校验器
	// 管理员登录 username 必须为 admin
	val.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "admin"
	})
	// 自定义翻译器
	// username 验证错误的提示
	val.RegisterTranslation("valid_username", trans,
		func(ut ut.Translator) error {
			return ut.Add("valid_username", "{0}输入错误", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("valid_username", fe.Field())
			return t
		})
}

// 管理员登录 password 格式校验
func passwordValid(val *validator.Validate, trans ut.Translator) {
	// 自定义校验器
	// 密码只能包含 字母、数字、下划线
	val.RegisterValidation("valid_password", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.Match(`^[a-zA-Z0-9_]+$`, []byte(fl.Field().String()))
		return matched
	})
	// 自定义翻译器
	// 密码验证错误的提示
	val.RegisterTranslation("valid_password", trans,
		func(ut ut.Translator) error {
			return ut.Add("valid_password", "{0}只能包含数字、字母和下划线", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("valid_password", fe.Field())
			return t
		})
}
