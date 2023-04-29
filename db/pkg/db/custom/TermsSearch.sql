-- name: TermsSearchRaw :many
WITH const AS (SELECT ? AS query, ? AS pop_log, ? AS pop_weight, ? AS len_log)
SELECT terms.*,
    1 / LOG(pop_log, CAST(popularity AS REAL) + pop_log - 1) * pop_weight AS pop_calc,
    1 / LOG(len_log, ABS(LENGTH(TRIM(text, '-')) - LENGTH(query)) + len_log) * (100 - pop_weight) AS len_calc
FROM terms, const WHERE text LIKE '%' || const.query ||'%' OR variants LIKE '%' || const.query ||'%'
ORDER BY pop_calc + len_calc DESC;