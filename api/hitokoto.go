package api

import (
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"io"
	"net/http"
)

type PersonalSignature struct {
	Hitokoto string `json:"hitokoto"`
}

func GeneratePersonalSignature() string {
	resp, err := http.Get("https://v1.hitokoto.cn")
	if err != nil {
		logger.Error(err)
	}
	defer resp.Body.Close()
	var personalSignature PersonalSignature
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}
	err = json.Unmarshal(bytes, &personalSignature)
	if err != nil {
		logger.Error(err)
	}
	return personalSignature.Hitokoto
}
