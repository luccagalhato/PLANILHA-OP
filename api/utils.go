package api

import (
	"crypto/rand"
	_ "embed" //EMB
	"encoding/base64"
	"io"
	sql "op/database"
	"op/models"
	"sync"
	"time"
)

var (
	//go:embed excel.xlsx
	op []byte
)
var mu sync.Mutex
var connectionLinx *sql.SQLStr
var files map[string]*models.Excel = map[string]*models.Excel{}

// genKey ...
func genKey() string {
	b := make([]byte, 12)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

//SetSQLConn ...
func SetSQLConn(l *sql.SQLStr) {
	connectionLinx = l
}

//GetFile ...
func GetFile(id string) *models.Excel {
	mu.Lock()
	defer mu.Unlock()
	return files[id]
}

func addExcel(name string, value io.Reader) string {
	key := genKey()
	mu.Lock()
	files[key] = &models.Excel{name, value}
	mu.Unlock()

	go func() {
		time.Sleep(time.Second * 5)
		mu.Lock()
		delete(files, key)
		mu.Unlock()
	}()

	return key
}
