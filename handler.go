package rdgo

type Handler interface {
	Input(Data) error
	GetLast() Data
	Destroy()
	OnDestroy(func())
}

type Data interface{}
