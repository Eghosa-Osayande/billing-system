-- name: GetUserProfileWhere :many
Select
    users.*,
    JSON_AGG(business.*) as business
from
    users
left JOIN
    business on users.id = business.owner_id
where
    (
        users.id = sqlc.narg('id') or sqlc.narg('id') IS NULL
    ) and
    (
        users.email = sqlc.narg('email') or sqlc.narg('email') IS NULL
    ) 
group by
    users.id
LIMIT sqlc.narg('limit');