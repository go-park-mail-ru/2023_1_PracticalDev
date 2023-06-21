CREATE USER pickpin_main WITH PASSWORD 'pickpinpswd';
CREATE USER pickpin_search WITH PASSWORD 'pickpinpswd';

GRANT SELECT, UPDATE, INSERT, DELETE ON
    users, users_id_seq,
    boards, boards_id_seq,
    pins, pins_id_seq,
    pin_likes,
    boards_pins,
    comments, comments_id_seq,
    followings,
    chats, chats_id_seq,
    messages, messages_id_seq,
    notifications, notifications_id_seq,
    new_pin_notifications,
    new_like_notifications,
    new_comment_notifications,
    new_follower_notifications 
    TO pickpin_main;

GRANT SELECT ON 
    users,
    pins,
    boards
    TO pickpin_search;

