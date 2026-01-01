package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCuOrderHandler_getPackageList(t *testing.T) {
	handler := &CuOrderHandler{}

	packages := handler.getPackageList()

	// 验证返回了3种套餐
	assert.Equal(t, 3, len(packages))

	// 验证试用版
	trial := packages[0]
	assert.Equal(t, "trial", trial["id"])
	assert.Equal(t, "试用版", trial["name"])
	assert.Equal(t, 0, trial["price"])
	assert.Equal(t, 1, trial["max_devices"])

	// 验证基础版
	basic := packages[1]
	assert.Equal(t, "basic", basic["id"])
	assert.Equal(t, "基础版", basic["name"])
	assert.Equal(t, 300, basic["price"])

	// 验证专业版
	professional := packages[2]
	assert.Equal(t, "professional", professional["id"])
	assert.Equal(t, "专业版", professional["name"])
	assert.Equal(t, 2000, professional["price"])
}
