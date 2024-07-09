FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o assets-scheduler

FROM alpine:3.20.1
WORKDIR /app
COPY --from=build /app/site /app/site
COPY --from=build /app/assets-scheduler /app/assets-scheduler
CMD ["/app/assets-scheduler"]


