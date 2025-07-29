package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/monoMonu/travel-itinerary-pdf/types"
	"github.com/monoMonu/travel-itinerary-pdf/utils"
)

func GeneratePDF(c *gin.Context) {
	var data types.BookingData
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	fileName, err := generatePDF(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}

	fileBase := strings.TrimPrefix(fileName, "./pdfs/")

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	pdfURL := fmt.Sprintf("%s://%s/pdfs/%s", scheme, c.Request.Host, fileBase)

	c.JSON(http.StatusOK, gin.H{
		"message": "PDF generated successfully",
		"url":     pdfURL,
	})
}

func generatePDF(data types.BookingData) (string, error) {
	error := os.RemoveAll("./pdfs")
	if error != nil {
		log.Println("Couldn't clean /pdfs dir")
	}

	err := os.MkdirAll("./pdfs", os.ModePerm)
	if err != nil {
		return "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)

	// PAGE 1:
	pdf.AddPage()

	addFooterToAllPages(pdf)

	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(0, 0, 210, 297, "F")

	pdf.SetTextColor(107, 70, 193)
	pdf.SetFont("Arial", "B", 28)
	pdf.SetY(25)
	pdf.CellFormat(0, 15, "vigovia", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 8, "PLAN.PACK.GO", "", 1, "C", false, 0, "")

	pdf.SetFillColor(107, 70, 193)
	pdf.RoundedRect(15, 55, 180, 50, 5, "1234", "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 18)
	pdf.SetY(65)
	pdf.CellFormat(0, 10, fmt.Sprintf("Hi, %s!", data.CustomerName), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 22)
	pdf.CellFormat(0, 12, fmt.Sprintf("%s Itinerary", data.Destination), "", 1, "C", false, 0, "")

	nights := utils.CalculateNights(data.DepartureDate, data.ReturnDate)
	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(0, 10, fmt.Sprintf("%d Days %d Nights", nights+1, nights), "", 1, "C", false, 0, "")

	pdf.SetFillColor(248, 250, 252)
	pdf.SetDrawColor(220, 220, 220)
	pdf.RoundedRect(15, 115, 180, 60, 5, "1234", "FD")

	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 10)

	details := [][]string{
		{"Departure From", data.DepartureFrom},
		{"Departure", utils.FormatDate(data.DepartureDate)},
		{"Arrival", utils.FormatDate(data.ReturnDate)},
		{"Destination", data.Destination},
		{"No. Of Travellers", fmt.Sprintf("%d", data.Travelers)},
	}

	y := 125.0
	for _, detail := range details {
		pdf.SetXY(25, y)
		pdf.SetFont("Arial", "B", 9)
		pdf.Cell(40, 6, detail[0]+":")
		pdf.SetFont("Arial", "", 9)
		pdf.Cell(80, 6, detail[1])
		y += 8
	}

	// PAGE 2
	pdf.AddPage()

	addPageHeader(pdf)

	pdf.SetY(40)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Daily Itinerary")
	pdf.Ln(20)

	pageHeight := 297.0
	topMargin := 15.0
	bottomMargin := 20.0
	usablePageHeight := pageHeight - topMargin - bottomMargin

	for i, day := range data.Days {
		estimatedHeight := 25 + len(day.Activities)*25

		currentY := pdf.GetY()

		log.Println(estimatedHeight, usablePageHeight, currentY)

		if currentY+float64(estimatedHeight) > usablePageHeight {
			pdf.AddPage()
			addPageHeader(pdf)
			pdf.SetY(40)
		}

		dayY := pdf.GetY()

		pdf.SetFillColor(63, 45, 123)
		pdf.Circle(25, dayY+15, 12, "F")

		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "B", 12)
		pdf.SetXY(18, dayY+10)
		pdf.Cell(10, 10, fmt.Sprintf("Day\n%d", i+1))

		pdf.SetXY(45, dayY+10)
		pdf.SetTextColor(55, 65, 81)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, utils.FormatDate(day.Date))
		pdf.Ln(6)
		pdf.SetX(45)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(0, 6, "Arrival in "+data.Destination+" & City Exploration")
		pdf.Ln(10)

		timelineX := 120.0
		activityY := dayY + 10

		for j, activity := range day.Activities {
			pdf.SetFillColor(107, 70, 193)
			pdf.Circle(timelineX, activityY, 3, "F")

			if j < len(day.Activities)-1 {
				pdf.SetDrawColor(200, 200, 200)
				pdf.Line(timelineX, activityY+3, timelineX, activityY+20)
			}

			// Activity details
			pdf.SetXY(timelineX+8, activityY-3)
			pdf.SetTextColor(55, 65, 81)
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(0, 5, activity.Time)
			pdf.Ln(5)
			pdf.SetX(timelineX + 8)
			pdf.SetFont("Arial", "", 8)
			pdf.MultiCell(65, 4, "- "+activity.Description, "", "L", false)

			activityY += 25
		}

		pdf.SetY(dayY + 60)
	}

	// PAGE 3
	pdf.AddPage()
	addPageHeader(pdf)

	pdf.SetY(40)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Flight Summary")
	pdf.Ln(15)

	for _, flight := range data.Flights {
		pdf.SetFillColor(248, 250, 252)
		pdf.RoundedRect(15, pdf.GetY(), 180, 15, 3, "1234", "F")

		pdf.SetY(pdf.GetY() + 3)
		pdf.SetX(25)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(100, 100, 100)
		pdf.Cell(40, 8, utils.FormatDate(flight.Date))

		pdf.SetTextColor(55, 65, 81)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 8, fmt.Sprintf("Fly %s From %s (%s) To %s (%s).",
			flight.Airline, flight.From, "DEL", flight.To, "SIN"))
		pdf.Ln(18)
	}

	pdf.Ln(5)
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(100, 100, 100)
	pdf.Cell(0, 5, "Note: All Flights Include Meals, Seat Choice (Excluding XL), And 20kg/25Kg Checked Baggage.")
	pdf.Ln(15)

	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Hotel Bookings")
	pdf.Ln(15)

	pdf.SetFillColor(63, 45, 123)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"City", "Check In", "Check Out", "Nights", "Hotel Name"}
	widths := []float64{25, 25, 25, 15, 90}

	x := 15.0
	for i, header := range headers {
		pdf.SetXY(x, pdf.GetY())
		pdf.CellFormat(widths[i], 8, header, "1", 0, "C", true, 0, "")
		x += widths[i]
	}
	pdf.Ln(8)

	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "", 8)
	for i, hotel := range data.Hotels {
		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		x = 15.0
		values := []string{
			hotel.City,
			utils.FormatDate(hotel.CheckIn),
			utils.FormatDate(hotel.CheckOut),
			fmt.Sprintf("%d", hotel.Nights),
			hotel.Name,
		}

		for j, value := range values {
			pdf.SetXY(x, pdf.GetY())
			pdf.CellFormat(widths[j], 6, value, "1", 0, "C", true, 0, "")
			x += widths[j]
		}
		pdf.Ln(6)
	}

	addNotesPage(pdf)
	addServiceScopePage(pdf)
	addActivityTablePage(pdf, data)
	addPaymentPage(pdf, data)

	fileName := fmt.Sprintf("./pdfs/%s_%s_itinerary_%d.pdf",
		utils.SanitizeFileName(data.CustomerName),
		utils.SanitizeFileName(data.Destination),
		time.Now().Unix())
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func addPageHeader(pdf *gofpdf.Fpdf) {
	pdf.SetTextColor(107, 70, 193)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetY(15)
	pdf.CellFormat(0, 8, "vigovia", "", 0, "L", false, 0, "")

	pdf.SetTextColor(100, 100, 100)
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 8, "PLAN.PACK.GO", "", 1, "R", false, 0, "")

	pdf.SetDrawColor(220, 220, 220)
	pdf.Line(15, 25, 195, 25)
}

