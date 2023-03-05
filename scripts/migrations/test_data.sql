INSERT INTO users(username, hashed_password, name, email, account_type)
VALUES ('geogreck', '$2a$10$Whoj5L2Bpn6qrvCxT5Ld8Oc4TOQXzlaqXdUFURPaA2/f10ij1Ffgm', 'George', 'geogreck@vk.com', 'personal'),
       ('kirill', '$2a$10$Whoj5L2Bpn6qrvCxT5Ld8Oc4TOQXzlaqXdUFURPaA2/f10ij1Ffgm', 'Kirill', 'figma@vk.com', 'personal'),
       ('slava', '$2a$10$Whoj5L2Bpn6qrvCxT5Ld8Oc4TOQXzlaqXdUFURPaA2/f10ij1Ffgm', 'Slava', 'iu7@vk.com', 'personal'),
       ('evgenii', '$2a$10$Whoj5L2Bpn6qrvCxT5Ld8Oc4TOQXzlaqXdUFURPaA2/f10ij1Ffgm', 'evgenii', 'test@vk.com', 'personal');

INSERT INTO boards(name, privacy, user_id)
VALUES ('Notes', 'public', 1),
       ('Saved', 'public', 1),
       ('Good images', 'secret', 1),
       ('Pictures', 'public', 2),
       ('My board', 'public', 2),
       ('ToDo', 'secret', 3);

INSERT INTO pins(title, board_id, media_source)
VALUES ('Road', 1, 'https://wg.grechkogv.ru/assets/pet7.webp'),
       ('Ice', 1, 'https://wg.grechkogv.ru/assets/armorChest4.webp'),
       ('Future', 1, 'https://wg.grechkogv.ru/assets/pet6.webp'),
       ('Color', 2, 'https://wg.grechkogv.ru/assets/pet8.webp'),
       ('Question', 2, 'https://wg.grechkogv.ru/assets/weapon5.webp'),
       ('Shops', 3, 'https://wg.grechkogv.ru/assets/weapon1.webp'),
       ('School', 4, 'https://wg.grechkogv.ru/assets/armorBeing3.webp');

INSERT INTO comments(description, pin_id, user_id)
VALUES ('Why?', 1, 2),
       ('It is good.', 1, 3),
       ('Normal', 2, 1),
       ('Ok', 2, 2),
       ('OK', 2, 3);
