package shared

const (
	MMEra_Feb9 = 0
	MMEra_010  = 1
	MMERA_100  = 2
)

const (
	// none = make no changes to current
	BluetoothEra_None = 0
	BluetoothEra_V2   = 1
	BluetoothEra_V3   = 2
	BluetoothEra_V4   = 3
	BluetoothEra_V5   = 4
)

type Personality struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	ID           string `json:"id"`
	MMEra        int    `json:"mmera"`
	BluetoothEra int    `json:"bluetoothera"`
	// needs an extra service file
	CustomWireProgram bool `json:"customwireprogram"`
}
