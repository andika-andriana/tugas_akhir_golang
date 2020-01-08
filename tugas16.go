package main

import "fmt"
import "database/sql"
import _ "mysql-master"

type student struct {
	ID          string
	Nama        string
	Tinggibadan int
}

func koneksi() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_mahasiswa")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func sqlQuery() {
	db, err := koneksi()
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

	var result []student

	for rows.Next() {
		var each = student{}
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

	fmt.Println("\n\t\t - - - Data Mahasiswa - - - \n")

	for _, each := range result {
		fmt.Printf("ID : %s\t | Nama Mahasiswa : %s\t | Tinggi Badan : %d\n", each.ID, each.Nama, each.Tinggibadan)
	}
	fmt.Println("\n\t\t - - - - - - - - - - - - - \n")
}

func main() {
	sqlQuery()
}
