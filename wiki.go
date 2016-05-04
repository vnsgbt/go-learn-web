package main 

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9])$")

func main() {

	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))

	http.HandleFunc("/", handler)

	//http.HandleFunc("/view/", viewHandler)
	//http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.ListenAndServe(":8080",nil)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

type Page struct {
	Title 	string
	Body 	[]byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/view/"):]
	//title, err := getTitle(w, r)
	//if err != nil {
	//	return
	//}

	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}

	//fmt.Fprintf(w, "<h1> %s </h1> <div> %s </div>", p.Title, p.Body)

	//t, _ := template.ParseFiles("view.html")
	//t.Execute(w, p)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/edit/"):]
	//title, err := getTitle(w, r)
	//if err != nil {
	//	return
	//}

	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	//fmt.Fprintf(w, "<h1> Editing %s </h1>" +
	//"<form action=\"/save/%s\" method=\"POST\">"+
	//"<textarea name=\"body\">%s</textarea><br>"+
	//"<input type=\"submit\" value=\"Save\">"+
	//"</form>", p.Title, p.Title, p.Body)

	//t, _ := template.ParseFiles("edit.html")
	//t.Execute(w, p)

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/save/"):]
	//title, err := getTitle(w, r)
	//if err != nil {
	//	return
	//}

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	//p.save()
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	//t, _ := template.ParseFiles(tmpl + ".html")
	//t, err := template.ParseFiles(tmpl + ".html")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//err = t.Execute(w, p)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}

	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//func viewHandler (w http.ResponseWriter, r *http.Request, title string)
//func editHandler (w http.ResponseWriter, r *http.Request, title string)
//func saveHandler (w http.ResponseWriter, r *http.Request, title string)

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}


