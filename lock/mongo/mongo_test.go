// +build mongo

package mongo_test

import (
	"testing"
	"time"

	"github.com/jessedearing/go-distributed/lock"
	_ "github.com/jessedearing/go-distributed/lock/mongo"
	"github.com/stretchr/testify/assert"
)

func TestBlockingLocks(t *testing.T) {
	l, err := lock.New("mongo", "mongodb://127.0.0.1:27017/mydb")
	assert.Nil(t, err)
	defer l.Close()
	l2, err := lock.New("mongo", "mongodb://127.0.0.1:27017/mydb")
	assert.Nil(t, err)
	defer l2.Close()

	var lock1Chan = make(chan struct{})
	var lock1Aquired = make(chan struct{}, 1)
	var lock2Chan = make(chan struct{})
	var lock2Aquired = make(chan struct{}, 1)

	go func() {
		<-lock1Chan
		l.Lock("test")
		lock1Aquired <- struct{}{}
	}()

	go func() {
		<-lock2Chan
		l2.Lock("test")
		lock2Aquired <- struct{}{}
	}()

	lock1Chan <- struct{}{}
	time.Sleep(500 * time.Millisecond)
	lock2Chan <- struct{}{}

	assert.Len(t, lock1Aquired, 1)
	assert.Len(t, lock2Aquired, 0)

	l.Unlock("test")

	time.Sleep(250 * time.Millisecond)
	assert.Len(t, lock2Aquired, 1)

	l2.Unlock("test")
}

func TestNonBlockingLocks(t *testing.T) {
	l, err := lock.New("mongo", "mongodb://localhost/mydb")
	assert.Nil(t, err)
	defer l.Close()
	l2, err := lock.New("mongo", "mongodb://localhost/mydb")
	assert.Nil(t, err)
	defer l2.Close()

	var lock1Chan = make(chan struct{})
	var lock1Result = make(chan bool)
	var lock2Chan = make(chan struct{})
	var lock2Result = make(chan bool)

	go func() {
		<-lock1Chan
		lock1Result <- l.NonBlockLock("test")
	}()

	go func() {
		<-lock2Chan
		lock2Result <- l2.NonBlockLock("test")
	}()

	lock1Chan <- struct{}{}
	time.Sleep(250 * time.Millisecond)
	lock2Chan <- struct{}{}

	assert.True(t, <-lock1Result)
	assert.False(t, <-lock2Result)
	l.Unlock("test")
	l2.Unlock("test")
}
