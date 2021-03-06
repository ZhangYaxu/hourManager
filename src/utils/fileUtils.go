/**
* @Project: hourManager
* @Package
* @Description: TODO
* @author : wj
* @date Date : 2018/12/14/ 9:52
* @version V1.0
 */
package utils

import (
	"github.com/beego/bee/logger"
	"github.com/xuri/excelize"
	"strconv"
	"strings"
	"time"
)

/**
 * 导出EXCEL
 * openFile 模板文件
 * details 详细列表 格式：${detail}
 * others 其他内容替换（单个） 格式：${XXX}
 * return fileName
 */
func ExportExcel(openFile *excelize.File, sheetName string, details []interface{}, others map[string]interface{}) (fileName string) {

	rows := openFile.GetRows(sheetName)

	for i, row := range rows {
		for j, colCell := range row {
			//设置单独字段
			for k, v := range others {
				if colCell == k {
					axis := intToAscii(int32(j+65)) + strconv.Itoa(i+1)
					beeLogger.Log.Info("替换列:" + axis + "||" + k)
					openFile.SetCellValue(sheetName, axis, v)
				}
			}
			//设置明细数据
			if colCell == "${detail}" {
				rowIndex := i + 1
				height := openFile.GetRowHeight(sheetName, rowIndex)
				for _, slice := range details {
					openFile.SetSheetRow(sheetName, "B"+strconv.Itoa(rowIndex), slice)
					openFile.SetRowHeight(sheetName, rowIndex, height)

					rowIndex++
					if rowIndex >= i+7 {
						openFile.InsertRow(sheetName, rowIndex)
					}
				}
			}
		}
	}

	fileName = "static/excel/" + time.Now().Format("20060102150405") + ".xlsx"

	return fileName
}

//转换成ASCII码 65-A  97-a
func intToAscii(value int32) (ascii string) {
	return strings.Replace(strconv.QuoteRuneToASCII(value), "'", "", -1)
}
