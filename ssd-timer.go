package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	cookieName = "TICKER_PASSWORD"
)

var password *string
var addr *string

var timerTempl *template.Template = parseTemplate("timer", "tmpl/timer.html")
var adminTempl *template.Template = parseTemplate("admin", "tmpl/admin.html")
var loginTempl *template.Template = parseTemplate("admin", "tmpl/login.html")

// parse template using delimiter '[[', ']]'
func parseTemplate(name, fname string) *template.Template {
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		panic("template not found!")
	}
	templ, _ := template.New(name).Delims("[[", "]]").Parse(string(f))
	return templ
}

func serveTimer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	timerTempl.Execute(w, r.Host)
}

func serveAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, _ := r.Cookie(cookieName)
	if cookie == nil || cookie.Value != *password {
		http.Redirect(w, r, "login", http.StatusTemporaryRedirect)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	adminTempl.Execute(w, r.Host)
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	loginTempl.Execute(w, cookieName)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func getArgs() {
	flag.Usage = usage
	password = flag.String("p", "", "admin interface password (required)")
	addr = flag.String("a", "localhost:8080", "http service address")
	flag.Parse()

	if len(*password) == 0 {
		fmt.Fprint(os.Stderr, "Required flag: -p\n", os.Args[0])
		usage()
	}
}

func main() {
	getArgs()

	go h.run()

	http.HandleFunc("/", serveTimer)
	http.HandleFunc("/timer", serveTimerWs)

	http.HandleFunc("/admin", serveAdmin)
	http.HandleFunc("/adminws", serveAdminWs)

	http.HandleFunc("/login", serveLogin)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	log.Printf("Listening on http://%s", *addr)
	log.Printf("Admin: http://%s/admin", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
