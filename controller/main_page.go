package controller

import (
	"context"
	"fmt"
	"main/database"
	"net/http"
	"text/template"
)

type frontEndData struct {
	Sutartys  []database.Sutarti
	Breziniai []database.Breziny
	Leidimai  []database.Leidima
}

func (s *Server) GetMainPage(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("permissions")
	if c == nil || err != nil {
		fmt.Println("get get main page first")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	cookie := c.Value
	dbconn := s.DBConn[cookie]
	data := frontEndData{}

	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("get main page second")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie == "1" {
		sutartys, err := queries.SelectSutartys(context.Background())
		if err != nil {
			fmt.Println("get main page third")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data.Sutartys = sutartys
	} else if cookie == "2" {
		breziniai, err := queries.SelectBreziniai(context.Background())
		if err != nil {
			fmt.Println("get main page fourth")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data.Breziniai = breziniai
	} else if cookie == "3" {
		leidimai, err := queries.SelectLeidimai(context.Background())
		if err != nil {
			fmt.Println("get main page fifth")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data.Leidimai = leidimai
	} else {
		fmt.Println("get main page sixth")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("./view/main_page/index.html", "./view/main_page/template.tmpl"))
	tmpl.Execute(w, data)
}
