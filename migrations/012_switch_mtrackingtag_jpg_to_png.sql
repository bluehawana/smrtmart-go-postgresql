-- Ensure any remaining references to mtrackingtag.jpg are updated to .png

UPDATE products
SET images = ARRAY['mtrackingtag.png']
WHERE images @> ARRAY['mtrackingtag.jpg'];

