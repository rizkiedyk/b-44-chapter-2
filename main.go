package main

import (
	"Personal-Web/connection"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Data = map[string]interface{}{
	"IsLogin": true,
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

	// untuk melakukan/membuat routing
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	// create
	router.HandleFunc("/blog", blog).Methods("GET")
	router.HandleFunc("/add-blog", addBlog).Methods("POST")
	// delete
	router.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	// update
	router.HandleFunc("/edit-blog/{id}", editBlog).Methods("GET")
	router.HandleFunc("/update-blog/{id}", updateBlog).Methods("POST")
	// read
	router.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	// register
	router.HandleFunc("/register", formRegister).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	// login
	router.HandleFunc("/login", formLogin).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")

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

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, description, technologies, image, start_date, end_date FROM public.tb_projects;")

	var result []Blog

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["Username"] = session.Values["Name"].(string)
	}

	for rows.Next() {
		var each = Blog{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Content, &each.Tech, &each.Image, &each.StartDate, &each.EndDate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "Kardi Bongkeng"
		each.Duration = countDuration(each.StartDate, each.EndDate)

		result = append(result, each)
	}

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

	// calculate distance startdate and enddate in milisecond
	// distance := timeEndDate.Sub(timeStartDate)

	image := "MyImage.png"

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, description, technologies, image, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)", title, content, tech, image, startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// append = seperti fungsi push di js
	// append(tujuan kirim data, data yg diambil)
	// Blogs = append(Blogs, newBlog)

	// fmt.Println(Blogs)
	fmt.Println(title)
	fmt.Println(content)
	fmt.Println(tech)
	fmt.Println(startDate)
	fmt.Println(endDate)

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

	fmt.Println(id)

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

	Id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var editBlog = Blog{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_projects WHERE id=$1", Id).Scan(&editBlog.Id, &editBlog.Title, &editBlog.Content, &editBlog.Tech, &editBlog.Image, &editBlog.StartDate, &editBlog.EndDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Blog": editBlog,
	}
	fmt.Println(editBlog.Content)
	fmt.Println(editBlog)
	tmpl.Execute(w, resp)
}

func updateBlog(w http.ResponseWriter, r *http.Request) {
	// untuk terima data
	err := r.ParseMultipartForm(10485760)
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

	image := "MyImage.png"
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err = connection.Conn.Exec(context.Background(), "UPDATE public.tb_projects SET name = $1, description = $2, technologies = $3, image = $4, start_date = $5, end_date = $6 WHERE id = $7", title, content, tech, image, StartDate, EndDate, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	fmt.Println(title)
	fmt.Println(content)
	fmt.Println(tech)
	fmt.Println(StartDate)
	fmt.Println(EndDate)

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

	BlogDetail := Blog{}
	// queryrow cuma untuk menampilkan data yg kita pengen ambil
	// query menampilkan semua kolom / data
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, technologies, image, start_date, end_date FROM tb_projects WHERE id=$1", id).Scan(&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Content, &BlogDetail.Tech, &BlogDetail.Image, &BlogDetail.StartDate, &BlogDetail.EndDate)

	BlogDetail.Duration = countDuration(BlogDetail.StartDate, BlogDetail.EndDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Blog": BlogDetail,
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

	tmpl.Execute(w, nil)
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

	tmpl.Execute(w, nil)
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
	session.Options.MaxAge = 10800

	session.AddFlash("Login succes", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
