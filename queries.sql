INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (1, 'OTVOREN', 'OPENED', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (2, 'POPUNJEN', 'FILLED', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (3, 'U TOKU', 'IN_PROGRESS', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (4, 'ZATVOREN', 'CLOSED', now(), now()) ON CONFLICT DO NOTHING;


INSERT INTO base_seminar_types (name, code, created_at, updated_at) VALUES ('OSNOVNA OBUKA', 'BASIC', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO base_seminar_types (name, code, created_at, updated_at) VALUES ('DODATNA OBUKA', 'ADDITIONAL', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO base_seminar_types (name, code, created_at, updated_at) VALUES ('PERIODIČNA OBUKA', 'CYCLE', now(), now()) ON CONFLICT DO NOTHING;

INSERT INTO seminar_themes (base_seminar_type_id, name, code) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '1', '1');
INSERT INTO seminar_themes (base_seminar_type_id, name, code) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '2', '2');
INSERT INTO seminar_themes (base_seminar_type_id, name, code) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '3', '3');
INSERT INTO seminar_themes (base_seminar_type_id, name, code) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '4', '4');
INSERT INTO seminar_themes (base_seminar_type_id, name, code) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '5', '5');