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

func (s *Server) PostDeleteSutartys(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cookie.Value != "1" {
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

	breziniai, err := queries.SelectDeleteBreziniai(context.Background(), int16(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, b := range breziniai {
		leidimai, err := queries.SelectDeleteLeidimai(context.Background(), int32(b.ID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, l := range leidimai {
			leidimasPath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", l.Leidimas)
			os.Remove(leidimasPath)
		}
		queries.DeleteLeidimai(context.Background(), int32(b.ID))
		brezinysPath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", b.Brezinys)
		os.Remove(brezinysPath)
		queries.DeleteBrezinys(context.Background(), int32(b.ID))
	}

	sutartis, err := queries.SelectSutartis(context.Background(), int16(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sutartisPath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", sutartis.Sutartis)
	os.Remove(sutartisPath)
	queries.DeleteSutartis(context.Background(), int16(id))

	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}

type updateSutartys struct {
	ID       int16
	FileName string
	Kaina    float64
}

func (s *Server) GetUpdateSutartys(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("get update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "1" {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("get update sutartys antras")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		fmt.Println("get update sutartys trecias")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sutartis, err := queries.SelectSutartis(context.Background(), int16(id))
	if err != nil {
		fmt.Println("get update sutartys ketvirtas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	kaina, err := strconv.ParseFloat(sutartis.Kaina, 64)
	if err != nil {
		fmt.Println("get update sutartys penktas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := updateSutartys{
		ID:       sutartis.ID,
		FileName: sutartis.Sutartis,
		Kaina:    kaina,
	}

	tmpl := template.Must(template.ParseFiles("./view/sutartys/edit.html"))
	tmpl.Execute(w, data)
}

func (s *Server) PostUpdateSutartys(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "1" {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("post update sutartys antras")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		fmt.Println("post update sutartys trecias")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := queries.SelectSutartis(context.Background(), int16(id))
	if err != nil {
		fmt.Println("post update sutartys ketvirtas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", data.Sutartis))

	// FILE UPLOADER
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("post update sutartys penktas")
		http.Redirect(w, r, "/data_list", http.StatusSeeOther)
		return
	}
	defer file.Close()
	fileName := header.Filename

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("post update sutartys sestas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	saveFile, err := os.Create(fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", fileName))
	if err != nil {
		fmt.Println("post update sutartys septintas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	saveFile.Write(fileBytes)
	// FILE UPLOADER

	r.ParseForm()
	kaina := r.Form["kaina"]
	sutartis := database.UpdateSutartisParams{
		ID:       int16(id),
		Sutartis: fileName,
		Kaina:    kaina[0],
	}

	queries.UpdateSutartis(context.Background(), sutartis)

	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}

func (s *Server) GetCreateSutartys(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "1" {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("./view/sutartys/create.html"))
	tmpl.Execute(w, nil)
}

func (s *Server) PostCreateSutartys(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "1" {
		fmt.Println("post update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("post update sutartys antras")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("get create sutartys trecias")
		http.Redirect(w, r, "/data_list", http.StatusSeeOther)
		return
	}
	defer file.Close()
	fileName := header.Filename

	filePath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", fileName)
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		fmt.Println("post create sutartys ketvirtas")
		http.Redirect(w, r, "/data_list", http.StatusSeeOther)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("post create sutartys penktas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	saveFile, err := os.Create(fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", fileName))
	if err != nil {
		fmt.Println("post create sutartys sestas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	saveFile.Write(fileBytes)

	r.ParseForm()
	kaina := r.Form["kaina"][0]
	sutartis := database.InsertSutartisParams{
		Sutartis: fileName,
		Kaina:    kaina,
	}
	queries.InsertSutartis(context.Background(), sutartis)

	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}