func addNotesPage(pdf *gofpdf.Fpdf) {
	pdf.AddPage()
	addPageHeader(pdf)

	pdf.SetY(40)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Important Notes")
	pdf.Ln(15)

	notes := [][]string{
		{"Airlines Standard Policy", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Flight/Hotel Cancellation", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Trip Insurance", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Hotel Check-in & Check Out", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
		{"Visa Rejection", "In Case Of Visa Rejection, Visa Fees Or Any Other Non Cancellable Component Cannot Be Reimbursed At Any Cost."},
	}

	pdf.SetFillColor(63, 45, 123)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 10, "Point", "1", 0, "C", true, 0, "")
	pdf.CellFormat(130, 10, "Details", "1", 1, "C", true, 0, "")

	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "", 9)

	const lineHeight = 6.0
	const paddingTop = 3.0
	const paddingBottom = 3.0
	const horizontalPadding = 4.0
	const minRowHeight = 12.0
	x := 15.0

	for i, note := range notes {
		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		y := pdf.GetY()
		width0, width1 := 50.0, 130.0

		lines0 := pdf.SplitLines([]byte(note[0]), width0-horizontalPadding)
		lines1 := pdf.SplitLines([]byte(note[1]), width1-horizontalPadding)
		maxLines := max(len(lines0), len(lines1))
		cellHeight := float64(maxLines)*lineHeight + paddingTop + paddingBottom
		if cellHeight < minRowHeight {
			cellHeight = minRowHeight
		}

		pdf.Rect(x, y, width0, cellHeight, "F")
		pdf.SetXY(x+horizontalPadding/2, y+paddingTop)
		pdf.MultiCell(width0-horizontalPadding, lineHeight, note[0], "", "C", false)

		pdf.Rect(x+width0, y, width1, cellHeight, "F")
		pdf.SetXY(x+width0+horizontalPadding/2, y+paddingTop)
		pdf.MultiCell(width1-horizontalPadding, lineHeight, note[1], "", "C", false)

		pdf.SetY(y + cellHeight)
	}
}

