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
	t.Parallel()
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
	lock2Chan <- struct{}{}
	select {
	case <-time.After(2 * time.Second):
		assert.Fail(t, "failed to aquire lock 1 after 2 seconds")
	case <-lock1Aquired:
	}

	l.Unlock("test")

	select {
	case <-time.After(2 * time.Second):
		assert.Fail(t, "failed to aquire lock 2 after 2 seconds")
	case <-lock2Aquired:
	}

	l2.Unlock("test")
}

func TestNonBlockingLocks(t *testing.T) {
	t.Parallel()
	l, err := lock.New("mongo", "mongodb://localhost/nbdb")
	assert.Nil(t, err)
	defer l.Close()
	l2, err := lock.New("mongo", "mongodb://localhost/nbdb")
	assert.Nil(t, err)
	defer l2.Close()

	var lock1Chan = make(chan struct{})
	var lock1Aquired = make(chan struct{})
	var lock1Result = make(chan bool, 1)
	var lock2Chan = make(chan struct{})
	var lock2Result = make(chan bool, 1)

	go func() {
		<-lock1Chan
		lock1Result <- l.NonBlockLock("non-block-test")
		lock1Aquired <- struct{}{}
	}()

	go func() {
		<-lock2Chan
		lock2Result <- l2.NonBlockLock("non-block-test")
	}()

	lock1Chan <- struct{}{}
	select {
	case <-time.After(2 * time.Second):
		assert.Fail(t, "Failed to aquire nonblocking lock after 2 seconds")
	case <-lock1Aquired:
	}
	lock2Chan <- struct{}{}

	assert.True(t, <-lock1Result)
	assert.False(t, <-lock2Result)
	l.Unlock("non-block-test")
	l2.Unlock("non-block-test")
}
