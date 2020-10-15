CREATE TABLE IF NOT EXISTS "user"
(
    id            serial       NOT NULL,
    nickname      varchar(20)  NOT NULL,
    car_no        varchar(20)  NOT NULL,
    avatar        varchar(255) NOT NULL DEFAULT '',
    sex           integer      NOT NULL DEFAULT 1,

    province_id   integer      NOT NULL DEFAULT 0,
    city_id       integer      NOT NULL DEFAULT 0,
    province_name varchar(20)  NOT NULL DEFAULT '',
    city_name     varchar(20)  NOT NULL DEFAULT '',

    created_at    timestamp    NOT NULL DEFAULT now(),
    updated_at    timestamp    NOT NULL DEFAULT now(),
    deleted_at    timestamp    NULL,

    CONSTRAINT user_pk PRIMARY KEY ("id")
);


CREATE TABLE IF NOT EXISTS "group"
(
    id         serial       NOT NULL,
    group_name varchar(255) NOT NULL DEFAULT '',
    type       integer      NOT NULL DEFAULT 0, --- discuss or main group
    owner      integer      NOT NULL DEFAULT 0,
    is_banned  boolean      NOT NULL DEFAULT false,
    notice     varchar(255) NOT NULL DEFAULT '',

    created_at timestamp    NOT NULL DEFAULT now(),
    updated_at timestamp    NOT NULL DEFAULT now(),
    deleted_at timestamp    NULL
);

CREATE TABLE IF NOT EXISTS "group_user"
(
    id       serial  NOT NULL,
    uid      integer NOT NULL DEFAULT 0,
    group_id integer NOT NULL DEFAULT 0,

    created_at timestamp    NOT NULL DEFAULT now(),
    updated_at timestamp    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "message"
(

    id         serial    NOT NULL,
    sender_id  integer NOT NULL DEFAULT 0,
    receiver_id integer NOT NULL DEFAULT 0,
    receiver_group_id integer NOT NULL DEFAULT 0,

    message_type integer NOT NULL DEFAULT 0, --- 单聊消息，讨论组/群聊消息
    message_media_type integer NOT NULL DEFAULT 0, --- 消息的具体类型: text,video,image,voice,link,
    message_system_type integer not null default 0, --- 消息的系统类型：system,user,
    content text NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp NULL
);


CREATE TABLE IF NOT EXISTS "chat_list"
(
    id                     serial       NOT NULL,
    title                  varchar(255) NOT NULL DEFAULT '',
    avatar                 varchar(255) NOT NULL DEFAULT '',
    type                   integer      NOT NULL DEFAULT 0, --- user/group
    user_id                integer      NOT NULL DEFAULT 0,

    link_user_id           integer      NOT NULL DEFAULT 0,
    link_group_id          integer      NOT NULL DEFAULT 0,
    unread_count           integer      NOT NULL DEFAULT 0,
    latest_message_desc    varchar(255) NOT NULL DEFAULT '',
    latest_message_updated timestamp    NOT NULL DEFAULT now(),

    created_at             timestamp    NOT NULL DEFAULT now(),
    updated_at             timestamp    NOT NULL DEFAULT now(),
    deleted_at             timestamp    NULL
);