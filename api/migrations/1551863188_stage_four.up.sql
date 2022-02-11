--
-- user
--
ALTER TABLE "user"
    DROP COLUMN IF EXISTS title ;

--
-- schedule
--

ALTER TABLE schedule
    DROP COLUMN IF EXISTS start_time ;

ALTER TABLE schedule
    DROP COLUMN IF EXISTS end_time ;

ALTER TABLE schedule
    ADD COLUMN start_time time NOT NULL;

ALTER TABLE schedule
    ADD COLUMN end_time time NOT NULL;

ALTER TABLE schedule
    DROP COLUMN IF EXISTS presenters ;

ALTER TABLE schedule
    DROP COLUMN IF EXISTS days ;

ALTER TABLE schedule
    ADD COLUMN presenter_id bigint NOT NULL ;

ALTER TABLE schedule
    ADD COLUMN day integer NOT NULL;

ALTER TABLE schedule
    ADD CONSTRAINT fk_schedule_presenter_id FOREIGN KEY(presenter_id) REFERENCES 
    presenter(id) ON DELETE CASCADE ON UPDATE CASCADE ;