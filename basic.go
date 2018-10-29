package rdgo

import "sync"

type BasicHandler struct {
	sync.Mutex
	destroylisteners []func()
	lastData         Data
}

func (h *BasicHandler) SetLast(data Data) {
	h.Lock()
	h.lastData = data
	h.Unlock()
}

func (h *BasicHandler) GetLast() Data {
	h.Lock()
	defer h.Unlock()
	return h.lastData
}

func (h *BasicHandler) Destroy() {
	h.Lock()
	for _, cb := range h.destroylisteners {
		cb()
	}
	h.Unlock()
}

func (h *BasicHandler) OnDestroy(cb func()) {
	h.Lock()
	h.destroylisteners = append(h.destroylisteners, cb)
	h.Unlock()
}
