package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/hryyan/b2"
)

type Login struct {
	b2.AuthResponse
	ExpiredAt int64 `json:"expiredAt"`
}

type Session struct {
	Login `json:"login"`
}

func login() *b2.B2 {
	var (
		accountId      = viper.GetString("B2_ACCOUNT_ID")
		applicationKey = viper.GetString("B2_APPLICATION_KEY")
		b              = &b2.B2{
			AccountId:      accountId,
			ApplicationKey: applicationKey,
		}
	)

	session := readSession()
	if session == nil {
		session = &Session{}
	}

	if time.Now().Before(time.Unix(session.ExpiredAt, 0)) {
		b.SetAuth(session.AuthResponse)
	} else {
		if err := b.Auth(); err != nil {
			fmt.Println("Auth error!")
			os.Exit(AUTH_ERROR_EXIT)
		} else {
			session.AuthResponse = b.GetAuth()
			session.ExpiredAt = time.Now().Add(24 * time.Hour).Unix()
			writeSession(session)
		}
	}

	return b
}

func writeSession(session *Session) {
	f, err := os.Create(sessionPath)
	if err != nil {
		fmt.Println("Create session file error!")
		os.Exit(CREATE_SESSION_ERROR_EXIT)
	}

	b, err := json.MarshalIndent(session, "", "    ")
	if err != nil {
		fmt.Println("Encode session file error!")
		os.Exit(ENCODE_SESSION_ERROR_EXIT)
	}

	if _, err = f.Write(b); err != nil {
		fmt.Println("Write session file error!")
		os.Exit(WRITE_SESSION_ERROR_EXIT)
	}
}

func readSession() *Session {
	f, err := os.Open(sessionPath)
	if err != nil {
		fmt.Println("Read session file error!")
		os.Exit(READ_SESSION_ERROR_EXIT)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Read session file error!")
		os.Exit(READ_SESSION_ERROR_EXIT)
	}

	var session Session
	if err = json.Unmarshal(b, &session); err != nil {
		return nil
	}

	return &session
}