func addServiceScopePage(pdf *gofpdf.Fpdf) {
	pdf.SetY(pdf.GetY() + 10)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Scope Of Service")
	pdf.Ln(15)

	services := [][]string{
		{"Flight Tickets And Hotel Vouchers", "Delivered 3 Days Post Full Payment"},
		{"Web Check-In", "Boarding Pass Delivery Via Email/WhatsApp"},
		{"Support", "Chat Support - Response Time: 4 Hours"},
		{"Cancellation Support", "Provided"},
		{"Trip Support", "Response Time: 5 Minutes"},
	}

	pdf.SetFillColor(63, 45, 123)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(60, 10, "Service", "1", 0, "C", true, 0, "")
	pdf.CellFormat(120, 10, "Details", "1", 1, "C", true, 0, "")

	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "", 9)

	const lineHeight = 6.0
	const paddingTop = 3.0
	const paddingBottom = 3.0
	const horizontalPadding = 4.0
	const minRowHeight = 12.0

	for i, service := range services {
		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		x := pdf.GetX()
		y := pdf.GetY()
		width0, width1 := 60.0, 120.0

		lines0 := pdf.SplitLines([]byte(service[0]), width0-horizontalPadding)
		lines1 := pdf.SplitLines([]byte(service[1]), width1-horizontalPadding)
		maxLines := max(len(lines0), len(lines1))
		rowHeight := float64(maxLines)*lineHeight + paddingTop + paddingBottom
		if rowHeight < minRowHeight {
			rowHeight = minRowHeight
		}

		pdf.Rect(x, y, width0, rowHeight, "F")
		pdf.SetXY(x+horizontalPadding/2, y+paddingTop)
		pdf.MultiCell(width0-horizontalPadding, lineHeight, service[0], "", "C", false)

		pdf.Rect(x+width0, y, width1, rowHeight, "F")
		pdf.SetXY(x+width0+horizontalPadding/2, y+paddingTop)
		pdf.MultiCell(width1-horizontalPadding, lineHeight, service[1], "", "C", false)

		pdf.SetXY(x, y+rowHeight)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func addActivityTablePage(pdf *gofpdf.Fpdf, data types.BookingData) {
	pdf.AddPage()
	addPageHeader(pdf)

	pdf.SetY(40)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Activity Table")
	pdf.Ln(15)

	headers := []string{"City", "Activity", "Type", "Time Required"}
	widths := []float64{35, 80, 35, 30}

	pdf.SetFillColor(63, 45, 123)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)

	x := 15.0
	for i, header := range headers {
		pdf.SetXY(x, pdf.GetY())
		pdf.CellFormat(widths[i], 8, header, "1", 0, "C", true, 0, "")
		x += widths[i]
	}
	pdf.Ln(8)

	var activities [][]string
	for _, day := range data.Days {
		for _, activity := range day.Activities {
			activityRow := []string{
				data.Destination,
				activity.Title,
				activity.Type,
				utils.FormatDuration(activity.Duration, activity.Time),
			}
			activities = append(activities, activityRow)
		}
	}

	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "", 9)

	currentY := pdf.GetY()

	const lineHeight = 5.0
	const paddingTop = 3.0
	const paddingBottom = 3.0
	const horizontalPadding = 4.0
	const minRowHeight = 12.0
	const pageBottomMargin = 270.0

	maxFloat := func(a, b float64) float64 {
		if a > b {
			return a
		}
		return b
	}

	for i, activity := range activities {

		heights := []float64{}
		for j, colText := range activity {
			wrappedWidth := widths[j] - horizontalPadding
			lines := pdf.SplitLines([]byte(colText), wrappedWidth)
			height := float64(len(lines))*lineHeight + paddingTop + paddingBottom
			heights = append(heights, height)
		}

		rowHeight := minRowHeight
		for _, h := range heights {
			rowHeight = maxFloat(rowHeight, h)
		}
		if currentY+rowHeight > pageBottomMargin {
			pdf.AddPage()
			addPageHeader(pdf)
			pdf.SetY(40)

			pdf.SetFillColor(63, 45, 123)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("Arial", "B", 10)
			x = 15.0
			for j, header := range headers {
				pdf.SetXY(x, pdf.GetY())
				pdf.CellFormat(widths[j], 8, header, "1", 0, "C", true, 0, "")
				x += widths[j]
			}
			pdf.Ln(8)

			pdf.SetTextColor(55, 65, 81)
			pdf.SetFont("Arial", "", 9)
			currentY = pdf.GetY()
		}

		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.SetXY(15, currentY)
		pdf.Rect(15, currentY, 180, rowHeight, "F")

		x = 15.0
		for j, value := range activity {
			pdf.SetXY(x, currentY)
			pdf.Rect(x, currentY, widths[j], rowHeight, "D")

			if j == 1 {
				pdf.SetXY(x+horizontalPadding/2, currentY+paddingTop)
				pdf.MultiCell(
					widths[j]-horizontalPadding,
					lineHeight,
					value,
					"",
					"C",
					false,
				)
			} else {
				textY := currentY + (rowHeight-lineHeight)/2
				pdf.SetXY(x, textY)
				pdf.CellFormat(widths[j], lineHeight, value, "", 0, "C", false, 0, "")
			}
			x += widths[j]
		}

		currentY += rowHeight
		pdf.SetY(currentY)
	}

	pdf.Ln(15)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Terms and Conditions")
	pdf.Ln(10)

	pdf.SetTextColor(107, 70, 193)
	pdf.SetFont("Arial", "U", 10)
	pdf.Cell(0, 8, "View all terms and conditions")
}

