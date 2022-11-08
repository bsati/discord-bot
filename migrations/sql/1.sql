CREATE TABLE IF NOT EXISTS birthdays (
    id INT GENERATED ALWAYS AS IDENTITY,
    user_id VARCHAR UNIQUE NOT NULL,
    date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_guild (
    user_id VARCHAR NOT NULL,
    guild_id VARCHAR NOT NULL,
    PRIMARY KEY (user_id, guild_id)
);

CREATE TABLE IF NOT EXISTS guild_bot_channel (
    id INT GENERATED ALWAYS AS IDENTITY,
    guild_id VARCHAR NOT NULL,
    channel_id VARCHAR NOT NULL
);