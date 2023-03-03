package tools

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type PersonalSignature struct {
	Hitokoto string `json:"hitokoto"`
}

func GeneratePersonalSignature() (string, error) {
	resp, err := http.Get("https://v1.hitokoto.cn")
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	defer resp.Body.Close()
	var personalSignature PersonalSignature
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	err = json.Unmarshal(bytes, &personalSignature)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return personalSignature.Hitokoto, nil
}
