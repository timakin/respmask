package respmask

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMaskData(t *testing.T) {
	tests := []struct {
		name             string
		initInputFunc    func() map[string]interface{}
		keysAndMaskFuncs map[string]MaskingFunc
		expected         map[string]interface{}
	}{
		{
			name: "default masking",
			initInputFunc: func() map[string]interface{} {
				return map[string]interface{}{
					"email":       "user@example.com",
					"password":    "supersecret",
					"credit_card": "1234567890123456",
					"phone":       "1234567890",
				}
			},
			keysAndMaskFuncs: map[string]MaskingFunc{
				"email":       DefaultMaskingRules[EmailMasking],
				"credit_card": DefaultMaskingRules[CreditCardMasking],
				"password":    DefaultMaskingRules[PasswordMasking],
				"phone":       DefaultMaskingRules[PhoneNumberMasking],
			},
			expected: map[string]interface{}{
				"email":       "u***@example.com",
				"credit_card": "************3456",
				"password":    "**********",
				"phone":       "******7890",
			},
		},
		{
			name: "Basic nested masking",
			initInputFunc: func() map[string]interface{} {
				return map[string]interface{}{
					"email": "user@example.com",
					"details": map[string]interface{}{
						"password": "supersecret",
					},
				}
			},
			keysAndMaskFuncs: map[string]MaskingFunc{
				"email":    DefaultMaskingRules[EmailMasking],
				"password": DefaultMaskingRules[PasswordMasking],
			},
			expected: map[string]interface{}{
				"email": "u***@example.com",
				"details": map[string]interface{}{
					"password": "**********",
				},
			},
		},
		{
			name: "Array and nested masking",
			initInputFunc: func() map[string]interface{} {
				return map[string]interface{}{
					"users": []interface{}{
						map[string]interface{}{
							"email": "user1@example.com",
						},
						map[string]interface{}{
							"email": "user2@example.com",
						},
					},
				}
			},
			keysAndMaskFuncs: map[string]MaskingFunc{
				"email": DefaultMaskingRules[EmailMasking],
			},
			expected: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"email": "u***@example.com",
					},
					map[string]interface{}{
						"email": "u***@example.com",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := tt.initInputFunc() // データを初期化
			MaskData(inputData, tt.keysAndMaskFuncs)
			if diff := cmp.Diff(inputData, tt.expected); diff != "" {
				t.Errorf("got %v, want %v", inputData, tt.expected)
			}
		})
	}
}
