package ws

const (
	OpAdd    int = 1
	OpRemove int = 2
	OpUpdate int = 3
	OpError  int = 4
)

// support multiple message types
type Message interface{}

type TempMes struct {
	Sub       string `json:"sub"`
	Op        int    `json:"op"`
	DeviceKey string `json:"device_key"`
	Temp      int    `json:"temp"`
	Date      int64  `json:"date"`
}

type ErrMes struct {
	Sub     string `json:"sub"`
	Op      int    `json:"op"`
	Message string `json:"message"`
}
