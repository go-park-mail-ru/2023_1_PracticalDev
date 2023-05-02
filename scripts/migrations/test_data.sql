INSERT INTO users(username, hashed_password, name, email, account_type, profile_image)
VALUES ('geogreck', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'George', 'geogreck@vk.com',
        'personal', 'https://pickpin.hb.bizmrg.com/default-user-icon-8-4024862977'),
       ('kirill', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'Kirill', 'kirill@vk.com',
        'personal', 'https://pickpin.hb.bizmrg.com/default-user-icon-8-4024862977'),
       ('slava', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'Slava', 'slava@vk.com',
        'personal', 'https://pickpin.hb.bizmrg.com/default-user-icon-8-4024862977'),
       ('evgenii', '$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS', 'Evgenii', 'evgenii@vk.com',
        'personal', 'https://pickpin.hb.bizmrg.com/default-user-icon-8-4024862977');

INSERT INTO boards(name, privacy, user_id)
VALUES ('Tea', 'secret', 1),
       ('Пейзажи', 'secret', 1),
       ('Good Images', 'secret', 1),
       ('Nature', 'secret', 2),
       ('Flowers', 'secret', 2),
       ('Океан', 'secret', 2),
       ('Природа', 'secret', 3),
       ('Airplanes', 'secret', 3),
       ('Cars', 'secret', 3),
       ('Basketball', 'secret', 4),
       ('Футбол', 'secret', 4);

