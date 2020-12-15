FROM golang:1.15-alpine3.12 AS build
WORKDIR /src
COPY . .
RUN go build -o /out/cookietest .

FROM alpine:3.12 AS bin
ARG JWT_SECRET
ENV JWT_SECRET ${JWT_SECRET:-"-secret"}
ENV ADDR=":5000"
ENV REMOTE="http://localhost:5000/redirpage"
COPY --from=build /out/cookietest /

EXPOSE 5000
ENTRYPOINT exec /cookietest -addr ${ADDR} -remote ${REMOTE}