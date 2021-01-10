// A simple kernel service wich provides access to a single github.com/etcd-io/bbolt
// object store
package bolt

import (
	"flag"
	"fmt"
	bbolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/golib/kernel"
	"os"
	"time"
)

type BoltService struct {
	FileName string
	db       *bbolt.DB
	dbFile   *string
}

func (s *BoltService) Name() string {
	return "bolt:" + s.FileName
}

func (s *BoltService) Init(k *kernel.Kernel) error {
	if s.FileName == "" {
		s.dbFile = flag.String("bucket-store", "", "The file to store all buckets")
	}
	return nil
}

func (s *BoltService) PostInit() error {
	if s.FileName == "" && s.dbFile != nil {
		s.FileName = *s.dbFile
	}

	if s.FileName == "" {
		s.FileName = os.Getenv("BUCKETSTORE")
	}
	if s.FileName == "" {
		return fmt.Errorf("No store provided by -bucket-store or BUCKETSTORE")
	}
	return nil
}

func (s *BoltService) Start() error {
	db, err := bbolt.Open(s.FileName, 0666, &bbolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *BoltService) Stop() {
	s.db.Close()
}
