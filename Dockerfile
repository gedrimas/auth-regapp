FROM golang:1.21.0

WORKDIR /project/auth-regapp/

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build -o /project/auth-regapp/build/auth .


EXPOSE 3001
ENTRYPOINT [ "/project/auth-regapp/build/auth" ]