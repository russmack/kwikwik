package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultPort  = 7171
	dataDir      = "notes/"
	templatesDir = "templates/"
	fileExt      = ".txt"
)

var (
	port = flag.Int("port", defaultPort, "Specify the listening port.")
)
var (
	templates = template.Must(template.ParseFiles(
		templatesDir+"index.html",
		templatesDir+"view.html",
		templatesDir+"edit.html",
		templatesDir+"error.html"))
	validPath = regexp.MustCompile(`^/((view|edit|save|styles|error)/([a-zA-Z0-9\.\-_]*))?$`)
)

type Model map[string]interface{}

type Page struct {
	Title string
	Body  string
}

func main() {
	flag.Parse()
	registerHandlers()
	fmt.Println("Running... ( port", *port, ")")
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

func registerHandlers() {
	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/styles/", makeHandler(styleHandler))
	http.HandleFunc("/error/", makeHandler(errorHandler))
}

func (p *Page) save() error {
	filename := p.Title + fileExt
	return ioutil.WriteFile(dataDir+filename, []byte(p.Body), 0660)
}

func (p *Page) load(title string) (*Page, error) {
	filename := title + fileExt
	filename = strings.Replace(filename, "%20", " ", -1)
	body, err := ioutil.ReadFile(dataDir + filename)
	if err != nil {
		return nil, err
	}
	p.Title = title
	p.Body = string(body)
	return p, nil
}

func buildModel(p *Page, asHtml bool) Model {
	b := p.Body
	if asHtml {
		b = strings.Replace(p.Body, "\n", "<br />", -1)
		link := regexp.MustCompile(`(^| )+([a-zA-Z0-9\-\._]+)\.txt`)
		b = link.ReplaceAllString(b, `$1<a href="view/$2">$2</a>`)
	}
	m := Model{
		"Title": p.Title,
		"Body":  template.HTML(b),
	}
	return m
}

func styleHandler(w http.ResponseWriter, r *http.Request, title string) {
	filename := "styles/" + title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("style handler err", err)
		return
	}
	w.Header().Set("content-type", "text/css")
	w.Write(body)
}

func indexHandler(w http.ResponseWriter, req *http.Request, title string) {
	title = "index"
	p, err := new(Page).load(title)
	if err != nil {
		http.Redirect(w, req, "/edit/"+title, http.StatusFound)
		return
	}
	model := buildModel(p, true)
	renderTemplate(w, "index", model)
}

func errorHandler(w http.ResponseWriter, req *http.Request, title string) {
	//title = "error"
	//p := &Page{Title: title}
	p := Page{}
	model := buildModel(p, true)
	renderTemplate(w, "error", model)
}

func viewHandler(w http.ResponseWriter, req *http.Request, title string) {
	if title == "index" {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
	p, err := new(Page).load(title)
	if err != nil {
		http.Redirect(w, req, "/edit/"+title, http.StatusFound)
		return
	}
	model := buildModel(p, true)
	renderTemplate(w, "view", model)
}

func editHandler(w http.ResponseWriter, req *http.Request, title string) {
	p, err := new(Page).load(title)
	if err != nil {
		p = &Page{Title: title}
	}
	model := buildModel(p, false)
	renderTemplate(w, "edit", model)
}

func saveHandler(w http.ResponseWriter, req *http.Request, title string) {
	body := req.FormValue("body")
	p := &Page{Title: title, Body: body}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, model Model) {
	err := templates.ExecuteTemplate(w, tmpl+".html", model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pathParts := validPath.FindStringSubmatch(req.URL.Path)
		if pathParts == nil {
			//u := &url.URL{Path: req.URL.Path}
			//encodedPath := u.String()
			http.Redirect(w, req, "/error/", http.StatusFound)
			return
		}
		fn(w, req, pathParts[3])
	}
}
