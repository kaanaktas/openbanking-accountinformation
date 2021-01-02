
INSERT INTO session_table(tid, tpp_id, user_id, aspsp_id, internal_access_token, reference_id, create_date_time, update_date_time)
    VALUES (1, 'Tpp_1', 'kaan', 'danske', '123456789', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9', '2020-12-05T15:17:50Z', '2020-12-05T15:17:50Z');
INSERT INTO session_table(tid, tpp_id, user_id, aspsp_id, internal_access_token, reference_id, create_date_time, update_date_time)
VALUES (2, 'Tpp_1', 'kaan', 'ozone', '987654321', 'test0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ91', '2020-12-05T15:17:50Z', '2020-12-05T15:17:50Z');

-- Consent with Token
--valid date
INSERT INTO consent_table(id, aspsp_id, consent_expiration_date_time, consent_id, consent_status, consent_status_update_date_time,
                          consent_transaction_from_date_time, consent_transaction_to_date_time, consent_type, create_date_time,
                          session_reference_id, tracking_id, update_date_time)
VALUES (1, 'danske', '2050-06-13T19:39:00Z', 'urn:accounts:v3:e14fb190-a7cd-4abf-8c87-ea73fde845e4', 'Authorised',
        '', '2018-06-13T19:39:00Z', '2019-06-13T19:39:00Z', 'ACCOUNT', '2050-12-05 15:17:50.928479',
        'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9', '163b39483659429edd305c8c9228a55456a9e42b',
        '2020-12-05 15:17:50.965');

INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token,
                                token_expiration_date_time, token_status, update_date_time, consent_tid)
VALUES ('AJANzM4MDVkMjVjNzQwODVlZjE1Y2RmO', '2020-12-05 15:17:50.938', '60',
        '11111111111111', '22222222222222222', '2050-12-05T15:17:50Z', 'Authorised', '2020-12-05 15:18:19.75', 1);

INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token,
                                token_expiration_date_time, token_status, update_date_time, consent_tid)
VALUES ('AJANzM4MDVkMjVjNzQwODVlZjE1Y2RmO', '2020-12-05 15:17:50.938', '60',
        '11111111111111', '22222222222222222', '2050-12-05T15:17:50Z', 'Revoked', '2020-12-05 15:18:19.75', 1);


--expiry date
INSERT INTO consent_table(id, aspsp_id, consent_expiration_date_time, consent_id, consent_status, consent_status_update_date_time,
                          consent_transaction_from_date_time, consent_transaction_to_date_time, consent_type, create_date_time,
                          session_reference_id, tracking_id, update_date_time)
VALUES (2, 'danske', '2018-06-13T19:39:00Z', 'urn:accounts:v3:e14fb190-a7cd-4abf-8c87-ea73fde845', 'Revoked',
        '', '2018-06-13T19:39:00Z', '2019-06-13T19:39:00Z', 'ACCOUNT', '2050-12-05 15:17:50.928479',
        'test0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ91', '163b39483659429edd305c8c9228a55456a9e4',
        '2020-12-05 15:17:50.965');

INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token,
                                token_expiration_date_time, token_status, update_date_time, consent_tid)
VALUES ('AJANzM4MDVkMjVjNzQwODVlZjE1Y2RmO', '2020-12-05 15:17:50.938', '60',
        '11111111111111', '22222222222222222', '2020-12-05T15:17:50Z', 'Expired', '2020-12-05 15:18:19.75', 2);

INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token,
                                token_expiration_date_time, token_status, update_date_time, consent_tid)
VALUES ('AJANzM4MDVkMjVjNzQwODVlZjE1Y2RmO', '2020-12-05 15:17:50.938', '60',
        '11111111111111', '22222222222222222', '2020-12-05T15:17:50Z', 'Expired', '2020-12-05 15:18:19.75', 2);

---------------------------
INSERT INTO consent_table(id, aspsp_id, consent_expiration_date_time, consent_id, consent_status, consent_status_update_date_time,
                          consent_transaction_from_date_time, consent_transaction_to_date_time, consent_type, create_date_time,
                          session_reference_id, tracking_id, update_date_time)
VALUES (3, 'danske', '2050-06-13T19:39:00Z', 'urn:accounts:v3:e14fb190-a7cd-4abf-8c87-ea', 'Authorised',
        '', '2018-06-13T19:39:00Z', '2019-06-13T19:39:00Z', 'ACCOUNT', '2050-12-05 15:17:50.928479',
        'test0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ91', '163b39483659429edd305c8c9228a55456a9',
        '2020-12-05 15:17:50.965');

INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token,
                                token_expiration_date_time, token_status, update_date_time, consent_tid)
VALUES ('AJANzM4MDVkMjVjNzQwODVlZjE1Y2RmO', '2020-12-05 15:17:50.938', '60',
        '11111111111111', '22222222222222222', '2020-12-05T15:17:50Z', 'Expired', '2020-12-05 15:18:19.75', 2);

INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token,
                                token_expiration_date_time, token_status, update_date_time, consent_tid)
VALUES ('AJANzM4MDVkMjVjNzQwODVlZjE1Y2RmO', '2020-12-05 15:17:50.938', '60',
        '11111111111111', '22222222222222222', '2020-12-05T15:17:50Z', 'Expired', '2020-12-05 15:18:19.75', 3);

-------
INSERT INTO consent_table(id, aspsp_id, consent_expiration_date_time, consent_id, consent_status, consent_status_update_date_time,
                          consent_transaction_from_date_time, consent_transaction_to_date_time, consent_type, create_date_time,
                          session_reference_id, tracking_id, update_date_time)
VALUES (4, 'danske', '2018-06-13T19:39:00Z', 'urn:accounts:v3:e14fb190-a7cd-4abf-8c87-ea73fde', 'Authorised',
        '', '2018-06-13T19:39:00Z', '2019-06-13T19:39:00Z', 'ACCOUNT', '2050-12-05 15:17:50.928479',
        'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9', '163b39483659429edd305c8c9228a5545',
        '2020-12-05 15:17:50.965');

INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token,
                                token_expiration_date_time, token_status, update_date_time, consent_tid)
VALUES ('AJANzM4MDVkMjVjNzQwODVlZjE1Y2RmO', '2020-12-05 15:17:50.938', '60',
        '11111111111111', '22222222222222222', '2020-12-05T15:17:50Z', 'Expired', '2020-12-05 15:18:19.75', 4);
