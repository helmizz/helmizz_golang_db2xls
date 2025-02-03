Berikut adalah contoh file `README.md` untuk library Anda. File ini memberikan informasi penting tentang library, cara instalasi, penggunaan, dan kontribusi.

---

# Helmizz DB2XLS

[![Go Reference](https://pkg.go.dev/badge/github.com/helmizz/helmizz_db2xls.svg)](https://pkg.go.dev/github.com/helmizz/helmizz_db2xls)

Helmizz DB2XLS adalah library Go yang memungkinkan Anda mengekspor data dari database ke file XLSX (Excel). Library ini mendukung berbagai jenis database seperti MySQL, PostgreSQL, SQLite, Oracle, dan SQL Server.

## Fitur
- Mendukung berbagai database: MySQL, PostgreSQL, SQLite, Oracle, dan SQL Server.
- Mengekspor hasil query ke file XLSX dengan format yang dapat disesuaikan.
- Mendukung beberapa query dalam satu konfigurasi, masing-masing ditulis ke sheet yang berbeda.
- Mudah digunakan dengan konfigurasi JSON.

## Instalasi
Untuk menggunakan library ini, tambahkan ke proyek Go Anda dengan perintah berikut:
```bash
go get -u github.com/helmizz/helmizz_db2xls
```

## Penggunaan
### 1. Konfigurasi JSON
Buat file konfigurasi JSON seperti berikut:
```json
{
  "type_db": "mysql",
  "host": "localhost",
  "port": "3306",
  "database": "testdb",
  "username": "root",
  "password": "password",
  "query": [
    {
      "sql": "SELECT * FROM users",
      "sheet_name": "Users"
    },
    {
      "sql": "SELECT * FROM products",
      "sheet_name": "Products"
    }
  ],
  "output_file": "output.xlsx"
}
```

### 2. Contoh Kode
Berikut adalah contoh cara menggunakan library ini dalam kode Go:
```go
package main

import (
	"log"
	"github.com/helmizz/helmizz_db2xls"
)

func main() {
	// Path ke file konfigurasi JSON
	configPath := "config/config.json"

	// Ekspor data dari database ke file XLSX
	err := helmizz_db2xls.ExportDBToXLSX(configPath)
	if err != nil {
		log.Fatalf("Error exporting data: %v", err)
	}

	log.Println("Data exported successfully!")
}
```

### 3. Output
File XLSX akan dihasilkan sesuai dengan konfigurasi yang diberikan. Setiap query akan ditulis ke sheet yang berbeda dalam file XLSX.

## Kontribusi
Kami sangat menghargai kontribusi Anda! Jika Anda ingin berkontribusi, silakan ikuti langkah-langkah berikut:
1. Fork repositori ini.
2. Buat branch baru (`git checkout -b feature/your-feature-name`).
3. Commit perubahan Anda (`git commit -m 'Add some feature'`).
4. Push ke branch (`git push origin feature/your-feature-name`).
5. Buat pull request.

## Lisensi
Library ini dilisensikan di bawah [MIT License](LICENSE).

## Kebutuhan Sistem
- Go 1.20 atau lebih tinggi.
- Driver database yang sesuai untuk jenis database yang digunakan:
  - MySQL: `github.com/go-sql-driver/mysql`
  - PostgreSQL: `github.com/lib/pq`
  - SQLite: `github.com/mattn/go-sqlite3`
  - Oracle: `github.com/godror/godror`
  - SQL Server: `github.com/denisenkom/go-mssqldb`

## Dukungan
Jika Anda memiliki pertanyaan atau masalah, silakan buka issue di repositori GitHub ini.

---

Dengan file `README.md` ini, pengguna akan memiliki pemahaman yang jelas tentang cara menggunakan library Anda, serta informasi tentang kontribusi dan lisensi.
