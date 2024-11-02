CREATE TABLE profiles (
    "user_id" INTEGER UNIQUE PRIMARY KEY NOT NULL,
    "first_name" text,
    "last_name" text,
    "description" text,
    "phone_number" text,
    "city" text,
    "photo_uuid" text
);
CREATE INDEX IF NOT EXISTS idx_user_id ON profiles (user_id);