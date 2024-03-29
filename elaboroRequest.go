package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

var sigma = configuration.Sigma

var eventi Jerks

func elaboroRequest(ctx context.Context, result []interface{}, device string) {

	serieppptime, seriepppvalue := estraiSerie(result)

	// Se non ci sono abbastanza valori per la serie esce
	if len(seriepppvalue) < 300 {
		log.Printf(
			"Error %s Non ci sono abbastanza dati per elaborare statistiche\n",
			device)
		return
	}

	mean, stdev, xdet, y := elaboraSerie(serieppptime, seriepppvalue)

	for i := 10; i < len(y); i++ {

		// Individuo se è avvenuto un Jerk
		if y[i] < mean-sigma*stdev {
			unixtimeUTC := time.Unix(int64(xdet[i]/1000), 0)
			// Serve per avere il timestamp di quando c'è stato il problema
			unixtimeinRFC3339 := unixtimeUTC.Format(time.RFC3339)

			// Devo verificare quali dopo il Jerk hanno avuto problemi
			numvalori := len(seriepppvalue)

			for l := 0; l <= 6; l++ {
				// Evita che si arrivi alla fine della serie di valori
				if i+l > numvalori-1 {
					break
				}
				// Verifica i valori dopo il jerk
				limite :=
					(seriepppvalue[i] - seriepppvalue[i+l]) / seriepppvalue[i]

				// Se il limite è negativo non ci interessa
				if limite < 0 {
					continue
				}

				if limite > configuration.Soglia {
					summary :=
						fmt.Sprintf(
							"abbassamento sessioni ppp superiore al %2.0f%%\n",
							configuration.Soglia*100)
					// Attenzione NON usare log.Print perchè serve printare il
					// timestamp non attuale ma di quando si è verificato
					// il problema
					fmt.Printf("%s Alert %s, %s\n", unixtimeinRFC3339,
						device, summary)

					evento := new(Jerk)
					evento.NasName = device
					evento.pppValue = seriepppvalue[i]
					evento.Timestamp = unixtimeUTC
					eventi = append(eventi, *evento)

					// Mandamail di notifica solo se siamo negli ultimi
					// 6 valori

				}
			}
		}
	}

}
