INSERT INTO users(username, hashed_password, name, email, account_type, profile_image)
VALUES ('geogreck', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'George', 'geogreck@vk.com',
        'personal', 'https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg'),
       ('kirill', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'Kirill', 'figma@vk.com', 'personal', 'https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg'),
       ('slava', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'Slava', 'iu7@vk.com', 'personal', 'https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg'),
       ('evgenii', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'evgenii', 'test@vk.com',
        'personal', 'https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg');

INSERT INTO boards(name, privacy, user_id)
VALUES ('Notes', 'public', 1),
       ('Saved', 'public', 1),
       ('Good images', 'secret', 1),
       ('Pictures', 'public', 2),
       ('My board', 'public', 2),
       ('ToDo', 'secret', 3);

INSERT INTO pins(title, author_id, media_source)
VALUES ('Ice', 1, 'https://pickpin.hb.bizmrg.com/68ed77e8-9002-491d-bd3d-d36938e054b8.jpg'),
       ('Future', 1, 'https://pickpin.hb.bizmrg.com/7c0a9457-6ea6-4c26-b8fd-4d7c0276c568.jpg'),
       ('Shops', 1, 'https://i.pinimg.com/564x/2f/93/56/2f9356b9346e82c14bf286c6a107bc7a.jpg'),
       ('Shops', 2, 'https://i.pinimg.com/564x/32/ff/71/32ff717c3cd3bd3d1886c775b59f0769.jpg'),
       ('Shops', 2, 'https://i.pinimg.com/564x/ce/e3/01/cee3011f3e19de4377dbf98f397c027b.jpg'),
       ('Shops', 3, 'https://i.pinimg.com/564x/a6/ba/55/a6ba553df2a0c0f3894ef328a86fb373.jpg'),
       ('Shops', 3, 'https://i.pinimg.com/564x/43/2d/3b/432d3b28d1661439245422e9218ffcce.jpg'),
       ('School', 4, 'https://i.pinimg.com/564x/98/9d/3f/989d3f5c158dcac7ca4d115bff866d84.jpg');


INSERT INTO comments(description, pin_id, user_id)
VALUES ('Why?', 1, 2),
       ('It is good.', 1, 3),
       ('Normal', 2, 1),
       ('Ok', 2, 2),
       ('OK', 2, 3);

INSERT INTO pin_likes(pin_id, author_id)
VALUES (1, 2),
       (1, 3),
       (2, 1),
       (2, 2),
       (2, 3);
