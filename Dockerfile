FROM alpine:latest
COPY main /main
EXPOSE 8080
CMD /main
