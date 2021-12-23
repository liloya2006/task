package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"time"
)

// HomePageHandler return the home page
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	html, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		logrus.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var tmpVar = TemplateVars{
		users,
		formError,
	}

	err = html.Execute(w, tmpVar)
	formError.Name = ""
	formError.Email = ""
}

// AddUser Add users in table
func AddUser(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method != "POST" {
		fmt.Println(r)
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	if !isNameValid(r.FormValue("name")) {
		formError.Name = "Not correct input " + r.FormValue("name")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if checkDuplicateEmail(r.FormValue("email")) {
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

func logRequest(r *http.Request) {
	fileName := "./logs/logs.txt"
	//file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0755)
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	_, err = file.WriteString(fmt.Sprintf("%v: %v\n", r.Method, r.RequestURI))
	if err != nil {
		return
	}
}

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

func checkDuplicateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return true
	}

	for _, user := range users {
		if user.Email == email {
			return true
		}
	}

	return false
}

func isNameValid(e string) bool {
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9.]+([a-zA-Z0-9.]( )[a-zA-Z0-9])*[a-zA-Z0-9]+$`)
	return nameRegex.MatchString(e)
}

