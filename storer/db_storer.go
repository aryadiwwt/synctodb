package storer

import (
	"context"
	"fmt"

	"github.com/aryadiwwt/synctodb/domain"
	customErrors "github.com/aryadiwwt/synctodb/errors"

	"github.com/jmoiron/sqlx"
)

type Wilayah struct {
	KodeProvinsi  string `db:"provinsi_id"`
	KodeKabupaten string `db:"kota_id"`
}

// Storer mendefinisikan kontrak untuk menyimpan data post.

type Storer interface {
	StoreOutputDetails(ctx context.Context, details []domain.OutputDetail) error
	// Diubah: Menerima slice kode provinsi untuk difilter
	GetWilayahByProvinsi(ctx context.Context, kodeProvinsi []string) ([]Wilayah, error)
}

// type Storer interface {
// 	StoreOutputDetails(ctx context.Context, details []domain.OutputDetail) error
// 	GetAllWilayah(ctx context.Context) ([]Wilayah, error)
// }

// Implementasi fungsi diubah untuk memfilter berdasarkan kd_prov
func (s *dbStorer) GetWilayahByProvinsi(ctx context.Context, kodeProvinsi []string) ([]Wilayah, error) {
	var wilayah []Wilayah

	// Query dasar
	baseQuery := `SELECT provinsi_id, kota_id FROM master_kota`

	var args []interface{}

	// Jika daftar provinsi diberikan, tambahkan klausa WHERE IN
	if len(kodeProvinsi) > 0 {
		baseQuery += ` WHERE provinsi_id IN (?)`
		args = append(args, kodeProvinsi)
	}

	baseQuery += ` ORDER BY provinsi_id, kota_id`

	// sqlx.In secara aman akan mengubah query (?) menjadi ($1, $2, ...)
	// dan menyesuaikan argumennya. Ini cara aman untuk klausa IN.
	query, args, err := sqlx.In(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat query IN: %w", err)
	}

	// Rebind query agar sesuai dengan placeholder PostgreSQL ($1, $2)
	query = s.db.Rebind(query)

	err = s.db.SelectContext(ctx, &wilayah, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data wilayah yang difilter: %w", err)
	}

	return wilayah, nil
}

// // Fungsi untuk mengambil semua wilayah dari DB
// func (s *dbStorer) GetAllWilayah(ctx context.Context) ([]Wilayah, error) {
// 	var wilayah []Wilayah
// 	query := `SELECT provinsi_id, kota_id FROM master_kota ORDER BY provinsi_id, kota_id`
// 	err := s.db.SelectContext(ctx, &wilayah, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal mengambil data wilayah: %w", err)
// 	}
// 	return wilayah, nil
// }

type dbStorer struct {
	db *sqlx.DB
}

func NewDBStorer(db *sqlx.DB) Storer {
	return &dbStorer{db: db}
}

const (
	// Query disimpan sebagai konstanta untuk menghindari 'magic strings'
	// dan memudahkan pengelolaan.
	upsertOutputDetailQuery = `INSERT INTO siskeudes_detail_output (
            tahun, kd_prov, nama_provinsi, kd_kab, nama_kabupaten,
            kd_kec, nama_kecamatan, kd_desa, nama_desa, id_keg,
            nama_kegiatan, kode_sumber, pagu, kode_output, no_id, nama_paket,
            lokasi, waktu, keluaran, uraian_output, volume, satuan, nilai,
            anggaran1, anggaran2, realisasi0, realisasi1, realisasi2, fisik0,
            fisik1, fisik2, namapptkd, nippptkd, jbtpptkd
        ) VALUES (
            :tahun, :kd_prov, :nama_provinsi, :kd_kab, :nama_kabupaten,
            :kd_kec, :nama_kecamatan, :kd_desa, :nama_desa, :id_keg,
            :nama_kegiatan, :kode_sumber, :pagu, :kode_output, :no_id, :nama_paket,
            :lokasi, :waktu, :keluaran, :uraian_output, :volume, :satuan, :nilai,
            :anggaran1, :anggaran2, :realisasi0, :realisasi1, :realisasi2, :fisik0,
            :fisik1, :fisik2, :namapptkd, :nippptkd, :jbtpptkd
        )
        ON CONFLICT (tahun, kd_prov, kd_kab, kd_kec, kd_desa, id_keg, no_id) DO UPDATE SET
            nama_provinsi = EXCLUDED.nama_provinsi,
            pagu = EXCLUDED.pagu,
            nilai = EXCLUDED.nilai,
            -- Tambahkan kolom lain yang ingin Anda update jika terjadi konflik
            nama_paket = EXCLUDED.nama_paket,
            realisasi1 = EXCLUDED.realisasi1,
            realisasi2 = EXCLUDED.realisasi2;`
)

func (s *dbStorer) StoreOutputDetails(ctx context.Context, details []domain.OutputDetail) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return &customErrors.ErrDBOperationFailed{Operation: "begin_transaction", Err: err}
	}
	defer tx.Rollback() // Aman untuk dipanggil meskipun sudah di-commit.

	for _, detail := range details {
		if _, err := tx.NamedExecContext(ctx, upsertOutputDetailQuery, detail); err != nil {
			return &customErrors.ErrDBOperationFailed{Operation: "upsert_post", Err: err}
		}
	}

	if err := tx.Commit(); err != nil {
		return &customErrors.ErrDBOperationFailed{Operation: "commit_transaction", Err: err}
	}

	return nil
}
