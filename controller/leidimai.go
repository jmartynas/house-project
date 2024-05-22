package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"main/database"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func (s *Server) PostDeleteLeidimai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "3" {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	leidimas, err := queries.SelectLeidimas(context.Background(), int64(id))
	filePath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", leidimas.Leidimas)
	os.Remove(filePath)
	queries.DeleteLeidimas(context.Background(), int64(id))
	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}

type updateLeidimai struct {
	ID       int64
	FileName string
	Brezinys []Breziniai
}

type Breziniai struct {
	Brezinys database.Breziny
	Selected bool
}

func (s *Server) GetUpdateLeidimai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("get update leidimai pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "3" {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("get update leidimai antras")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		fmt.Println("get update leidimai trecias")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	leidimas, err := queries.SelectLeidimas(context.Background(), int64(id))
	if err != nil {
		fmt.Println("get update leidimai ketvirtas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	breziniai, err := queries.SelectBreziniai(context.Background())
	if err != nil {
		fmt.Println("get update leidimai penktas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := updateLeidimai{}
	data.ID = leidimas.ID
	data.FileName = leidimas.Leidimas
	for _, b := range breziniai {
		brez := Breziniai{
			Brezinys: b,
		}
		if b.ID == leidimas.FkBrezinysID {
			brez.Selected = true
		} else {
			brez.Selected = false
		}
		data.Brezinys = append(data.Brezinys, brez)
	}

	tmpl := template.Must(template.ParseFiles("./view/leidimai/edit.html"))
	tmpl.Execute(w, data)
}

func (s *Server) PostUpdateLeidimai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("post update leidimai pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "3" {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("post update leidimai antras")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		fmt.Println("post update leidimai trecias")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := queries.SelectLeidimas(context.Background(), int64(id))
	if err != nil {
		fmt.Println("post update leidimai ketvirtas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", data.Leidimas))

	// FILE UPLOADER
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("post update leidimai penktas")
		http.Redirect(w, r, "/data_list", http.StatusSeeOther)
		return
	}
	defer file.Close()
	fileName := header.Filename

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("post update leidimai sestas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	saveFile, err := os.Create(fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", fileName))
	if err != nil {
		fmt.Println("post update leidimai septintas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	saveFile.Write(fileBytes)
	// FILE UPLOADER

	r.ParseForm()
	brezinysID := r.Form["breziniai"]
	bid, err := strconv.ParseInt(brezinysID[0], 0, 0)
	leidimas := database.UpdateLeidimasParams{
		ID:           id,
		Leidimas:     fileName,
		FkBrezinysID: int32(bid),
	}

	queries.UpdateLeidimas(context.Background(), leidimas)

	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}

type CreateLeidimas struct {
	Breziniai []database.Breziny
	Error     string
}

func (s *Server) GetCreateLeidimai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("1")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "3" {
		fmt.Println("2")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("3")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	breziniai, err := queries.SelectBreziniai(context.Background())
	if err != nil {
		fmt.Println("4")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	leidimas := CreateLeidimas{
		Breziniai: breziniai,
		Error:     "",
	}

	tmpl := template.Must(template.ParseFiles("./view/leidimai/create.html"))
	tmpl.Execute(w, leidimas)
}

func (s *Server) PostCreateLeidimai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("1")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "3" {
		fmt.Println("2")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("3")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("4")
		http.Redirect(w, r, "/data_list", http.StatusSeeOther)
		return
	}
	defer file.Close()
	fileName := header.Filename

	filePath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", fileName)
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		fmt.Println("5")
		// http.Redirect(w, r, "/data_list", http.StatusSeeOther)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("6")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	saveFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("8")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	saveFile.Write(fileBytes)

	r.ParseForm()
	brezinys, err := strconv.ParseInt(r.Form["breziniai"][0], 0, 0)
	if err != nil {
		fmt.Println("9")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	leidimas := database.InsertLeidimasParams{
		Leidimas:     fileName,
		FkBrezinysID: int32(brezinys),
	}
	queries.InsertLeidimas(context.Background(), leidimas)

	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}
