-- Seed data for the listings table
-- Requires a matching seller_id in the users table (e.g. insert uid_demo_seller_01 into users beforehand)
INSERT INTO listings (
    id,
    seller_id,
    title,
    description,
    images,
    price,
    quantity,
    status,
    item_condition
)
VALUES
    (
        'lst_00000000000000000000000001',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Vintage Leather Backpack',
        'Full-grain leather backpack with multiple compartments and sturdy hardware.',
        '[{"url": "https://images.unsplash.com/photo-1447933601403-0c6688de566e?auto=format&fit=crop&w=1200&q=80"}, {"url": "https://images.unsplash.com/photo-1489515217757-5fd1be406fef?auto=format&fit=crop&w=1200&q=80"}]',
        12500,
        3,
        'active',
        'excellent'
    ),
    (
        'lst_00000000000000000000000002',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Refurbished Mechanical Keyboard',
        'Tenkeyless mechanical keyboard refurbished with new switches and keycaps.',
        '[{"url": "https://images.unsplash.com/photo-1507878866276-a947ef722fee?auto=format&fit=crop&w=800&q=80"}]',
        9200,
        5,
        'draft',
        'good'
    ),
    (
        'lst_00000000000000000000000003',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Handmade Ceramic Mug Set',
        'Set of four handmade ceramic mugs with a matte glaze finish.',
        '[{"url": "https://images.unsplash.com/photo-1489515217757-5fd1be406fef?auto=format&fit=crop&w=1200&q=80"}, {"url": "https://images.unsplash.com/photo-1447933601403-0c6688de566e?auto=format&fit=crop&w=1200&q=80"}]',
        4800,
        12,
        'sold',
        'new'
    );
