package processor

import "strconv"

type CommonData struct {
	Battery     int `json:"battery"`
	Voltage     int `json:"voltage"`
	LinkQuality int `json:"linkquality"`
}

/*{
  "device": {
    "applicationVersion": 2,
    "dateCode": "20160516",
    "friendlyName": "0x00158d0003a0026c",
    "hardwareVersion": 30,
    "ieeeAddr": "0x00158d0003a0026c",
    "manufacturerID": 4151,
    "manufacturerName": "LUMI",
    "model": "WSDCGQ01LM",
    "networkAddress": 7545,
    "powerSource": "Battery",
    "softwareBuildID": "3000-0001",
    "stackVersion": 2,
    "type": "EndDevice",
    "zclVersion": 1
  },
  "humidity": 33.66,
  "linkquality": 159,
  "temperature": 25.98
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
	Contact     *bool    `json:"contact"`
	Temperature *float64 `json:"temperature"`
}

func (data *WindowsSwitch) Present() bool {
	return data.Contact != nil && data.Temperature != nil
}

func (data *WindowsSwitch) Messages(root string) (msgs []*Message) {
	if !data.Present() {
		return nil
	}
	msg := &Message{Topic: root + "_contact/state"}
	if *data.Contact {
		msg.Payload = []byte("CLOSED")
	} else {
		msg.Payload = []byte("OPEN")
	}
	msgs = append(msgs, msg)
	msg = &Message{
		Topic:   root + "_temperature/state",
		Payload: []byte(strconv.FormatFloat(*data.Temperature, 'f', 2, 64)),
	}
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
