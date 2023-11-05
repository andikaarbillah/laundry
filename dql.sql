CREATE DATABASE submit;--membuat database
CREATE TABLE role(id varchar (4) PRIMARY KEY,--membuat table role
                                 rl_nm varchar(5)
                                 );
CREATE TABLE user_(id varchar(5) PRIMARY KEY, --membuat table pengguna
                                  id_rl varchar(4),
                                  ur_nm varchar(20),
                                  pswd varchar(8),
                                  FOREIGN KEY(id_rl) REFERENCES role(id)
                                  );
CREATE TABLE customer(id varchar(5) PRIMARY KEY, --membuat table customer
					 nm_cs varchar(20),
					 no_tp varchar(14),
					 date_in date,
					 date_out date,
					 tambahan bigint,
					 ttl_byr bigint
					 );
CREATE TABLE layanan(id varchar(6) PRIMARY KEY,--membuat table layanan yang berisi jenis pelayanan
					jns_lyn varchar (20),
					satuan varchar(20),
					harga bigint
					);
CREATE TABLE transaksi(id serial PRIMARY KEY, --membuat table transaksi 
					  ur_id varchar(5),
					  cs_id varchar(5),
					  lyn_id varchar (5),
					  quantity int,
					  FOREIGN KEY(ur_id)REFERENCES user_(id),
					  FOREIGN KEY(cs_id)REFERENCES customer(id),
					  FOREIGN KEY(lyn_id)REFERENCES layanan(id)
					  );
CREATE TABLE transaksi_dtl(id serial PRIMARY KEY,--membuat table transaksi detail untuk menampung total harga
					  tr_id serial,
					  total_harga bigint,
					  FOREIGN KEY(tr_id)REFERENCES transaksi(id)
					  );
