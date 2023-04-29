-- name: TermsSearchRaw :many
WITH const AS (SELECT ? AS query)
SELECT terms.*, (
                1 / LOG(100, CAST(popularity AS REAL) + 99) * 40 + -- 1/log(100, popularity + 99 <base # - 1>)
                1 / LOG(3, ABS(LENGTH(TRIM(text, '-')) - LENGTH(const.query)) + 3) * 60 -- 1/log(log_base, len_diff + log_base)
    ) AS calc_rank
FROM terms, const WHERE text LIKE '%' || const.query || '%' OR variants LIKE '%' || const.query || '%'
ORDER BY calc_rank DESC;