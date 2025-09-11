-- Update Smart Tracking Card (by UUID) with English content and new image

UPDATE products
SET 
    name = 'Smart Tracking Card',
    description = 'Ultra-slim Bluetooth tracking card that fits perfectly in your wallet. Works with Apple Find My network for easy location tracking of your valuables.',
    price = 299.0,
    compare_price = NULL,
    category = 'accessories',
    tags = ARRAY['bluetooth', 'tracking', 'card', 'wallet', 'find my', 'apple'],
    images = ARRAY['monecard.png'],
    seo = jsonb_build_object(
        'title', 'Smart Tracking Card',
        'description', 'Ultra-slim Bluetooth tracking card that fits perfectly in your wallet. Works with Apple Find My network for easy location tracking of your valuables.',
        'keywords', to_jsonb(ARRAY['bluetooth','tracking','card','wallet','find my','apple'])
    ),
    updated_at = NOW()
WHERE id = '6f2e79da-9591-4b0c-82c5-ea8efadbd35d'::uuid;

