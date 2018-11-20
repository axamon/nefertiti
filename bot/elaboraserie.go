package main

import (
	"fmt"
	"regexp"
	"strconv"

	//"github.com/spf13/viper"

	ma "github.com/mxmCherry/movavg"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	//"compress/gzip"

	"log"
	"os"
)

func elaboraserie(lista []float64, device, interfaccia, metrica string) {

	//Inviare o no immagine del grafico a Teleram?
	var sendimage bool

	//Finita la funzione notifica il waitgroup
	defer wg.Done()

	speeds := lista

	if len(speeds) < 120 {
		log.Println("Non ci sono abbastanza dati:", device, interfaccia)
		return
	}

	// re := regexp.MustCompile("(ICR-.[0-9]+/[0-9]+)")
	// subnames := re.FindStringSubmatch(interfaccia)
	// if len(subnames) >= 0 {
	// 	nameICR := subnames[0]
	// }

	//Crea un nome per l'immagine che sia più contenuto del nomee interfaccia
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	nomeimmagine := reg.ReplaceAllString(interfaccia, "")

	//Creazione medie mobili di interesse
	sma3 := ma.ThreadSafe(ma.NewSMA(3))     //creo una moving average a 3
	sma7 := ma.ThreadSafe(ma.NewSMA(7))     //creo una moving average a 7
	sma20 := ma.ThreadSafe(ma.NewSMA(20))   //creo una moving average a 20
	sma100 := ma.ThreadSafe(ma.NewSMA(100)) //creo una moving average a 100

	n := len(speeds)
	//fmt.Println(n)

	//Crea contenitori parametrizzati al numero n di elementi in entrata
	x, dx := 0.0, 0.01
	xary := make([]float64, 0, n)
	yaryOrig := make([]float64, 0, n)
	ma3 := make([]float64, 0, n)
	ma7 := make([]float64, 0, n)
	ma20 := make([]float64, 0, n)
	ma100 := make([]float64, 0, n)
	ma20Upperband := make([]float64, 0, n)
	ma20Lowerband := make([]float64, 0, n)

	//
	// plot
	//

	//Inzializza il grafico
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//inzializza il puntatore
	var i int

	//aryOrig := speeds
	for i = 0; i < n; i++ {
		//y := math.Sin(x) + 0.1*(rand.NormFloat64()-0.5)
		y := speeds[i]
		//s.Set(0, i, y)
		x += dx

		xary = append(xary, x)

		ma3 = append(ma3, sma3.Add(y))       //aggiung alla media mobile il nuovo valore e storo la media
		ma7 = append(ma7, sma7.Add(y))       //aggiung alla media mobile il nuovo valore e storo la media
		ma20 = append(ma20, sma20.Add(y))    //aggiung alla media mobile il nuovo valore e storo la media
		ma100 = append(ma100, sma100.Add(y)) //aggiung alla media mobile il nuovo valore e storo la media
		//yaryOrig = append(yaryOrig, y-ma20[i])
		yaryOrig = append(yaryOrig, y)

		var devstdBands float64
		if i >= 100 {
			devstdBands = stat.StdDev(speeds[i-99:i], nil)
		}
		// ma20Upperband = append(ma20Upperband, sma20.Avg()+3*devstdBands)
		// ma20Lowerband = append(ma20Lowerband, sma20.Avg()-3*devstdBands)

		ma20Upperband = append(ma20Upperband, sma20.Avg()+sigma*devstdBands)
		ma20Lowerband = append(ma20Lowerband, sma20.Avg()-sigma*devstdBands)

		//Verifica anomalie
		if i > len(speeds)-3 { //Confronto solo gli ultimi3 valori per un ROPL di 15 minuti
			if yaryOrig[i] > ma20Upperband[i] {
				log.Printf("Violata soglia alta %s %s. Intf: %s, valore: %.2f", device, metrica, interfaccia, yaryOrig[i])
				//alert := fmt.Sprintf("Violata soglia alta %s %s. Intf: %s, valore: %.2f", device, metrica, interfaccia, yaryOrig[i])
				//msg <- alert
				sendimage = true
				//TODO inviare alert

			}

			if yaryOrig[i] < ma20Lowerband[i] {
				log.Printf("Violata soglia bassa %s %s. Intf: %s, valore: %.2f", device, metrica, interfaccia, yaryOrig[i])
				//alert := fmt.Sprintf("Violata soglia bassa %s %s. Intf: %s, valore: %.2f", device, metrica, interfaccia, yaryOrig[i])
				//TODO inviare alert
				//msg <- alert
				sendimage = true

			}
		}
	}

	//filtered := filter.Filter(s)
	//yaryFilt := mat64.Row(nil, 0, filtered)

	//salva sigma come string
	sigmastring := strconv.FormatFloat(sigma, 'f', 1, 64)

	err = plotutil.AddLinePoints(p,

		//"Filtered", generatePoints(xary, yaryFilt[len(yaryFilt)-120:]),
		//"MA3", generatePoints(xary, ma3),
		//"MA7", generatePoints(xary, ma7),

		"Up "+sigmastring+" sigma", generatePoints(xary, ma20Upperband[len(ma20Upperband)-120:len(ma20Upperband)-1]),
		"Original", generatePoints(xary, yaryOrig[len(yaryOrig)-120:len(yaryOrig)-1]),
		"Media mobile 20", generatePoints(xary, ma20[len(ma20)-120:]),
		"Media mobile 100", generatePoints(xary, ma100[len(ma20)-120:]),
		"Low "+sigmastring+" sigma", generatePoints(xary, ma20Lowerband[len(ma20Lowerband)-120:len(ma20Lowerband)-1]),
	)
	if err != nil {
		log.Println(err)
	}

	// Save the plot to a PNG file.

	//imposta su due righe del grafico nome apparato e interfaccia
	p.Title.Text = device + "\n " + interfaccia + "\n" + metrica

	path1 := "./grafici"
	path2 := path1 + "/" + device
	path3 := path2 + "/" + metrica

	paths := []string{path1, path2, path3}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 664)
		}
	}
	fmt.Println(path3) //debug

	//SALVA IL GRAFICO
	if err := p.Save(8*vg.Inch, 4*vg.Inch, path3+"/"+nomeimmagine+".png"); err != nil {
		panic(err)
	}

	if sendimage == true {
		image <- path3 + "/" + nomeimmagine + ".png"
	}

	return
}
