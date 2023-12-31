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
var ttl_byr, tambahan, total_harga, id_trks, harga, id_trks_dtl int64

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
		fmt.Print("\nLanjut Login [Y/T]: ")
		var lanjut string
		fmt.Scanln(&lanjut)
		if lanjut == "Y" || lanjut == "y" {
			goto a
		}
	}
}

// menu
func menu() {
b:
	fmt.Println(strings.Repeat("-", 56))
	fmt.Println("\t |SILAHKAN PILIH OPSI YANG DIINGINKAN|")
	fmt.Println(strings.Repeat("-", 56))
	fmt.Print("\n1. pelayanan\n2. Lihat data\n3. Update\n4. Delete\n5. Exit\n\n============: ")
	var menu int
	fmt.Scanln(&menu)
	switch menu {
	case 1:
	a:
		fmt.Println()
		fmt.Println(strings.Repeat("-", 16))
		fmt.Println("|MENU PELAYANAN|")
		fmt.Println(strings.Repeat("-", 16))

		viewCustomer()
		inputanCustomer()
	f:
		fmt.Print("\n[CAUTION!! MAX PELAYANAN 4]")
		fmt.Print("\nMasukkan jumlah pelayanan : ")
		var pel int
		fmt.Scanln(&pel)
		if pel > 4 {
			fmt.Println("Ulangi Masukkan jumlah pelayanan")
			fmt.Scanln()
			goto f
		}
		for i := 0; i < pel; i++ {
			fmt.Print("\n")
			fmt.Println(strings.Repeat("-", 14))
			fmt.Printf("|Pelayan ke-%d|\n", i+1)
			fmt.Println(strings.Repeat("-", 14))
			fmt.Print("\n")
			inputanTransaksi()
		}
		viewNotaHead(cs_id)
		viewNotaBody(cs_id)
		var pil string
		fmt.Print("\n\nApakah ada transaksi lain [Y/T]: ")
		fmt.Scanln(&pil)
		if pil == "y" || pil == "Y" {
			goto a
		} else {
			goto b
		}
	case 2:
	c:
		fmt.Println()
		fmt.Println(strings.Repeat("-", 21))
		fmt.Println("|MENU TAMPILKAN DATA|")
		fmt.Println(strings.Repeat("-", 21))

		fmt.Print("\nMenampilkan data :\n1. Nota\n2. customer\n3. transaksi\n4. layanan\n5. back to main menu\n6. exit\n==============: ")
		var view int
		fmt.Scanln(&view)
		switch view {
		case 1:
			nota()
			fmt.Print("\nApakah ingin melihat data yang lain [Y/T]: ")
			var pil string
			fmt.Scanln(&pil)
			if pil == "y" || pil == "Y" {
				goto c
			} else {
				goto b
			}
		case 2:
			viewCustomer()
			fmt.Print("\nApakah ingin melihat data yang lain [Y/T]: ")
			var pil string
			fmt.Scanln(&pil)
			if pil == "y" || pil == "Y" {
				goto c
			} else {
				goto b
			}
		case 3:
			viewTransaksi()
			viewTransaksidtl()
			fmt.Print("\nApakah ingin melihat data yang lain [Y/T]: ")
			var pil string
			fmt.Scanln(&pil)
			if pil == "y" || pil == "Y" {
				goto c
			} else {
				goto b
			}
		case 4:
			viewLayanan()
			fmt.Print("\nApakah ingin melihat data yang lain [Y/T]: ")
			var pil string
			fmt.Scanln(&pil)
			if pil == "y" || pil == "Y" {
				goto c
			} else {
				goto b
			}
		case 5:
			goto b
		default:
		}
	case 3:
	d:
		fmt.Println()
		fmt.Println(strings.Repeat("-", 18))
		fmt.Println("|MENU UPDATE DATA|")
		fmt.Println(strings.Repeat("-", 18))

		fmt.Print("\nUpdate Data : \n1. Customer\n2. Layanan\n3. kembali ke menu awal\n4. exit \n================: ")
		var pil int
		fmt.Scanln(&pil)
		switch pil {
		case 1:
			Updatecs()
			fmt.Print("\nApakah ingin mengupdate data yang lain [Y/T]: ")
			var pil string
			fmt.Scanln(&pil)
			if pil == "y" || pil == "Y" {
				goto d
			} else {
				goto b
			}
		case 2:
			updatelyn()
			fmt.Print("\nApakah ingin mengupdate data yang lain [Y/T]: ")
			var pil string
			fmt.Scanln(&pil)
			if pil == "y" || pil == "Y" {
				goto d
			} else {
				goto b
			}
		case 3:
			goto b
		default:

		}

	case 4:
	g:
		fmt.Println()
		fmt.Println(strings.Repeat("-", 16))
		fmt.Println("|MENU DELETE CSUTOMER|")
		fmt.Println(strings.Repeat("-", 16))
		deletecs()
		fmt.Print("\nApakah ingin menghapus data yang lain [Y/T]: ")
		var pil string
		fmt.Scanln(&pil)
		if pil == "y" || pil == "Y" {
			goto g
		} else {
			goto b
		}
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
	insertTransaksidtl(data.Transaksi_dtl{Id_trks: id_trks, Ttl_harga: total_harga + tambahan}, tx)
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
	fmt.Println("[CAUTION!!] Tuliskan ID berikutnya yang belum terdaftar")
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
	fmt.Print("\nMenu Paket:\n\n1. Paket Kilat    | 1 hari |+ 10 K fee|\n2. Paket Badai    | 2 hari |+ 05 K fee|\n3. Paket Gerimis  | 3 hari |+ 00 K fee|\n\n=======================================: ")
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
	fmt.Print("\nMenu Pelayanan:\n")
	viewLayananT()
	fmt.Print("\n============================: ")
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

// // select total harga
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
	updateQuery := "UPDATE customer SET ttl_byr = (SELECT SUM(total_harga) FROM transaksi_dtl JOIN transaksi ON transaksi_dtl.tr_id = transaksi.id WHERE transaksi.cs_id = $1) + tambahan WHERE id = $1;"
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
	fmt.Println(strings.Repeat("-", 75))
	fmt.Print("\t\t\t     ENIGMA AUNDRY\n")
	fmt.Println(strings.Repeat("-", 75))

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

	fmt.Printf("%-15s %-25s %s %s\n", "No Transaksi   : ", cs_id, "Nama customer : ", cs_nm)
	fmt.Printf("%-15s %-25s %s %s\n", "Tanggal Masuk  : ", date_in.Format("2006-01-02"), "No telepon    : ", no_tlp)
	fmt.Printf("%-15s %-25s %s %s\n", "Tanggal Selesai: ", date_out.Format("2006-01-02"), "Paket Laundry : ", paket)
	fmt.Printf("%-15s %-25s %s %d\n", "Di Terima oleh : ", ur_nm, "Fee paket     : ", tambahan)
	fmt.Println(strings.Repeat("-", 75))
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
	fmt.Printf("%-2s %-25s %-8s %-9s %-10s %-11s", "No|", "Pelayanan", "|Jumlah", "|Satuan", "|Harga", "|Total harga\n")
	fmt.Println(strings.Repeat("-", 75))
	for rows.Next() {
		rows.Scan(&no_nota, &jenis_layanan, &quantity, &satuan, &harga, &total_harga)
		fmt.Printf("%-3d %-26s %-8d %-9s %-10d %-9d\n", no_nota, jenis_layanan, quantity, satuan, harga, total_harga)
	}
	fmt.Print("\n\n\n")
	fmt.Print(strings.Repeat("-", 75))
	fmt.Printf("\n%+63s %-8d %s\n", "||  Toatal Pembayaran : ", ttl_byr, "||")
	fmt.Println(strings.Repeat("-", 75))

}

// view table customer
func viewCustomer() {
	db := connectDB()
	defer db.Close()

	query5 := "SELECT * FROM customer;"
	rows, err := db.Query(query5)
	if err != nil {
		fmt.Println("LIST MASIH BELUM TERISI")
	}
	defer rows.Next()
	fmt.Print("\nTable Customer:\n")
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("%-5s %-14s %-14s %-11s %-11s %-8s %s\n", "ID", "|Nama", "|No telepon", "|Tanggal Masuk", "|Tanggal Selesai", "|Tambahan", "|Total Bayar")
	fmt.Println(strings.Repeat("-", 90))

	for rows.Next() {
		rows.Scan(&cs_id, &cs_nm, &no_tlp, &date_in, &date_out, &tambahan, &ttl_byr)
		fmt.Printf("%-6s %-14s %-14s %-14s %-16s %-9d %d\n", cs_id, cs_nm, no_tlp, date_in.Format("2006-01-02"), date_out.Format("2006-01-02"), tambahan, ttl_byr)
	}
	fmt.Println(strings.Repeat("-", 90))

}

// view table transaksi
func viewTransaksi() {
	db := connectDB()
	defer db.Close()

	query6 := "SELECT * FROM transaksi;"
	rows, err := db.Query(query6)
	if err != nil {
		fmt.Println("LIST MASIH BELUM TERISI")
	}
	defer rows.Next()
	fmt.Print("\nTable Transaksi:\n")
	fmt.Println(strings.Repeat("-", 43))
	fmt.Printf("%-3s %-8s %-8s %-8s %-9s\n", "ID", "|User ID", "|CS ID", "|Layanan ID", "|Quantity")
	fmt.Println(strings.Repeat("-", 43))

	for rows.Next() {
		rows.Scan(&id_trks, &ur_id, &cs_id, &lyn_id, &quantity)
		fmt.Printf("%-4d %-8s %-8s %-11s %-9d\n", id_trks, ur_id, cs_id, lyn_id, quantity)
	}

}

// view table transaksi detail
func viewTransaksidtl() {
	db := connectDB()
	defer db.Close()

	query7 := "SELECT * FROM transaksi_dtl;"
	rows, err := db.Query(query7)
	if err != nil {
		fmt.Println("LIST MASIH BELUM TERISI!!")
	}
	defer rows.Close()
	fmt.Print("\nTable transaksi detail:\n")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("%-3s %-12s %-10s\n", "ID", "|Transaksi ID", "|Total harga")
	fmt.Println(strings.Repeat("-", 30))

	for rows.Next() {
		rows.Scan(&id_trks_dtl, &id_trks, &total_harga)
		fmt.Printf("%-4d %-13d %-10d\n", id_trks_dtl, id_trks, total_harga)
	}

}

// view  table layanan
func viewLayanan() {
	db := connectDB()
	defer db.Close()

	query8 := "SELECT * FROM layanan;"
	rows, err := db.Query(query8)
	if err != nil {
		fmt.Println("LIST MASIH BELUM TERISI!!")
	}
	defer rows.Close()
	fmt.Print("\nTable Layanan:\n")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("%-4s %-18s %-7s %s\n", "ID", "|Jenis Layanan", "|Satuan", "|Harga")
	fmt.Println(strings.Repeat("-", 40))

	for rows.Next() {
		rows.Scan(&lyn_id, &jenis_layanan, &satuan, &harga)
		fmt.Printf("%-5s %-18s %-7s %d\n", lyn_id, jenis_layanan, satuan, harga)
	}

}
func viewLayananT() {
	db := connectDB()
	defer db.Close()

	query8 := "SELECT ROW_NUMBER() OVER(ORDER BY id),jns_lyn,satuan,harga from layanan ORDER BY id ASC;"
	rows, err := db.Query(query8)
	if err != nil {
		fmt.Println("LIST MASIH BELUM TERISI!!")
	}
	defer rows.Close()
	fmt.Print("\nTable Layanan:\n")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("%-4s %-18s %-7s %s\n", "ID", "|Jenis Layanan", "|Satuan", "|Harga")
	fmt.Println(strings.Repeat("-", 40))

	for rows.Next() {
		rows.Scan(&lyn_id, &jenis_layanan, &satuan, &harga)
		fmt.Printf("%-5s %-18s %-7s %d\n", lyn_id, jenis_layanan, satuan, harga)
	}
	fmt.Println(strings.Repeat("-", 40))

}

// update customer bagian paket pelayanan
func enrollupdT(customer data.UpdateandDelete) {
	db := connectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	UpdateCustomer(customer, tx)
	updateTtlByr(tx, customer.Cs_id)
	err = tx.Commit()
	if err != nil {
		panic(err)
	} else {

		fmt.Print(strings.Repeat("-", 23))
		fmt.Print("\n|Transaction Commited!|\n")
		fmt.Println(strings.Repeat("-", 23))
	}
}

func UpdateCustomer(update data.UpdateandDelete, tx *sql.Tx) {
	sqlStatement := "UPDATE customer SET date_in = $1, date_out = $2, tambahan = $3 WHERE id = $4;"
	_, err := tx.Exec(sqlStatement, date_in, date_out, tambahan, cs_id)
	if err != nil {
		fmt.Println("ID cs tidak ada!")
	}
	validate(err, "Update Changes in customer!", tx)
}

func Updatecs() {
	viewCustomer()
	fmt.Print("\nMasukkan Id customer      : ")
	fmt.Scanln(&cs_id)
	var pil int
	fmt.Print("\nMenu Paket:\n\n1. Paket Kilat    | 1 hari |+ 10 K fee|\n2. Paket Badai    | 2 hari |+ 05 K fee|\n3. Paket Gerimis  | 3 hari |+ 00 K fee|\n\n=======================================: ")
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
	update := data.UpdateandDelete{Cs_id: cs_id, Date_in: date_in, Date_out: date_out, Tambahan: tambahan}
	enrollupdT(update)
}

// delete customer
func deletecs() {
	viewCustomer()
	fmt.Print("\nMasukkan ID cs yang ingin dihapus : ")
	fmt.Scanln(&cs_id)
	deletecs := data.UpdateandDelete{Cs_id: cs_id}
	enrollDEletecs(deletecs)
}

func enrollDEletecs(delete data.UpdateandDelete) {
	db := connectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	deleteCustomer(delete.Cs_id, tx)
	err = tx.Commit()
	if err != nil {
		panic(err)
	} else {

		fmt.Print(strings.Repeat("-", 23))
		fmt.Print("\n|Transaction Commited!|\n")
		fmt.Println(strings.Repeat("-", 23))
	}
}

func deleteCustomer(id string, tx *sql.Tx) {
	// Hapus data dari tabel transaksi_detail
	deleteTrksDtl := "DELETE FROM transaksi_dtl WHERE tr_id IN (SELECT id FROM transaksi WHERE cs_id = $1);"
	_, err := tx.Exec(deleteTrksDtl, id)
	if err != nil {
		fmt.Println("Error deleting transaksi detail:", err)
		tx.Rollback()
		return
	}

	// Hapus data dari tabel transaksi
	deleteTrks := "DELETE FROM transaksi WHERE cs_id = $1;"
	_, err = tx.Exec(deleteTrks, id)
	if err != nil {
		fmt.Println("Error deleting transaksi:", err)
		tx.Rollback()
		return
	}

	// Hapus data dari tabel customer
	deleteCs := "DELETE FROM customer WHERE id = $1;"
	_, err = tx.Exec(deleteCs, id)
	if err != nil {
		fmt.Println("Error deleting customer:", err)
		tx.Rollback()
		return
	}
	validate(err, "Delete Customer!", tx)
}

// update layanan
func updatelyn() {
	viewLayanan()
	fmt.Print("\nMasukkan ID layanan yang ingin di Update harganya : ")
	fmt.Scanln(&lyn_id)
	fmt.Print("\nTentukan Harga terbaru : ")
	fmt.Scanln(&harga)
	updatelyn := data.UpdateandDelete{Lyn_id: lyn_id, Harga: harga}
	enrollupdatelayanan(updatelyn)
}

func enrollupdatelayanan(update data.UpdateandDelete) {
	db := connectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	updateLayanan(update.Harga, update.Lyn_id, tx)
	err = tx.Commit()
	if err != nil {
		panic(err)
	} else {
		fmt.Print(strings.Repeat("-", 23))
		fmt.Print("\n|Transaction Commited!|\n")
		fmt.Println(strings.Repeat("-", 23))
	}

}
func updateLayanan(harga int64, id_lyn string, tx *sql.Tx) {
	query9 := "UPDATE layanan SET harga = $1 WHERE id = $2;"
	_, err := tx.Exec(query9, harga, id_lyn)
	if err != nil {
		fmt.Println("ID Layanan Belum terdaftar atau tidak ada!")
		tx.Rollback()
	}
}
