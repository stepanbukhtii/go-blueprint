FROM alpine:latest

COPY ./build/example .

CMD ["./example"]