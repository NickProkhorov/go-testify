package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

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
