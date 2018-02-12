package storage

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mchetelat/bazo_miner/protocol"
	"log"
	"os"
	"strings"
	"time"
)

var (
	db                 *bolt.DB
	logger             *log.Logger
	State              = make(map[[32]byte]*protocol.Account)
	RootKeys           = make(map[[32]byte]*protocol.Account)
	txMemPool          = make(map[[32]byte]protocol.Transaction)
	AllClosedBlocksAsc []*protocol.Block
)

const (
	ERROR_MSG = "Initiate storage aborted: "
)

//Entry function for the storage package
func Init(connTuple string, dbname string) {
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	var err error
	db, err = bolt.Open(dbname, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		logger.Fatal(ERROR_MSG, err)
	}

	//Check if db file is empty for all non-bootstraping miners
	ipport := strings.Split(connTuple, ":")
	if ipport[1] != "8000" {
		err := db.View(func(tx *bolt.Tx) error {
			err := tx.ForEach(func(name []byte, bkt *bolt.Bucket) error {
				err := bkt.ForEach(func(k, v []byte) error {
					if k != nil && v != nil {
						return errors.New("Non-empty database given.")
					}
					return nil
				})
				return err
			})
			return err
		})

		if err != nil {
			logger.Fatal(ERROR_MSG, err)
		}
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("openblocks"))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("closedblocks"))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("closedfunds"))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("closedaccs"))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("closedstakes"))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("closedconfigs"))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("lastclosedblock"))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
}

func TearDown() {
	db.Close()
}