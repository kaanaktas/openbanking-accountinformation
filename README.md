![openbanking-accountinformation-ci](https://github.com/kaanaktas/openbanking-accountinformation/workflows/openbanking-accountinformation-ci/badge.svg)

### **IMPORTANT**: DON'T FORGET TO HANDLE TOKENS, SECRETS, CERTIFICATES(public and private keys), CLIENT_IDS AND PASSWORDS SECURELY.

Introduction
------------

openbanking-accountinformation is a complete application to process account information journeys of Open Banking UK.
To be able to run these journeys, user needs;

1- to complete dynamic or manual registration process with each aspsp(configs added for **Danske** and **OzoneBank**)

2- to collect clientId and password(optional) and update in **config_table**


# Installation and usage

## Environment Properties


You can either use an .env file or export env variables to pass below values. It is recommended to keep keys in a secure environment rather than public access.

**.env** file:

```dotenv
#Port for account service
PORT=8080
#Port for callback service
PORT_CALLBACK=8081
#shared secret key to be used for signing internal jwts
INTERNAL_SIGN_KEY=<internal_signing.key>
#this key is used to sign JWT. It can be obtained from OpenBanking Portal  
OB_SIGN_KEY=<client_signing.key>
#KID value will be the same with your sign public cert. This can be obtained from 
KID=<CLIENT_SIGNING_KID>
#This can be the issuer of Open Banking certificate chain
CLIENT_CA_CERT_PEM=<ob_issuer.cer>
#Key pairs of TPP's transport and used for TLS-MA. This can be obtained from Open Banking portal. 
CLIENT_CERT_PEM=<client_transport.pem>
CLIENT_KEY_PEM=<client_transport.key>
#database driver
DRIVER_NAME=postgres
#database connection url
DATASOURCE_URL=postgres://postgres:postgres@localhost:5432/accountinformation?sslmode=disable
#database migration properties. version should be logical with sql files prefixes; 1_,2_ etc
MIGRATE_VERSION=2
MIGRATE_SCRIPT_URL=file://scripts/postgresql
MIGRATE_DATABASE_URL=postgres://postgres:postgres@localhost:5432/accountinformation?sslmode=disable
```

## Application Setup

### Database Migration

We can use either an in memory database or a relational database to store our data. Here, primarily Postgresql is used and for fast testing, sqlite is preferred.
Once you make ready your database, then migration process should be triggered.
- If you want to complete it on your local; run `/cmd/migrate/migrate.go`. Make sure you have correct values for env variables; **MIGRATE_VERSION, MIGRATE_SCRIPT_URL, MIGRATE_DATABASE_URL**.
  Once this is successfull, you need to see message `migration has been completed`, then check the database and the tables.


- If you want to run in your docker environment, docker-compose handles it directly. Please check `docker-compose.yml` for more information.


- You can also check the github action flow to understand how it is handled during the deployment flow.


### Account Service

Account service is the main service which exposes APIs and functions which interact with Open Banking UK and ASPSPs. It picks **PORT** from environment variables and can be built as a standalone application.

- to run on your local, go to /cmd/account/account.go and run/debug the go file.
- to build via command line, you can run `go build -o account ./cmd/account`
- or you can run it directly by calling `go run ./cmd/account`

### Callback Service

Callback service is used to handle redirected response from ASPSPs to complete consent authorization journey. It picks **PORT_CALLBACK** from environment variables and can be built as a standalone application.

- to run on your local, go to /cmd/callback/callback.go and run/debug the go file.
- to build via command line, you can run `go build -o callback ./cmd/callback`
- or you can run it directly by calling `go run ./cmd/callback`

# Docker

As there are 2 different applications in the platform as **account** and **callback**, we need to handle this either with in a single Dockerfile or separate file for each service.

### Docker in CI/CD pipeline

Here, we used one single configuration file, and we didn't use any command, so we need run each separate container for each service by passing its own name from the same image.
After the CI/CD pipeline builds account and callback in the **Build step**, the binaries is picked up and put into the Docker image in the next step.

Once your image is ready, the containers can be created with;

`docker run -p 8080:8080 <repo-domain-name>/openbanking-accountinformation ./account`

`docker run -p 8081:8081 <repo-domain-name>/openbanking-accountinformation ./callback`

### Dockerfile Standalone

We can use **Dockerfile.standalone** file, if we want to build the image in our local environment. Once the process is completed successfully, we can use the same commands above to run the containers.

### Docker Compose

Docker compose provides us a compact environment, which includes all dependencies with Postgresql and Redis. Once it runs successfully, we can see two services up as **openbanking-accountinformation** and **openbanking-accountinformation_callback**. Another difference here is, all environment variables pass to the system in the docker-compose.yml file instead of .env file 

# Services

### Initiate Session

Initiate Session is used to create a session on the application, and the reference number is returned with a valid token for 60 minutes by default. After the token expires, it is necessary to create a new token to access the system.

###### **Endpoint**

`{url}/internal/access/initiate/{user_name}/{tpp}/{tid}`

**`url`**: should be pointing your application's domain name and port number.

**`user_name`**: can be preferred if the same tpp has multiple users. Otherwise, this may be set with a static value.

**`tpp`**: points the tpp name. An organization can manage more than one tpp and this needs to be set with the name of each specific tpp which makes request to the system.

**`tid`**: is the unique number to retrieve a new token and reference number from the application. Tt is combined with **`tpp`** to provide uniqueness on the table, so the same tid can be used by different tpps.

**Example**

###### **Request**

>curl -v localhost:8080/internal/access/initiate/**MyUser**/**MyTpp**/**e12ed6e0-aec2- 489a-bee6-d24e4f2da3e0**

###### **Response**

>{
"**internal_access_token**": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk0Mjc3MDEsImlhdCI6MTYwOTQyNDEwMSwidGlkIjoiZTEyZWQ2ZTAtYWVjMi00ODlhLWJlZTYtZDI0ZTRmMmRhM2UwIiwidHBwSWQiOiJUcHBfMSJ9.1gGKp007QqcGA1ZMjogK1GL0koikqTfeRyFj57FmisU",
"**reference_id**": "c5e1d69000758b4ab3a368c4b446488cf2d96905"
}

### Create Consent

The first thing users need to do to manage their accounts is to allow TPP for the respective accounts. For this, the TPP has to call the Create Permission service and show the next ASPSP web page to its user. Thus, the user can log in through the web page and complete the requested permissions.


**Endpoint**

`{url}/{aspspId}/account-access-consents/reference/{reference}`

**`url`**: should be pointing your application's domain name and port number.

**`aspspId`**: needs to be set with the aspspId which service will be called from.

**`reference`**: is an internal **trackingId** and should be unique for each consent request. Once consentId is created on the system, authorization flow will try to match with this ID after redirection.


 **Example**

###### **Request**

>curl -d '{
"Data": {
"Permissions": [
"ReadAccountsBasic"
],
"ExpirationDateTime": "2022-12-31T21:35:00Z",
"TransactionFromDateTime": "2018-06-13T19:39:00Z",
"TransactionToDateTime": "2019-06-13T19:39:00Z"
},
"Risk": {}
}' -H 'Content-Type: application/json' -H 'Authorization: Bearer <internal_access_token>'
http://localhost:8080/danske/account-access-consents/reference/6541561651516165

###### **Response**

>https://sandbox-obp-web.danskebank.com/ui/authorize?original-url=https://sandbox-obp-api.danskebank.com...

The url which is retrieved after the request needs to be handled by the TPP in a way. The url will open a web page for the user to give consent for their accounts to the TPP.

### Retrieve Active Consents

This service retrieves authenticated consents for the selected ASPSP. Then, the consent reference(consentTid) of the consent needs to be selected and passed for each API calls.

**Endpoint**

`{url}/{aspspId}/internal/consent/active`

**`url`**: should be pointing your application's domain name and port number.

**`aspspId`**: needs to be set with the aspspId which service will be called from.


**Example**

###### **Request**

>curl -v -H 'Authorization: Bearer <internal_access_token>' http://localhost:8080/danske/internal/consent/active

###### **Response**

>[{\"consentTid\":10,\"aspspId\":\"danske\"}]

**ConsentId** parameter will be sent as a parameter in Open Banking API calls.


### Get Consent Details

This service returns the detail of selected consent from ASPSP directly.

**Endpoint**

`{url}/{aspspId}/account-access-consents/{consentTid}`

**`url`**: should be pointing your application's domain name and port number.

**`aspspId`**: needs to be set with the aspspId which service will be called from.

**`consentTid`**: needs to be set with a consentTid which is retrieved from Retrieve Active Consents.

**Example**

###### **Request**

>curl -v -H 'Authorization: Bearer <internal_access_token>' http://localhost:8080/danske/account-access-consents/10

###### **Response**

>"{Data":{"ConsentId":"urn:accounts:v3:d4c55300-8690-4521-8a53-8071dbc4b0a7","Status":"Authorised","CreationDateTime":"2021-01-01T11:34:45.265611+02:00","Permissions":["ReadAccountsBasic"],"ExpirationDateTime":"2021-01-05T23:35:00+02:00","TransactionFromDateTime":"2018-06-13T22:39:00+03:00","TransactionToDateTime":"2019-06-13T22:39:00+03:00","StatusUpdateDateTime":"2021-01-01T11:35:39.591842+02:00"},"Risk":{},"Links":{"Self":"https://sandbox-obp-api.danskebank.com/sandbox-open-banking/v3.1/aisp/account-access-consents/urn:accounts:v3:d4c55300-8690-4521-8a53-8071dbc4b0a7"},"Meta":{}}"

**ConsentId** needs to be sent as a parameter with each Open Banking API call.


### Accounts

This service returns the detail of attached accounts to the consent. If accountId isn't specified, all accounts will be appeared in th response.

**Endpoint**

`{url}/{aspspId}/accounts/{accountId}/cid/{consentTid}`

`{url}/{aspspId}/accounts/cid/{consentTid}`

**`url`**: should be pointing your application's domain name and port number.

**`aspspId`**: needs to be set with the aspspId which service will be called from.

**`cid`**: needs to be set with a consentTid which is retrieved from Retrieve Active Consents.

**`accountId`**: needs to be set with an accountId, if it is intended to retrieve a specific one.

**Example**

###### **Request**

>curl -v -H 'Authorization: Bearer <internal_access_token>' http://localhost:8080/danske/accounts/cid/10

###### **Response**

>{"Data":{"Account":[{"AccountId":"654f239a-5b8b-4f06-8aa6-8622d41d518b","Currency":"GBP","AccountType":"Personal","AccountSubType":"Savings","Nickname":"Sandbox Test Nickname","SwitchStatus":"UK.CASS.SwitchCompleted"},{"AccountId":"6c62c3dc-2763-4f70-9f4a-ffbc65dbcfaa","Currency":"GBP","AccountType":"Personal","AccountSubType":"Savings","Nickname":"Sandbox Test Nickname","SwitchStatus":"UK.CASS.SwitchCompleted"}]},"Links":{"Self":"https://sandbox-obp-api.danskebank.com/sandbox-open-banking/v3.1/aisp/accounts/"},"Meta":{}}

###### **Request**

>curl -v -H 'Authorization: Bearer <internal_access_token>' http://localhost:8080/danske/accounts/6c62c3dc-2763-4f70-9f4a-ffbc65dbcfaa/cid/10

###### **Response**

>{"Data":{"Account":[{"AccountId":"6c62c3dc-2763-4f70-9f4a-ffbc65dbcfaa","Currency":"GBP","AccountType":"Personal","AccountSubType":"Savings","Nickname":"Sandbox Test Nickname","SwitchStatus":"UK.CASS.NotSwitched"}]},"Links":{"Self":"https://sandbox-obp-api.danskebank.com/sandbox-open-banking/accounts/6c62c3dc-2763-4f70-9f4a-ffbc65dbcfaa/"},"Meta":{}}

