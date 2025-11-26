package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

var (
	listeCourses []string
	mutex        sync.Mutex
)

func ajouterArticle(liste []string, article string) []string {
	for _, a := range liste {
		if a == article {
			return liste
		}
	}
	return append(liste, article)
}

func compterArticles(liste []string) int {
	return len(liste)
}

func supprimerArticle(liste []string, article string) []string {
	nouvelleListe := []string{}
	for _, a := range liste {
		if a != article {
			nouvelleListe = append(nouvelleListe, a)
		}
	}
	return nouvelleListe
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		if r.Form.Get("action") == "delete" {
			article := r.Form.Get("delete_article")
			mutex.Lock()
			listeCourses = supprimerArticle(listeCourses, article)
			mutex.Unlock()
		} else {
			article := r.Form.Get("article")
			mutex.Lock()
			listeCourses = ajouterArticle(listeCourses, article)
			mutex.Unlock()
		}
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	mutex.Lock()
	defer mutex.Unlock()
	tmpl.Execute(w, struct {
		Liste []string
		Count int
	}{
		Liste: listeCourses,
		Count: compterArticles(listeCourses),
	})
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Serveur Go lancé : accès sur http://<IP_DE_TA_MACHINE>:3000 depuis le réseau")
	// Écoute sur toutes les interfaces réseau, port 3000
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Erreur lors du lancement du serveur : %v", err)
	}
	// importLog importe le package log si besoin
	//
	//	func importLog() {
	//		// Hack pour forcer l'import du package log
	//		_ = log.Flags
	//	}
}
