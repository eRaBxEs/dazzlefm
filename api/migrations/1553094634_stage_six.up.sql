--
--- schedule
--

ALTER TABLE schedule
    ADD COLUMN copresenter_id bigint NULL ;

ALTER TABLE schedule
    ADD CONSTRAINT fk_schedule_copresenter_id FOREIGN KEY(copresenter_id) REFERENCES 
    presenter(id) ON DELETE CASCADE ON UPDATE CASCADE ;

INSERT INTO presenter VALUES (0, '{}', 'None', 'None', 'None');


-- ALTER TABLE schedule
--     DROP COLUMN IF EXISTS start_time ;

-- ALTER TABLE schedule
--     DROP COLUMN IF EXISTS end_time ;

ALTER TABLE schedule
    ALTER COLUMN start_time TYPE time with time zone;

ALTER TABLE schedule
    ALTER COLUMN end_time TYPE time with time zone;