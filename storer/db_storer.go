package storer

import (
	"context"

	"github.com/aryadiwwt/synctodb/domain"
	customErrors "github.com/aryadiwwt/synctodb/errors"

	"github.com/jmoiron/sqlx"
)

// Storer mendefinisikan kontrak untuk menyimpan data post.
type Storer interface {
	StoreOutputDetails(ctx context.Context, details []domain.OutputDetail) error
}

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
