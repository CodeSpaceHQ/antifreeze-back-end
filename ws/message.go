package ws

const (
	Add    int = 1
	Remove int = 2
	Update int = 3
)

// support multiple message types
type Message interface{}

type TempMes struct {
	Sub       string
	Op        int
	DeviceKey string
	Temp      int
	Date      int64
}

type ErrMes struct {
	Message string
}
