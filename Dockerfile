# Build Stage
FROM golang:1.21.3 AS BuildStage
WORKDIR /
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Deploy Stage
FROM scratch
WORKDIR /
COPY --from=BuildStage /app /
COPY --from=BuildStage /database/main-data.sql /
EXPOSE 8080
ENTRYPOINT ["./app"]