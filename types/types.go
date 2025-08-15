package types

type BookingData struct {
	CustomerName  string   `json:"customerName"`
	Destination   string   `json:"destination"`
	DepartureFrom string   `json:"departureFrom"`
	DepartureDate string   `json:"departureDate"`
	ReturnDate    string   `json:"returnDate"`
	Travelers     int      `json:"travelers"`
	Days          []Day    `json:"days"`
	Flights       []Flight `json:"flights"`
	Hotels        []Hotel  `json:"hotels"`
	TotalAmount   float64  `json:"totalAmount"`
	Installment1  float64  `json:"installment1"`
	Installment2  float64  `json:"installment2"`
}

type Day struct {
	Date       string     `json:"date"`
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
	Date      string `json:"date"`
	Airline   string `json:"airline"`
	From      string `json:"from"`
	To        string `json:"to"`
	Arrival   string `json:"arrival"`
	Departure string `json:"departure"`
}

type Hotel struct {
	City     string `json:"city"`
	CheckIn  string `json:"checkIn"`
	CheckOut string `json:"checkOut"`
	Nights   int    `json:"nights"`
	Name     string `json:"name"`
}
