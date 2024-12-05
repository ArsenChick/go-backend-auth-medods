-- create a table
CREATE TABLE users (
    guid UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    refresh_token_hash TEXT
);

-- add test data
INSERT INTO users (guid, email)
VALUES
    ('525ee10f-e1e0-4509-b233-6b37821a4fef', 'mak3ntosh@hatanet.network'),
    ('a914c512-5ab1-4116-9064-895ec838e58b', 'clebel@yotomail.com'),
    ('5d63da45-c8fe-4176-a9da-788eb42dd96b', 'pemanning@available-home.com'),
    ('40341b70-c2f3-413c-9d8f-5e0572991de4', 'madisongracie@gmailbrt.com'),
    ('2279209e-3de6-484d-b369-527a7aec5b67', 'nguyenj@cumfoto.com');