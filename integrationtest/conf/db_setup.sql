CREATE DATABASE streamlistner;
CREATE ROLE streamlistner WITH LOGIN PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE streamlistner TO streamlistner;
