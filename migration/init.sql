CREATE TABLE
    "performer" (
        "id" UUID UNIQUE DEFAULT gen_random_uuid (),
        "username" TEXT NOT NULL UNIQUE,
        "youtube_channel" TEXT,
        "twitch_channel" TEXT,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX "performer_index_0" ON "performer" (
    "id",
    "username",
    "youtube_channel",
    "twitch_channel"
);

CREATE TABLE
    "request" (
        "id" UUID UNIQUE DEFAULT gen_random_uuid (),
        "channel_id" TEXT NOT NULL,
        "music" TEXT NOT NULL,
        "artist" TEXT,
        "status" TEXT NOT NULL,
        "requester" TEXT NOT NULL,
        "requester_platform" TEXT NOT NULL,
        "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "moditied_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX "request_index_0" ON "request" ("id", "channel_id");