package main

// package yang diperlukan untuk API dan DB mysql
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)
import _ "mysql-master"

// deklarari tipe struct
type mahasiswa struct {
	ID          string
	Nama        string
	Tinggibadan string
}

// variabel penampung
var data = []mahasiswa{}

// koneksi ke mysql dan database
func koneksi() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_mahasiswa")

	if err != nil {
		return nil, err
	}
	return db, nil
}

// memasukan tabel kedalam variabel penampung
func get_mahasiswa() {
	db, err := koneksi()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	row, err := db.Query("select * from tbl_mahasiswa")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	for row.Next() {
		var each = mahasiswa{}
		var err = row.Scan(&each.ID, &each.Nama, &each.Tinggibadan)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data = append(data, each)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

// func API POST untuk mendapatkan isi tabel mahasiswa
func tampil_mahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	if r.Method == "POST" {
		var result, err = json.Marshal(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

// func API POST untuk mendapatkan detail isi tabel mahasiswa
func detail_mahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var id = r.FormValue("id")
		var result []byte
		var err error

		for _, each := range data {
			if each.ID == id {
				result, err = json.Marshal(each)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Write(result)
				return
			}
		}

		http.Error(w, "Detail Mahasiswa, Gagal Ditampilkan", http.StatusBadRequest)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

// func main sekaligus rooting dari API dan server
func main() {
	get_mahasiswa()
	http.HandleFunc("/mahasiswa", tampil_mahasiswa)
	http.HandleFunc("/detail_mahasiswa", detail_mahasiswa)

	fmt.Println("Berhasil: Menjalankan Web Server Pada http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
