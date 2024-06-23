package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Definindo a struct Pessoa
type Pessoa struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// Simulando um banco de dados com um slice
var pessoas []Pessoa
var nextID int

// Função para inicializar os dados
func init() {
	pessoas = make([]Pessoa, 0)
	nextID = 1
}

// Função para buscar todas as pessoas
func getPessoas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)
}

// Função para buscar uma pessoa por ID
func getPessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, pessoa := range pessoas {
		if pessoa.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(pessoa)
			return
		}
	}

	http.Error(w, "Pessoa não encontrada", http.StatusNotFound)
}

// Função para criar uma nova pessoa
func createPessoa(w http.ResponseWriter, r *http.Request) {
	var novaPessoa Pessoa
	err := json.NewDecoder(r.Body).Decode(&novaPessoa)
	if err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	novaPessoa.ID = nextID
	nextID++

	pessoas = append(pessoas, novaPessoa)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(novaPessoa)
}

// Função para deletar uma pessoa por ID
func deletePessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for i, pessoa := range pessoas {
		if pessoa.ID == id {
			pessoas = append(pessoas[:i], pessoas[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Pessoa não encontrada", http.StatusNotFound)
}

func main() {
	// Inicializando o roteador
	r := mux.NewRouter()

	// Definindo as rotas
	r.HandleFunc("/pessoa", getPessoas).Methods("GET")
	r.HandleFunc("/pessoa/{id}", getPessoa).Methods("GET")
	r.HandleFunc("/pessoa", createPessoa).Methods("POST")
	r.HandleFunc("/pessoa/{id}", deletePessoa).Methods("DELETE")

	// Iniciando o servidor
	fmt.Println("Servidor iniciado na porta 8000")
	http.ListenAndServe(":8000", r)
}
