-- sqlc can't read this table, see matching table in schema_implied.sql
CREATE VIRTUAL TABLE terms USING fts5(
  text,
  variants,
  part_of_speech UNINDEXED,
  common_level UNINDEXED,
  translations UNINDEXED,
  popularity UNINDEXED
);