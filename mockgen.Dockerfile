FROM golang:1.23.3-alpine3.20

# RUN useradd -ms /bin/bash mockgen
# USER mockgen
WORKDIR /home/mockgen/src

RUN ls 

RUN go install go.uber.org/mock/mockgen@latest

# RUN go work use ./pkg/repository/factory

# RUN go generate ./...

ENTRYPOINT ["go", "generate", "-v", "./..."]