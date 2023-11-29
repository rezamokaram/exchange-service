# QExchange

## Introduction
**QExchange** is an advanced digital currency exchange platform designed to streamline the process of trading cryptocurrencies. Built with the latest industry standards, this platform emphasizes clean, secure coding practices to ensure a safe trading environment, crucial for maintaining financial integrity and trust.

## Technology Stack
- Docker
- Golang 1.21.3
- Echo 4.11.3
- Gorm 1.25.5
- PostgreSQL 16.1
- echo-swagger 1.4.1
- testify
- sqlite

## Setup
First install the dependencies, then:
1. **Clone the Repository**
   ```bash
   git clone https://github.com/Quera-Go-Zilla/QExchange-System.git
   ```

2. **Navigate to the Project Directory**
   ```bash
   cd QExchange-System
   ```

3. **Start the Project**
   ```bash
   sudo docker-compose up --build
   ```

4. **Access the endpoints** at `localhost:8080`.

## Documentation
To generate the docs run this in the root directory:
   ```bash
   swag init --parseDependency
   ```
Then start the project using docker-compose and access the docs at this url:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
   
## Testing
Run the following command in the root directory of the project:
```bash
go test ./test
```