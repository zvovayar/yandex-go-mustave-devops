package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	inst "internal/storage"
	"log"
	"net/http"
)

func (m inst.Monitor) SendMetrics() {

	//Just encode to json and print
	b, _ := json.Marshal(m)
	log.Println("SendMetrics -> " + string(b))
	var body = []byte(b)

	// internal.Gauge type send
	for key, element := range inst.Gmetricnames {
		var url = "http://" + inst.ServerAddress + "/update/gauge/" + key + "/" + fmt.Sprint(m.Gmetrics[element])
		log.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			log.Println(err)
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}
		// отправляем запрос

		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		log.Println(resp)
	}

	// counter type send
	for key, element := range inst.Cmetricnames {
		var url = "http://" + inst.ServerAddress + "/update/counter/" + key + "/" + fmt.Sprint(m.Cmetrics[element])
		log.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			log.Println(err)
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}
		// отправляем запрос

		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		log.Println(resp)
	}
}
