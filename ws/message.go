package ws

const (
	OpAdd     int = 1
	OpRemove  int = 2
	OpUpdate  int = 3
	OpError   int = 4
	OpSuccess int = 5
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

type AlarmMes struct {
	Sub       string `json:"sub"`
	Op        int    `json:"op"`
	DeviceKey string `json:"device_key"`
	Alarm     *int   `json:"alarm"`
}

type DeviceMes struct {
	Sub       string `json:"sub"`
	Op        int    `json:"op"`
	DeviceKey string `json:"device_key"`
	Name      string `json:"name"`
	Alarm     *int   `json:"alarm"`
}

type ErrMes struct {
	Sub     string `json:"sub"`
	Op      int    `json:"op"`
	Message string `json:"message"`
}

type SuccessMes struct {
	Sub string `json:"sub"`
	Op  int    `json:"op"`
}