func addPaymentPage(pdf *gofpdf.Fpdf, data types.BookingData) {
	pdf.AddPage()
	addPageHeader(pdf)

	pdf.SetY(40)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Payment Plan")
	pdf.Ln(20)

	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(15, pdf.GetY(), 180, 15, 3, "1234", "F")
	pdf.SetY(pdf.GetY() + 4)
	pdf.SetX(25)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 8, "Total Amount")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Rs.  %.0f For %d Pax (Inclusive Of GST)", data.TotalAmount, data.Travelers))
	pdf.Ln(20)

	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(15, pdf.GetY(), 180, 15, 3, "1234", "F")
	pdf.SetY(pdf.GetY() + 4)
	pdf.SetX(25)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 8, "TCS")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Not Collected")
	pdf.Ln(25)

	pdf.SetFillColor(63, 45, 123)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)
	headers := []string{"Installment", "Amount", "Due Date"}
	widths := []float64{60, 60, 60}

	x := 15.0
	for i, header := range headers {
		pdf.SetXY(x, pdf.GetY())
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", true, 0, "")
		x += widths[i]
	}
	pdf.Ln(10)

	installments := [][]string{
		{"Installment 1", fmt.Sprintf("Rs. %.0f", data.Installment1), "Initial Payment"},
		{"Installment 2", fmt.Sprintf("Rs. %.0f", data.Installment2), "Post Visa Approval"},
		{"Installment 3", "Remaining", "20 Days Before Departure"},
	}

	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "", 10)
	for i, installment := range installments {
		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		x = 15.0
		for j, value := range installment {
			pdf.SetXY(x, pdf.GetY())
			pdf.CellFormat(widths[j], 10, value, "1", 0, "C", true, 0, "")
			x += widths[j]
		}
		pdf.Ln(10)
	}

	pdf.Ln(15)
	pdf.SetTextColor(55, 65, 81)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Visa Details")
	pdf.Ln(15)

	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(15, pdf.GetY(), 180, 35, 5, "1234", "F")

	pdf.SetY(pdf.GetY() + 8)
	pdf.SetX(25)

	pdf.SetFont("Arial", "B", 11)

	pdf.Cell(40, 6, "Visa Type:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 6, "Tourist")

	pdf.Ln(8)

	pdf.SetX(25)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(30, 6, "Validity:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 6, "30 Days")

	pdf.Ln(8)

	pdf.SetX(25)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(40, 6, "Processing Date :")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 6, "14/06/2025")

	pdf.Ln(20)
	pdf.SetTextColor(63, 45, 123)
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 15, "PLAN.PACK.GO!", "", 1, "C", false, 0, "")

	rectHeight := 15.0
	rectY := pdf.GetY() + 5

	pdf.SetFillColor(63, 45, 123)
	pdf.RoundedRect(75, rectY, 60, rectHeight, 8, "1234", "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 12)

	textY := rectY + (rectHeight / 2) - 4

	pdf.SetY(textY)
	pdf.CellFormat(0, 8, "Book Now", "", 1, "C", false, 0, "")

}

func addFooterToAllPages(pdf *gofpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-20)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(100, 100, 100)

		pdf.SetX(15)
		pdf.Cell(60, 5, "Vigovia Tech Pvt. Ltd")
		pdf.Ln(4)
		pdf.SetX(15)
		pdf.Cell(60, 5, "Registered Office: Hd-109 Cinnabar Hills,")
		pdf.Ln(4)
		pdf.SetX(15)
		pdf.Cell(60, 5, "Links Business Park, Karnataka, India.")

		pdf.SetXY(120, -20)
		pdf.Cell(0, 5, "Phone: +91-99X9999999")
		pdf.SetXY(120, -16)
		pdf.Cell(0, 5, "Email ID: contact@Vigovia.Com")

		pdf.SetXY(170, -18)
		pdf.SetTextColor(107, 70, 193)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 5, "vigovia")
		pdf.SetXY(170, -14)
		pdf.SetTextColor(100, 100, 100)
		pdf.SetFont("Arial", "", 6)
		pdf.Cell(0, 5, "PLAN.PACK.GO")
	})
}
