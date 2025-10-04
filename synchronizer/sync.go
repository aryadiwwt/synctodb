package synchronizer

import (
	"context"
	"fmt"
	"log"

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

// SynchronizePosts menjalankan seluruh alur kerja sinkronisasi.
func (ps *OutputDetailSynchronizer) Synchronize(ctx context.Context) error {
	ps.log.Println("Starting post synchronization...")

	details, err := ps.fetcher.FetchOutputDetails(ctx)
	if err != nil {
		return fmt.Errorf("synchronization failed during fetch phase: %w", err)
	}
	ps.log.Printf("Successfully fetched %d details.", len(details))

	if len(details) == 0 {
		ps.log.Println("No new details to synchronize.")
		return nil
	}

	if err := ps.storer.StoreOutputDetails(ctx, details); err != nil {
		return fmt.Errorf("synchronization failed during store phase: %w", err)
	}
	ps.log.Println("Successfully stored details to the database.")

	ps.log.Println("detail synchronization finished successfully.")
	return nil
}
