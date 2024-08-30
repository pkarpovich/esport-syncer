-- +migrate Up
ALTER TABLE events ADD COLUMN team_id INTEGER DEFAULT NULL;
ALTER TABLE events ADD COLUMN game_type TEXT DEFAULT NULL;
UPDATE events SET team_id = 1669, game_type = 'dota2';

-- +migrate Down
ALTER TABLE events DROP COLUMN team_id;
ALTER TABLE events DROP COLUMN game_type;
