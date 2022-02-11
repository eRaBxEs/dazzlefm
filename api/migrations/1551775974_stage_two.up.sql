--
-- upload
--
CREATE TABLE upload (
	id    serial NOT NULL,
	path  varchar(400) NOT NULL,
	size  bigint NOT NULL,


	CONSTRAINT pk_upload_id PRIMARY KEY (id)
);


--
-- music
--
ALTER TABLE music
    ALTER COLUMN music_file TYPE varchar(400);

ALTER TABLE music
    RENAME COLUMN music_file TO file_path;


--
-- schedule
--
ALTER TABLE schedule
    DROP CONSTRAINT IF EXISTS fk_schedule_presenter_id;

ALTER TABLE schedule
    DROP COLUMN IF EXISTS presenter_id;

ALTER TABLE schedule
    ADD COLUMN presenters jsonb NOT NULL DEFAULT '{}';

