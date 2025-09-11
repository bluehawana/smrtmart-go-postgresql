-- Update Smart Tracking Card product details and image filename
-- Assumptions: product is identified by existing image filename or name pattern

UPDATE products
SET 
    name = 'Smart Tracking Card',
    description = 'Ultra-slim Bluetooth tracking card that fits perfectly in your wallet. Works with Apple Find My network for easy location tracking of your valuables.',
    price = 299.0,
    compare_price = NULL,
    category = 'accessories',
    tags = ARRAY['bluetooth', 'tracking', 'card', 'wallet', 'find my', 'apple'],
    images = ARRAY['mtrackingtag.png'],
    seo = jsonb_build_object(
        'title', 'Smart Tracking Card',
        'description', 'Ultra-slim Bluetooth tracking card that fits perfectly in your wallet. Works with Apple Find My network for easy location tracking of your valuables.',
        'keywords', to_jsonb(ARRAY['bluetooth','tracking','card','wallet','find my','apple'])
    ),
    updated_at = NOW()
WHERE images @> ARRAY['mtrackingtag.jpg']
   OR images @> ARRAY['smart tracking card.jpg']
   OR images @> ARRAY['mtrackingtag.png']
   OR (LOWER(name) LIKE '%tracking%' AND LOWER(name) LIKE '%card%');
