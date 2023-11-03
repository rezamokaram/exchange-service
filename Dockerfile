# Build Stage
FROM golang:1.21.3 AS BuildStage
WORKDIR /
COPY . .
RUN go mod download
RUN go build -o /app
ENTRYPOINT ["/app"]

# Deploy Stage
# FROM scratch
# WORKDIR /
# COPY --from=BuildStage /app /
# EXPOSE 8080
# USER nonroot:nonroot
# ENTRYPOINT ["/app"]