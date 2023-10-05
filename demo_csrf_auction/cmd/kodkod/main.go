package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	auction "github.com/connorkuljis/kodkod/internal/auctions"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	_ "github.com/mattn/go-sqlite3"
)

// globals
var db *sql.DB

const (
	DatabaseName = "my_database.db"
	CookieName   = "session"
)

// Provides cookie and filesystem sessions
var store = sessions.NewCookieStore([]byte("verySecretValue"))

// Data object for templates to access information about a session.
type PageData struct {
	Username       string
	IsLoggedIn     bool
	AuctionData    []auction.Auction
	CurrentAuction auction.Auction
	CsrfField      template.HTML
	Message        string
}

// Entry point to our application
func main() {
	fmt.Println("💿 serving on http://localhost:8080/")

	instantiateDBConnection()

	router := mux.NewRouter()

	// https://stackoverflow.com/questions/27945310/why-do-i-need-to-use-http-stripprefix-to-access-my-static-files
	s := http.StripPrefix("/web/static/", setupStaticFileServerHandler())
	router.PathPrefix("/web/static/").Handler(s)

	// Map actions to handlers.
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/users/login", loginHandler)
	router.HandleFunc("/users/logout", logoutHandler)
	router.HandleFunc("/users/register", registerHandler)

	router.HandleFunc("/auction/{auctionID}", auctionHandler)

	router.HandleFunc("/bid", bidHandler)

	router.HandleFunc("/xsrf", xsrfHandler)

	var csrfMiddleware = csrf.Protect(
		[]byte("your-secret-key"), // Replace with your own secret key
		csrf.Secure(false),        // Set to true in production to enforce HTTPS
	)
	http.ListenAndServe(":8080", csrfMiddleware(router))

	db.Close()
}

func instantiateDBConnection() {
	var err error
	db, err = sql.Open("sqlite3", DatabaseName)
	if err != nil {
		log.Fatal("Error opening database.")
	}
}

// A handler that serves static assets from the file system.
func setupStaticFileServerHandler() http.Handler {
	resourcePath := "./web/static/"
	return http.FileServer(http.Dir(resourcePath))
}

// Stores session information in data object.
func setPageData(session *sessions.Session) PageData {
	pageData := PageData{}
	if session.Values["logged_in"] != nil {
		pageData.IsLoggedIn = true
	}
	if session.Values["username"] != nil {
		pageData.Username = session.Values["username"].(string)
	}

	pageData.AuctionData = auction.LoadMockAuctionData(db)

	return pageData
}

// Prints http request and cookie data to stdout.
func logHttpRequest(r *http.Request) {
	log.Println(r.Method, r.RequestURI)

	cookie, err := r.Cookie(CookieName)
	if err != nil {
		fmt.Println("Error, could not read cookie: " + CookieName + " " + err.Error())
		return
	}
	fmt.Printf("\033[33m%s\033[0m\n", cookie)
}

// --- HTTP Handler Callback Functions ---
// Each handler generally needs to:
// 1. Log the request information.
// 2. Create or access the existing session store using a cookie.
// 3. Render a html document by composing html templates.

// indexHandler renders the index.html page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	logHttpRequest(r)
	session, _ := store.Get(r, CookieName) // will create a new cookie if not exists
	data := setPageData(session)

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// loginHander renders the login.html page.
// Form data contains login information. The session state is updated when
// authenticated and redirected to the default handler.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, CookieName)
	data := PageData{}
	data.CsrfField = csrf.TemplateField(r)
	if r.Method == "GET" {
		logHttpRequest(r)
		tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	} else {
		r.ParseForm()
		// TODO: validate passwords with users table
		fmt.Println("username:", r.Form["username"][0])
		fmt.Println("password:", r.Form["password"][0])
		// Set some session values.
		session.Values["logged_in"] = true
		session.Values["username"] = r.Form["username"][0]
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// registerHandler renders the register.html page.
func registerHandler(w http.ResponseWriter, r *http.Request) {
	logHttpRequest(r)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/register.html")
		if err != nil {

		}
		tmpl.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

// logoutHandler expires the session and redirects to
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	logHttpRequest(r)
	session, _ := store.Get(r, CookieName)
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func auctionHandler(w http.ResponseWriter, r *http.Request) {
	logHttpRequest(r)
	session, _ := store.Get(r, CookieName)
	data := setPageData(session)
	data.CsrfField = csrf.TemplateField(r)

	vars := mux.Vars(r)
	s := vars["auctionID"]
	auctionID, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get from db
	var a auction.Auction
	a, err = auction.GetAuctionById(db, auctionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.CurrentAuction = a

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/auction.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func bidHandler(w http.ResponseWriter, r *http.Request) {
	logHttpRequest(r)
	session, _ := store.Get(r, CookieName)
	data := setPageData(session)

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check the HTTP request method, we only want to parse the body for POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST methods are allowed.", http.StatusBadRequest)
		return
	}

	// Parse the form data from the request body
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Access the form values by their keys
	auctionIDStr := r.FormValue("auction")
	amountStr := r.FormValue("amount")

	if auctionIDStr == "" || amountStr == "" {
		http.Error(w, "Error parsing form data:", http.StatusMethodNotAllowed)
		return
	}

	var amount float64
	amount, err = strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	res, err := db.Exec(`UPDATE auctions SET price = ? WHERE id = ?`, amount, auctionIDStr)
	if err != nil {
		http.Error(w, "Error updating auction price:", http.StatusMethodNotAllowed)
		return
	}

	log.Println(res)

	str := fmt.Sprintf("'%s' placed a bid on item %s for amount: $%.2f\n", data.Username, auctionIDStr, amount)
	fmt.Println(str)
	data.Message = str

	tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/bid.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func xsrfHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/xsrf.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
