package main

import(
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var alarma string

type Objeto struct{
	Fecha string `bson: "fecha" json: "fecha"`
	Texto string `bson: "texto" json: "texto"`
}


func newObjeto(fecha, texto string) Objeto {
	return Objeto{
		Fecha : fecha,
		Texto : texto,
	}
}

type Payload struct{
	Mensaje string `json:message`
}

var isDropMe = true

func main() {
	
	router := mux.NewRouter()	
	router.HandleFunc("/AlarmaArduino/{id}", Index).Methods("GET")
	router.HandleFunc("/AlarmaAndroid", TodoIndex).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

//Funcion ARDUINO 
func Index(w http.ResponseWriter, r *http.Request){

	params := mux.Vars(r)
	num, err := strconv.Atoi(params["id"])
	
	
	fmt.Println(num)
	fmt.Println(err)

	if num > 350 && err == nil {
		client := http.Client{}
		request, err := http.NewRequest("POST", "https://graph.facebook.com/v3.1/2040235742664167/feed?message=SU%20CASA%20SE%20INCENDIA&access_token=EAAew0cc1Jz4BACINeKF15TxPblVHpb7wYCsLma4F80cZC6LHmmitZC9zp8etzi88reKx88N7cN1RosxyJ6LdhZBvAulkUffZAeLWoZAmAEvz2nQ2NL78nemlxUPzcpvqpn1XD3w03YPa5ZAuZAoMeCViUWRBiAY5BsLQRStCRu0fGmXIdIPjwqh74UPupJBLUcZBKEZBtFZCyxXgZDZD", nil)
		
		if err != nil{
			log.Fatalln(err)
		}
		
		resp, err := client.Do(request)
		if err != nil{
			log.Fatalln(err)
			fmt.Println(resp)
		}
		
		fmt.Fprintf(w, "Hay un incendio!!!")
		alarma = "true"
			
	}else{		
		fmt.Fprintf(w, "No hay incendio!!!")
		alarma = "false"
	}
	
}

//
func TodoIndex(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, alarma)
}
