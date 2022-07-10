CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	media_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
ALTER TABLE ONLY comments
ADD CONSTRAINT comments_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE ONLY comments
ADD CONSTRAINT comments_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE SET NULL;
