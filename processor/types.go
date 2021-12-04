package processor

import "strconv"

type CommonData struct {
	Battery     int `json:"battery"`
	Voltage     int `json:"voltage"`
	LinkQuality int `json:"linkquality"`
}

func (data *CommonData) Messages(root string) (msgs []*Message) {
	msg := &Message{Topic: root + "_battery/state/set"}
	msg.Payload = []byte(strconv.Itoa(data.Battery))
	msgs = append(msgs, msg)
	msg = &Message{Topic: root + "_voltage/state/set"}
	msg.Payload = []byte(strconv.Itoa(data.Battery))
	msgs = append(msgs, msg)
	msg = &Message{Topic: root + "_linkquality/state/set"}
	msg.Payload = []byte(strconv.Itoa(data.Battery))
	msgs = append(msgs, msg)
	return
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
	Contact *bool `json:"contact"`
}

func (data *WindowsSwitch) Messages(root string) (msgs []*Message) {
	if data.Contact == nil {
		return
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

func (data *PushButton) Messages(root string) (msgs []*Message) {
	if data.Action == nil {
		return
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

type Humidity struct {
	Humidity *float64 `json:"humidity"`
}

/* Датчик движения
{
    "illuminance": 6,
    "illuminance_lux": 6,
    "linkquality": 150,
    "occupancy": true
}
*/

func (data *Humidity) Present() bool {
	return data.Humidity != nil
}

func (data *Humidity) Messages(root string) (msgs []*Message) {
	if data.Humidity == nil {
		return
	}
	msg := &Message{Topic: root + "_humidity/state/set"}
	msg.Payload = []byte(strconv.FormatFloat(*data.Humidity, 'f', 2, 64))
	msgs = append(msgs, msg)
	return
}

type Temperature struct {
	Temperature *float64 `json:"temperature"`
}

func (data *Temperature) Messages(root string) (msgs []*Message) {
	if data.Temperature == nil {
		return
	}
	msg := &Message{Topic: root + "_temperature/state/set"}
	msg.Payload = []byte(strconv.FormatFloat(*data.Temperature, 'f', 2, 64))
	msgs = append(msgs, msg)
	return
}

type Motion struct {
	Illuminance    *int  `json:"illuminance"`
	IlluminanceLux *int  `json:"illuminance_lux"`
	Occupancy      *bool `json:"occupancy"`
}

func (data *Motion) Message(root string) (msgs []*Message) {
	if data.Illuminance != nil {
		msg := &Message{Topic: root + "_illuminance/state/set"}
		msg.Payload = []byte(strconv.Itoa(*data.Illuminance))
		msgs = append(msgs, msg)
	}
	return
}

type DeviceData struct {
	CommonData
	WindowsSwitch
	PushButton
	Humidity
	Temperature
	Motion
}

func (data *DeviceData) Messages(root string) (msgs []*Message) {
	msgs = append(msgs, data.WindowsSwitch.Messages(root)...)
	msgs = append(msgs, data.PushButton.Messages(root)...)
	msgs = append(msgs, data.Humidity.Messages(root)...)
	msgs = append(msgs, data.Temperature.Messages(root)...)
	return
}
