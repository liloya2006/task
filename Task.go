package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"time"
)

type Member struct {
	Id               int
	Name             string
	Email            string
	RegistrationData time.Time
}

type FormError struct {
	Name  string
	Email string
}

type TemplateVars struct {
	Users []Member
	Err   FormError
}

var users []Member
var formError FormError

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO Write to file info about The application logs incoming requests and returned responses.
func logRequest(r *http.Request) {
	fileName := "./logs/logs.txt"
	//file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0755)
	file, err := os.Create(fileName)
	if err != nil {
		check(err)
		return
	}
	_, err2 := file.WriteString(fmt.Sprintf("%v: %v\n", r.Method, r.RequestURI))
	if err2 != nil {
		check(err2)
		return
	}
}

// Open home-page
func home_page(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	html, err := template.ParseFiles("templates/home_page.html")
	check(err)
	var tmpVar = TemplateVars{
		users,
		formError,
	}

	err = html.Execute(w, tmpVar)
	formError.Name = ""
	formError.Email = ""

	check(err)
}

// AddUser Add users in table
func AddUser(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method != "POST" {
		fmt.Println(r)
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	if !IsNameValid(r.FormValue("name")) {
		formError.Name = "Not correct input " + r.FormValue("name")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if CheckDuplicateEmail(r.FormValue("email")) {
		formError.Email = "Duplicate email " + r.FormValue("email")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var member Member
	member.Id = len(users) + 1
	member.Name = r.FormValue("name")
	member.Email = r.FormValue("email")
	member.RegistrationData = time.Now()

	users = append(users, member)
	http.Redirect(w, r, "/", http.StatusFound)
}

// check duplicate email
func CheckDuplicateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	check(err)
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func IsNameValid(e string) bool {
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9.]+([a-zA-Z0-9.]( )[a-zA-Z0-9])*[a-zA-Z0-9]+$`)
	return nameRegex.MatchString(e)
}

func handleRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/AddUser", AddUser)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func main() {
	handleRequest()
}
