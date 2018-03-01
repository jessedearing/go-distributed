// Package mongo provides the MongoDB locking capability
package mongo

import (
	"time"

	"github.com/jessedearing/go-distributed/lock"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	lock.Lockers["mongo"] = newMongoLocker
}

type mongoLock struct {
	session  *mgo.Session
	database string
}

func newMongoLocker(connectionURL string) (lock.DistributedLocker, error) {
	sess, err := mgo.Dial(connectionURL)
	if err != nil {
		return nil, err
	}

	sess.SetMode(mgo.Primary, true)

	return &mongoLock{session: sess}, nil
}

func (m *mongoLock) Lock(lockname string) {
	m.lock(lockname, false)
}

func (m *mongoLock) NonBlockLock(lockname string) bool {
	return m.lock(lockname, true)
}

func (m *mongoLock) lock(lockname string, nonblocking bool) bool {
	t := time.NewTicker(50 * time.Millisecond)
	var errcount int
	for range t.C {
		i, err := m.session.DB("jesse").C("locks").Upsert(bson.M{"name": lockname}, bson.M{"name": lockname})
		if err != nil {
			if errcount < 5 {
				panic(err)
			} else {
				t.Stop()
				break
			}
			continue
		}

		// Got the lock
		if i.UpsertedId != nil {
			t.Stop()
			return true
		}

		if nonblocking {
			t.Stop()
			break
		}
	}
	return false
}

func (m *mongoLock) Unlock(lockname string) {
	_, err := m.session.DB("jesse").C("locks").RemoveAll(bson.M{"name": lockname})
	if err != nil {
		panic(err)
	}
}

func (m *mongoLock) Close() {
	m.session.Close()
}
