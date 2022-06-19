package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Persona struct {
	IdPersona int    `json:Id_persona`
	Nombre    string `json:Nombre`
	Apellido  string `json:Apellido`
	Email     string `json:Email`
	Genero    string `json:Genero`
}

type allPersonas []Persona

var ListaPersonas = allPersonas{
	{
		IdPersona: 1,
		Nombre:    "Toinette",
		Apellido:  "Whitewood",
		Email:     "twhitewood0@homestead.com",
		Genero:    "Female",
	}, {
		IdPersona: 2,
		Nombre:    "Alick",
		Apellido:  "Antonik",
		Email:     "aantonik1@homestead.com",
		Genero:    "Male",
	}, {
		IdPersona: 3,
		Nombre:    "Theo",
		Apellido:  "Trawin",
		Email:     "ttrawin2@mail.ru",
		Genero:    "Male",
	},
}

func main() {
	//fmt.Println("hola mundo")

	//EL url debe tener una ruta valida
	router := mux.NewRouter().StrictSlash(true)

	//ruta inicial
	router.HandleFunc("/", index)

	//Responder a una ruta y por un get la lista de Personas
	router.HandleFunc("/listaPersona", getPersonas).Methods("GET")

	//Ruta para Agregar una nuevaPersona
	router.HandleFunc("/crearPersona", SavePersonas).Methods("POST")

	//Buscar una Persona por su Id
	router.HandleFunc("/buscarPersona/{IdPersona}", buscarPersonaId).Methods("GET")

	//Buscar una Persona por su Id
	router.HandleFunc("/eliminarPersona/{IdPersona}", eliminarPersonaID).Methods("DELETE")

	//Buscar una Persona por su Id
	router.HandleFunc("/actualizarPersona/{IdPersona}", updatePersonaID).Methods("PUT")

	//iniciar el Servidor en el  puerto 8080
	servidor := http.ListenAndServe(":8080", router)
	//Log
	log.Fatal(servidor)

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola Miki")
}

func getPersonas(w http.ResponseWriter, r *http.Request) {
	//Envia o Responde de Tipo Json
	w.Header().Set("Content-Type", "application/json")
	//Enviamos un Codigo de estado correcto
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ListaPersonas)
}

func SavePersonas(w http.ResponseWriter, r *http.Request) {
	var nuevaPersona Persona
	//Permite manejar las entradas y salidad del Servidor
	datos, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte Datos Validos")
	}
	//Si los datos son validos
	//se asigna a la nuevaPersona los datos recibidos
	// para poder manipular
	json.Unmarshal(datos, &nuevaPersona)

	/*
		Agregamos el id de manera manual,
		Se basa en el tamano del arreglo para asignar el idPersona al nuevo dato que se agrega
	*/
	nuevaPersona.IdPersona = len(ListaPersonas) + 1

	//Agregar una nueva Persona a la lista
	ListaPersonas = append(ListaPersonas, nuevaPersona)

	//Enviamos la nueva Persona que se agrego al arreglo en formato
	w.Header().Set("Content-Type", "application/json")
	//Enviamos un Codigo de estado correcto
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(nuevaPersona)

}
func buscarPersonaId(w http.ResponseWriter, r *http.Request) {
	//extrae los valores del url
	vars := mux.Vars(r)
	//Convierte un String a un Entero
	//Puede o no ocurrir un error
	dato, err := strconv.Atoi(vars["IdPersona"])
	if err != nil {
		fmt.Fprintf(w, "ID Invalido")
		return
	}

	// _ es el indice de la tarea i=1 i=2
	for _, Persona := range ListaPersonas {
		if Persona.IdPersona == dato {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Persona)
		}
	}
}

func eliminarPersonaID(w http.ResponseWriter, r *http.Request) {
	//extrae los valores del url
	vars := mux.Vars(r)
	//Convierte un String a un Entero
	//Puede o no ocurrir un error
	dato, err := strconv.Atoi(vars["IdPersona"])
	if err != nil {
		fmt.Fprintf(w, "ID Invalido")
		return
	}
	// guardo el valor del indice
	for i, Per := range ListaPersonas {
		if Per.IdPersona == dato {
			/*
				Concatena todo lo que esta antes del indice y todo lo que esta despues del indice
				[1,2,3,4,5] elimina el valor 3 se posiciona  en i=2 el valor 2
				Retorna [1,2,4,5]
			*/
			ListaPersonas = append(ListaPersonas[:i], ListaPersonas[i+1:]...)
			fmt.Fprintf(w, "La Persona con ID %v fue eliminada", dato)
		}
	}
}

func updatePersonaID(w http.ResponseWriter, r *http.Request) {
	//obtiene los datos de la url
	vars := mux.Vars(r)
	//
	var actualizarPersona Persona
	//Convierte un String a un Entero
	//Puede o no ocurrir un error
	dato, err := strconv.Atoi(vars["IdPersona"])
	if err != nil {
		fmt.Fprintf(w, "ID Invalido")
		return
	}

	//Obtener los Datos del r.Body de la Peticion
	nuevosDatos, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Datos Invalidos")
	}

	//asignar los datos del body a la variable actualizarPersona para manipularlso
	json.Unmarshal(nuevosDatos, &actualizarPersona)

	for i, Per := range ListaPersonas {
		if Per.IdPersona == dato {
			//eliminar la Persona que coincide con la lista y remplazar
			ListaPersonas = append(ListaPersonas[:i], ListaPersonas[i+1:]...)
			actualizarPersona.IdPersona = Per.IdPersona
			ListaPersonas = append(ListaPersonas, actualizarPersona)

			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(updatedTask)
			fmt.Fprintf(w, "La Persona %v ha sido Actualizada", dato)
		}
	}
}
