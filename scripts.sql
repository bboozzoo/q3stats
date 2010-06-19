select count(alias_id), from match_player_stats group by match_player_stats.match_id ;
select * from match_player_stats where match_id = 1;
select matches.id, matches.date, matches.map, count(match_player_stats.alias_id) from match_player_stats inner join matches on match_player_stats.match_id = matches.id group by matches.id order by matches.date desc;
select alias_id, score, player_aliases.alias from match_player_stats inner join player_aliases on match_player_stats.alias_id = player_aliases.id where match_id = 33 order by score desc limit 1;
select * from match_player_stats where alias_id = 9;
select * from player_aliases where id = 9;
select match_player_weapon_stats.type, match_player_weapon_stats.hits, match_player_weapon_stats.shots, match_player_weapon_stats.kills 
from match_player_weapon_stats 
inner join match_player_stats on match_player_stats.id = match_player_weapon_stats.match_player_stats_id where alias_id = 9 group by type;
select match_player_weapon_stats.type, match_player_weapon_stats.hits, match_player_weapon_stats.shots, match_player_weapon_stats.kills 
from match_player_weapon_stats 
inner join match_player_stats on match_player_stats.id = match_player_weapon_stats.match_player_stats_id where alias_id = 9;

-- select weapon stats by player
select match_player_weapon_stats.type, sum(match_player_weapon_stats.hits), sum(match_player_weapon_stats.shots), sum(match_player_weapon_stats.kills) 
from match_player_weapon_stats where match_player_weapon_stats.match_player_stats_id in (select id from match_player_stats where alias_id = 1)
group by match_player_weapon_stats.type;

-- select items stats by player
select match_player_items_stats.type, sum(match_player_items_stats.pickups), sum(match_player_items_stats.time)
from match_player_items_stats where match_player_items_stats.match_player_stats_id in (select id from match_player_stats where alias_id = 1)
group by match_player_items_stats.type;

select id from match_player_stats where alias_id = 9;