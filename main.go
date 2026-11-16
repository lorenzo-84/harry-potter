package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

const apiURL = "https://potterapi-fedeperin.vercel.app/fr/characters"

// Personnage décrit un personnage de Harry Potter.
// L'API renvoie un tableau [ ... ] : on décode donc dans un []Personnage.
type Personnage struct {
	FullName      string `json:"fullName"`
	Nickname      string `json:"nickname"`
	HogwartsHouse string `json:"hogwartsHouse"`
	InterpretedBy string `json:"interpretedBy"`
	Image         string `json:"image"`
	Birthdate     string `json:"birthdate"`
}

// getPersonnages appelle l'API et range la réponse dans nos structs.
func getPersonnages() ([]Personnage, error) {
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu : %s", resp.Status)
	}

	var personnages []Personnage
	if err := json.NewDecoder(resp.Body).Decode(&personnages); err != nil {
		return nil, err
	}
	return personnages, nil
}

func main() {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	// Fichiers statiques (le CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		personnages, err := getPersonnages()
		if err != nil {
			http.Error(w, "Impossible de joindre l'API : "+err.Error(), http.StatusBadGateway)
			return
		}
		if err := tmpl.Execute(w, personnages); err != nil {
			log.Println(err)
		}
	})

	fmt.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
