FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY nemesis-backend .
EXPOSE 8080
CMD ["./nemesis-backend"]