INSERT INTO pins(title, description, media_source, author_id)
VALUES ('Самолет', '#самолет #путешествия', 'https://pickpin.hb.bizmrg.com/680caf09-00dc-4e6c-ad09-f455e4e9b8e8.jpg',
        3),
       ('Two Planes', '', 'https://pickpin.hb.bizmrg.com/9b3a1f23-cf9b-41ed-9871-282fad93841a.jpg', 3),
       ('Plane in city', '', 'https://pickpin.hb.bizmrg.com/5230f56b-ed84-42de-9e50-76adbd665ef2.jpg', 3),

       ('Basketball', 'Our basketball team', 'https://pickpin.hb.bizmrg.com/e7b00a68-d644-496c-9745-e367d665268a.jpg',
        4),
       ('Ball in basket', 'Our basketball team',
        'https://pickpin.hb.bizmrg.com/d4a54111-cc28-4ce2-9366-6044f6135e96.jpg', 4),
       ('Баскетбол', 'Our basketball team', 'https://pickpin.hb.bizmrg.com/a69fe86f-b857-45af-b5d6-0e4974538cb2.jpg',
        4),

       ('Футбол', 'Наша футбольная команда', 'https://pickpin.hb.bizmrg.com/6e209d7b-651b-4cbd-97aa-d8594a00fc66.jpg',
        4),
       ('Football', 'Our football team', 'https://pickpin.hb.bizmrg.com/1a682e14-672b-48ad-9d46-63db96c11693.jpg', 4),
       ('Football', 'Our football team', 'https://pickpin.hb.bizmrg.com/d14dad98-a533-4be9-8edc-cf35cd908c79.jpg', 4),

       ('Some tea', '#tasty #tea', 'https://pickpin.hb.bizmrg.com/f6e0346a-0dec-40df-b70f-ea1e945ed65f.jpg', 1),
       ('Tea', '#tasty #tea', 'https://pickpin.hb.bizmrg.com/ad3afd89-8beb-469d-880c-2f0baa3b09b5.jpg', 1),
       ('Good tea', '#tasty #tea', 'https://pickpin.hb.bizmrg.com/214763ac-b62e-489f-adec-e26cc9641915.jpg', 1),

       ('Beautiful Car', '', 'https://pickpin.hb.bizmrg.com/e4dd748e-091e-4810-be98-58e630128574.jpg', 3),
       ('Car', '', 'https://pickpin.hb.bizmrg.com/95cc56c2-dbf2-47dc-ae78-1659d899655e.jpg', 3),
       ('Tesla Car', '', 'https://pickpin.hb.bizmrg.com/d98a7ab9-827c-4855-aa4c-ce55528ccce7.jpg', 3),
       ('Formula 1', '', 'https://pickpin.hb.bizmrg.com/2383c62b-ec7f-4b3c-8c30-edd51d73f455.jpg', 3),
       ('Future Train', '#future', 'https://pickpin.hb.bizmrg.com/9742d1a9-f5b5-46af-9e00-4c9bdba49c4e.jpg', 3),

       ('Moscow City', '', 'https://pickpin.hb.bizmrg.com/2bf53183-438f-4379-b798-9cb6d8611f58.jpg', 1),
       ('Москва', '', 'https://pickpin.hb.bizmrg.com/2e6ab359-b627-4597-977f-ba361a8e8ea3.jpg', 1),

       ('Nature', '', 'https://pickpin.hb.bizmrg.com/15ea0fa6-0064-40af-b3b7-a5f5455119ae.jpg', 2),
       ('Nature', '', 'https://pickpin.hb.bizmrg.com/bbe7bb81-15fd-45a9-b614-3e4ca3e00cef.jpg', 2),
       ('Облака', '', 'https://pickpin.hb.bizmrg.com/1bc0410f-9977-4097-bffe-d16889a388d5.jpg', 2),
       ('Ocean', '', 'https://pickpin.hb.bizmrg.com/97d80f59-c0dd-43b5-b491-ccce0cf0c2c5.jpg', 2),
       ('Природа', '', 'https://pickpin.hb.bizmrg.com/2f1e32b2-51ab-48e1-8ca2-aa64781a4d98.jpg', 2),
       ('Природа', '', 'https://pickpin.hb.bizmrg.com/a6eb885c-e845-4aba-a51b-020a2306d21f.jpg', 2),
       ('Beautiful Ocean', '', 'https://pickpin.hb.bizmrg.com/7ddacf84-8835-4e11-9df5-f0363c2449fa.jpg', 2),
       ('Nature', '', 'https://pickpin.hb.bizmrg.com/84706766-c5e2-4649-a39f-feff2f68bd35.jpg', 2),
       ('Boat and River', '', 'https://pickpin.hb.bizmrg.com/6dfb24dd-40d9-48d5-9bd3-d54913b26c00.jpg', 2),
       ('Flowers', '', 'https://pickpin.hb.bizmrg.com/519ad2c9-27e7-4772-b3da-5230a0d6a3fd.jpg', 1),
       ('Rainbow', '', 'https://pickpin.hb.bizmrg.com/782584df-02c1-43fd-bf59-61c369c4411d.jpg', 1),
       ('Красивый гриб', '', 'https://pickpin.hb.bizmrg.com/85cc1542-424d-4625-b86b-09e50b50d5c0.jpg', 1),
       ('Океан', '', 'https://pickpin.hb.bizmrg.com/46415b8b-3482-4b5d-bb98-fee423f967f9.jpg', 2),

       ('Апельсины', '', 'https://pickpin.hb.bizmrg.com/3ba25a33-d262-4099-a9f8-bd46a6001f8a.jpg', 3),
       ('Сова', '', 'https://pickpin.hb.bizmrg.com/6360462f-ea43-4589-87cd-e3d365427c7c.jpg', 3),
       ('Поле', '', 'https://pickpin.hb.bizmrg.com/2a09e3e4-0ee2-422b-8693-422c48f8c8a1.jpg', 3),
       ('Мак', '', 'https://pickpin.hb.bizmrg.com/7c888297-1e25-4e0f-9f83-9df5251ebcf9.jpg', 3),

       ('Красивый пейзаж', '', 'https://pickpin.hb.bizmrg.com/62690b32-5c69-4a29-a2b1-8ec7031f92f2.jpg', 2),
       ('Зима Пейзаж', '', 'https://pickpin.hb.bizmrg.com/be03308f-916b-49a1-bbd6-67428d4ef154.jpg', 2),
       ('Пейзаж "Кит"', '', 'https://pickpin.hb.bizmrg.com/eca6ca6a-841f-4098-a8ef-8530210f8aec.jpg', 2),
       ('Гранат', '', 'https://pickpin.hb.bizmrg.com/9ec4ce1a-ff5c-4f28-b4b4-b5d93882367e.jpg', 2);

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

INSERT INTO boards_pins(board_id, pin_id)
VALUES (8, 1),
       (8, 2),
       (8, 3),

       (10, 4),
       (10, 5),
       (10, 6),

       (11, 7),
       (11, 8),
       (11, 9),

       (1, 10),
       (1, 11),
       (1, 12),

       (9, 13),
       (9, 14),
       (9, 15),
       (9, 16),
       (9, 17),

       (3, 18),
       (3, 19),

       (4, 19),
       (4, 20),
       (4, 21),
       (4, 22),
       (4, 23),
       (4, 24),
       (4, 25),
       (4, 26),
       (4, 27),
       (4, 28),
       (4, 29);

INSERT INTO followings (followee_id, follower_id)
VALUES (1, 2),
       (2, 3),
       (3, 4),
       (4, 1);

INSERT INTO chats (user1_id, user2_id)
VALUES (1, 2), -- Гоша - Кирилл
       (2, 3), -- Кирилл - Слава
       (3, 4); -- Слава - Женя
