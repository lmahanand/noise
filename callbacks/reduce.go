package callbacks

import "sync"

type ReduceCallbackManager struct {
	*sync.Mutex

	callbacks []*reduceCallback
	reverse   bool
}

func NewReduceCallbackManager(mu *sync.Mutex) *ReduceCallbackManager {
	r := &ReduceCallbackManager{
		reverse: false,
	}

	if mu == nil {
		r.Mutex = new(sync.Mutex)
	} else {
		r.Mutex = mu
	}

	return r
}

func (m *ReduceCallbackManager) Reverse() *ReduceCallbackManager {
	m.reverse = true
	return m
}

func (m *ReduceCallbackManager) RegisterCallback(c reduceCallback) {
	m.Lock()

	m.UnsafeRegisterCallback(c)

	m.Unlock()
}

func (m *ReduceCallbackManager) UnsafeRegisterCallback(c reduceCallback) {
	if m.reverse {
		m.callbacks = append([]*reduceCallback{&c}, m.callbacks...)
	} else {
		m.callbacks = append(m.callbacks, &c)
	}
}

// RunCallbacks runs all callbacks on a variadic parameter list, and de-registers callbacks
// that throw an error.
func (m *ReduceCallbackManager) RunCallbacks(in interface{}, params ...interface{}) (res interface{}, errs []error) {
	m.Lock()

	var remaining []*reduceCallback
	var err error

	for _, c := range m.callbacks {
		if in, err = (*c)(m.UnsafeRegisterCallback, in, params...); err != nil {
			if err != DeregisterCallback {
				errs = append(errs, err)
			}
		} else {
			remaining = append(remaining, c)
		}
	}
	m.callbacks = remaining

	m.Unlock()

	return in, errs
}

// MustRunCallbacks runs all callbacks on a variadic parameter list, and de-registers callbacks
// that throw an error. Errors are ignored.
func (m *ReduceCallbackManager) MustRunCallbacks(in interface{}, params ...interface{}) interface{} {
	out, _ := m.RunCallbacks(in, params...)
	return out
}
