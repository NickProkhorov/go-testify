package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// тест_1: Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandler_ValidRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?city=moscow&count=2", nil)
	w := httptest.NewRecorder()

	mainHandle(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, body)
	assert.Contains(t, cafeList["moscow"], strings.Split(body, ",")[0]) // просто для примера
}

// тест_2: Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandler_UnsupportedCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?city=spb&count=2", nil)
	w := httptest.NewRecorder()

	mainHandle(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "wrong city value", body)
}

// тест_3: Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?city=moscow&count=10", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	resp := responseRecorder.Result()
	body := responseRecorder.Body.String()

	expected := strings.Join(cafeList["moscow"][:totalCount], ",")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expected, body)

}
