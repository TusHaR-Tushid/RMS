package utilities

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Key string

const (
	UserContextKey Key = "values"
)

func Decoder(r *http.Request, inter interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&inter)
	if err != nil {
		logrus.Printf("decoderr error:%v", err)
		return err
	}
	return nil
}

func Encoder(w http.ResponseWriter, inter interface{}) error {
	err := json.NewEncoder(w).Encode(&inter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("encoder error:%v", err)
		return err
	}
	return nil
}
