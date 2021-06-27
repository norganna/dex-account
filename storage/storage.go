package storage

import (
	"sync"
	"time"
)

// StorageFactory is supplied to register a storage class.
type StorageFactory interface {
	New() Storage
	Name() string
}

// Storage can add or get codes.
type Storage interface {
	AddCode(email string, exists bool) *Credential
	GetCode(code string) (cred *Credential, ok bool)
	Config() interface{}
}

// NewCredendtial creates a new credential with the supplied values.
func NewCredential(code, email string, exists bool) *Credential {
	return &Credential{
		code:   code,
		email:  email,
		exists: exists,
		until:  time.Now().Add(4*time.Hour),
		used:   false,
	}
}

var useMutex sync.Mutex

// Credential is returned when adding or retrieving codes from storage.
type Credential struct {
	code   string
	email  string
	exists bool
	until  time.Time
	used   bool
}

// Code returns the credential's code.
func (c *Credential) Code() string {
	return c.code
}

// Email returns the credential's email address.
func (c *Credential) Email() string {
	return c.email
}

// Exists returns whether the credential is for an existing account.
func (c *Credential) Exists() bool {
	return c.exists
}

// Valid returns whether the credential can be used.
func (c *Credential) Valid() bool {
	return c.ValidAt(time.Now())
}

// ValidAt returns whether the credential can be used at given time.
func (c *Credential) ValidAt(at time.Time) bool {
	return !c.used && c.until.After(at)
}

// Use consumes the credential (uses a global mutex lock to ensure no race).
func (c *Credential) Use() (ok bool) {
	useMutex.Lock()
	if c.Valid() {
		ok = true
		c.used = true // Stop it from being used again.
		c.until = time.Time{} // Mark for deletion at next clean.
	}
	useMutex.Unlock()
	return ok
}

var registry = struct {
	sync.RWMutex
	factories map[string]StorageFactory
} {
	factories: map[string]StorageFactory{},
}

// NewStore returns a storage of the given class or nil if not found.
func NewStore(name string) Storage {
	registry.RLock()
	f, ok := registry.factories[name]
	registry.RUnlock()

	if !ok {
		return nil
	}
	return f.New()
}

// RegisterStorage adds the given factory to the storage types.
func RegisterStorage(factory StorageFactory)  {
	registry.Lock()
	registry.factories[factory.Name()] = factory
	registry.Unlock()
}