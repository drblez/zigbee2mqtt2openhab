package processor

type CommonData struct {
	Battery     int `json:"battery"`
	Voltage     int `json:"voltage"`
	LinkQuality int `json:"linkquality"`
}

/*{
  "battery": 100,
  "contact": true,
  "device": {
    "applicationVersion": 3,
    "dateCode": "20161128",
    "friendlyName": "oh/office_window",
    "hardwareVersion": 2,
    "ieeeAddr": "0x00158d0001b91e14",
    "manufacturerID": 4151,
    "manufacturerName": "LUMI",
    "model": "MCCGQ11LM",
    "networkAddress": 54897,
    "powerSource": "Battery",
    "softwareBuildID": "3000-0001",
    "stackVersion": 2,
    "type": "EndDevice",
    "zclVersion": 1
  },
  "linkquality": 156,
  "temperature": 30,
  "voltage": 3035
}
*/

/* Дверь/окно
{
    "battery": 100,
    "contact": false,
    "temperature": 28,
    "voltage": 3035,
    "linkquality": 144
}
*/

type WindowsSwitch struct {
	Temperature interface{} `json:"temperature"`
	Contact     *bool       `json:"contact"`
}

func (data *WindowsSwitch) Present() bool {
	return data.Contact != nil //&& data.Temperature != nil
}

func (data *WindowsSwitch) Messages(root string) (msgs []*Message) {
	if !data.Present() {
		return nil
	}
	msg := &Message{Topic: root + "_contact/state/set"}
	if *data.Contact {
		msg.Payload = []byte("CLOSED")
	} else {
		msg.Payload = []byte("OPEN")
	}
	msgs = append(msgs, msg)
	return
}

/* Кнопка
{
    "battery": 100,
    "linkquality": 168,
    "voltage": 3035,
    "action": "single"
}
*/

type PushButton struct {
	Action *string `json:"action"`
}

func (data *PushButton) Present() bool {
	return data.Action != nil
}

func (data *PushButton) Messages(root string) (msgs []*Message) {
	if !data.Present() {
		return nil
	}
	topic := root + "_" + *data.Action + "/state/set"
	msgs = append(msgs, &Message{
		Topic:   topic,
		Payload: []byte("ON"),
	})
	msgs = append(msgs, &Message{
		Topic:   topic,
		Payload: []byte("OFF"),
	})
	return
}

/* Температура/влажность
{
    "linkquality": 141,
    "temperature": 25.4,
    "humidity": 39.79
}
*/

type HumidityTemperature struct {
	Humidity    *float64 `json:"humidity"`
	Temperature *float64 `json:"temperature"`
}

/* Датчик движения
{
    "illuminance": 6,
    "illuminance_lux": 6,
    "linkquality": 150,
    "occupancy": true
}
*/

type Motion struct {
	Illuminance    *int `json:"illuminance"`
	IlluminanceLux *int `json:"illuminance_lux"`
	Occupancy      *int `json:"occupancy"`
}

type DeviceData struct {
	CommonData
	WindowsSwitch
	PushButton
	HumidityTemperature
	Motion
}

func (data *DeviceData) Messages(root string) (msgs []*Message) {
	msgs = append(msgs, data.WindowsSwitch.Messages(root)...)
	msgs = append(msgs, data.PushButton.Messages(root)...)
	return
}
