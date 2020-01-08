package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)
import _ "mysql-master"

var baseURL = "http://localhost:8080"

type mahasiswa struct {
	ID          string
	Nama        string
	Tinggibadan string
}

func http_request() ([]mahasiswa, error) {
	var err error
	var client = &http.Client{}
	var data []mahasiswa

	request, err := http.NewRequest("POST", baseURL+"/mahasiswa", nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func http_request_form(ID string) (mahasiswa, error) {
	var err error
	var client = &http.Client{}
	var data mahasiswa

	var param = url.Values{}
	param.Set("id", ID)
	var payload = bytes.NewBufferString(param.Encode())

	request, err := http.NewRequest("POST", baseURL+"/detail_mahasiswa", payload)
	if err != nil {
		return data, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func baca_api() {
	var mahasiswa, err = http_request()
	if err != nil {
		fmt.Println("Error!", err.Error())
		return
	}

	fmt.Println("\n> [Baca] Data dari API...")
	fmt.Println("---------------------------------------------")

	for _, each := range mahasiswa {
		fmt.Printf("ID: %s\t Nama: %s\t Tinggi Badan: %s\n", each.ID, each.Nama, each.Tinggibadan)
	}
}

func cari_api() {
	var mahasiswa, err = http_request_form("001")
	if err != nil {
		fmt.Println("Menu tidak Tersedia", err.Error())
		return
	}
	fmt.Println("\n> [Cari] Data dari API id = 001...")
	fmt.Println("---------------------------------------------")
	fmt.Printf("ID: %s\t Nama: %s\t Tinggi Badan: %s\n", mahasiswa.ID, mahasiswa.Nama, mahasiswa.Tinggibadan)
}

func main() {
	baca_api()
	cari_api()
	sql_crud()
	fmt.Println("\nSelesai...\n")
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_mahasiswa")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func table_view() {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from tbl_mahasiswa")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	var result []mahasiswa

	for rows.Next() {
		var each = mahasiswa{}
		var err = rows.Scan(&each.ID, &each.Nama, &each.Tinggibadan)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("---------------------------------------------")

	for _, each := range result {
		fmt.Printf("ID: %s\t Nama Menu: %s\t Harga: %s\n", each.ID, each.Nama, each.Tinggibadan)
	}
}

func sql_crud() {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec("insert into tbl_mahasiswa values (?, ?, ?)", "004", "Joni", 160)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("\n> [Menambahkan] Item (Joni)")
	table_view()

	_, err = db.Exec("update tbl_mahasiswa set Tinggibadan = ? where ID = ?", 165, "004")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("\n> [Mengubah] Item (Joni, Tinggi 160 > 165)")
	table_view()

	_, err = db.Exec("delete from tbl_mahasiswa where id = ?", "004")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("\n>[Menghapus] Item (Joni)")
	table_view()
}
