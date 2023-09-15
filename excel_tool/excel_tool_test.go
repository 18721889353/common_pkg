package excel_tool

import (
	"testing"
)

func TestExcel(t *testing.T) {
	headers := []string{"用户名", "性别", "年龄"}
	values := [][]interface{}{
		{"测试", "1", "男"},
		{
			"ggr1", "1", "男",
		},
	}
	f, err := ExportExcel("sheet1", headers, values)
	if err != nil {
		t.Log(err)
	}
	err = f.SaveAs("./static/export/" + "test1.xlsx")
	if err != nil {
		t.Log(err)
	}
}
