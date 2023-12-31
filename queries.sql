INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (1, 'OTVOREN', 'OPENED', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (2, 'POPUNJEN', 'FILLED', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (3, 'U TOKU', 'IN_PROGRESS', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (id, name, code, created_at, updated_at) VALUES (4, 'ZATVOREN', 'CLOSED', now(), now()) ON CONFLICT DO NOTHING;


INSERT INTO base_seminar_types (name, code, created_at, updated_at) VALUES ('OSNOVNA OBUKA', 'BASIC', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO base_seminar_types (name, code, created_at, updated_at) VALUES ('DODATNA OBUKA', 'ADDITIONAL', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO base_seminar_types (name, code, created_at, updated_at) VALUES ('PERIODIČNA OBUKA', 'CYCLE', now(), now()) ON CONFLICT DO NOTHING;

INSERT INTO seminar_themes (base_seminar_type_id, name, code, number_of_days) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '1', '1', 1);
INSERT INTO seminar_themes (base_seminar_type_id, name, code, number_of_days) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '2', '2', 1);
INSERT INTO seminar_themes (base_seminar_type_id, name, code, number_of_days) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '3', '3', 1);
INSERT INTO seminar_themes (base_seminar_type_id, name, code, number_of_days) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '4', '4', 1);
INSERT INTO seminar_themes (base_seminar_type_id, name, code, number_of_days) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'CYCLE'), '5', '5', 1);
INSERT INTO seminar_themes (base_seminar_type_id, name, code, number_of_days) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'BASIC'), 'OSNOVNA OBUKA', 'BASIC', 21);
INSERT INTO seminar_themes (base_seminar_type_id, name, code, number_of_days) VALUES ((SELECT id FROM base_seminar_types WHERE code = 'ADDITIONAL'), 'DODATNA OBUKA', 'ADDITIONAL', 5);

INSERT INTO roles (id, name, code, created_at, updated_at) VALUES (1, 'Administrator', 'ADMINISTRATOR', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO roles (id, name, code, created_at, updated_at) VALUES (2, 'Predavač', 'PREDAVAC', now(), now()) ON CONFLICT DO NOTHING;

ALTER TABLE client_seminars ADD UNIQUE (client_id, seminar_id);

ALTER TABLE seminar_class_names ADD CONSTRAINT unique_theme_seminar_day_number_class_number UNIQUE (seminar_theme_id, day_number, class_number);
ALTER TABLE seminar_day_theme_names ADD CONSTRAINT unique_theme_seminar_day_number UNIQUE (seminar_theme_id, day_number);