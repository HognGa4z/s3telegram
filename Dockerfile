FROM golang:latest
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
    
WORKDIR /app
COPY . .

RUN go build -ldflags "-s -w"
RUN rm -rf client bot util config s3
RUN rm -f *.go
RUN rm -f go.mod go.sum
RUN rm env README.md
RUN rm -f Dockerfile
RUN rm *.log
ENTRYPOINT ["/app/s3telegram"]
