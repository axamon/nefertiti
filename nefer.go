package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func nefer(ctx context.Context, device, intf string) {

	// Riceve il contesto padre e aggiunge un timeout.
	// massimo per terminare la richiesta dati.
	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)

	// Esegue cancel a fine procedura.
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Printf(
				"Error %s Superato tempo massimo per raccolta dati\n")
			return

		default:
			// Attendo un tempo random per evitare troppe query insieme.
			// se sono attive le goroutines.
			randomdelay := rand.Intn(100)
			time.Sleep(
				time.Duration(randomdelay) * time.Millisecond)

			// Ripulisce eventiali impostazioni di proxy
			// a livello di sistema.
			// os.Setenv("HTTP_PROXY", "")
			// os.Setenv("HTTPS_PROXY", "")
			// fmt.Println(os.Getenv("HTTP_PROXY"))
			// fmt.Println(os.Getenv("HTTPS_PROXY"))

			//fmt.Println(device)

			// Recupera le credenziali per IPDOM.
			username := configuration.IPDOMUser
			password := configuration.IPDOMPassword

			// Ricompongo la URL di IPDOM con il nome del NAS all'interno.
			url := IPWUrlRicerca + device + IPWUrlRicercaMiddle + intf + IPWUrlRicercaFooter

			// Avvia la richiesta web.
			result := clientRequest(ctx, url, username, password)

			// Elabora il risulatato della richiesta web.
			//elaboroRequest(ctx, result, device)

			fmt.Println(string(result))
			return
		}
	}
}
