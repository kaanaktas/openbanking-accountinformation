FROM scratch

COPY ./certs ./certs
COPY ./.env ./

# Copy the binary file
COPY ./account ./
COPY ./callback ./

EXPOSE 8080 8081

#services can be run with following commands
#docker run -p 8080:8080 kaktas/openbanking-accountinformation ./account
#docker run -p 8081:8081 kaktas/openbanking-accountinformation ./callback