CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              hash VARCHAR(100) UNIQUE,
                              from_addr VARCHAR(100),
                              to_addr VARCHAR(100),
                              value VARCHAR(255),
                              block_number BIGINT
);