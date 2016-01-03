package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const (
	pubLen int = 12
)

func checkOTP(w http.ResponseWriter, r *http.Request, dal *Dal) {
	if r.URL.Query()["otp"] == nil || r.URL.Query()["nonce"] == nil || r.URL.Query()["id"] == nil {
		reply(w, "", "", "", MISSING_PARAMETER, "", dal)
		return
	}
	otp := r.URL.Query()["otp"][0]
	nonce := r.URL.Query()["nonce"][0]
	id := r.URL.Query()["id"][0]
	name := ""

	if len(otp) < pubLen {
		reply(w, otp, name, nonce, BAD_OTP, id, dal)
		return
	}
	pub := otp[:pubLen]

	k, err := dal.GetKey(pub)
	if err != nil {
		reply(w, otp, name, nonce, BAD_OTP, id, dal)
		return
	} else {
		k, err = Gate(k, otp)
		if err != nil {
			reply(w, otp, name, nonce, err.Error(), id, dal)
			return
		} else {
			err = dal.UpdateKey(k)
			if err != nil {
				log.Println("fail to update key counter/session")
				return
			}

			name := k.Name
			reply(w, otp, name, nonce, OK, id, dal)
			return
		}
	}
}

func Sign(values []string, key []byte) []byte {
	payload := ""
	for _, v := range values {
		payload += v + "&"
	}
	payload = payload[:len(payload)-1]

	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}

func loadKey(id string, dal *Dal) ([]byte, error) {
	i, err := dal.GetApp(id)
	if err != nil {
		return []byte{}, errors.New(NO_SUCH_CLIENT)
	}

	return i, nil
}

func reply(w http.ResponseWriter, otp, name, nonce, status, id string, dal *Dal) {
	values := []string{}
	key := []byte{}
	err := errors.New("")

	values = append(values, "nonce="+nonce)
	values = append(values, "otp="+otp)
	if status != MISSING_PARAMETER {
		key, err = loadKey(id, dal)
		if err == nil {
			values = append(values, "name="+name)
			values = append(values, "status="+status)
		} else {
			values = append(values, "status="+err.Error())
		}
	} else {
		values = append(values, "status="+status)
	}
	values = append(values, "t="+time.Now().Format(time.RFC3339))
	if status != MISSING_PARAMETER {
		values = append(values, "h="+base64.StdEncoding.EncodeToString(Sign(values, key)))
	}

	ret := ""
	for _, v := range values {
		ret += v + "\n"
	}

	w.Write([]byte(ret))
}

func runAPI(dal *Dal, host, port string) {
	r := mux.NewRouter()

	r.HandleFunc("/wsapi/2.0/verify", func(w http.ResponseWriter, r *http.Request) {
		checkOTP(w, r, dal)
	}).Methods("GET")

	http.Handle("/", r)
	log.Printf("Listening on: %s:%s...", host, port)
	http.ListenAndServe(host+":"+port, nil)
}
