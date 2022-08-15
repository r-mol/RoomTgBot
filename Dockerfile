FROM golang as build
WORKDIR /app

# avoid redownloading without dependencies changing
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# copy sources
COPY ./cmd ./
# build standalone application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./

FROM busybox as production
WORKDIR /app
# copy built from previous stage
COPY --from=build /app/main ./
# define entypoint
CMD ["/app/main"]

