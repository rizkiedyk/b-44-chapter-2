package main

import (
	"Personal-Web/connection"
	"Personal-Web/middleware"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	Username  string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

// untuk mendifinisikan isi dan type data objek yg akan dibuat
type Blog struct {
	Id        int
	Image     string
	Title     string
	Post_date time.Time
	Author    string
	StartDate time.Time
	EndDate   time.Time
	Duration  string
	Content   string
	Tech      []string
}

func main() {
	router := mux.NewRouter()

	connection.DatabaseConnect()

	// untuk menyimpan/menampilkan file assets
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads"))))

	// untuk melakukan/membuat routing
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	// create
	router.HandleFunc("/blog", blog).Methods("GET")
	router.HandleFunc("/add-blog", middleware.UploadFile(addBlog)).Methods("POST")
	// delete
	router.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	// update
	router.HandleFunc("/edit-blog/{id}", editBlog).Methods("GET")
	router.HandleFunc("/update-blog/{id}", middleware.UploadFile(updateBlog)).Methods("POST")
	// read
	router.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	// register
	router.HandleFunc("/register", formRegister).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	// login
	router.HandleFunc("/login", formLogin).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	// log out
	router.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", router)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var result []Blog

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false

		rows, errQuery := connection.Conn.Query(context.Background(), "SELECT id, title,description, technologies, image , start_date, end_date FROM tb_projects")
		if errQuery != nil {
			fmt.Println("message : " + errQuery.Error())
			return
		}

		for rows.Next() {
			each := Blog{}

			err := rows.Scan(&each.Id, &each.Title, &each.Content, &each.Tech, &each.Image, &each.StartDate, &each.EndDate)
			if err != nil {
				fmt.Println("message : " + err.Error())
				return
			}

			each.Duration = countDuration(each.StartDate, each.EndDate)

			result = append(result, each)
		}
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)

		user := session.Values["Id"]

		rows, errQuery := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, title, description, technologies, image, start_date, end_date, tb_user.name as author FROM tb_projects LEFT JOIN tb_user ON tb_projects.author_id = tb_user.id WHERE tb_projects.author_id = $1", user)
		if errQuery != nil {
			fmt.Println("Message : " + errQuery.Error())
			return
		}

		for rows.Next() {
			var each = Blog{}

			var err = rows.Scan(&each.Id, &each.Title, &each.Content, &each.Tech, &each.Image, &each.StartDate, &each.EndDate, &each.Author)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			each.Duration = countDuration(each.StartDate, each.EndDate)

			result = append(result, each)
		}

	}

	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	resp := map[string]interface{}{
		"Data": Data,
		"Blog": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func blog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	// untuk terima data
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	// sama seperti document.getElementById tapi pakai name bukan id
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("description")
	tech := r.Form["technologies"]

	// duration start
	const timeFormat = "2006-01-02" //declare format date
	startDate, _ := time.Parse(timeFormat, r.PostForm.Get("start-date"))
	endDate, _ := time.Parse(timeFormat, r.PostForm.Get("end-date"))
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	author := session.Values["Id"].(int)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(title, description, technologies, image, start_date, end_date, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", title, content, tech, image, startDate, endDate, author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// make duration
func countDuration(start time.Time, end time.Time) string {
	distance := end.Sub(start)

	// Convert time to month, week and day
	monthDistance := int(distance.Hours() / 24 / 30)
	dayDistance := int(distance.Hours() / 24)

	var duration string

	if monthDistance >= 1 && dayDistance%30 >= 1 {
		duration = strconv.Itoa(monthDistance) + " Months " + strconv.Itoa(dayDistance%30) + " Days"
	} else if monthDistance >= 1 {
		duration = strconv.Itoa(monthDistance) + " Months"
	} else if monthDistance < 1 && dayDistance >= 0 {
		duration = strconv.Itoa(dayDistance) + " Days"
	} else {
		duration = "0 Days"
	}

	return duration
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text-html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/edit-blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	Id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var editBlog = Blog{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_projects WHERE id=$1", Id).Scan(&editBlog.Id, &editBlog.Title, &editBlog.Content, &editBlog.Tech, &editBlog.Image, &editBlog.StartDate, &editBlog.EndDate, &editBlog.Author)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Blog": editBlog,
		"Data": Data,
	}

	tmpl.Execute(w, resp)
}

func updateBlog(w http.ResponseWriter, r *http.Request) {
	// untuk terima data
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	// sama seperti document.getElementById tapi pakai name bukan id
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("description")
	tech := r.Form["technologies"]

	const timeFormat = "2006-01-02" //declare format date
	StartDate, _ := time.Parse(timeFormat, r.PostForm.Get("start-date"))
	EndDate, _ := time.Parse(timeFormat, r.PostForm.Get("end-date"))

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)
	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	author := session.Values["Id"].(int)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err = connection.Conn.Exec(context.Background(), "UPDATE public.tb_projects SET title = $1, description = $2, technologies = $3, image = $4, start_date = $5, end_date = $6, author_id = $7 WHERE id = $8", title, content, tech, image, StartDate, EndDate, author, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	BlogDetail := Blog{}
	// queryrow cuma untuk menampilkan data yg kita pengen ambil
	// query menampilkan semua kolom / data
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_projects WHERE id=$1", id).Scan(&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Content, &BlogDetail.Tech, &BlogDetail.Image, &BlogDetail.StartDate, &BlogDetail.EndDate, &BlogDetail.Author)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	BlogDetail.Duration = countDuration(BlogDetail.StartDate, BlogDetail.EndDate)

	dataTime := map[string]interface{}{
		"dStart": BlogDetail.StartDate.Format("2006-01-02"),
		"dEnd":   BlogDetail.EndDate.Format("2006-01-02"),
	}

	resp := map[string]interface{}{
		"Blog": BlogDetail,
		"Data": Data,
		"Time": dataTime,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data": Data,
	}

	tmpl.Execute(w, resp)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	// encrypt password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_user(name, email, password) VALUES ($1, $2, $3);", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.AddFlash("Register success!", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	resp := map[string]interface{}{
		"Data": Data,
	}

	tmpl.Execute(w, resp)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800

	session.AddFlash("Login succes", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
