--
-- music
--
ALTER TABLE music
    DROP COLUMN IF EXISTS file_path;

ALTER TABLE music
    ADD COLUMN music_file jsonb NOT NULL DEFAULT '{}';

ALTER TABLE music
    ADD COLUMN upload_id bigint NOT NULL ;

ALTER TABLE music
    ADD CONSTRAINT fk_music_upload_id FOREIGN KEY(upload_id) REFERENCES 
    upload(id) ON DELETE CASCADE ON UPDATE CASCADE;
    -- ADD CONSTRAINT fk_music_upload_id UNIQUE (upload_id);

