package respmask

import (
	"strings"
)

type MaskingMode = string

const (
	ExactMode     MaskingMode = "ExactMode"
	RecursiveMode MaskingMode = "RecursiveMode"
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

func traverseAndMaskExact(node interface{}, funcs map[string]MaskingFunc, nestedKeys ...string) {
	switch v := node.(type) {
	case map[string]interface{}:
		for key, value := range v {
			currentPath := append(nestedKeys, key)
			joinedPath := strings.Join(currentPath, ".")

			if maskFunc, ok := funcs[joinedPath]; ok {
				if strVal, ok := value.(string); ok {
					v[key] = maskFunc(strVal)
				}
			}
			traverseAndMaskExact(value, funcs, currentPath...)
		}
	case []interface{}:
		for index, item := range v {
			if subMap, ok := item.(map[string]interface{}); ok {
				traverseAndMaskExact(subMap, funcs, nestedKeys...)
			} else if strVal, ok := item.(string); ok {
				joinedPath := strings.Join(nestedKeys, ".")
				if maskFunc, ok := funcs[joinedPath]; ok {
					v[index] = maskFunc(strVal)
				}
			}
		}
	}
}

func traverseAndMaskAllLevels(data map[string]interface{}, funcs map[string]MaskingFunc) {
	for k, v := range data {
		if maskFunc, exists := funcs[k]; exists {
			if strVal, ok := v.(string); ok {
				data[k] = maskFunc(strVal)
			}
		}

		switch value := v.(type) {
		case map[string]interface{}:
			traverseAndMaskAllLevels(value, funcs)
		case []interface{}:
			for i, item := range value {
				if submap, ok := item.(map[string]interface{}); ok {
					traverseAndMaskAllLevels(submap, funcs)
					value[i] = submap
				}
			}
		}
	}
}

func Mask(data map[string]interface{}, funcs map[string]MaskingFunc, mode string) {
	switch mode {
	case ExactMode:
		traverseAndMaskExact(data, funcs)
	case RecursiveMode:
		traverseAndMaskAllLevels(data, funcs)
	default:
		panic("invalid masking mode")
	}
}
