package pkg

import (
	"encoding/json"
	"os"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 存在且是文件
func IsFile(filePath string) bool {
	s, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// WriteJSON 将对象写入json文件
func WriteJSON(v interface{}, path string) error {
	file, e := os.Create(path)
	if e != nil {
		return e
	}
	defer file.Close()
	jsonEncode := json.NewEncoder(file)
	jsonEncode.SetIndent("", "\t")
	return jsonEncode.Encode(v)
}
