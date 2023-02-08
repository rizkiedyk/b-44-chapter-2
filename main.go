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
)

var Data = map[string]interface{}{
	"Title":   "Personal Web",
	"IsLogin": false,
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

var Blogs = []Blog{
	{
		Title: "Ini title",
		// Post_date: "12 Juli 2023 | 22.30 WIB",
		Author:   "Rizki",
		Duration: "3 Bulan",
		Content:  "Ini deskripsi",
		Tech:     []string{"reactjs", "typescript", "nodejs", "nextjs"},
	},
}

func main() {
	router := mux.NewRouter()

	connection.DatabaseConnect()

	// untuk menyimpan/menampilkan file assets
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	// untuk melakukan/membuat routing
	router.HandleFunc("/hello", hello).Methods("GET")
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/blog", blog).Methods("GET")
	router.HandleFunc("/add-blog", addBlog).Methods("POST")
	router.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	router.HandleFunc("/edit-blog/{id}", editBlog).Methods("GET")
	router.HandleFunc("/update-blog/{id}", updateBlog).Methods("POST")
	router.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", router)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("aku siapa?"))
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
	for rows.Next() {
		var each = Blog{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Content, &each.Tech, &each.Image, &each.StartDate, &each.EndDate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "Kardi Bongkeng"
		each.Duration = each.Post_date.Format("Jan 21, 2000")

		result = append(result, each)
	}

	resp := map[string]interface{}{
		"Title": Data,
		"Blogs": result,
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
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	content := r.PostForm.Get("description")
	tech := r.Form["technologies"]

	// duration start
	const timeFormat = "2006-01-02"                       //declare format date
	timeStartDate, _ := time.Parse(timeFormat, startDate) //change start date timeformat
	timeEndDate, _ := time.Parse(timeFormat, endDate)     // change end date timeformat

	// calculate distance startdate and enddate in milisecond
	distance := timeEndDate.Sub(timeStartDate)

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

	var newBlog = Blog{
		Title: title,
		// Post_date: time.Now().String(),
		Author: "Rizki",
		// StartDate: startDate,
		// EndDate:   endDate,
		Duration: duration,
		Content:  content,
		Tech:     tech,
	}

	// append = seperti fungsi push di js
	// append(tujuan kirim data, data yg diambil)
	Blogs = append(Blogs, newBlog)

	fmt.Println(Blogs)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	fmt.Println(id)

	Blogs = append(Blogs[:id], Blogs[id+1:]...)

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

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var editBlog = Blog{}

	for i, data := range Blogs {
		if i == id {
			editBlog = Blog{
				Title:     data.Title,
				Post_date: data.Post_date,
				Author:    data.Author,
				Content:   data.Content,
				StartDate: data.StartDate,
				EndDate:   data.EndDate,
				Duration:  data.Duration,
				Tech:      data.Tech,
			}
		}
	}

	resp := map[string]interface{}{
		"Data": Data,
		"Id":   id,
		"Blog": editBlog,
	}
	fmt.Println(editBlog.Content)
	fmt.Println(editBlog)
	tmpl.Execute(w, resp)
}

func updateBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// untuk terima data
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	// sama seperti document.getElementById tapi pakai name bukan id
	title := r.PostForm.Get("title")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	content := r.PostForm.Get("description")
	tech := r.Form["technologies"]

	// duration start
	const timeFormat = "2006-01-02"                       //declare format date
	timeStartDate, _ := time.Parse(timeFormat, startDate) //change start date timeformat
	timeEndDate, _ := time.Parse(timeFormat, endDate)     // change end date timeformat

	// calculate distance startdate and enddate in milisecond
	distance := timeEndDate.Sub(timeStartDate)

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

	var newBlog = Blog{
		Title: title,
		// Post_date: time.Now().String(),
		Author: "Rizki",
		// StartDate: startDate,
		// EndDate:   endDate,
		Duration: duration,
		Content:  content,
		Tech:     tech,
	}

	Blogs[id] = newBlog
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

	for i, data := range Blogs {
		if i == id {
			BlogDetail = Blog{
				Title:     data.Title,
				Post_date: data.Post_date,
				Author:    data.Author,
				Content:   data.Content,
				StartDate: data.StartDate,
				EndDate:   data.EndDate,
				Duration:  data.Duration,
				Tech:      data.Tech,
			}
		}
	}

	resp := map[string]interface{}{
		"Data": Data,
		"Id":   id,
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
