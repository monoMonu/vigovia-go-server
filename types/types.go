package types

type BookingData struct {
	CustomerName  string   `json:"customerName"`
	Destination   string   `json:"destination"`
	DepartureFrom string   `json:"departureFrom"`
	DepartureDate string   `json:"departureDate"` // changed to string
	ReturnDate    string   `json:"returnDate"`    // changed to string
	Travelers     int      `json:"travelers"`
	Days          []Day    `json:"days"`
	Flights       []Flight `json:"flights"`
	Hotels        []Hotel  `json:"hotels"`
	TotalAmount   float64  `json:"totalAmount"`
	Installment1  float64  `json:"installment1"`
	Installment2  float64  `json:"installment2"`
}

type Day struct {
	Date       string     `json:"date"` // changed to string
	Activities []Activity `json:"activities"`
}

type Activity struct {
	Time        string `json:"time"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Type        string `json:"type"`
}

type Flight struct {
	Date      string `json:"date"` // changed to string
	Airline   string `json:"airline"`
	From      string `json:"from"`
	To        string `json:"to"`
	Arrival   string `json:"arrival"`   // changed to string
	Departure string `json:"departure"` // changed to string
}

type Hotel struct {
	City     string `json:"city"`
	CheckIn  string `json:"checkIn"`  // changed to string
	CheckOut string `json:"checkOut"` // changed to string
	Nights   int    `json:"nights"`
	Name     string `json:"name"`
}
