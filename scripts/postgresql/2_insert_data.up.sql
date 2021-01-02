-- Danske
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_OAUTH2', 'https://sandbox-obp-api.danskebank.com/sandbox-open-banking/private/oauth2/token', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_DOMESTIC_CONSENT', '', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_AUTHORIZE', 'https://sandbox-obp-auth.danskebank.com/sandbox-open-banking/private/oauth2/authorize', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_DOMESTIC_PAYMENT', 'https://sandbox.hsbc.com/psd2/obie/v3.1/domestic-payments', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_ACCOUNT_ACCESS_CONSENT', 'https://sandbox-obp-api.danskebank.com/sandbox-open-banking/v3.1/aisp/account-access-consents', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_ACCOUNTS', 'https://sandbox-obp-api.danskebank.com/sandbox-open-banking/v3.1/aisp/accounts', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('FAPI_FINANCIAL_ID', '0015800000jf7AeAAI', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('AUD', 'https://sandbox-obp-api.danskebank.com/sandbox-open-banking/private', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ISS', '', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('TOKEN_ENDPOINT_AUTH_METHOD', 'tls_client_auth', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('AUTHENTICATION_TOKEN', '', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('REDIRECT_URL', 'https://<domain>/callback', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('APPLICATION_TYPE', 'web', 'danske');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('CLIENT_ID', '', 'danske');

--Ozone
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_OAUTH2', 'https://ob19-auth1.o3bank.co.uk:4201/token', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_DOMESTIC_CONSENT', '', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_AUTHORIZE', 'https://ob19-auth1-ui.o3bank.co.uk/auth', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_DOMESTIC_PAYMENT', '', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_ACCOUNT_ACCESS_CONSENT', 'https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/account-access-consents', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ENDPOINT_ACCOUNTS', 'https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/accounts', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('FAPI_FINANCIAL_ID', '0015800001041RHAAY', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('AUD', '0015800001041RHAAY', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('ISS', '', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('TOKEN_ENDPOINT_AUTH_METHOD', 'tls_client_auth', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('AUTHENTICATION_TOKEN', '', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('REDIRECT_URL', 'https://<domain>/callback', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('APPLICATION_TYPE', 'web', 'ozone');
INSERT INTO config_table(config_name, config_value, aspsp_id) VALUES ('CLIENT_ID', '', 'ozone');
