-- Replace legacy CloudFront URL for Smart Tracking Card with monecard.png

UPDATE products
SET images = ARRAY['monecard.png']
WHERE images @> ARRAY['https://d10qehs4k3bdf9.cloudfront.net/smart-tracking-card.jpg'];

