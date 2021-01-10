package db

import (
	"database/sql"
	"log"
)

func (tx *Tx) SetEncoding(enc string) error {
	_, err := tx.Exec("SET CLIENT_ENCODING TO " + enc)
	return err
}

// Set's the connection encoding to "UTF8" - the default in the database.
// Use if SetEncodingWIN1252 or SetEncoding has been used.
func (tx *Tx) SetEncodingUTF8() error {
	return tx.SetEncoding("UTF8")
}

// Set's the connection encoding to "WIN1252" - use with csv style imports
// which have latin characters rather than UTF8
func (tx *Tx) SetEncodingWIN1252() error {
	return tx.SetEncoding("WIN1252")
}

// DeleteFrom is a wrapper around deleting from a table
func (tx *Tx) DeleteFrom(table string) (sql.Result, error) {
	log.Println("Deleting from", table)

	result, err := tx.Exec("DELETE FROM " + table)
	if err != nil {
		return nil, err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	log.Println("Deleted", ra, "from", table)

	return result, err
}

func (tx *Tx) OnCommitVacuum(table string) {
	tx.OnCommit(func() error {
		log.Println("VACUUM", table)
		_, err := tx.db.Exec("VACUUM " + table)
		return err
	})
}

func (tx *Tx) OnCommitVacuumFull(table string) {
	tx.OnCommit(func() error {
		log.Println("VACUUM FULL", table)
		_, err := tx.db.Exec("VACUUM " + table)
		return err
	})
}

func (tx *Tx) OnCommitCluster(table, column string) {
	tx.OnCommit(func() error {
		log.Println("VACUUM", table)
		_, err := tx.db.Exec("VACUUM " + table)
		if err != nil {
			return err
		}

		log.Println("Clustering", table, "on", column)
		_, err = tx.db.Exec("CLUSTER " + table + " USING " + column)
		if err != nil {
			return err
		}

		log.Println("ANALYZE", table)
		_, err = tx.db.Exec("ANALYZE " + table)
		return err
	})
}
