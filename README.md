# Oura Dashboard - API v2

A fully interactive web application to visualize and analyze your health data from the Oura Ring using the **Oura API v2**.

## Features

- **Multi-Tab Dashboard**: Switch between Sleep, Activity, Readiness, and Heart Rate views
- **Interactive Charts**: Line charts, bar charts, and radar charts powered by Chart.js
- **Comprehensive Data**: Access all Oura API v2 endpoints in one dashboard
- **Date Range Selection**: Custom dates or quick select buttons (7, 30, 90 days)
- **Modern Dark Theme**: Responsive design that looks great on all devices
- **Color-Coded Scores**: Visual indicators for score quality (excellent, good, fair, poor)

## Data Types

| Tab | Oura API v2 Endpoint | Data Included |
|-----|---------------------|---------------|
| **Sleep** | `/daily_sleep`, `/sleep` | Sleep scores, contributors, sleep stages (deep, light, REM), duration |
| **Activity** | `/daily_activity` | Activity scores, steps, calories, activity levels, contributors |
| **Readiness** | `/daily_readiness` | Readiness scores, HRV, recovery index, temperature deviation |
| **Heart Rate** | `/heartrate` | 5-minute interval heart rate data, min/max/avg BPM |

## Installation

### Prerequisites
- Go 1.16 or later (for embed support)
- An Oura Ring account with API access
- Active Oura Membership (required for Gen3/Ring 4 users)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/piercecohen1/Oura-Client.git
cd Oura-Client
```

2. Set your Oura API token as an environment variable:
```bash
export OURA_TOKEN="your_personal_access_token"
```

You can get your Personal Access Token from the [Oura Developer Portal](https://cloud.ouraring.com/personal-access-tokens).

3. Build the application:
```bash
go build -o oura_client .
```

4. Run the server:
```bash
./oura_client
```

5. Open your browser and navigate to:
```
http://localhost:8080
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OURA_TOKEN` | Your Oura API Personal Access Token | Required |
| `PORT` | The port the server listens on | `8080` |

## API Endpoints

The server exposes several API endpoints that proxy to the Oura API v2:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Serves the web dashboard |
| `/api/dashboard` | GET | Returns all data types for the date range |
| `/api/sleep` | GET | Returns daily sleep scores |
| `/api/sleep/periods` | GET | Returns detailed sleep sessions |
| `/api/activity` | GET | Returns daily activity data |
| `/api/readiness` | GET | Returns daily readiness scores |
| `/api/heartrate` | GET | Returns heart rate data |
| `/api/personal` | GET | Returns user personal info |
| `/api/health` | GET | Health check endpoint |

### Query Parameters

All data endpoints accept:
- `start_date` (optional): Start date in YYYY-MM-DD format (defaults to 30 days ago)
- `end_date` (optional): End date in YYYY-MM-DD format (defaults to today)

## Technology Stack

- **Backend**: Go (standard library only)
- **Frontend**: HTML5, CSS3, Vanilla JavaScript
- **Charts**: Chart.js
- **API**: Oura Ring API v2

## Oura API v2 Documentation

For more information about the Oura API v2, visit:
- [Oura Developer Documentation](https://developer.ouraring.com/docs/api/v2/oura-api-documentation)
- [Oura API Help](https://support.ouraring.com/hc/en-us/articles/4415266939155-The-Oura-API)

## License

MIT License - See [LICENSE](LICENSE) for details.
