package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"akhir/data"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "26527y25tw"
	dbname   = "submit"
)

var psqlInfo = fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)

var ur_id, cs_id, lyn_id, pswd, ur_nm, cs_nm, no_tlp string
var date_in, date_out time.Time
var quantity int
var ttl_byr, tambahan, total_harga, id_trks int64

func main() {
	login()
	menu()
}

// connect db
func connectDB() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// validate
func validate(err error, message string, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		fmt.Print("\nTransaction Rollback!\n")
	} else {
		fmt.Print("\nSuccessfully " + message + " Data!\n")
	}
}

// login
func login() {
a:
	db := connectDB()
	defer db.Close()

	fmt.Println(strings.Repeat("=", 56))
	fmt.Println("\t\t\tLOGIN")
	fmt.Println(strings.Repeat("=", 56))

	fmt.Print("Masukkan ID User       : ")
	fmt.Scanln(&ur_id)
	fmt.Print("Masukkan Password User : ")
	fmt.Scanln(&pswd)
	var pswdConnect string
	err := db.QueryRow("SELECT pswd,ur_nm from user_ WHERE id = $1;", ur_id).Scan(&pswdConnect, &ur_nm)
	if err != nil {
		panic(err)
	}
	if pswd == pswdConnect {
		fmt.Printf("\nLogin berhasil !, Selamat datang %s\n\n", ur_nm)
	} else {
		fmt.Println("\nLogin Gagal!, Id atau sandi salah!!")
		fmt.Println("\nTekan 1 untuk Login ulang dan 0 untuk Exit")
		fmt.Print(strings.Repeat("=", 24))
		fmt.Print(" : ")
		var lanjut int
		fmt.Scanln(&lanjut)
		switch lanjut {
		case 1:
			goto a
		case 0:
		default:
		}
	}
}

// menu
func menu() {
	fmt.Println(strings.Repeat("-", 56))
	fmt.Println("\t |SILAHKAN PILIH OPSI YANG DIINGINKAN|")
	fmt.Println(strings.Repeat("-", 56))
	fmt.Print("\n1. pelayanan\n2. Lihat data\n3. Nota \n4. Update\n5. Delete\n6. Exit\n\n============: ")
	var menu int
	fmt.Scanln(&menu)
	switch menu {
	case 1:
		inputanCustomer()
		inputanTransaksi()
	default:

	}

}

// enrol
func enrollC(customer data.Customer) {
	db := connectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	insertCustomer(customer, tx)
	err = tx.Commit()
	if err != nil {
		panic(err)
	} else {

		fmt.Print(strings.Repeat("-", 23))
		fmt.Print("\n|Transaction Commited!|\n")
		fmt.Println(strings.Repeat("-", 23))
	}
}

func enrollT(transaksi data.Transaksi) {
	db := connectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	insertTransaksi(transaksi, tx)
	total_harga := selectData(transaksi.Id, tx)
	insertTransaksidtl(data.Transaksi_dtl{Id_trks: id_trks, Ttl_harga: total_harga}, tx)
	err = tx.Commit()
	if err != nil {
		panic(err)
	} else {

		fmt.Print(strings.Repeat("-", 23))
		fmt.Print("\n|Transaction Commited!|\n")
		fmt.Println(strings.Repeat("-", 23))
	}
}

// insert Data
func insertCustomer(customer data.Customer, tx *sql.Tx) {
	insertC := "INSERT INTO customer(id,nm_cs,no_tp,date_in,date_out,tambahan,ttl_byr) VALUES ($1,$2,$3,$4,$5,$6,$7);"
	_, err := tx.Exec(insertC, customer.Id, customer.Nama, customer.No_tlp, customer.Date_in, customer.Date_out, customer.Tambahan, customer.Total_bayar)
	validate(err, "Insert customer", tx)
}

func insertTransaksi(transaksi data.Transaksi, tx *sql.Tx) {
	insertTransaksi := "INSERT INTO transaksi(id,ur_id,cs_id,lyn_id,quantity) VALUES ($1,$2,$3,$4,$5);"
	_, err := tx.Exec(insertTransaksi, transaksi.Id, transaksi.Ur_id, transaksi.Cs_id, transaksi.Lyn_id, transaksi.Quantity)
	validate(err, "Transaksi", tx)
}

// inputan user_
func inputanCustomer() {
	fmt.Print("\nMasukkan Id customer      : ")
	fmt.Scanln(&cs_id)
	fmt.Print("Masukkan Nama customer    : ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cs_nm = scanner.Text()
	fmt.Print("Masukkan no hp customer   : ")
	fmt.Scanln(&no_tlp)
	ttl_byr = 0
	var pil int
	fmt.Print("\nMenu Paket:\n\n1. Paket Kilat    | 1 hari |+ 10 K fee|\n2. Paket Badai    | 2 hari |+ 05 K fee|\n3. Paket Gerimis  | 3 hari |+ 00 K fee|\n\n=============: ")
	fmt.Scanln(&pil)
	switch pil {
	case 1:
		date_in = time.Now()
		date_out = time.Now().Add(1 * 24 * time.Hour)
		tambahan = 10000
	case 2:
		date_in = time.Now()
		date_out = time.Now().Add(2 * 24 * time.Hour)
		tambahan = 5000
	case 3:
		date_in = time.Now()
		date_out = time.Now().Add(3 * 24 * time.Hour)
		tambahan = 0
	default:
	}

	customerInsert := data.Customer{Id: cs_id, Nama: cs_nm, No_tlp: no_tlp, Date_in: date_in, Date_out: date_out, Tambahan: tambahan, Total_bayar: ttl_byr}
	enrollC(customerInsert)
}

func inputanTransaksi() {
	db := connectDB()
	defer db.Close()

	db.QueryRow("SELECT id FROM transaksi ORDER BY id DESC LIMIT 1;").Scan(&id_trks)
	id_trks = id_trks + 1

	var pill int
	fmt.Print("\nMenu Pelayanan:\n1. Cuci + gosok     | 07k /KG |\n2. Laundry Bedcover |50k/Buah | \n3. Laundry Boneka   |25k/Buah |\n4. laundry ambal    |20k/meter|\n=============: ")
	fmt.Scanln(&pill)
	switch pill {
	case 1:
		lyn_id = "L11"
	case 2:
		lyn_id = "L12"
	case 3:
		lyn_id = "L13"
	case 4:
		lyn_id = "L14"
	default:
	}
	fmt.Print("\nmasukkan banyak barang : ")
	fmt.Scanln(&quantity)
	transaksiInsert := data.Transaksi{Id: id_trks, Ur_id: ur_id, Cs_id: cs_id, Lyn_id: lyn_id, Quantity: quantity}
	enrollT(transaksiInsert)
}

// // select
func selectData(id_trks int64, tx *sql.Tx) int64 {
	selectD := "SELECT SUM(transaksi.quantity * layanan.harga) from transaksi join layanan on transaksi.lyn_id = layanan.id WHERE transaksi.id = $1;"
	total_harga = 0
	err := tx.QueryRow(selectD, id_trks).Scan(&total_harga)
	validate(err, "Select data", tx)
	return total_harga
}

func insertTransaksidtl(transaksi_dtl data.Transaksi_dtl, tx *sql.Tx) {
	insertTd := "INSERT INTO transaksi_dtl(tr_id,total_harga) VALUES ($1,$2);"
	_, err := tx.Exec(insertTd, id_trks, total_harga)
	validate(err, "Insert transaksi detail", tx)
}
