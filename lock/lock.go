package lock

// DistributedLocker defines the functions for aquiring distributed locks
type DistributedLocker interface {
	Lock(string)
	NonBlockLock(string) bool
	Unlock(string)
	Close()
}

type newFunc func(string) (DistributedLocker, error)

var Lockers = make(map[string]newFunc)

// New returns a distributed locker based on the type of locker passed in
//
// Currently go-distributed supports mysql, mongo, and postgres
func New(lockerType, connectionString string) (DistributedLocker, error) {
	f := Lockers[lockerType]
	locker, err := f(connectionString)
	if err != nil {
		return nil, err
	}

	return locker, nil
}
