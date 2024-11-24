CREATE TABLE
    "performer" (
        "id" UUID NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        "username" TEXT NOT NULL UNIQUE,
        "youtube_channel" TEXT,
        "twitch_channel" TEXT,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX "performer_index_0" ON "performer" ("id", "username");

CREATE TABLE
    "performance" (
        "id" UUID NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        "performer_id" UUID NOT NULL,
        "status" TEXT NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX "performance_index_0" ON "performance" ("id", "performer_id");

CREATE TABLE
    "request" (
        "id" UUID NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        "performance_id" UUID NOT NULL,
        "music" TEXT NOT NULL,
        "artist" TEXT,
        "status" TEXT NOT NULL,
        "requester" TEXT NOT NULL,
        "requester_platform" TEXT NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX "request_index_0" ON "request" ("id", "performance_id");

ALTER TABLE "performance" ADD FOREIGN KEY ("performer_id") REFERENCES "performer" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE "request" ADD FOREIGN KEY ("performance_id") REFERENCES "performance" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;