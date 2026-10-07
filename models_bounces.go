package haskimail

// Bounce — отказ доставки.
type Bounce struct {
	ID        int64  `json:"ID"`
	MessageID string `json:"MessageID"`
	ServerID  int64  `json:"ServerID"`
	Type      string `json:"Type"`
	// TypeCode — код причины возврата; API отдаёт его строкой.
	TypeCode      string `json:"TypeCode"`
	BouncedAt     *Time  `json:"BouncedAt"`
	DumpAvailable bool   `json:"DumpAvailable"`
	Email         string `json:"Email"`
	From          string `json:"From"`
	Subject       string `json:"Subject"`
	Content       string `json:"Content"`
	Inactive      bool   `json:"Inactive"`
	CanActivate   bool   `json:"CanActivate"`
	Name          string `json:"Name"`
	Tag           string `json:"Tag"`
	Description   string `json:"Description"`
	Details       string `json:"Details"`
	RecordType    string `json:"RecordType"`
	MessageStream string `json:"MessageStream"`
}

// Bounces — страница списка отказов.
type Bounces struct {
	TotalCount int      `json:"TotalCount"`
	Bounces    []Bounce `json:"Bounces"`
}

// BounceType — счётчик отказов одного типа.
type BounceType struct {
	Type  string `json:"Type"`
	Name  string `json:"Name"`
	Count int    `json:"Count"`
}

// DeliveryStats — сводка по доставке.
type DeliveryStats struct {
	InactiveMails int          `json:"InactiveMails"`
	Bounces       []BounceType `json:"Bounces"`
}
