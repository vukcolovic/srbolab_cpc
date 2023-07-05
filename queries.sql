INSERT INTO seminar_statuses (name, code, created_at, updated_at) VALUES ('U TOKU', 'IN_PROGRESS', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (name, code, created_at, updated_at) VALUES ('NA ČEKANJU', 'PENDING', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_statuses (name, code, created_at, updated_at) VALUES ('ZAVRŠEN', 'FINISHED', now(), now()) ON CONFLICT DO NOTHING;

INSERT INTO seminar_types (name, code, created_at, updated_at) VALUES ('OSNOVNA OBUKA', 'BASIC', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_types (name, code, created_at, updated_at) VALUES ('DODANA OBUKA', 'ADDITIONAL', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_types (name, code, created_at, updated_at) VALUES ('PERIODICNA OBUKA 1', 'CYCLE_1', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_types (name, code, created_at, updated_at) VALUES ('PERIODICNA OBUKA 2', 'CYCLE_2', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_types (name, code, created_at, updated_at) VALUES ('PERIODICNA OBUKA 3', 'CYCLE_3', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_types (name, code, created_at, updated_at) VALUES ('PERIODICNA OBUKA 4', 'CYCLE_4', now(), now()) ON CONFLICT DO NOTHING;
INSERT INTO seminar_types (name, code, created_at, updated_at) VALUES ('PERIODICNA OBUKA 5', 'CYCLE_5', now(), now()) ON CONFLICT DO NOTHING;

