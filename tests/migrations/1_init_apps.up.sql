INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'super-secret')
ON CONFLICT DO NOTHING;