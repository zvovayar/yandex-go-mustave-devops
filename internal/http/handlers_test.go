package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

const (
	GetAllMetricsHTMLText = `<table border="1"><tr><td>type gauge</td><td> Alloc</td><td> #0 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> BuckHashSys</td><td> #1 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> FreeMemory</td><td> #29 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Frees</td><td> #2 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> GCCPUFraction</td><td> #3 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> GCSys</td><td> #4 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapAlloc</td><td> #5 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapIdle</td><td> #6 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapInuse</td><td> #7 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapObjects</td><td> #8 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapReleased</td><td> #9 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> HeapSys</td><td> #10 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> LastGC</td><td> #11 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Lookups</td><td> #12 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MCacheInuse</td><td> #13 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MCacheSys</td><td> #14 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MSpanInuse</td><td> #15 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> MSpanSys</td><td> #16 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Mallocs</td><td> #17 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> NextGC</td><td> #18 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> NumForcedGC</td><td> #19 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> NumGC</td><td> #20 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> OtherSys</td><td> #21 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> PauseTotalNs</td><td> #22 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> RandomValue</td><td> #27 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> StackInuse</td><td> #23 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> StackSys</td><td> #24 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> Sys</td><td> #25 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> TotalAlloc</td><td> #26 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> TotalMemory</td><td> #28 =</td><td> 0.000000 </td></tr><tr><td>type gauge</td><td> testSetGet134</td><td> #30 =</td><td> 0.000000 </td></tr><tr><td>type counter</td><td> PollCount</td><td> #0 =</td><td> 0</td></tr><tr><td>type counter</td><td> testSetGet33</td><td> #1 =</td><td> 0</td></tr></table>`
)

func TestMain(m *testing.M) {
	fmt.Print("TestMain run\n")

	// config test database URI
	inst.DatabaseDSN = "postgres://postgres:qweasd@localhost:5432/yandex?sslmode=disable"
	inst.StoreMonitor.OpenDB()
	// inst.StoreMonitor.LoadData()
	os.Exit(m.Run())

}

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
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	assert.Equal(t, GetAllMetricsHTMLText, rr.Body.String())
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

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

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

	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

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

	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

	// bad json
	expected = "invalid character 'b' looking for beginning of value\n"
	expectedStatus = 400

	b = bytes.NewBuffer([]byte("bad json"))

	req, err = http.NewRequest("POST", "/value", b)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetMvalueJSON)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
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

	assert.Equal(t, http.StatusNotImplemented, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

}

func TestPingStorage(t *testing.T) {

	// init httptest parameters
	handler := http.HandlerFunc(PingStorage)

	expected := `<h1>Ping database OK</h1>DSN=postgres://postgres:qweasd@localhost:5432/yandex?sslmode=disable`
	expected = "<h1>Ping database fail</h1>DSN=postgres://postgres:qweasd@localhost:5432/yandex?sslmode=disable<br>dial tcp [::1]:5432: connect: connection refused"
	expected = "<h1>Ping database fail</h1>DSN=postgres://postgres:qweasd@localhost:5432/yandex?sslmode=disable<br>dial tcp 127.0.0.1:5432: connect: connection refused"

	expectedStatus := http.StatusOK
	expectedStatus = http.StatusInternalServerError

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, req)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

}

func TestUpdateMetricJSON(t *testing.T) {

	// init httptest parameters
	// test counter
	handler := http.HandlerFunc(UpdateMetricJSON)
	v := inst.Metrics{
		ID:    "testSetGet33",
		MType: "counter",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	*v.Delta = 554

	expected := `<h1>Counter metric</h1>testSetGet33`
	expectedStatus := http.StatusOK

	buf, _ := json.Marshal(v)
	b := bytes.NewBuffer(buf)

	req, err := http.NewRequest("POST", "/update", b)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, req)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

	// test gauge
	handler = http.HandlerFunc(UpdateMetricJSON)
	v = inst.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	*v.Value = 0.987654321

	expected = `<h1>Gauge metric</h1>RandomValue`
	expectedStatus = http.StatusOK

	buf, _ = json.Marshal(v)
	b = bytes.NewBuffer(buf)

	req, err = http.NewRequest("POST", "/update", b)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, req)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

	// bad mtype
	handler = http.HandlerFunc(UpdateMetricJSON)
	v = inst.Metrics{
		ID:    "RandomValue",
		MType: "gaugex",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	*v.Value = 0.987654321

	expected = `<h1>Unknown metric type</h1>gaugex`
	expectedStatus = 400

	buf, _ = json.Marshal(v)
	b = bytes.NewBuffer(buf)

	req, err = http.NewRequest("POST", "/update", b)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, req)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

}

