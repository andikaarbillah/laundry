package data

import "time"

type User struct {
	Id    string
	Id_rl string
	Nama  string
	Pswd  string
}

type Role struct {
	Id        string
	Role_name string
}

type Customer struct {
	Id          string
	Nama        string
	No_tlp      string
	Date_in     time.Time
	Date_out    time.Time
	Tambahan    int64
	Total_bayar int64
}

type Layanan struct {
	Id     string
	Lyn_nm string
	Satuan string
	Harga  int64
}

type Transaksi struct {
	Id       int64
	Ur_id    string
	Cs_id    string
	Lyn_id   string
	Quantity int
}

type Transaksi_dtl struct {
	Id_trks   int64
	Ttl_harga int64
}

type UpdateandDelete struct {
	Cs_id    string
	Date_in  time.Time
	Date_out time.Time
	Tambahan int64
	Lyn_id   string
	Quantity int
	Id_trks  int
	Harga    int64
}
