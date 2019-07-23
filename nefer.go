package main

import (
	// "github.com/axamon/nefertiti/algoritmi"
	"context"

	"github.com/gonum/matrix/mat64"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"fmt"
	"log"
	"math/rand"
	"time"

	kalman "github.com/ryskiwt/go-kalman"
)

func nefer(ctx context.Context, device, intf string) {

	// Riceve il contesto padre e aggiunge un timeout.
	// massimo per terminare la richiesta dati.
	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)

	// Esegue cancel a fine procedura.
	defer cancel()

	sstd := 0.000001
	ostd := 0.1

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

			fmt.Println("VolumeIN: ", result.NetVolumeIn.Device.Intf.Data)
			in := result.NetVolumeIn.Device.Intf.Data
			for n := range in {
				fmt.Println(in[n].Time)
				fmt.Println(in[n].Value)
			}
			time.Sleep(2 * time.Second)
			fmt.Println("VolumeOUT: ", result.NetVolumeOut.Device.Intf.Data)

			out := result.NetVolumeOut.Device.Intf.Data
			l := len(out)
			fmt.Println(l)
			for n := range out {
				fmt.Println(out[n].Time)
				fmt.Println(out[n].Value)
			}
			time.Sleep(2 * time.Second)

			// trend model
			filter, err := kalman.New(&kalman.Config{
				F: mat64.NewDense(2, 2, []float64{3, -1, 1, 0}),
				G: mat64.NewDense(2, 1, []float64{1, 0}),
				Q: mat64.NewDense(1, 1, []float64{sstd}),
				H: mat64.NewDense(1, 2, []float64{2, 0}),
				R: mat64.NewDense(1, 1, []float64{ostd}),
			})

			if err != nil {
				panic(err)
			}

			s := mat64.NewDense(1, l, nil)
			xary := make([]float64, 0, l)
			xaryfiltered := make([]float64, 0, l)

			yaryOrig := make([]float64, 0, l)

			diff := make([]float64, 0, l)

			for i := 0; i < l; i++ {
				y := float64(out[i].Value)
				s.Set(0, i, y)
				x := float64(out[i].Time)
				xfiletered := float64(out[i].Time)
				xaryfiltered = append(xary, xfiletered)
				xary = append(xary, x)

				yaryOrig = append(yaryOrig, y)
			}

			// xdet, ydet := algoritmi.Detrend(xary, yaryOrig)

			// yderived, err := algoritmi.Derive3(ydet)
			// if err != nil  {
			// 	log.Println(err.Error())
			// }

			filtered := filter.Filter(s)
			yaryFilt := mat64.Row(nil, 0, filtered)

			for n, v := range yaryFilt {
				diff = append(diff ,v - yaryOrig[n])
			}

			//
			// plot
			//

			p, err := plot.New()
			if err != nil {
				panic(err)
			}

			err = plotutil.AddLinePoints(p,
				"Original", generatePoints(xary, yaryOrig),
				"Filtered", generatePoints(xaryfiltered, yaryFilt),
				"Diff", generatePoints(xary, diff),
				//"detrend", generatePoints(xdet,yderived),
			)
			if err != nil {
				panic(err)
			}

			// Save the plot to a PNG file.
			if err := p.Save(16*vg.Inch, 4*vg.Inch, "sample.png"); err != nil {
				panic(err)
			}
		}

		return
	}
}

func generatePoints(x []float64, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))

	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	return pts
}
