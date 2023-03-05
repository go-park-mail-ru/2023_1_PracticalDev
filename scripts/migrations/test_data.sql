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

INSERT INTO pins(title, board_id)
VALUES ('Road', 1),
       ('Ice', 1),
       ('Future', 1),
       ('Color', 2),
       ('Question', 2),
       ('Shops', 3),
       ('School', 4);

INSERT INTO comments(description, pin_id, user_id)
VALUES ('Why?', 1, 2),
       ('It is good.', 1, 3),
       ('Normal', 2, 1),
       ('Ok', 2, 2),
       ('OK', 2, 3);
