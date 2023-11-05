INSERT INTO role(id,rl_nm)VALUES('rl01','owner'),--mengisi table role
('rl02','admin'),
('rl03','kasir');

INSERT INTO user_(id,id_rl,ur_nm,pswd)VALUES --mengisi user yang berhak akses
('U111','rl01','andika','230104'),
('U112','rl02','mirna','123456'),
('U113','rl03','genjie','000111');

INSERT INTO layanan(id,jns_lyn,satuan,harga)VALUES -- mengisi table layanan yang tersedia
('L11','Cuci + Gosok','Kg',7000),
('L12','Laundry Bedcover','Buah',50000),
('L13','Laundry Boneka','Buah',25000),
('L14','Laundry Amabal','Meter',20000);