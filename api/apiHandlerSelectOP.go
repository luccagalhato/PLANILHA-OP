package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	sql "op/database"
	"op/models"
)

type NewOP struct {
	Cod string `json:"cod,omitempty"`
}

func SelectOP(w http.ResponseWriter, r *http.Request) {
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
	if NewOP.Cod == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Bad Request",
			Error:  "",
			Data:   "Favor Inserir uma OP válida",
		})
	}

	data, err := connectionLinx.SelectOPDatabase(NewOP.Cod)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Bad Request",
			Error:  "",
			Data:   err.Error(),
		})
		return
	}
	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Not Found",
			Error:  "",
			Data:   fmt.Sprintf("op '%s' nao encontrada", NewOP.Cod),
		})
		return
	}

	if hasNullGTIN(data) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Invalid Code",
			Error:  "",
			Data:   fmt.Sprintf("existem produtos da op '%s' sem GTIN cadastrado", NewOP.Cod),
		})
		return
	}

	excels, err := excelOP(data).WriteToBuffer()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status: "internal error",
			Error:  "",
			Data:   "Erro ao gerar Excel",
		})
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/json")
	b, _ := ioutil.ReadAll(excels)
	json.NewEncoder(w).Encode(models.ResponseExcel{
		Status: "OK",
		Error:  "",
		Data:   data,
		Excel:  b,
	})
}

func hasNullGTIN(ops []sql.Op) bool {
	for i := 0; i < len(ops); i++ {
		if ops[i].Ean == nil {
			return true
		}
	}
	return false
}
