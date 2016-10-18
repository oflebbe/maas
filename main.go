package main

import (
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/oflebbe/mandel"
)

// Interesting URIs
// http://localhost:8080/?pr=-1.5&pi=-1.&sz=2.&max=2048
// http://localhost:8080/?pr=-0.759081&pi=-0.071950&sz=0.001611&max=2048
// http://localhost:8080/?pr=-0.758287&pi=-0.070960&sz=0.000020&max=8000
// http://localhost:8080/?pr=-0.580820&pi=-0.652497&sz=0.000003&max=13000

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pr, err := strconv.ParseFloat(r.FormValue("pr"), 64)
		pi, err2 := strconv.ParseFloat(r.FormValue("pi"), 64)
		if err != nil || err2 != nil {
			pr = -1.3
			pi = 0.
		}
		sz, _ := strconv.ParseFloat(r.FormValue("sz"), 64)
		//res, _ := strconv.ParseInt(r.FormValue("res"), 10, 32)
		max, _ := strconv.ParseInt(r.FormValue("max"), 10, 32)
		w.Header().Set("Content-type", "image/png")
		m := mandel.NewMandel(int(max))
		im := m.Pic(complex(pr, pi), sz, 800)
		png.Encode(w, im)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
