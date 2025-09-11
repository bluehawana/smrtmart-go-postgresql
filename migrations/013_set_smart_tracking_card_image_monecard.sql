-- Point Smart Tracking Card image to monecard.png

UPDATE products
SET images = ARRAY['monecard.png']
WHERE images @> ARRAY['mtrackingtag.png']
   OR images @> ARRAY['mtrackingtag.jpg']
   OR images @> ARRAY['smart tracking card.jpg']
   OR (LOWER(name) LIKE '%tracking%' AND LOWER(name) LIKE '%card%');

