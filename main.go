package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/tkanos/gonfig"
)

// IPWUrlRicerca è la parte di url dove recuperare le iformazioni
const (
 IPWUrlRicerca = "https://ipw.telecomitalia.it/ipwweb/api/0.16.0/metrics/?metrics[]=net.volume.in&metrics[]=net.volume.out&devices[]="
 IPWUrlRicercaMiddle = "&slots[]="
 IPWUrlRicercaFooter =  "&start=172800s-ago&end=1s-ago"
)
// Versione attuale di Nefertiti.
var version = "version: 1.1"

// wg è un Waitgroup che gestisce quante richieste contemporanee fare a IPDOM
var wg = sizedwaitgroup.New(5)

// Crea variabile con le configurazioni del file passato come argomento
var configuration Configuration

// Crea delle mappe a tempo per storicizzare avvenimenti
var antistorm = NewTTLMap(24 * time.Hour)
var violazioni = NewTTLMap(24 * time.Hour)
var nientedatippp = NewTTLMap(12 * time.Hour)

var username, password string

func main() {

	// Creo il contesto inziale che verrà propagato alle go-routine
	// con la funzione cancel per uscire dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//username = os.Getenv("nefertitiuser")
	// password = os.Getenv("nefertitipass")

	// Creo il canale c di buffer 1 per gestire i segnali di tipo CTRT+.
	c := make(chan os.Signal, 1)

	// Notifica al canale c i tipi di segnali di interrupt.
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c) // Ipedisce a c di ricevere ulteriori segnalazioni.
		cancel()       // Avvia la funzione di chiusura.
	}()

	// Avvia una go-routine in background in ascolto sul canale c.
	go func() {
		select {
		case <-c: // Se arriva qualche segnale in c.
			fmt.Println()
			fmt.Println("Spengo Nefertiti, docilmente...")
			log.Println("INFO Invio mail per comunicare spegnimento Nefertiti")
			// Prima di terminare la funzione invia una mail
			mandamail(
				configuration.SmtpFrom,
				configuration.SmtpTo,
				"Chiusura",
				eventi)
			cancel()
			os.Exit(0)
		case <-ctx.Done(): // Se il contesto ctx viene terminato.
		}
	}()

	// Scrive su standard output la versione di Nefertiti.
	log.Printf("Avvio Nefertiti %s\n", version)

	// file è il file di configurazione.
	file := os.Args[1] // TODO: creare flag

	// Recupera valori dal file di configurazione passato come argomento.
	err := gonfig.GetConf(file, &configuration)

	// Se ci sono errori nel recuperare le informazioni
	// l'applicazione viene chiusa.
	if err != nil {
		log.Printf(
			"Error Impossibile recupere valori da %s: %s\n",
			file,
			err.Error())
		os.Exit(1)
	}

	// GatherInfo recupera informazioni di sevizio sul funzionamento dell'APP.
	//GatherInfo(ctx)

	log.Printf("INFO Inizio recupero informazioni su IPDOM\n")

	// Dorme per 3 secondi.
	time.Sleep(3 * time.Second)

	device := "rgamt001"
	intf := "GigabitEthernet100/0/0/26"

	nefer(ctx, device, intf)

}