func TestUpdateMetricBatch(t *testing.T) {
	// init httptest parameters
	// test counter
	handler := http.HandlerFunc(UpdateMetricBatch)

	v := []inst.Metrics{
		{
			ID:    "RandomValue",
			MType: "gauge",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
		{
			ID:    "testSetGet33",
			MType: "counter",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
		{
			ID:    "testSetGet33132",
			MType: "counter",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
	}

	*v[0].Value = 5.5555

	expected := ``
	expectedStatus := http.StatusOK

	if err := inst.StoreMonitor.PingSQLserver(context.Background()); err != nil {
		expected = "dial tcp [::1]:5432: connect: connection refused\n"
		expected = "dial tcp 127.0.0.1:5432: connect: connection refused\n"
		expected = "database unreachable\n"
		expectedStatus = 200
	} else {
		expected = ``
		expectedStatus = http.StatusOK
	}

	buf, _ := json.Marshal(v)
	b := bytes.NewBuffer(buf)

	req, err := http.NewRequest("POST", "/updates", b)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, req)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetGMvalue(t *testing.T) {
	// init httptest parameters
	// test good request
	handler := http.HandlerFunc(GetGMvalue)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("GMname", "RandomValue")

	expected := `5.5555`
	expectedStatus := http.StatusOK

	r, _ := http.NewRequest("GET", "/value/gauge/{GMname}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

	// test bad request
	handler = http.HandlerFunc(GetGMvalue)

	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("GMname", "RandomValue1")

	expected = `<h1>404 Gauge metric not found</h1>`
	expectedStatus = http.StatusNotFound

	r, _ = http.NewRequest("GET", "/value/gauge/{GMname}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr = httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

}

func TestGetCMvalue(t *testing.T) {
	// test good request
	handler := http.HandlerFunc(GetCMvalue)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("CMname", "testSetGet33")

	expected := `554`
	expectedStatus := http.StatusOK

	r, _ := http.NewRequest("GET", "/value/counter/{CMname}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

	// test bad request
	handler = http.HandlerFunc(GetCMvalue)

	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("CMname", "RandomValue13")

	expected = `<h1>404 Counter metric not found</h1>`
	expectedStatus = http.StatusNotFound

	r, _ = http.NewRequest("GET", "/value/counter/{CMname}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr = httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

}

func TestUpdateGaugeMetric(t *testing.T) {
	// init httptest parameters
	// test good request
	handler := http.HandlerFunc(UpdateGaugeMetric)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("GMname", "RandomValue")
	rctx.URLParams.Add("GMvalue", "0.55555")

	expected := `<h1>Gauge metric</h1>RandomValue0.55555`
	expectedStatus := http.StatusOK

	r, _ := http.NewRequest("POST", "/update/gauge/{GMname}/{GMvalue}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

	// test bad request
	handler = http.HandlerFunc(UpdateGaugeMetric)

	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("GMname", "RandomValue13")
	rctx.URLParams.Add("GMvalue", "0.55555")

	expected = `<h1>Gauge metric</h1>RandomValue130.55555`
	expectedStatus = http.StatusOK

	r, _ = http.NewRequest("POST", "/update/gauge/{GMname}/{GMvalue}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr = httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
}

func TestUpdateCounterMetric(t *testing.T) {
	// init httptest parameters
	// test good request
	handler := http.HandlerFunc(UpdateCounterMetric)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("CMname", "testSetGet33")
	rctx.URLParams.Add("CMvalue", "55555")

	expected := `<h1>Counter metric</h1>testSetGet3355555`
	expectedStatus := http.StatusOK

	r, _ := http.NewRequest("POST", "/update/counter/{CMname}/{CMvalue}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())

	// test bad request
	handler = http.HandlerFunc(UpdateCounterMetric)

	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("CMname", "RandomValue1")
	rctx.URLParams.Add("CMvalue", "55555")

	expected = `<h1>Counter metric</h1>RandomValue155555`
	expectedStatus = http.StatusOK

	r, _ = http.NewRequest("POST", "/update/counter/{CMname}/{CMvalue}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr = httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	// results
	assert.Equal(t, expectedStatus, rr.Code)
	assert.Equal(t, expected, rr.Body.String())
}

func BenchmarkGetAllMetrics(b *testing.B) {

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMetrics)

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkUpdateMetricBatch(b *testing.B) {

	// init httptest parameters
	// test counter
	handler := http.HandlerFunc(UpdateMetricBatch)

	v := []inst.Metrics{
		{
			ID:    "RandomValue",
			MType: "gauge",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
		{
			ID:    "testSetGet33",
			MType: "counter",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
		{
			ID:    "testSetGet33",
			MType: "counter",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
	}

	*v[0].Value = 5.12345
	*v[1].Delta = 555551
	*v[2].Delta = 555551

	sbuf, _ := json.Marshal(v)
	buf := bytes.NewBuffer(sbuf)

	// run test
	for i := 0; i < b.N; i++ {

		b.StopTimer()
		req, _ := http.NewRequest("POST", "/updates", buf)
		rr := httptest.NewRecorder()
		b.StartTimer()

		handler.ServeHTTP(rr, req)
	}

}

func BenchmarkUpdateGaugeMetric(b *testing.B) {
	// init httptest parameters
	// test good request
	handler := http.HandlerFunc(UpdateGaugeMetric)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("GMname", "RandomValue")
	rctx.URLParams.Add("GMvalue", "0.55555")

	r, _ := http.NewRequest("POST", "/update/gauge/{GMname}/{GMvalue}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, r)
	}

}

func BenchmarkUpdateMetricJSON(b *testing.B) {
	handler := http.HandlerFunc(UpdateMetricJSON)
	v := inst.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	*v.Value = 0.087654321

	sbuf, _ := json.Marshal(v)
	buf := bytes.NewBuffer(sbuf)

	// run test
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		req, err := http.NewRequest("POST", "/update", buf)
		if err != nil {
			fmt.Print(err)
		}

		rr := httptest.NewRecorder()

		b.StartTimer()
		handler.ServeHTTP(rr, req)
	}
}

func ExampleUpdateMetricJSON() {
	handler := http.HandlerFunc(UpdateMetricJSON)
	v := inst.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}
	*v.Value = 0.087654321

	sbuf, _ := json.Marshal(v)
	buf := bytes.NewBuffer(sbuf)

	// run test
	req, err := http.NewRequest("POST", "/update", buf)
	if err != nil {
		fmt.Print(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// <h1>Gauge metric</h1>RandomValue
}

func ExampleGetAllMetrics() {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		fmt.Println(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMetrics)

	handler.ServeHTTP(rr, req)

	fmt.Println(rr.Code)

	// Output:
	// 200
}

func ExampleGetMvalueJSON() {

	v := inst.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Delta: new(int64),
		Value: new(float64),
		Hash:  "",
	}

	buf, _ := json.Marshal(v)
	b := bytes.NewBuffer(buf)

	req, err := http.NewRequest("POST", "/value", b)
	if err != nil {
		fmt.Println(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMvalueJSON)

	handler.ServeHTTP(rr, req)

	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// {"id":"Alloc","type":"gauge","delta":0,"value":0}
}

func ExampleUpdateMetricBatch() {
	handler := http.HandlerFunc(UpdateMetricBatch)

	v := []inst.Metrics{
		{
			ID:    "RandomValue",
			MType: "gauge",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
		{
			ID:    "testSetGet33",
			MType: "counter",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
		{
			ID:    "testSetGet33132",
			MType: "counter",
			Delta: new(int64),
			Value: new(float64),
			Hash:  "",
		},
	}

	*v[0].Value = 5.5555

	buf, _ := json.Marshal(v)
	b := bytes.NewBuffer(buf)

	req, err := http.NewRequest("POST", "/updates", b)
	if err != nil {
		fmt.Println(err)
	}

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, req)

	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
}

func ExampleGetGMvalue() {
	handler := http.HandlerFunc(GetGMvalue)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("GMname", "RandomValue")

	r, _ := http.NewRequest("GET", "/value/gauge/{GMname}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// 5.5555

}

func ExampleGetCMvalue() {
	handler := http.HandlerFunc(GetCMvalue)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("CMname", "testSetGet33")

	r, _ := http.NewRequest("GET", "/value/counter/{CMname}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)

	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// 56109
}

func ExampleUpdateGaugeMetric() {
	handler := http.HandlerFunc(UpdateGaugeMetric)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("GMname", "RandomValue")
	rctx.URLParams.Add("GMvalue", "0.55555")

	r, _ := http.NewRequest("POST", "/update/gauge/{GMname}/{GMvalue}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// <h1>Gauge metric</h1>RandomValue0.55555
}

func ExampleUpdateCounterMetric() {
	// init httptest parameters
	// test good request
	handler := http.HandlerFunc(UpdateCounterMetric)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("CMname", "testSetGet33")
	rctx.URLParams.Add("CMvalue", "55555")

	r, _ := http.NewRequest("POST", "/update/counter/{CMname}/{CMvalue}", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// run test
	handler.ServeHTTP(rr, r)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// <h1>Counter metric</h1>testSetGet3355555

}
