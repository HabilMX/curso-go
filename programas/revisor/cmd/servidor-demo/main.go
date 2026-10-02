// servidor-demo levanta cuatro rutas que se comportan como los cuatro casos
// que el revisor tiene que saber reportar: uno sano, uno lento, uno que falla
// y uno que nunca contesta. Sirve para probar el revisor de principio a fin
// sin depender de que haya servicios de verdad arriba.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	puerto := flag.String("puerto", "8090", "puerto donde escuchar")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/lento", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	mux.HandleFunc("/nunca-contesta", func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done() // no responde nunca por su cuenta; solo cede cuando el cliente se cansa
	})

	servidor := &http.Server{Addr: ":" + *puerto, Handler: mux}

	fmt.Printf("servidor-demo escuchando en :%s (/ok, /lento, /error, /nunca-contesta)\n", *puerto)
	if err := servidor.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
