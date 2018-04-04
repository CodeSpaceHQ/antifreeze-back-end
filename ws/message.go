package ws

const (
	Add    int = 1
	Remove int = 2
	Update int = 3
)

// support multiple message types
type Message interface {
	GetSub() string
	GetOp() int
}

type Temperature struct {
	Sub       string
	Op        int
	DeviceKey string
	Temp      int
	Date      int64
}

func (v Temperature) GetSub() string { return v.Sub }
func (v Temperature) GetOp() int     { return v.Op }
