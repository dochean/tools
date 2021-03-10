package bicrypto

import (
	"errors"
	"github.com/dochean/tools/chain"
	"github.com/gofrs/uuid"
	"sync"
	"time"
)

var (
	EXPIRE_DURATION = 30 * time.Minute

	NoneExistId = errors.New("No this ID in record.")
)

type BiManager struct {
	// map[id]bicrypto
	bm map[string]BiCryptor
	mu sync.RWMutex
	tm *chain.TimeChain
	d time.Duration
}

func NewBiManager() *BiManager {
	bi := &BiManager{
		bm: make(map[string]BiCryptor),
		tm: chain.NewTimeChain(),
		d: EXPIRE_DURATION,
	}
	go func() {
		for {
			select {
			case <-time.Tick(5 * time.Second):
				bi.Check()
			}
		}
	}()
	return bi
}

func (b *BiManager) Check() {
	b.mu.Lock()
	for b.tm.Len() > 0 && b.tm.PeekTime().Before(time.Now()) {
		id := b.tm.DeleteHead().(string)
		delete(b.bm, id)
	}
	b.mu.Unlock()
}

func (b *BiManager) Set(d time.Duration) {
	b.d = d
}

// client public key
func (b *BiManager) Add(pkey []byte) (string, []byte) {
	id := uuid.Must(uuid.NewV4()).String()
	bi := &Bicrypto{Id: id}
	bi.SetClientPublicKey(pkey)

	b.mu.Lock()
	defer b.mu.Unlock()
	b.bm[id] = bi
	b.tm.Add(b.d, id)
	serverKey := b.bm[id].MarshalServerPublicKey()
	return id, serverKey
}

func (b *BiManager) Encrypt(id string, data []byte) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	bi, ok := b.bm[id]
	if !ok {
		return nil, NoneExistId
	}
	res := bi.EncryptClientMessage(data)

	return res, nil
}

func (b *BiManager) Decrypt(id string, cipher []byte) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	bi, ok := b.bm[id]
	if !ok {
		return nil, NoneExistId
	}
	data := bi.DecryptServerMessage(cipher)

	return data, nil
}