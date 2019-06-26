FROM golang:1.12.3 AS BUILD

RUN mkdir /wfs-tiler
WORKDIR /wfs-tiler

ADD go.mod .
ADD go.sum .
RUN go mod download

#now build source code
ADD . ./
RUN go build -o /go/bin/wfs-tiler
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/wfs-tiler .


FROM golang:1.12.3

ENV WFS3_API_URL ''

COPY --from=BUILD /go/bin/* /bin/
ENTRYPOINT /bin/wfs-tiler

EXPOSE 3000

