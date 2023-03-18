package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type RandomUserResult struct {
	Gender string `json:"gender"`
	Name   struct {
		Title string `json:"title"`
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Location struct {
		Street struct {
			Number int    `json:"number"`
			Name   string `json:"name"`
		} `json:"street"`
		City        string      `json:"city"`
		State       string      `json:"state"`
		Country     string      `json:"country"`
		Postcode    interface{} `json:"postcode"`
		Coordinates struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"coordinates"`
		Timezone struct {
			Offset      string `json:"offset"`
			Description string `json:"description"`
		} `json:"timezone"`
	} `json:"location"`
	Email string `json:"email"`
	Login struct {
		Uuid     string `json:"uuid"`
		Username string `json:"username"`
		Password string `json:"password"`
		Salt     string `json:"salt"`
		Md5      string `json:"md5"`
		Sha1     string `json:"sha1"`
		Sha256   string `json:"sha256"`
	} `json:"login"`
	Dob struct {
		Date string `json:"date"`
		Age  int    `json:"age"`
	} `json:"dob"`
	Registered struct {
		Date string `json:"date"`
		Age  int    `json:"age"`
	} `json:"registered"`
	Phone string `json:"phone"`
	Cell  string `json:"cell"`
	Id    struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"id"`
	Picture struct {
		Large     string `json:"large"`
		Medium    string `json:"medium"`
		Thumbnail string `json:"thumbnail"`
	} `json:"picture"`
	Nat string `json:"nat"`
}

type RandomUserResponse struct {
	Results []RandomUserResult `json:"results"`
}

func user(w http.ResponseWriter, r *http.Request) {
	tmpl2 := template.Must(template.ParseFiles("HTML/user.html"))
	tmpl2.Execute(w, nil)
}
func contact(w http.ResponseWriter, r *http.Request) {
	tmpl3 := template.Must(template.ParseFiles("HTML/contact.html"))
	tmpl3.Execute(w, nil)
}

var newName string

func handler(w http.ResponseWriter, r *http.Request) {

	genre := r.FormValue("genre")
	fmt.Printf("Le genre est %s\n", genre)

	newName := r.FormValue("name")
	fmt.Printf("Le nouveau nom est %s\n", newName)

	link := fmt.Sprintf("https://randomuser.me/api/?gender=%s", strings.ToLower(genre))

	resp, err := http.Get(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		http.Error(w, readErr.Error(), http.StatusInternalServerError)
		return
	}

	result := RandomUserResponse{}

	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	if newName != "" {
		result.Results[0].Name.First = newName
	}

	tmpl, err := template.ParseFiles("HTML/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, result.Results[0])
	return

	http.ServeFile(w, r, "HTML/index.html")
}

func main() {
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/CSS/", http.StripPrefix("/CSS/", fs))
	http.HandleFunc("/", handler)
	http.HandleFunc("/user", user)
	http.HandleFunc("/contact", contact)
	http.ListenAndServe(":8080", nil)
}
