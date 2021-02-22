package domain

type Device struct {
	ID           int64 `db:"id" goqu:"defaultifempty"`
	Brand        Brand
	Name         string `db:"name"`
	Kind         DeviceKind
	ModelYear    int
	Size         float32
	PanelType    string
	Resolution   string
	PixelDensity int
}
