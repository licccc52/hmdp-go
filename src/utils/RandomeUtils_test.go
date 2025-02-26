package utils

import "testing"

func TestGenerateRandNumber(t *testing.T) {
	randomStr := RandomUtil.GenerateVerifyCode()
	t.Log(randomStr)
}

func TestGenerateRandStr(t *testing.T) {
	randomStr := RandomUtil.GenerateRandomStr(10)
	t.Log(randomStr)
}
