# OneStepGPS / backend 
## Backend Installation

### Prerequisites
1.  Install  [Go](https://golang.org/dl/).
2.  Install PostgreSQL and set up a database.

### Install Dependencies
    go mod tidy
### Set Up Environment Variables
    ONESTEPGPS_API_KEY=your_onestepgps_api_key
    GOOGLE_MAPS_API_KEY=your_google_maps_api_key

    <!-- Be sure to use quotes here for the DSN -->
    DSN="host=my-postgres-db.lettersandwords.us-east-1.rds.amazonaws.com user=yourusername dbname=yourdatabase password=yourpass port=5432 sslmode=require TimeZone=America/Los_Angeles"
## Database Setup

The application uses  **GORM**  to manage database schema. Ensure your PostgreSQL database is up and running. The schema for  `UserPreferences`  will be migrated automatically on application start.
## Usage

### Run the Server

Start the server:

    go run main.go
    
By default, the server runs on  `http://localhost:8080`.

## API Endpoints

### **V1 API**

#### Get User Preferences

    `GET /api/v1/preferences/:userId`

-   **Description**: Retrieve preferences for a specific user.
-   **Response**: JSON object containing user preferences.

#### Update User Preferences

    `PUT /api/v1/preferences/:userId`

-   **Description**: Update preferences for a specific user.
-   **Request Body**: JSON object with preference updates.

#### Get Devices

    `GET /api/v1/devices`

-   **Description**: Fetch a list of devices.
-   **Response**: JSON array of devices.

### **V3 API**

#### Get Device Info

    `GET /v3/api/device-info`

-   **Description**: Retrieve detailed device information.
-   **Response**: JSON object with device details.

#### Drive Stop Route

    `GET /v3/api/route/drive-stop`

-   **Description**: Fetch driving and stopping route information.
-   **Response**: JSON object with route details.

## CORS Configuration

CORS middleware is configured to allow requests from specific origins for development and production. Adjust the  `AllowOrigins`  array in  `main.go`  to include your own domains if needed.

## File Structure
```
backend/
├── handlers/
│   └── handlers.go    # HTTP request handlers for devices and preferences
├── models/
│   └── models.go      # Data structures for devices and preferences
├── services/
│   └── services.go    # OneStepGPS API client and preferences storage
└── main.go           # Server setup and routing
```
## Development Tips

-   Use  `log.Printf`  for debugging and monitoring.
-   For custom migrations or database interactions, extend the  `models`  and  `services`  packages.
-   Update  `.env`  values as needed for production.
