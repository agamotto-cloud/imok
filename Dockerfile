FROM ubuntu

ARG serviceName=imok

ENV GIN_MODE=release                

EXPOSE 8080

WORKDIR /opt/${serviceName}

CMD [ "/opt/${serviceName}" ]