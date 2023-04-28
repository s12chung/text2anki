-- name: TermsSearch :many
SELECT *, rank, (
        1 / LOG(100, CAST(popularity AS REAL) + 99) * 50 + -- 1/log(100, popularity + 99 <base # - 1>)
        rank / -12 * 30 +
        1 / LOG(2, ABS(LENGTH(TRIM(text, '-')) - ?) + ? + 1) * 20 -- 1/log(2, len_diff + 1 <so don't log zero> + 1 <base # - 1>)
    ) AS calc_rank
FROM terms WHERE terms MATCH ? AND rank MATCH 'bm25(1.0, 0.5)'
ORDER BY calc_rank DESC;