package logging

import (
	"fmt"
	"runtime"
	"testing"
)

func FuncA() string {
	// 当前函数
	pc, _, _, _ := runtime.Caller(0)
	FuncAName := FuncName(pc)
	// 上层调用函数
	ppc, _, _, _ := runtime.Caller(1)
	PFuncAName := FuncName(ppc)
	callInfo := fmt.Sprint(PFuncAName, " called ", FuncAName)
	return callInfo
}

func TestFuncName(t *testing.T) {
	expect := "TestFuncName called FuncA"
	r := FuncA()
	if r != expect {
		t.Error(r)
	}
}
