-- Update Sony WH-1000XM5 to WH-1000XM6 with new price and description
UPDATE products
SET
    name = 'Sony WH-1000XM6 Headphones',
    price = 4888,
    description = 'Latest generation industry-leading noise canceling headphones with exceptional sound quality, 40-hour battery life, and enhanced comfort. Features advanced V2 processor, multipoint connection, and premium build quality for the ultimate listening experience.',
    images = ARRAY['sony.jpg'],  -- Keep existing image until new XM6 image is provided
    updated_at = NOW()
WHERE numeric_id = 1;