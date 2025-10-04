package domain

// OutputDetail merepresentasikan struktur data detail output kegiatan.
type OutputDetail struct {
	Tahun         string  `json:"tahun" db:"tahun"`
	KodeProvinsi  string  `json:"kd_prov" db:"kd_prov"`
	NamaProvinsi  string  `json:"nama_provinsi" db:"nama_provinsi"`
	KodeKabupaten string  `json:"kd_kab" db:"kd_kab"`
	NamaKabupaten string  `json:"nama_kabupaten" db:"nama_kabupaten"`
	KodeKecamatan string  `json:"kd_kec" db:"kd_kec"`
	NamaKecamatan string  `json:"nama_kecamatan" db:"nama_kecamatan"`
	KodeDesa      string  `json:"kd_desa" db:"kd_desa"`
	NamaDesa      string  `json:"nama_desa" db:"nama_desa"`
	IDKegiatan    string  `json:"id_keg" db:"id_keg"`
	NamaKegiatan  string  `json:"nama_kegiatan" db:"nama_kegiatan"`
	KodeSumber    string  `json:"kode_sumber" db:"kode_sumber"`
	Pagu          float64 `json:"pagu,string" db:"pagu"`
	KodeOutput    string  `json:"kode_output" db:"kode_output"`
	NoID          string  `json:"no_id" db:"no_id"`
	NamaPaket     string  `json:"nama_paket" db:"nama_paket"`
	Lokasi        string  `json:"lokasi" db:"lokasi"`
	Waktu         string  `json:"waktu" db:"waktu"`
	Keluaran      string  `json:"keluaran" db:"keluaran"`
	UraianOutput  string  `json:"uraian_output" db:"uraian_output"`
	Volume        float64 `json:"volume,string" db:"volume"`
	Satuan        string  `json:"satuan" db:"satuan"`
	Nilai         float64 `json:"nilai,string" db:"nilai"`
	Anggaran1     float64 `json:"anggaran1,string" db:"anggaran1"`
	Anggaran2     float64 `json:"anggaran2,string" db:"anggaran2"`
	Realisasi0    float64 `json:"realisasi0,string" db:"realisasi0"`
	Realisasi1    float64 `json:"realisasi1,string" db:"realisasi1"`
	Realisasi2    float64 `json:"realisasi2,string" db:"realisasi2"`
	Fisik0        float64 `json:"fisik0,string" db:"fisik0"`
	Fisik1        float64 `json:"fisik1,string" db:"fisik1"`
	Fisik2        float64 `json:"fisik2,string" db:"fisik2"`
	NamaPPTKD     string  `json:"namapptkd" db:"namapptkd"`
	NIPPPTKD      string  `json:"nippptkd" db:"nippptkd"`
	JabatanPPTKD  string  `json:"jbtpptkd" db:"jbtpptkd"`
}
