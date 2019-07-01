package main

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

func clientRequest(
	ctx context.Context, url, username, password string) (
	result []byte) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error Get new request %s\n", err.Error())
	}

	// Costringe il client ad accettare anche certificati https non validi
	// o scaduti.
	transCfg := &http.Transport{
		// Ignora certificati SSL scaduti.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.SetBasicAuth(username, password)
	req.Header.Add("cache-control", "no-cache")
	req.WithContext(ctx)

	client := &http.Client{Transport: transCfg}

	res, err := client.Do(req)

	if err != nil {
		log.Printf(
			"Error HTTP Client Do impossibile raggiungere: %s\n",
			err.Error())
		return nil
	}

	if res.StatusCode > 300 {
		log.Printf(
			"Error Ricevuto un errore http: %d\n",
			res.StatusCode)
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf(
			"Error impossibile ricevere HTTP body per %s, %s\n",
			err.Error())
		// os.Exit(1)
	}
	defer res.Body.Close()
	// fmt.Println(res)
	// fmt.Println(string(body))
	result := new(RichiestaNefer)
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf(
			"Error unmarshal impossibile per %s, %s\n",
			intf,
			err.Error())
	}

	return result

}
