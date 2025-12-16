# Oura Sleep Dashboard

A fully interactive web application to visualize and analyze your sleep data from the Oura Ring API.

## Features

- **Interactive Dashboard**: Modern, responsive web interface with dark theme
- **Sleep Score Trends**: Line chart showing your sleep score over time
- **Contributor Analysis**: Radar chart displaying average sleep quality contributors
- **Daily Sleep Cards**: Detailed breakdown of each night's sleep with visual progress bars
- **Date Range Selection**: Pick custom date ranges or use quick select buttons (7, 30, 90 days)
- **Statistics Overview**: Average score, best night, lowest night, and days tracked
- **Color-Coded Scores**: Visual indicators for sleep quality (excellent, good, fair, poor)

## Screenshots

The dashboard includes:
- Summary statistics cards
- Sleep score trend line chart
- Average contributors radar chart
- Daily sleep cards with contributor breakdowns

## Installation

### Prerequisites
- Go 1.16 or later (for embed support)
- An Oura Ring account with API access

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

### Example

```bash
export OURA_TOKEN="your_token_here"
export PORT="3000"
./oura_client
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Serves the web dashboard |
| `/api/sleep` | GET | Returns sleep data for date range |
| `/api/health` | GET | Health check endpoint |

### Query Parameters for `/api/sleep`

- `start_date` (optional): Start date in YYYY-MM-DD format (defaults to 30 days ago)
- `end_date` (optional): End date in YYYY-MM-DD format (defaults to today)

## Technology Stack

- **Backend**: Go (standard library only)
- **Frontend**: HTML5, CSS3, Vanilla JavaScript
- **Charts**: Chart.js
- **API**: Oura Ring API v2

## License

MIT License - See [LICENSE](LICENSE) for details.
