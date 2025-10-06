package synchronizer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aryadiwwt/synctodb/domain"
	"github.com/aryadiwwt/synctodb/fetcher"
	"github.com/aryadiwwt/synctodb/storer"
)

// PostSynchronizer mengorkestrasi proses sinkronisasi data post.
type OutputDetailSynchronizer struct {
	fetcher fetcher.Fetcher
	storer  storer.Storer
	log     *log.Logger
}

func NewOutputDetailSynchronizer(f fetcher.Fetcher, s storer.Storer, l *log.Logger) *OutputDetailSynchronizer {
	return &OutputDetailSynchronizer{
		fetcher: f,
		storer:  s,
		log:     l,
	}
}

// Synchronize mengurutkan alur kerja, sekarang dengan langkah transformasi.
func (s *OutputDetailSynchronizer) Synchronize(ctx context.Context) error {
	s.log.Println("Starting output detail synchronization...")

	// 1. FETCH: Ambil data mentah dari API
	details, err := s.fetcher.FetchOutputDetails(ctx)
	if err != nil {
		return fmt.Errorf("synchronization failed during fetch phase: %w", err)
	}
	s.log.Printf("Successfully fetched %d output details.", len(details))

	if len(details) == 0 {
		s.log.Println("No new data to synchronize.")
		return nil
	}

	// 2. TRANSFORM: Ubah data sesuai aturan bisnis
	s.log.Println("Transforming data...")
	transformedDetails := transformDetails(details)
	s.log.Println("Data transformation complete.")

	// 3. STORE: Simpan data yang sudah ditransformasi
	if err := s.storer.StoreOutputDetails(ctx, transformedDetails); err != nil {
		return fmt.Errorf("synchronization failed during store phase: %w", err)
	}
	s.log.Println("Successfully stored transformed data to the database.")

	s.log.Println("Synchronization finished successfully.")
	return nil
}

// FUNGSI BARU: transformDetails berisi logika untuk mengubah data
func transformDetails(details []domain.OutputDetail) []domain.OutputDetail {
	// Loop melalui setiap record dan modifikasi nilainya
	for i := range details {
		// Simpan nilai asli untuk digunakan dalam penggabungan
		prov := details[i].KodeProvinsi
		kab := details[i].KodeKabupaten
		kec := details[i].KodeKecamatan
		desa := details[i].KodeDesa

		// Aturan 1: kd_kab = kd_prov.kd_kab
		details[i].KodeKabupaten = fmt.Sprintf("%s.%s", prov, kab)

		// Aturan 2: kd_kec = kd_prov.kd_kab.kd_kec
		details[i].KodeKecamatan = fmt.Sprintf("%s.%s.%s", prov, kab, kec)

		// Aturan 3: kd_desa = kd_prov.kd_kab.kd_desa
		// Membersihkan titik di akhir kd_desa dari API jika ada
		cleanedDesa := strings.TrimSuffix(desa, ".")
		details[i].KodeDesa = fmt.Sprintf("%s.%s.%s", prov, kab, cleanedDesa)
	}

	return details
}
