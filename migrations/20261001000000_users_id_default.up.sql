-- Generate IDs for new users; existing IDs remain unchanged.
-- Also upgrades databases that applied create_user before its UUID default was added.
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();
