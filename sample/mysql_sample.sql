CREATE TABLE all_types_table (
    id INT PRIMARY KEY,
    
    -- Integer types
    col_tinyint TINYINT,
    col_smallint SMALLINT,
    col_int INT,
    col_bigint BIGINT,
    
    -- String types
    col_varchar VARCHAR(50),
    col_char CHAR(10),
    col_text TEXT,
    
    -- Boolean
    col_bool BOOLEAN,
    
    -- Date and Time
    col_date DATE,
    col_datetime DATETIME,
    col_timestamp TIMESTAMP,
    created_at TIMESTAMP, -- Common audit column
    
    -- Floating point types
    col_float FLOAT,
    col_double DOUBLE,
    col_decimal DECIMAL(10,2)
);
