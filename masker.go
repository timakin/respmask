package respmask

import (
	"strings"
)

type MaskingFunc func(input string) string

type MaskingRuleType = string

const (
	EmailMasking       MaskingRuleType = "EMAIL_MASKING"
	PasswordMasking    MaskingRuleType = "PASSWORD_MASKING"
	CreditCardMasking  MaskingRuleType = "CREDIT_CARD_MASKING"
	PhoneNumberMasking MaskingRuleType = "PHONE_NUMBER_MASKING"

	PasswordLength = 10
)

var DefaultMaskingRules = map[MaskingRuleType]MaskingFunc{
	EmailMasking: func(input string) string {
		return input[0:1] + "***" + input[strings.LastIndex(input, "@"):]
	},
	PasswordMasking: func(input string) string {
		return strings.Repeat("*", PasswordLength)
	},
	CreditCardMasking: func(input string) string {
		return strings.Repeat("*", len(input)-4) + input[len(input)-4:]
	},
	PhoneNumberMasking: func(input string) string {
		return strings.Repeat("*", len(input)-4) + input[len(input)-4:]
	},
}

func traverseAndMask(node interface{}, keysAndMaskFuncs map[string]MaskingFunc) {
	switch v := node.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if maskFunc, ok := keysAndMaskFuncs[key]; ok {
				if strVal, ok := value.(string); ok {
					v[key] = maskFunc(strVal)
				}
			}
			traverseAndMask(value, keysAndMaskFuncs)
		}
	case []interface{}:
		for _, item := range v {
			traverseAndMask(item, keysAndMaskFuncs)
		}
	}
}

func MaskData(data map[string]interface{}, keysAndMaskFuncs map[string]MaskingFunc) {
	traverseAndMask(data, keysAndMaskFuncs)
}
