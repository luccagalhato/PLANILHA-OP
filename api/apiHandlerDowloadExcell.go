package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	sql "op/database"
	"op/models"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func DownloadExcell(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/octet-stream")
	NewOP := NewOP{}
	if err := json.NewDecoder(r.Body).Decode(&NewOP); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Bad Request",
			Error:  "",
			Data:   err.Error(),
		})
		return
	}

	if file := GetFile(NewOP.Cod); file != nil {
		w.Header().Set("File-Name", file.Name)
		io.Copy(w, file.Value)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

var f *excelize.File

func gerarExcel(rst []sql.Op, op string) string {
	var keys string

	f = excelOP(rst)

	if b, err := f.WriteToBuffer(); err == nil {
		keys = addExcel(fmt.Sprintf("%s_%s.xlsx", op, time.Now().Format("20060102_150405")), b)
	}
	return keys
	// pretty.Print(rst)
}

func excelOP(OPS []sql.Op) *excelize.File {
	f, err := excelize.OpenReader(bytes.NewReader(op))
	if err != nil {
		log.Println(err)
	}
	sheet := f.GetSheetName(1)
	for i, Op := range OPS {
		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+2), Op.Ref)
		f.SetCellStr(sheet, fmt.Sprintf("B%d", i+2), Op.Ean)
		f.SetCellStr(sheet, fmt.Sprintf("C%d", i+2), Op.Nome)
		f.SetCellStr(sheet, fmt.Sprintf("D%d", i+2), Op.Cor)
		f.SetCellStr(sheet, fmt.Sprintf("E%d", i+2), Op.Tamanho)
		f.SetCellStr(sheet, fmt.Sprintf("F%d", i+2), Op.Uni)
		f.SetCellStr(sheet, fmt.Sprintf("G%d", i+2), Op.Quanti)
		f.SetCellStr(sheet, fmt.Sprintf("H%d", i+2), Op.Ex1)
		f.SetCellStr(sheet, fmt.Sprintf("I%d", i+2), Op.Ex2)
		f.SetCellStr(sheet, fmt.Sprintf("J%d", i+2), Op.Ex20)
		f.SetCellStr(sheet, fmt.Sprintf("K%d", i+2), Op.Grupo)
	}
	return f
}
