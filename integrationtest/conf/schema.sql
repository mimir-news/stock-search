CREATE TABLE stock (
  symbol VARCHAR(20) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  is_active BOOLEAN,
  total_count INTEGER,
  updated_at TIMESTAMP
);

CREATE TABLE tweet (
    id VARCHAR(50) PRIMARY KEY,
    text VARCHAR(500),
    language VARCHAR(10),
    author_id VARCHAR(50),
    author_followers INTEGER,
    created_at TIMESTAMP
);

CREATE TABLE tweet_link (
    id INTEGER PRIMARY KEY,
    url VARCHAR(200),
    tweet_id VARCHAR(50) REFERENCES tweet(id)
);

CREATE TABLE tweet_symbol (
    id INTEGER PRIMARY KEY,
    symbol VARCHAR(20) REFERENCES stock(symbol),
    tweet_id VARCHAR(50) REFERENCES tweet(id)
);

INSERT INTO stock(symbol, name, is_active, total_count, updated_at) VALUES
    ('TWTR', 'Twitter, Inc.', TRUE, 0, CURRENT_TIMESTAMP),
    ('T', 'AT&T, Inc.', TRUE, 0, CURRENT_TIMESTAMP);

INSERT INTO tweet(id, text) VALUES 
    ('0', 'TWTR tweet 1'),
    ('1', 'T tweet 1'),
    ('2', 'TWTR tweet 2'),
    ('3', 'TWTR tweet 3'),
    ('4', 'T tweet 2'),
    ('5', 'BOTH tweet 2');

INSERT INTO tweet_symbol(id, symbol, tweet_id) VALUES
    (1, 'TWTR', '0'),
    (2, 'T', '1'),
    (3, 'TWTR', '2'),
    (4, 'TWTR', '3'),
    (5, 'T', '4'),
    (6, 'TWTR', '5'),
    (7, 'T', '5');