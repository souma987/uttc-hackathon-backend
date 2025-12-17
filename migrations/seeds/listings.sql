DELETE
FROM listings
WHERE TRUE;

INSERT INTO listings (id,
                      seller_id,
                      title,
                      description,
                      images,
                      price,
                      quantity,
                      status,
                      item_condition)
VALUES ('lst_00000000000000000000000001',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Vintage Leather Backpack',
        'Full-grain leather backpack with multiple compartments and sturdy hardware.',
        '[
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2Ftest%2FIMG_20251020_140312.jpg?alt=media"
          },
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2FK1yWuMX4phYOHanNk9pMuYTvsyo2%2F123d21bb-87fe-4cef-8f5a-7128c37619ad.png?alt=media"
          }
        ]',
        12500,
        3,
        'active',
        'excellent'),
       ('lst_00000000000000000000000002',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Refurbished Mechanical Keyboard',
        'Tenkeyless mechanical keyboard refurbished with new switches and keycaps.',
        '[
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2Ftest%2FIMG_20251020_140312.jpg?alt=media"
          }
        ]',
        9200,
        5,
        'draft',
        'good'),
       ('lst_00000000000000000000000003',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Handmade Ceramic Mug Set',
        'Set of four handmade ceramic mugs with a matte glaze finish.',
        '[
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2FK1yWuMX4phYOHanNk9pMuYTvsyo2%2F123d21bb-87fe-4cef-8f5a-7128c37619ad.png?alt=media"
          },
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2Ftest%2FIMG_20251020_140312.jpg?alt=media"
          }
        ]',
        4800,
        12,
        'sold',
        'new'),
       ('lst_00000000000000000000000004',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Vintage Leather Backpack',
        'Full-grain leather backpack with multiple compartments and sturdy hardware.',
        '[
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2Ftest%2FIMG_20251020_140312.jpg?alt=media"
          },
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2FK1yWuMX4phYOHanNk9pMuYTvsyo2%2F123d21bb-87fe-4cef-8f5a-7128c37619ad.png?alt=media"
          }
        ]',
        12500,
        3,
        'active',
        'excellent'),
       ('lst_00000000000000000000000005',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Vintage Leather Backpack',
        'Full-grain leather backpack with multiple compartments and sturdy hardware.',
        '[
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2Ftest%2FIMG_20251020_140312.jpg?alt=media"
          },
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2FK1yWuMX4phYOHanNk9pMuYTvsyo2%2F123d21bb-87fe-4cef-8f5a-7128c37619ad.png?alt=media"
          }
        ]',
        12500,
        3,
        'active',
        'excellent'),
       ('lst_00000000000000000000000006',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Vintage Leather Backpack',
        'Full-grain leather backpack with multiple compartments and sturdy hardware.',
        '[
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2Ftest%2FIMG_20251020_140312.jpg?alt=media"
          },
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2FK1yWuMX4phYOHanNk9pMuYTvsyo2%2F123d21bb-87fe-4cef-8f5a-7128c37619ad.png?alt=media"
          }
        ]',
        12500,
        3,
        'active',
        'excellent'),
       ('lst_00000000000000000000000007',
        'K1yWuMX4phYOHanNk9pMuYTvsyo2',
        'Vintage Leather Backpack',
        'Full-grain leather backpack with multiple compartments and sturdy hardware.',
        '[
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2Ftest%2FIMG_20251020_140312.jpg?alt=media"
          },
          {
            "url": "https://firebasestorage.googleapis.com/v0/b/term8-souma-nagano.firebasestorage.app/o/users%2FK1yWuMX4phYOHanNk9pMuYTvsyo2%2F123d21bb-87fe-4cef-8f5a-7128c37619ad.png?alt=media"
          }
        ]',
        12500,
        3,
        'active',
        'excellent')
;
