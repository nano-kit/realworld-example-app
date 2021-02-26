FROM alpine
ADD realworld-service /realworld-service
ENTRYPOINT [ "/realworld-service" ]
