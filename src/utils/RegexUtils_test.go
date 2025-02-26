package utils

import "testing"

func TestPhoneValid(t *testing.T) {
	phone1 := "15871773042"
	res := RegexUtil.IsPhoneValid(phone1)
	if res {
		t.Log("测试成功...")
	} else {
		t.Log("测试失败")
	}
}

