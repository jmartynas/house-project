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

func (s *Server) PostDeleteBreziniai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cookie.Value != "2" {
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

	vars := mux.Vars(r)
	brezinysid := vars["id"]
	bid, err := strconv.ParseInt(brezinysid, 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	leidimai, err := queries.SelectDeleteLeidimai(context.Background(), int32(bid))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, l := range leidimai {
		filePath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", l.Leidimas)
		os.Remove(filePath)
	}
	brezinys, err := queries.SelectBrezinys(context.Background(), int32(bid))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filePath := fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", brezinys.Brezinys)
	os.Remove(filePath)
	queries.DeleteLeidimai(context.Background(), int32(id))

	queries.DeleteBrezinys(context.Background(), int32(id))
	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}

type updateBreziniai struct {
	ID       int32
	FileName string
	Sutartis []Sutartys
}

type Sutartys struct {
	Sutartis database.Sutarti
	Selected bool
}

func (s *Server) GetUpdateBreziniai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("get update breziniai pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "2" {
		fmt.Println("get update sutartys pirmas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("get update breziniai antras")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		fmt.Println("get update breziniai trecias")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	brezinys, err := queries.SelectBrezinys(context.Background(), int32(id))
	if err != nil {
		fmt.Println("get update breziniai ketvirtas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sutartys, err := queries.SelectSutartys(context.Background())
	if err != nil {
		fmt.Println("get update breziniai penktas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := updateBreziniai{}
	data.ID = brezinys.ID
	data.FileName = brezinys.Brezinys
	for _, s := range sutartys {
		sut := Sutartys{
			Sutartis: s,
			Selected: false,
		}
		if s.ID == sut.Sutartis.ID {
			sut.Selected = true
		}
		data.Sutartis = append(data.Sutartis, sut)
	}

	tmpl := template.Must(template.ParseFiles("./view/breziniai/edit.html"))
	tmpl.Execute(w, data)
}

func (s *Server) PostUpdateBreziniai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("post update breziniai pirmas duxas")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "2" {
		fmt.Println("post update sutartys what")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		fmt.Println("post update breziniai antras")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		fmt.Println("post update breziniai trecias")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := queries.SelectBrezinys(context.Background(), int32(id))
	if err != nil {
		fmt.Println("post update breziniai ketvirtas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", data.Brezinys))

	// FILE UPLOADER
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("post update breziniai penktas")
		http.Redirect(w, r, "/data_list", http.StatusSeeOther)
		return
	}
	defer file.Close()
	fileName := header.Filename

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("post update breziniai sestas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	saveFile, err := os.Create(fmt.Sprintf("/home/fyzikas/git-repos/dokumentai/house-project/files/%v", fileName))
	if err != nil {
		fmt.Println("post update breziniai septintas")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	saveFile.Write(fileBytes)
	// FILE UPLOADER

	r.ParseForm()
	sutartisID := r.Form["sutartys"]
	sid, err := strconv.ParseInt(sutartisID[0], 0, 0)
	brezinys := database.UpdateBrezinysParams{
		ID:           int32(id),
		Brezinys:     fileName,
		FkSutartisID: int16(sid),
	}

	queries.UpdateBrezinys(context.Background(), brezinys)

	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}

type CreateBrezinys struct {
	Sutartys []database.Sutarti
	Error    string
}

func (s *Server) GetCreateBreziniai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "2" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbconn := s.DBConn[cookie.Value]
	queries, err := dbconn.DBConn(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sutartys, err := queries.SelectSutartys(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	brezinys := CreateBrezinys{
		Sutartys: sutartys,
		Error:    "",
	}

	tmpl := template.Must(template.ParseFiles("./view/breziniai/create.html"))
	tmpl.Execute(w, brezinys)
}

func (s *Server) PostCreateBreziniai(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("permissions")
	if err != nil {
		fmt.Println("1")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if cookie.Value != "2" {
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
		http.Redirect(w, r, "/data_list", http.StatusSeeOther)
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
	sutartis, err := strconv.ParseInt(r.Form["sutartys"][0], 0, 0)
	if err != nil {
		fmt.Println("9")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	brezinys := database.InsertBrezinysParams{
		Brezinys:     fileName,
		FkSutartisID: int16(sutartis),
	}
	queries.InsertBrezinys(context.Background(), brezinys)

	http.Redirect(w, r, "/data_list", http.StatusSeeOther)
}
