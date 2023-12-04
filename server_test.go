package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_GetLogin(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/login", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if strings.Contains(body, "新規登録") == false {
		t.Errorf("Body does not contain Sign in")
	}
}
