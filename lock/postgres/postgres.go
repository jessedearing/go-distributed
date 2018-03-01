package postgres

import (
	"database/sql"
	"hash/crc64"

	"github.com/jessedearing/go-distributed/lock"
	// importing pq for use by database/sql
	_ "github.com/lib/pq"
)

type postgresLocker struct {
	db *sql.DB
}

func init() {
	lock.Lockers["postgres"] = newPostgresLocker
}

func newPostgresLocker(connectionString string) (lock.DistributedLocker, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	plock := &postgresLocker{
		db: db,
	}

	return plock, nil
}

func (p *postgresLocker) Lock(lockname string) {
	lockid := getLockID(lockname)
	_, err := p.db.Exec("select pg_advisory_lock($1)", lockid)
	if err != nil {
		panic(err)
	}
}

func (p *postgresLocker) NonBlockLock(lockname string) bool {
	lockid := getLockID(lockname)
	rows, err := p.db.Query("select pg_try_advisory_lock($1)", lockid)
	if err != nil {
		panic(err)
	}

	var gotLock bool

	rows.Next()

	var errors []error
	errors = append(errors, rows.Scan(&gotLock))

	rows.Next()
	errors = append(errors, rows.Close())
	errors = append(errors, rows.Err())

	if len(errors) > 0 {
		panic(errors)
	}

	return gotLock
}

func (p *postgresLocker) Unlock(lockname string) {
	_, err := p.db.Exec("select pg_advisory_unlock_all()")
	if err != nil {
		panic(err)
	}
}

func (p *postgresLocker) Close() {
	p.db.Close()
}

func getLockID(lockname string) uint64 {
	return crc64.Checksum([]byte(lockname), crc64.MakeTable(crc64.ECMA))
}
