CREATE TABLE file_metadata (
    id VARCHAR(48) PRIMARY KEY
    ,user_id VARCHAR(48)
    ,filename TEXT NOT NULL
    ,file_size INTEGER
    ,file_type VARCHAR(100)
    ,filehash VARCHAR(255)
    ,chunks INTEGER
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
---
CREATE TABLE users (
	id VARCHAR(48) PRIMARY KEY
	,email TEXT NOT NULL
	,hashedpass TEXT
	,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
---
CREATE TABLE user_files (
	user_id VARCHAR(48)
	,file_id VARCHAR(48) NOT NULL
	,chunk_id INTEGER NOT NULL
	,chunk_blob_url TEXT NOT NULL
	,chunk_hash TEXT NOT NULL
	,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
---