CREATE TABLE transactions (
                        id SERIAL PRIMARY KEY,
                        hash VARCHAR(50),
                        from_addr VARCHAR(50),
                        to_addr VARCHAR(50),
                        value VARCHAR(50),
                        contract_addr VARCHAR(50),
                        block_number BIGINT
);