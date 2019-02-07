package callbacks

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type SequentialCallbackManager struct {
	sync.Mutex

	callbacks *[]callbackState
	reverse   bool
}

type callbackState struct {
	cb             callback
	pendingRemoval uint32
}

func NewSequentialCallbackManager() *SequentialCallbackManager {
	callbacks := make([]callbackState, 0)

	return &SequentialCallbackManager{
		reverse:   false,
		callbacks: &callbacks,
	}
}

func (m *SequentialCallbackManager) UnsafelySetReverse() *SequentialCallbackManager {
	m.reverse = true
	return m
}

func (m *SequentialCallbackManager) loadCallbacks() []callbackState {
	return *(*[]callbackState)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.callbacks))))
}

func (m *SequentialCallbackManager) storeCallbacks(callbacks []callbackState) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&m.callbacks)), unsafe.Pointer(&callbacks))
}

func (m *SequentialCallbackManager) Trim() {
	m.Lock()
	newCallbacks := make([]callbackState, 0)
	for _, cb := range m.loadCallbacks() {
		if cb.pendingRemoval == 0 {
			newCallbacks = append(newCallbacks, cb)
		}
	}
	m.storeCallbacks(newCallbacks)
	m.Unlock()
}

func (m *SequentialCallbackManager) RegisterCallback(c callback) {
	m.Lock()
	m.storeCallbacks(append(m.loadCallbacks(), callbackState{
		cb: c,
	}))
	m.Unlock()
}

// RunCallbacks runs all callbacks on a variadic parameter list, and de-registers callbacks
// that throw an error.
func (m *SequentialCallbackManager) RunCallbacks(params ...interface{}) (errs []error) {
	callbacks := m.loadCallbacks()
	if m.reverse {
		for i := len(callbacks) - 1; i >= 0; i-- {
			c := &callbacks[i]
			if err := c.run(params...); err != nil {
				errs = append(errs, err)
			}
		}
	} else {
		for i := 0; i < len(callbacks); i++ {
			c := &callbacks[i]
			if err := c.run(params...); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (c *callbackState) run(params ...interface{}) error {
	if atomic.LoadUint32(&c.pendingRemoval) == 0 {
		err := c.cb(params...)
		if err != nil {
			atomic.StoreUint32(&c.pendingRemoval, 1)
			if err != DeregisterCallback {
				return err
			}
		}
	}

	return nil
}
