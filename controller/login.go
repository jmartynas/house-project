package controller

import (
	"context"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	user := r.PostFormValue("user")
	dbconn := NewDBConn(user, "localhost", "5432")
	cookie := http.Cookie{
		Name:    "permissions",
		Expires: time.Now().Add(time.Duration(time.Minute * 30)),
	}
	counter := 0

	q, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("post login pirmas")
		tmpl := template.Must(template.ParseFiles("./view/login/index.html"))
		tmpl.Execute(w, "")
		return
	}
	_, err = q.CheckSutartys(context.Background())
	if err == nil {
		counter++
		cookie.Value = "1"
		cookie.SameSite = 1
		http.SetCookie(w, &cookie)
		s.DBConn["1"] = NewDBConn(user, "localhost", "5432")
	}
	if counter == 0 {
		_, err = q.CheckBreziniai(context.Background())
		if err == nil {
			counter++
			cookie.Value = "2"
			cookie.SameSite = 1
			http.SetCookie(w, &cookie)
			s.DBConn["2"] = NewDBConn(user, "localhost", "5432")
		}
	}
	if counter == 0 {
		_, err = q.CheckLeidimai(context.Background())
		if err == nil {
			counter++
			cookie.Value = "3"
			cookie.SameSite = 1
			http.SetCookie(w, &cookie)
			s.DBConn["3"] = NewDBConn(user, "localhost", "5432")
		}
	}
	if counter == 0 || counter > 1 {
		fmt.Println("post login antras")
		tmpl := template.Must(template.ParseFiles("./view/login/index.html"))
		tmpl.Execute(w, "")
		return
	}
	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}

func (s *Server) GetLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./view/login/index.html"))
	tmpl.Execute(w, nil)
}

// projekto vadovas - sutartis
// architektas - brezinys
// projektuotojas - leidimai
