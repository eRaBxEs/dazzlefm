--
-- Database : dazzlefmdb
--

--
-- page
--
CREATE TABLE "page" (
	id                serial NOT NULL,
  name              varchar(50) NOT NULL,
	route              varchar(50) NOT NULL,
	attributes        jsonb NOT NULL DEFAULT '{}',

	CONSTRAINT pk_page_id PRIMARY KEY (id)
);

INSERT INTO page VALUES (1, 'Home', '/home', '{}');
INSERT INTO page VALUES (2, 'About', '/about', '{}');
INSERT INTO page VALUES (3, 'Presenter', '/presenter', '{}');
INSERT INTO page VALUES (4, 'Schedule', '/schedule', '{}');
INSERT INTO page VALUES (5, 'News', '/news', '{}');
INSERT INTO page VALUES (6, 'Gallery', '/gallery', '{}');


--
-- presenter
--
CREATE TABLE presenter (
	id                serial NOT NULL,
	image             jsonb NOT NULL,
	name        varchar(150) NOT NULL,
	email             varchar(100) NOT NULL,
	description          text NOT NULL,

	CONSTRAINT pk_person_id PRIMARY KEY (id)
);

--
-- news
--
CREATE TABLE news (
	id             serial NOT NULL,
	title          varchar(200) NOT NULL,
	image        	 jsonb NOT NULL,
	date           date NOT NULL,
	author          varchar(150) NOT NULL,
	content        text NOT NULL,

	CONSTRAINT pk_news_id PRIMARY KEY (id)
);

--
-- schedule
--
CREATE TABLE schedule (
	id             serial NOT NULL,
	title        	 varchar(100) NOT NULL,
	presenter_id   bigint NOT NULL,
	days           jsonb NOT NULL,
	start_time          date NOT NULL,
	end_time          date NOT NULL,

	CONSTRAINT pk_schedule_id PRIMARY KEY (id),
	CONSTRAINT fk_schedule_presenter_id FOREIGN KEY(presenter_id)
		REFERENCES presenter (id) MATCH FULL
		ON DELETE NO ACTION ON UPDATE NO ACTION
);

--
-- gallery
--
CREATE TABLE gallery (
	id   serial NOT NULL,
	image jsonb NOT NULL,

	CONSTRAINT pk_gallery_id PRIMARY KEY (id)
);

--
-- contact
--
CREATE TABLE contact (
    id      serial NOT NULL,
    name    varchar(150) NOT NULL,
    email   varchar(250) NOT NULL,
    subject varchar(150) NOT NULL,
    date    timestamp NOT NULL DEFAULT localtimestamp,
    message text NOT NULL,

    CONSTRAINT pk_contact_id PRIMARY KEY (id)
);

--
-- user
--
CREATE TABLE "user" (
	id                serial NOT NULL,
  role              smallint NOT NULL,
  user_name         varchar(50) NOT NULL,
	first_name        varchar(50) NOT NULL,
	middle_name       varchar(50) NOT NULL,
	surname           varchar(50) NOT NULL,
  email             varchar(200) NOT NULL,
	title             varchar(100) NOT NULL,
  password          text NOT NULL,

	CONSTRAINT pk_user_id PRIMARY KEY (id),
  CONSTRAINT uq_user_user_name UNIQUE (user_name)
);


-- password 1234
INSERT INTO "user" 
	VALUES (1, 1, 'deled', 'dele', 'o', 'dada', 'deled@mailinator.com', 'mr', '$2a$14$oDmsATMGXCZLov.jSyjpDObT3.x0DeuuloY.spb0tXRqa8acy.dBK');


--
-- music
--
CREATE TABLE music (
	id             serial NOT NULL,
	title          varchar(200) NOT NULL,
	artist          varchar(150) NOT NULL,
	music_file       jsonb NOT NULL,
	rank           integer NOT NULL,
	

	CONSTRAINT pk_music_id PRIMARY KEY (id)
);

--
-- event
--
CREATE TABLE event (
    id           serial NOT NULL,
    title        varchar(200) NOT NULL,
    venue        text NOT NULL,
    status       integer NOT NULL,
    date         date NOT NULL,
    link         text NULL,

    CONSTRAINT pk_event_id PRIMARY KEY (id)
);