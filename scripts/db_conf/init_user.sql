CREATE USER pickpin_main WITH PASSWORD 'pickpinpswd';
CREATE USER pickpin_search WITH PASSWORD 'pickpinpswd';

GRANT SELECT, UPDATE, INSERT, DELETE ON
    users,
    boards,
    pins,
    pin_likes,
    boards_pins,
    comments,
    followings,
    chats,
    messages,
    notifications,
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

