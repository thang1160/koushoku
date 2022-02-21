CREATE TABLE IF NOT EXISTS artist (
  id   BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS artist_slug_uindex ON artist(slug);
CREATE UNIQUE INDEX IF NOT EXISTS artist_name_uindex ON artist(name);

CREATE TABLE IF NOT EXISTS circle (
  id BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS circle_slug_uindex ON circle(slug);
CREATE UNIQUE INDEX IF NOT EXISTS circle_name_uindex ON circle(name);

CREATE TABLE IF NOT EXISTS tag (
  id   BIGSERIAL PRIMARY KEY,
  slug VARCHAR(32) NOT NULL DEFAULT NULL,
  name VARCHAR(32) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS tag_slug_uindex ON tag(slug);
CREATE UNIQUE INDEX IF NOT EXISTS tag_name_uindex ON tag(name);

CREATE TABLE IF NOT EXISTS magazine (
  id BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS magazine_slug_uindex ON magazine(slug);
CREATE UNIQUE INDEX IF NOT EXISTS magazine_name_uindex ON magazine(name);

CREATE TABLE IF NOT EXISTS parody (
  id BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS parody_slug_uindex ON parody(slug);
CREATE UNIQUE INDEX IF NOT EXISTS parody_name_uindex ON parody(name);

CREATE TABLE IF NOT EXISTS archive (
  id                BIGSERIAL PRIMARY KEY,
  path              TEXT NOT NULL DEFAULT NULL,
  
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  published_at      TIMESTAMP,

  title             VARCHAR(1024) NOT NULL DEFAULT NULL,
  slug              VARCHAR(1024) NOT NULL DEFAULT NULL,
  pages             SMALLINT NOT NULL DEFAULT NULL,
  size              VARCHAR(16) NOT NULL DEFAULT NULL,

  circle_id         BIGINT DEFAULT NULL REFERENCES circle(id) ON DELETE SET NULL,
  magazine_id       BIGINT DEFAULT NULL REFERENCES magazine(id) ON DELETE SET NULL,
  parody_id         BIGINT DEFAULT NULL REFERENCES parody(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS archive_path_uindex ON archive(path);
CREATE INDEX IF NOT EXISTS archive_title_index ON archive(title);
CREATE INDEX IF NOT EXISTS archive_slug_index ON archive(slug);
CREATE INDEX IF NOT EXISTS archive_created_at_index ON archive(created_at);
CREATE INDEX IF NOT EXISTS archive_updated_at_index ON archive(updated_at);
CREATE INDEX IF NOT EXISTS archive_published_at_index ON archive(published_at);
CREATE INDEX IF NOT EXISTS archive_title_index ON archive(title);
CREATE INDEX IF NOT EXISTS archive_circle_id_index ON archive(circle_id);
CREATE INDEX IF NOT EXISTS archive_magazine_id_index ON archive(magazine_id);
CREATE INDEX IF NOT EXISTS archive_parody_id_index ON archive(parody_id);

CREATE TABLE IF NOT EXISTS archive_artists (
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  artist_id BIGINT NOT NULL DEFAULT NULL REFERENCES artist(id) ON DELETE CASCADE,
  PRIMARY KEY(archive_id, artist_id)
);

CREATE INDEX IF NOT EXISTS archive_artists_archive_id_index ON archive_artists(archive_id);
CREATE INDEX IF NOT EXISTS archive_artists_artist_id_index ON archive_artists(artist_id);

CREATE TABLE IF NOT EXISTS archive_tags (
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  tag_id    BIGINT NOT NULL DEFAULT NULL REFERENCES tag(id) ON DELETE CASCADE,
  PRIMARY KEY(archive_id, tag_id)
);

CREATE INDEX IF NOT EXISTS archive_tags_archive_id_index ON archive_tags(archive_id);
CREATE INDEX IF NOT EXISTS archive_tags_tag_id_index ON archive_tags(tag_id);
