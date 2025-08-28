# recordtocsv

`recordtocsv` adalah library sederhana untuk **mencatat data (log/record) ke file CSV** dengan format nama file yang otomatis menyesuaikan waktu (harian, bulanan, atau tahunan).  
Library ini cocok digunakan untuk mencatat log request/response, transaksi, atau data lain yang perlu disimpan secara rapi dalam format CSV.

---

## ✨ Features
- Simpan data ke file CSV dengan **header otomatis**.
- Mendukung format nama file berdasarkan waktu:
  - **Daily** → `filename_YYYY_MM_DD.csv`
  - **Monthly** → `filename_YYYY_MM.csv`
  - **Yearly** → `filename_YYYY.csv`
- Secara otomatis membuat direktori tujuan jika belum ada.
- Payload fleksibel: bisa berupa struct atau map.

---

## 📦 Installation

```bash
go get github.com/ojipoji/recordtocsv
```

### Inisialisasi service

```go
import "github.com/username/recordtocsv"

columns := []string{"id", "request", "response"}

service := recordtocsv.NewRecordToCSV(
    "files/record",            // Folder penyimpanan
    "agoda_booking_record",    // Nama dasar file
    columns,                   // Kolom CSV
    "daily",                   // Jenis record: daily, monthly, yearly
)
```

### Mencatat record

```go
payload := map[string]interface{}{
    "id":       1,
    "request":  "GET /booking/123",
    "response": "200 OK",
}

if err := service.Record(payload); err != nil {
    panic(err)
}
```

### Contoh secara keseluruhan

```go
package main

import (
	"fmt"
	"github.com/ojipoji/recordtocsv"
)

type BookingRecord struct {
	ID       string `json:"id"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

func main() {
	// Inisialisasi service
	service := recordtocsv.NewRecordToCSV(
		"files/record",                     // Direktori tujuan
		"booking_record",                   // Nama file dasar
		[]string{"id", "request", "response"}, // Header CSV
		"daily",                            // Tipe record: daily | monthly | yearly
	)

	// Data yang ingin dicatat
	record := BookingRecord{
		ID:       "123",
		Request:  `{"room":"Deluxe"}`,
		Response: `{"status":"confirmed"}`,
	}

	// Simpan ke CSV
	if err := service.Record(record); err != nil {
		panic(err)
	}

	fmt.Println("Record berhasil dicatat ke CSV!")
}
```

### Struktur File yang Dihasilkan

```bash
files/record/
└── booking_record_2025_08_26.csv
```

Isi file:

```bash
id,request,response
123,"{""room"":""Deluxe""}","{""status"":""confirmed""}"
```

Jika kamu menggunakan monthly, maka nama file akan menjadi:

```bash
files/record/agoda_booking_record_2025_08.csv
```

---

### ⚠️ Notes

- Pastikan kolom (Column) sesuai dengan field JSON pada struct/map yang Anda kirimkan.
- Lokasi waktu default adalah Asia/Jakarta.
- Mendukung berbagai tipe data sederhana (string, int, float, dll.).



