package main

import (
	"encoding/hex"
	"errors"
	"github.com/conformal/yubikey"
)

const (
	OK                = "OK"
	REPLAYED_OTP      = "REPLAYED_OTP"
	MISSING_PARAMETER = "MISSING_PARAMETER"
	BAD_OTP           = "BAD_OTP"
	BAD_SIGNATURE     = "BAD_SIGNATURE"
	NO_SUCH_CLIENT    = "NO_SUCH_CLIENT"
)

func Gate(key *Key, otp string) (*Key, error) {
	priv, err := getSecretKey(key.Secret)
	if err != nil {
		return nil, err
	}
	token, err := getToken(otp, priv)
	if err != nil {
		return nil, errors.New(BAD_OTP)
	}

	if token.Ctr < uint16(key.Counter) {
		return nil, errors.New(REPLAYED_OTP)
	} else if token.Ctr == uint16(key.Counter) && token.Use <= uint8(key.Session) {
		return nil, errors.New(REPLAYED_OTP)
	} else {
		key.Counter = int(token.Ctr)
		key.Session = int(token.Use)
	}

	return key, nil
}

func getSecretKey(key string) (*yubikey.Key, error) {
	b, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	priv := yubikey.NewKey(b)

	return &priv, nil
}

func getToken(otpString string, priv *yubikey.Key) (*yubikey.Token, error) {
	_, otp, err := yubikey.ParseOTPString(otpString)
	if err != nil {
		return nil, err
	}

	t, err := otp.Parse(*priv)
	if err != nil {
		return nil, err
	}
	return t, nil
}
