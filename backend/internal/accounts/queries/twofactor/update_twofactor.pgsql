UPDATE talvi.twofactor
SET
    (enabled) = ($2)
WHERE
    parent_account_hash = $1;