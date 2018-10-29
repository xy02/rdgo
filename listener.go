package rdgo

type Listener struct {
	BasicHandler
	OnData func(Data)
}

func (h *Listener) Input(data Data) error {
	if h.OnData == nil {
		return nil
	}
	h.SetLast(data)
	h.OnData(data)
	return nil
}
