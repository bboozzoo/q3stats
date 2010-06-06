CREATE TABLE "player_aliases" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "player_id" INTEGER NOT NULL,
    "alias" TEXT NOT NULL UNIQUE
);
CREATE TABLE "match_player_items_stats" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "match_player_stats_id" INTEGER NOT NULL,
    "type" TEXT NOT NULL,
    "pickups" INTEGER NOT NULL,
    "time" INTEGER NOT NULL
);
CREATE TABLE "match_player_stats"
(
	"id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	"match_id" INTEGER NOT NULL,
	"alias_id" INTEGER NOT NULL,
	"score" INTEGER NOT NULL,
	"kills" INTEGER NOT NULL,
	"deaths" INTEGER NOT NULL,
	"suicides" INTEGER NOT NULL,
	"net" INTEGER NOT NULL,
	"damage_given" INTEGER NOT NULL,
	"damage_taken" INTEGER NOT NULL,
	"health_total" INTEGER NOT NULL,
	"armor_total" INTEGER NOT NULL
);
CREATE TABLE "match_player_weapon_stats" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "match_player_stats_id" INTEGER NOT NULL,
    "type" TEXT NOT NULL,
    "hits" INTEGER NOT NULL,
    "shots" INTEGER NOT NULL,
    "kills" INTEGER NOT NULL
);
CREATE TABLE "matches" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "date" INTEGER NOT NULL,
    "type" TEXT NOT NULL,
    "map" TEXT NOT NULL,
    "source_file_digest" TEXT NOT NULL UNIQUE
);
CREATE TABLE "players" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "name" TEXT NOT NULL UNIQUE
);
