package haskimail

// Suppression — запись стоп-листа.
type Suppression struct {
	EmailAddress      string `json:"EmailAddress"`
	SuppressionReason string `json:"SuppressionReason"`
	Origin            string `json:"Origin"`
	CreatedAt         *Time  `json:"CreatedAt"`
}

// Suppressions — список записей стоп-листа.
type Suppressions struct {
	Suppressions []Suppression `json:"Suppressions"`
}

// SuppressionEntry — адрес для добавления/удаления.
type SuppressionEntry struct {
	EmailAddress string `json:"EmailAddress"`
}

// SuppressionEntries — тело запроса на изменение стоп-листа.
type SuppressionEntries struct {
	Suppressions []SuppressionEntry `json:"Suppressions"`
}

// NewSuppressionEntries собирает тело запроса из списка адресов.
func NewSuppressionEntries(emails ...string) SuppressionEntries {
	e := SuppressionEntries{Suppressions: make([]SuppressionEntry, 0, len(emails))}
	for _, addr := range emails {
		e.Suppressions = append(e.Suppressions, SuppressionEntry{EmailAddress: addr})
	}
	return e
}

// SuppressionStatus — результат операции для одного адреса.
type SuppressionStatus struct {
	EmailAddress string `json:"EmailAddress"`
	Status       string `json:"Status"`
	Message      string `json:"Message"`
}

// SuppressionStatuses — результаты операции над стоп-листом.
type SuppressionStatuses struct {
	Suppressions []SuppressionStatus `json:"Suppressions"`
}
