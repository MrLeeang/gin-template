package utils

import (
	"gin-template/pkg/utils"
	"testing"
)

func TestGetTimeStr(t *testing.T) {
	// 测试逻辑
	timeStr := utils.GetTimeStr()

	if timeStr == "" {
		t.Errorf("Expected utils.GetTimeStr to string")
	}

}
