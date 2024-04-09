SELECT
    *
FROM
    talvi.twofactor
WHERE
    parent_account_hash = $1
LIMIT
    1;