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

var ur_id, cs_id, lyn_id, pswd, ur_nm, cs_nm, no_tlp, jenis_layanan, satuan string
var date_in, date_out time.Time
var quantity, no_nota int
var ttl_byr, tambahan, total_harga, id_trks, harga int64

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

// login user sekalian input
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
		fmt.Print("\nMasukkan jumlah pelayanan : ")
		var pel int
		fmt.Scanln(&pel)
		for i := 0; i < pel; i++ {
			fmt.Print("\n")
			fmt.Println(strings.Repeat("-", 14))
			fmt.Printf("|Pelayan ke-%d|\n", i+1)
			fmt.Println(strings.Repeat("-", 14))
			fmt.Print("\n")

			inputanTransaksi()
		}
	case 2:
	case 3:
		nota()
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
	updateTtlByr(tx, transaksi.Cs_id)
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

// trks pleyanan
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

// insert ttlharga
func insertTransaksidtl(transaksi_dtl data.Transaksi_dtl, tx *sql.Tx) {
	insertTd := "INSERT INTO transaksi_dtl(tr_id,total_harga) VALUES ($1,$2);"
	_, err := tx.Exec(insertTd, id_trks, total_harga)
	validate(err, "Insert transaksi detail", tx)
}

// update ttl byar
func updateTtlByr(tx *sql.Tx, csID string) {
	updateQuery := "UPDATE customer SET ttl_byr = (SELECT SUM(transaksi_dtl.total_harga + customer.tambahan) FROM transaksi_dtl join transaksi on transaksi_dtl.tr_id = transaksi.id join customer on transaksi.cs_id = customer.id WHERE transaksi_dtl.tr_id IN (SELECT id FROM transaksi WHERE cs_id = $1)) WHERE id = $1;"
	_, err := tx.Exec(updateQuery, csID)
	if err != nil {
		validate(err, "Update ttl_byr", tx)
	}
}

// view nota
func nota() {
	fmt.Print("Masukkan No ID pelanggan : ")
	fmt.Scan(&cs_id)
	viewNotaHead(cs_id)
	viewNotaBody(cs_id)
}
func viewNotaHead(id_cs string) {
	db := connectDB()
	defer db.Close()

	query := "SELECT customer.id, customer.date_in, customer.date_out, user_.ur_nm, customer.nm_cs, customer.no_tp, customer.tambahan, customer.ttl_byr FROM customer JOIN transaksi ON customer.id = transaksi.cs_id JOIN user_ ON transaksi.ur_id = user_.id WHERE customer.id = $1 LIMIT 1;"

	rows, err := db.Query(query, id_cs)
	if err != nil {
		fmt.Println("Id pelanggan tidak ditemukan !!")
	}
	defer rows.Close()

	fmt.Print("\nNota : \n")
	fmt.Println(strings.Repeat("=", 75))
	fmt.Print("\t\t\t     ENIGMA LAUNDRY\n")
	fmt.Println(strings.Repeat("=", 75))

	for rows.Next() {
		rows.Scan(&cs_id, &date_in, &date_out, &ur_nm, &cs_nm, &no_tlp, &tambahan, &ttl_byr)
	}
	var paket string
	if tambahan == 10000 {
		paket = "Paket Kilat"
	} else if tambahan == 5000 {
		paket = "Paket Badai"
	} else if tambahan == 0 {
		paket = "Paket Gerimis"
	}
	//kepala
	fmt.Printf("%-15s %-25s %s %s\n", "No Transaksi   : ", cs_id, "Nama customer : ", cs_nm)
	fmt.Printf("%-15s %-25s %s %s\n", "Tanggal Masuk  : ", date_in.Format("2006-01-02"), "No telepon    : ", no_tlp)
	fmt.Printf("%-15s %-25s %s %s\n", "Tanggal Selesai: ", date_out.Format("2006-01-02"), "Paket Laundry : ", paket)
	fmt.Printf("%-15s %-25s %s %d\n", "Di Terima oleh : ", ur_nm, "Fee paket     : ", tambahan)
	fmt.Println(strings.Repeat("=", 75))
}

func viewNotaBody(id_cs string) {
	db := connectDB()
	defer db.Close()

	query2 := "SELECT ROW_NUMBER() OVER (ORDER BY layanan.jns_lyn), layanan.jns_lyn, transaksi.quantity, layanan.satuan, layanan.harga, transaksi_dtl.total_harga FROM customer JOIN transaksi ON customer.id = transaksi.cs_id JOIN layanan ON layanan.id = transaksi.lyn_id JOIN transaksi_dtl ON transaksi.id = transaksi_dtl.tr_id where customer.id = $1;"

	rows, err := db.Query(query2, id_cs)
	if err != nil {
		fmt.Println("Error querying data :", err)
		return
	}
	defer rows.Close()
	//body
	fmt.Printf("%-2s %-19s %-8s %-9s %-7s %-9s %-11s", "No|", "Pelayanan", "|Jumlah", "|Satuan", "|Harga", "|Tambahan", "|Total harga\n")
	fmt.Println(strings.Repeat("=", 75))
	for rows.Next() {
		rows.Scan(&no_nota, &jenis_layanan, &quantity, &satuan, &harga, &total_harga)
		fmt.Printf("%-3d %-20s %-8d %-9s %-8d %-10d %-11d\n", no_nota, jenis_layanan, quantity, satuan, harga, tambahan, total_harga)
	}
	fmt.Print("\n\n\n")
	fmt.Print(strings.Repeat("=", 75))
	fmt.Printf("\n%+63s %d %+4s\n", "||  Toatal Pembayaran : ", ttl_byr, "||")
	fmt.Println(strings.Repeat("=", 75))

}
