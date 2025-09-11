-- Normalize mtrackingtag image to filename for frontend resolver
-- Run this against your production database (Heroku Postgres)

-- If product currently stores the filename, keep it; otherwise convert from any R2 URL
UPDATE products
SET images = ARRAY['mtrackingtag.png']
WHERE images @> ARRAY['mtrackingtag.jpg']
   OR images @> ARRAY['mtrackingtag.png']
   OR images @> ARRAY['smart tracking card.jpg']
   OR images @> ARRAY['https://2a35af424f8734e497a5d707344d79d5.r2.cloudflarestorage.com/smrtmart/mtrackingtag.jpg'];
