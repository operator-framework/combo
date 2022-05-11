FROM gcr.io/distroless/static:debug

WORKDIR /

COPY combo ./bin/combo

EXPOSE 8080

ENTRYPOINT ["/bin/combo"]