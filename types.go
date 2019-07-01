package main

import "time"

// Jerk contiene le informazioni sulle cadute di sessioni PPP
// che si sono verificate.
type Jerk struct {
	NasName   string
	Timestamp time.Time
	pppValue  float64
}

// Jerks è un insieme di Jerk
type Jerks []Jerk

// Configuration tiene gli elementi di configurazione
type Configuration struct {
	IPDOMUser          string
	IPDOMPassword      string
	IPDOMUrlRicerca    string
	IPDOMSnmpReceiver  string
	IPDOMSnmpPort      uint16
	IPDOMSnmpCommunity string
	Sigma              float64
	Soglia             float64
	NasInventory       string
	NasDaIgnorare      string
	URLSessioniPPP     string
	URLTail7d          string
	SmtpPort           int
	SmtpServer         string
	SmtpUser           string
	SmtpPassword       string
	SmtpSender         string
	SmtpFrom           string
	SmtpTo             string
}

type RichiestaNefer struct {
	NetVolumeIn struct {
		Device struct {
			Intf struct {
				Data []struct {
					Time  int64 `json:"time"`
					Value int   `json:"value"`
				} `json:"data"`
			} `json:"GigabitEthernet100/0/0/26"`
		} `json:"rgamt001"`
	} `json:"net.volume.in"`
	NetVolumeOut struct {
		Device struct {
			Intf struct {
				Data []struct {
					Time  int64 `json:"time"`
					Value int   `json:"value"`
				} `json:"data"`
			} `json:"GigabitEthernet100/0/0/26"`
		} `json:"rgamt001"`
	} `json:"net.volume.out"`
}
