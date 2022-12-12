# Oura-Client
I wrote this program to view my sleep data for any range of dates, using the Oura API. It was originally written in Python, but I wanted to learn Go, so I rewrote it in Go. In the future, I plan to run statistical analysis on the data.

## Installation
Clone the git repo:
```
git clone https://github.com/piercecohen1/Oura-Client.git
```

Replace <token> with your personal access token (PAT), on the following line:
```
"Authorization": "Bearer <token>",
```


Build Oura Client:
```
cd Oura-Client
go build ./oura_client.go
```

To run:
```
./oura_client
```
