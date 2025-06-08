-- docker exec -it GO_project-RestAPI psql -U myuser -d GO_project -c "SHOW database"

-- Создание базы, если ее нет
CREATE DATABASE my_restapi;

-- Подключение к новой БД
\c my_restapi

-- Создание таблиц
CREATE TABLE public.users (
	id int NOT NULL,
	"name" varchar(128) NOT NULL,
	age smallint NOT NULL,
	is_student bool NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);


-- Начальные данные
INSERT INTO users ("name", age, is_student) 
VALUES 
    ('admin', 14, true),
    ('user1', 21, false)
ON CONFLICT (username) DO NOTHING;