package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

func TestGetAllMetrics(t *testing.T) {
	// Создаем запрос с указанием нашего хендлера. Нам не нужно
	// указывать параметры, поэтому вторым аргументом передаем nil
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMetrics)

	// Наш хендлер соответствует интерфейсу http.Handler, а значит
	// мы можем использовать ServeHTTP и напрямую указать
	// Request и ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Проверяем код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := `<table border="1"><tr><td>type gauge</td><td> Alloc</td><td> #0 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> BuckHashSys</td><td> #1 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> FreeMemory</td><td> #29 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Frees</td><td> #2 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> GCCPUFraction</td><td> #3 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> GCSys</td><td> #4 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapAlloc</td><td> #5 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapIdle</td><td> #6 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapInuse</td><td> #7 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapObjects</td><td> #8 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapReleased</td><td> #9 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapSys</td><td> #10 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> LastGC</td><td> #11 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Lookups</td><td> #12 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MCacheInuse</td><td> #13 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MCacheSys</td><td> #14 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MSpanInuse</td><td> #15 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MSpanSys</td><td> #16 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Mallocs</td><td> #17 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> NextGC</td><td> #18 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> NumForcedGC</td><td> #19 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> NumGC</td><td> #20 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> OtherSys</td><td> #21 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> PauseTotalNs</td><td> #22 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> RandomValue</td><td> #27 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> StackInuse</td><td> #23 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> StackSys</td><td> #24 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Sys</td><td> #25 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> TotalAlloc</td><td> #26 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> TotalMemory</td><td> #28 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> testSetGet134</td><td> #30 =</td><td> 0.000000 </td></tr><tr><td>type counter</td><td> PollCount</td><td> #0 =</td><td> 0</td></tr><tr><td>type counter</td><td> testSetGet33</td><td> #1 =</td><td> 0</td></tr></table>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetMvalueJSON(t *testing.T) {

	// good gauge test
	v := inst.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	expected := `{"id":"Alloc","type":"gauge","delta":0,"value":0}`

	buf, _ := json.Marshal(v)
	b := bytes.NewBuffer(buf)

	req, err := http.NewRequest("POST", "/value", b)
	if err != nil {
		t.Fatal(err)
	}

	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMvalueJSON)

	// Наш хендлер соответствует интерфейсу http.Handler, а значит
	// мы можем использовать ServeHTTP и напрямую указать
	// Request и ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Проверяем код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// good counter test
	v = inst.Metrics{
		ID:    "testSetGet33",
		MType: "counter",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	expected = `{"id":"testSetGet33","type":"counter","delta":0,"value":0}`

	buf, _ = json.Marshal(v)
	b = bytes.NewBuffer(buf)

	req, err = http.NewRequest("POST", "/value", b)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetMvalueJSON)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// bad counter test
	v = inst.Metrics{
		ID:    "testSetGet33a",
		MType: "counter",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	expected = `<h1>404 metric type or name not found</h1>`
	expectedStatus := http.StatusNotFound

	buf, _ = json.Marshal(v)
	b = bytes.NewBuffer(buf)

	req, err = http.NewRequest("POST", "/value", b)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetMvalueJSON)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// bad gauge test
	v = inst.Metrics{
		ID:    "testSetGet33a",
		MType: "gauge",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	expected = `<h1>404 metric type or name not found</h1>`
	expectedStatus = http.StatusNotFound

	buf, _ = json.Marshal(v)
	b = bytes.NewBuffer(buf)

	req, err = http.NewRequest("POST", "/value", b)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetMvalueJSON)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// bad json
	expected = `invalid character 'b' looking for beginning of value`
	expectedStatus = 400

	b = bytes.NewBuffer([]byte("bad json"))

	req, err = http.NewRequest("POST", "/value", b)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetMvalueJSON)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func TestNotImplemented(t *testing.T) {
	// Создаем запрос с указанием нашего хендлера. Нам не нужно
	// указывать параметры, поэтому вторым аргументом передаем nil
	req, err := http.NewRequest("GET", "/wrewrwe", nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := `<h1>Not implemented</h1> length=`

	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotImplemented)

	// Наш хендлер соответствует интерфейсу http.Handler, а значит
	// мы можем использовать ServeHTTP и напрямую указать
	// Request и ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Проверяем код
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func BenchmarkGetAllMetrics(b *testing.B) {

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMetrics)

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
	}
}
