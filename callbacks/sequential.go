package callbacks

import "sync"

type SequentialCallbackManager struct {
	*sync.Mutex

	callbacks []*callback
	reverse   bool
}

func NewSequentialCallbackManager(mu *sync.Mutex) *SequentialCallbackManager {
	s := &SequentialCallbackManager{
		reverse: false,
	}
	if mu == nil {
		s.Mutex = new(sync.Mutex)
	} else {
		s.Mutex = mu
	}
	return s
}

func (m *SequentialCallbackManager) Reverse() *SequentialCallbackManager {
	m.reverse = true
	return m
}

func (m *SequentialCallbackManager) RegisterCallback(c callback) {
	m.Lock()

	m.UnsafeRegisterCallback(c)

	m.Unlock()
}

func (m *SequentialCallbackManager) UnsafeRegisterCallback(c callback) {
	if m.reverse {
		m.callbacks = append([]*callback{&c}, m.callbacks...)
	} else {
		m.callbacks = append(m.callbacks, &c)
	}
}

// RunCallbacks runs all callbacks on a variadic parameter list, and de-registers callbacks
// that throw an error.
func (m *SequentialCallbackManager) RunCallbacks(params ...interface{}) (errs []error) {
	m.Lock()

	var remaining []*callback
	var err error

	for _, c := range m.callbacks {
		if err = (*c)(m.UnsafeRegisterCallback, params...); err != nil {
			if err != DeregisterCallback {
				errs = append(errs, err)
			}
		} else {
			remaining = append(remaining, c)
		}
	}
	m.callbacks = remaining

	m.Unlock()

	return
}
