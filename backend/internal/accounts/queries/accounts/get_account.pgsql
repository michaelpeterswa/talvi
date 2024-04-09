SELECT
    *
FROM
    talvi.accounts
WHERE
    email_provider_hash = $1
LIMIT
    1;