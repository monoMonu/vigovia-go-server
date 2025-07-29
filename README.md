# Developement Setup Guide

## Prerequisites

- **Go 1.24.5 or later** - [Click for installation guide](https://go.dev/doc/install)
- **Git** - For cloning the repository

## Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/monoMonu/vigovia-go-server.git
   cd vigovia-go-server
   ```

2. **Download dependencies**
   ```bash
   go mod download
   ```

3. **Verify installation**
   ```bash
   go mod verify
   ```


4. **Install all dependencies (Cleanup + add missing)**
   ```bash
   go mod tidy
   ```

## Running the Application

### Development Mode

1. **Start the server**
   ```bash
   go run main.go
   ```

2. **Verify the server is running**
   
   Go to - http://localhost:3002/

   Expected response: `"Hello World"`

## API Documentation - Use Postman or any other tool to test below APIs

### Endpoints

#### Health Check
- **GET** `/` - Returns "Hello World" to verify server status

#### Generate PDF
- **POST** `/generate-itinerary` - Generates a PDF from travel data

**Request Body Example:**
```json
{
  "customerName": "John Doe",
  "destination": "Paris, France",
  "departureFrom": "New York, NY",
  "departureDate": "2024-06-15",
  "returnDate": "2024-06-22",
  "travelers": 2,
  "totalAmount": 2500.00,
  "installment1": 1250.00,
  "installment2": 1250.00,
  "days": [
    {
      "date": "2024-06-15",
      "activities": [
        {
          "time": "09:00",
          "title": "Arrival in Paris",
          "description": "Land at Charles de Gaulle Airport",
          "duration": 120,
          "type": "travel"
        }
      ]
    }
  ],
  "flights": [
    {
      "airline": "Air France",
      "flightNumber": "AF123",
      "departure": "JFK",
      "arrival": "CDG",
      "date": "2024-06-15",
      "time": "08:00"
    }
  ],
  "hotels": [
    {
      "name": "Hotel Le Marais",
      "address": "123 Rue de Rivoli, Paris",
      "checkIn": "2024-06-15",
      "checkOut": "2024-06-22",
      "nights": 7
    }
  ]
}
```
- Returns :
```
{
    "message": "PDF generated successfully",
    "url": "link to statically served pdf"
}
```

#### Static Files
- **GET** `/pdfs/*filepath` - Serves generated PDF files

#