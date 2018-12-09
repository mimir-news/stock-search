CREATE ROLE stocksearch WITH LOGIN PASSWORD 'password';
GRANT CONNECT ON DATABASE streamlistner TO stocksearch;
GRANT USAGE ON SCHEMA public TO stocksearch;
GRANT SELECT ON tweet_symbol TO stocksearch;
GRANT INSERT, UPDATE, SELECT ON stock TO stocksearch;