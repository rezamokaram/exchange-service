# Exchange Service

## Introduction
**Exchange Service** is a monolithic REST API project implemented in Go using the fiber framework for the project layout.


## Setup
### Docker
First install the dependencies, then:
- docker
- Clone via Git (or provide the source code in some other way).
1. **Clone the Repository**
   ```bash
   git clone https://github.com/rezamokaram/exchange-service.git
   ```

2. **Navigate to the Project Directory**
   ```bash
   cd exchange-service
   ```

3. **Start the Project**
   ```bash
   sudo docker-compose up -d --build
   ```

4. **Access the endpoints** at `localhost:8080`.

### Helm Chart
coming soon!

## Documentation
Api documentations provided in `docs` directory.
   
## Testing

Run the following command in the root directory of the project:
```bash
go test ./test
```