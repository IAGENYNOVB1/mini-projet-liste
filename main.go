package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Item struct {
	ID   int64
	Name string
}

var (
	listeCourses []Item
	mutex        sync.Mutex
	nextID       int64 = 1
)

func ajouterArticle(liste []Item, article string) []Item {
	for _, a := range liste {
		if a.Name == article {
			return liste
		}
	}
	item := Item{ID: nextID, Name: article}
	nextID++
	return append(liste, item)
}

func compterArticles(liste []Item) int {
	return len(liste)
}

func supprimerArticleByName(liste []Item, article string) []Item {
	nouvelleListe := []Item{}
	for _, a := range liste {
		if a.Name != article {
			nouvelleListe = append(nouvelleListe, a)
		}
	}
	return nouvelleListe
}

func supprimerArticleByID(liste []Item, id int64) []Item {
	nouvelleListe := []Item{}
	for _, a := range liste {
		if a.ID != id {
			nouvelleListe = append(nouvelleListe, a)
		}
	}
	return nouvelleListe
}

func editArticleByID(liste []Item, id int64, newVal string) []Item {
	// replace first matching id; if newVal exists already, remove old
	exists := false
	for _, it := range liste {
		if it.Name == newVal {
			exists = true
			break
		}
	}
	for i, it := range liste {
		if it.ID == id {
			if exists {
				return supprimerArticleByID(liste, id)
			}
			liste[i].Name = newVal
			return liste
		}
	}
	return liste
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		action := r.Form.Get("action")
		switch action {
		case "delete":
			// prefer id if provided
			if idstr := r.Form.Get("delete_id"); idstr != "" {
				if id, err := strconv.ParseInt(idstr, 10, 64); err == nil {
					mutex.Lock()
					listeCourses = supprimerArticleByID(listeCourses, id)
					mutex.Unlock()
					break
				}
			}
			article := r.Form.Get("delete_article")
			mutex.Lock()
			listeCourses = supprimerArticleByName(listeCourses, article)
			mutex.Unlock()
		case "edit":
			// prefer id
			if idstr := r.Form.Get("id"); idstr != "" {
				if id, err := strconv.ParseInt(idstr, 10, 64); err == nil {
					newVal := r.Form.Get("new_article")
					mutex.Lock()
					listeCourses = editArticleByID(listeCourses, id, newVal)
					mutex.Unlock()
					break
				}
			}
			old := r.Form.Get("old_article")
			newVal := r.Form.Get("new_article")
			// fallback: edit by name
			mutex.Lock()
			for i, a := range listeCourses {
				if a.Name == old {
					// ensure not duplicating
					dup := false
					for _, b := range listeCourses {
						if b.Name == newVal {
							dup = true
							break
						}
					}
					if !dup {
						listeCourses[i].Name = newVal
					} else {
						listeCourses = supprimerArticleByName(listeCourses, old)
					}
					break
				}
			}
			mutex.Unlock()
		default:
			article := r.Form.Get("article")
			if article != "" {
				mutex.Lock()
				listeCourses = ajouterArticle(listeCourses, article)
				mutex.Unlock()
			}
		}
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	mutex.Lock()
	defer mutex.Unlock()
	tmpl.Execute(w, struct {
		Liste []Item
		Count int
	}{
		Liste: listeCourses,
		Count: compterArticles(listeCourses),
	})
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Serveur Go lancé : accès sur http://127.0.0.1:3000")
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
