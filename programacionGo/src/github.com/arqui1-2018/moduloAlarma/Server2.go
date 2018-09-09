package main

import(
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
)

const(
	MongoDBHOST = "arqui1db-2018.documents.azure.com"
	AuthDatabase = "admin"
	AuthUserName = "arqui1db-2018"
	AuthPassword = "jY4INnddURHmEEJDL05qGHEYGZVQgGvf4EmomytTCqhf3wsxhuxbUPN9CJAzkJWyKvt9MLrfx1TOdxOwhNT1Xw=="
	TestDatabase = "admin"
		
)

func ErrorWithJSON(w http.ResponseWriter, message string, code int){
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "(message: %q)", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

type Objeto struct{
	ID string `json:"id"`
	Fecha string `json:"fecha"`
	Texto string `json:"texto"`
}

func main(){

	MongoDBDialInfo := &mgo.DialInfo{
		Addrs: []string{MongoDBHOST},
		//Timeout: 60* time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}
	
	session, err := mgo.DialWithInfo(MongoDBDialInfo)
	if err != nil{
		panic(err)
	}
	defer session.Close()
	
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)
	
	router := mux.NewRouter()
	router.HandleFunc("/AlarmaArduino/{id}", add(session)).Methods("GET")
	//router.HandleFunc("/AlarmaAndroid", TodoIndex).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func ensureIndex(s *mgo.Session){
	session := s.Copy()
	defer session.Close()
	
	c:= session.DB("admin").C("Alarma")
	
	index := mgo.Index{
		Key: []string{"id"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}
	err := c.EnsureIndex(index)
	if err != nil{
		panic(err)
	}
	
}

func insert(db *mgo.Session, objeto *Objeto) error{
	c := db.DB("admin").C("Alarma")
	count, err := c.Find(bson.M{"ID": objeto.ID}).Limit(1).Count()
	if err != nil{
		return err
	}
	if count > 0{
		return fmt.Errorf("El ID ya existe")
	}
	return c.Insert(objeto)
}

func add(s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	t := time.Now()
	fecha := t.String()
	texto := "Su casa se incendia"
	objeto2 := &Objeto{
		Fecha: fecha,
		Texto: texto,
	}
	err := insert(s, objeto2)	
	if err != nil{
		log.Println("No se pudo guardar en la base de datos ", err)
	}

	return func(w http.ResponseWriter, r *http.Request){
		session := s.Copy()
		defer session.Close()
		
		var objeto Objeto
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&objeto)
		if err != nil{
			ErrorWithJSON(w, "Cuerpo incorrecto", http.StatusBadRequest)
			return
		}
		
		c := session.DB("admin").C("Alarma")
		
		err = c.Insert(objeto)
		if err != nil{
			if mgo.IsDup(err){
				ErrorWithJSON(w, "Ya existe este objeto", http.StatusBadRequest)
				return
			}
			
			ErrorWithJSON(w, "Error en la Bases de Datos", http.StatusInternalServerError)
			log.Println("Fallo la insercion", err)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+objeto.ID)
		w.WriteHeader(http.StatusCreated)
	}
}
