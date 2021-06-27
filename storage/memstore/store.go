package memstore

import (
	"sync"
	"time"

	"github.com/norganna/dex-account-api/rand"
	"github.com/norganna/dex-account-api/storage"
)

// Type aliases from storage
type (
	Storage = storage.Storage
	Credential = storage.Credential
)

func init() {
	storage.RegisterStorage(&factory{})
}

type factory struct {}

func (f *factory) New() storage.Storage {
	store := &memStore{
		cfg: &config{},
		codes: map[string]*Credential{},
	}

	go store.clean()

	return store
}

func (f *factory) Name() string {
	return "memstore"
}

type config struct {
	// Store config is under .store.*
	Store struct {
		// Class is the name of the desired storage class (us).
		Class string `mapstructure:"class"`
		// Other config vars go here.
	} `mapstructure:"store"`
}

type memStore struct {
	sync.RWMutex

	cfg *config
	codes map[string]*Credential
}

func (m *memStore) Config() interface{} {
	return m.cfg
}

func (m *memStore) AddCode(email string, exists bool) *Credential {
	m.Lock()
	var code string
	for ok := true; ok;  _, ok = m.codes[code] {
		code = rand.SafeString(9)
	}

	cred := storage.NewCredential(code, email, exists)
	m.codes[code] = cred
	m.Unlock()

	return cred
}

func (m *memStore) GetCode(code string) (cred *Credential, ok bool) {
	m.RLock()
	cred, ok = m.codes[code]
	m.RUnlock()

	if ok && !cred.Valid() {
		return nil, false
	}
	return cred, ok
}

func (m *memStore) clean() {
	timer := time.NewTicker(time.Hour)
	defer timer.Stop()
	for now := range timer.C {
		m.RLock()
		var clean []string
		for k, c := range m.codes {
			if !c.ValidAt(now) {
				clean = append(clean, k)
			}
		}
		m.RUnlock()
		if len(clean) > 0 {
			m.Lock()
			for _, k := range clean {
				delete(m.codes, k)
			}
			m.Unlock()
		}
	}
}
