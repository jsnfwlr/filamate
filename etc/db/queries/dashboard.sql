-- name: GetStorageStats :many
-- GetStorageStats returns the storage statistics for each location, including the maximum capacity, used capacity, and free capacity. It also includes a total row that sums up the statistics for all locations. The results are ordered by location ID.
SELECT locations.id, tally.label::text, max::bigint, used::bigint, free::bigint
FROM (
SELECT locations.label, locations.capacity as max, count(spools.id) as used, (locations.capacity - count(spools.id)) as free FROM locations LEFT JOIN spools ON spools.location_id = locations.id AND spools.deleted_at IS NULL WHERE locations.tally = true GROUP BY locations.label, locations.capacity, locations.tally
UNION
SELECT 'TOTAL' as label, SUM(max) as max, SUM(used) as used, sum(free) as free FROM (
SELECT locations.label, locations.capacity as max, count(spools.id) as used, (locations.capacity - count(spools.id)) as free FROM locations LEFT JOIN spools ON spools.location_id = locations.id AND spools.deleted_at IS NULL WHERE locations.tally = true GROUP BY locations.label, locations.capacity, locations.tally
)
) as tally
LEFT JOIN locations on locations.label = tally.label
ORDER BY locations.id;


-- name: GetUsageStats :many
-- GetUsageStats returns the most used color and material combinations, sorted by the number of empty spools and total spools. It filters out combinations that have only one spool or where the number of empty spools is less than half of the total spools.
SELECT color::text, material, SUM(empty)::bigint as used, count(*)::bigint as ordered
FROM (
SELECT STRING_AGG(colors.label, ', ') as color,
materials.class as material,
CASE WHEN spools.emptied_at IS NOT NULL THEN 1 ELSE 0 END as empty
FROM spools
JOIN spool_colors ON spool_id = spools.id
JOIN colors ON color_id = colors.id
JOIN materials  ON material_id = materials.id
GROUP BY spools.id, materials.class, (spools.emptied_at IS NOT NULL)
)
GROUP BY color,material
HAVING count(*) != 1 AND ((SUM(empty) * 2) >= count(*) OR count(*) > 1)
ORDER BY SUM(empty) desc, count(*) DESC;

-- name: GetMaterialChartData :many
-- GetMaterialChartData returns the data for a material chart: class, material, brand, and count of spools for each combo.
-- It includes a total row for each class, class and material, and a grand total row. The results are ordered by material label and class, with null values treated as 'All'.
SELECT (CASE WHEN materials.class IS NULL THEN '' ELSE materials.class END)::text as class,
(CASE WHEN materials.label IS NULL THEN '' ELSE materials.label END)::text as material,
(CASE WHEN brands.label IS NULL THEN '' ELSE brands.label END)::text as brand,
count(spools.id) as count
FROM spools
JOIN materials ON material_id = materials.id
JOIN brands ON brand_id = brands.id
WHERE spools.deleted_at IS NULL AND spools.emptied_at IS NULL
GROUP BY ROLLUP (materials.class, materials.label, brands.label)
ORDER BY CASE WHEN materials.class IS NULL THEN '' ELSE materials.class END DESC,
CASE WHEN materials.label IS NULL THEN '' ELSE materials.label END DESC,
CASE WHEN brands.label IS NULL THEN '' ELSE brands.label END DESC;








