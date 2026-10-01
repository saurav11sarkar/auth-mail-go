-- Roll back the application to supplying IDs before removing this default.
ALTER TABLE users ALTER COLUMN id DROP DEFAULT;
