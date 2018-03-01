// Package mysql provides the MySQL locking capability
package mysql

import (
	"database/sql"

	// need to pull in mysql driver for database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessedearing/go-distributed/lock"
)

func init() {
	lock.Lockers["mysql"] = newMySQLLocker
}

type mySQLLock struct {
	db *sql.DB
}

func newMySQLLocker(connectionString string) (lock.DistributedLocker, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	my := &mySQLLock{
		db: db,
	}

	return my, nil
}

func (m *mySQLLock) Lock(lockname string) {
	m.lock(lockname, false)
}

func (m *mySQLLock) NonBlockLock(lockname string) bool {
	return m.lock(lockname, true)
}

func (m *mySQLLock) Unlock(lockname string) {
}

func (m *mySQLLock) Close() {
	m.db.Close()
}

func (m *mySQLLock) lock(lockname string, nonblocking bool) bool {
	var timeout int
	if !nonblocking {
		timeout = -1
	}
	rows, err := m.db.Query("SELECT GET_LOCK(?,?)", lockname, timeout)
	if err != nil {
		panic(err)
		return false
	}
	defer rows.Close()

	var isLocked *int
	rows.Next()

	err = rows.Scan(&isLocked)
	if err != nil {
		panic(err)
	}

	// Got the lock
	if isLocked != nil && *isLocked == 1 {
		return true
	}

	return false
}
