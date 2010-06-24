#!/usr/bin/env python

import cgi
import cgitb
import ConfigParser
import sqlite3
import string
import os
import Cookie
import datetime

# to see backtraces
cgitb.enable()

SCRIPT_NAME=None

class StatsError(Exception):
    '''exception wrapper'''
    def __init__(self, what, error = True):
        self.what = what
    def __str__(self):
        return str(self.what)
        
class Stats:
    '''wrapper for all relevant config files/template access/db access'''
    CONFIG_FILE='/etc/q3stat.conf'
    TEMPLATES = { 'main-page' : 'main-template.xhtml',
                  'player-stats' : 'player-stats-template.xhtml',
                  'player-stats-weapon-stats-table-entry' : 'player-stats-weapon-stats-table-entry-template.xhtml',
                  'player-stats-item-stats-table-entry' : 'player-stats-item-stats-table-entry-template.xhtml',
                  'recent-match-list-table-entry': 'recent-match-list-table-entry-template.xhtml',
                  'player-list-table-entry': 'player-list-table-entry-template.xhtml',
                  'alias-list-table-entry' : 'alias-list-table-entry-template.xhtml',
                  'add-player' : 'add-player-template.xhtml',
                  'modify-player-alias-checkbox' : 'modify-player-alias-checkbox-template.xhtml'
                  }
    IMAGES = { 'G' : 'iconw_gauntlet.png',
               'SG' : 'iconw_shotgun.png',
               'RL' : 'iconw_rocket.png',
               'RG' : 'iconw_railgun.png',
               'LG' : 'iconw_lightning.png',
               'PG' : 'iconw_plasma.png',
               'BFG' : 'iconw_bfg.png',
               'GL' : 'iconw_grenade.png',
               'MG' : 'iconw_machinegun.png',
               'MH' : 'iconh_mega.png',
               'RA' : 'iconr_red.png',
               'YA' : 'iconr_yellow.png',
               'GA' : 'iconr_green.png',
               'Quad' : 'quad.png',
               'Invis' : 'invis.png',
               'BattleSuit' : 'envirosuit.png',
               'Regen' : 'regen.png',
               'Flight' : 'flight.png'
               }
    
    def __init__(self):
        self.load_config()
        self.__db_conn = None

    def load_config(self):
        cp = ConfigParser.SafeConfigParser()
        cp.read(self.CONFIG_FILE)
        
        self.__db_path = cp.get('database', 'path')
        self.__document_root = cp.get('general', 'document-root')
        self.__templates_path = cp.get('general', 'templates-path')
        
    def db_get(self):
        if not self.__db_conn:
            self.__db_conn = sqlite3.connect(self.__db_path)
            self.__db_conn.row_factory = sqlite3.Row
        return self.__db_conn
            
    def db_close(self):
        if self.__db_conn:
            self.__db_conn.close()
    
    def template_get(self, which):
        if not self.TEMPLATES.has_key(which):
            raise StatsError('template %s not found' % (which))
        try:
            s = open(self.__templates_path + os.sep + self.TEMPLATES[which]).read()
        except Except, e:
            raise StatsError('error reading template, reason: %s' % (str(e)))
        return s

    def image_path_get(self, name):
        if not self.IMAGES.has_key(name):
            raise StatsError('image %s not found' % (name))
        return self.__document_root + self.IMAGES[name]
                  
                             
def get_aliases_table(stats):
    '''return HTML formatted player aliases table'''
    conn = stats.db_get()
    html_aliases_table = ''
    c = conn.execute('select id, alias, player_id from player_aliases order by alias asc')
    tmpl = string.Template(stats.template_get('alias-list-table-entry'))
    for row in c:
        alias_name = row[1]
        alias_id = row[0]
        player_id = row[2]
        html_aliases_table += tmpl.substitute(script_name = SCRIPT_NAME,
                                              player_id = player_id,
                                              alias_name = alias_name)
    return html_aliases_table

def get_players_table(stats):
    '''return HTML formatted players table'''
    conn = stats.db_get()
    html_players_table = ''
    c = conn.execute('select id, name from players order by name asc')
    tmpl = string.Template(stats.template_get('player-list-table-entry'))
    for row in c:
        player_name = row[1]
        player_id = row[0]
        html_players_table += tmpl.substitute(script_name = SCRIPT_NAME,
                                              player_id = player_id,
                                              player_name = player_name)
    if len(html_players_table) == 0:
        html_players_table = '<tr><td>None</td></tr>'
    return html_players_table

def get_players_list(stats, match_id):
    '''return list of players ordered by score'''
    players_list = []
    conn = stats.db_get()
    c = conn.execute('''select alias_id, score, 
                        player_aliases.alias 
                        from 
                        match_player_stats 
                        inner join 
                        player_aliases 
                        on match_player_stats.alias_id = player_aliases.id 
                        where match_id = ? 
                        order by score desc''', (match_id,))
    for row in c:
        players_list.append(row[2])
    return players_list

def get_recent_matches_table(stats):
    '''return HTML formatted table of recent matches'''
    conn = stats.db_get()
    html_matches_table = ''
    c = conn.execute('''select matches.id, matches.date, 
                        matches.map, count(match_player_stats.alias_id),
                        matches.winner_alias_id
                        from 
                        match_player_stats 
                        inner join 
                        matches 
                        on match_player_stats.match_id = matches.id 
                        group by matches.id 
                        order by matches.date desc''')
    tmpl = string.Template(stats.template_get('recent-match-list-table-entry'))
    for row in c:
        dt = datetime.datetime(1900,1,1)
        dt = dt.fromtimestamp(row[1])
        when = dt.strftime('%A, %d. %B %Y %I:%M%p')
        match_map = row[2]
        players_count = row[3]
        match_id = row[0]
        who_won = get_player_name(stats, row[4])
        html_matches_table += tmpl.substitute(script_name = SCRIPT_NAME,
                                              match_id = match_id,
                                              when = when,
                                              map = match_map,
                                              who_won = who_won,
                                              how_many_players = players_count)
    return html_matches_table

def get_player_name(stats, player_id):
    '''return string - player name given player alias id'''
    conn = stats.db_get()
    c = conn.execute('select alias from player_aliases where id = ?', (player_id,))
    row = c.fetchone()
    return row[0]

def get_player_matches_count(stats, player_id):
    '''return string - count of matches played by player of given id'''
    conn = stats.db_get()
    c = conn.execute('select count(*) from match_player_stats where alias_id = ?', (player_id, ))
    row = c.fetchone()
    return row[0]

def get_player_wins_count(stats, player_id):
    '''return string - count of player wins'''
    conn = stats.db_get()
    c = conn.execute('select count(*) from matches where winner_alias_id = ?', (player_id,))
    row = c.fetchone()
    return row[0]

def get_player_frags_deaths_suicides_count(stats, player_id):
    '''return tuple of strings (frags, deaths, suicides) for given player id'''
    conn = stats.db_get()
    c = conn.execute('''select 
                      sum(kills), 
                      sum(deaths), 
                      sum(suicides) 
                      from 
                      match_player_stats 
                      where alias_id = ?''', (player_id,))
    row = c.fetchone()
    return row[0], row[1], row[2]

def get_player_hits_shots_accuracy(stats, player_id):
    '''return tupe of strings(hits, shots, accuracy) for given player id'''
    hits = 0
    shots = 0
    accuracy = 'N/A'
    conn = stats.db_get()
    c = conn.execute('''select 
                        sum(hits), 
                        sum(shots) 
                        from 
                        match_player_weapon_stats 
                        where 
                        match_player_stats_id 
                        in 
                        (select id from match_player_stats where alias_id = ?)''', (player_id,))
    row = c.fetchone()
    hits = row[0]
    shots = row[1]
    if int(shots) > 0:
        accuracy = '%d%%' % (int(100.0 * float(hits) / float(shots)))
    return hits, shots, accuracy

def get_player_vital_stats(stats, player_id):
    '''output dictionary containing string indexed player stats fields (also used as key names):
    - matches
    - wins
    - losses === (matches - wins)
    - kills
    - deaths
    - suicides
    - hits
    - shots
    - accuracy
    '''
    vstats = {}
    for s in ['matches', 'wins', 'losses', 'kills', 'deaths', 'suicides', 'hits', 'shots', 'accuracy']:
        vstats[s] = '0'
        
    vstats['matches'] = get_player_matches_count(stats, player_id)
    vstats['wins'] = get_player_wins_count(stats, player_id)
    vstats['losses'] = '%d' % (int(vstats['matches']) - int(vstats['wins']))
    vstats['kills'], vstats['deaths'], vstats['suicides'] = get_player_frags_deaths_suicides_count(stats, player_id)
    vstats['hits'], vstats['shots'], vstats['accuracy'] = get_player_hits_shots_accuracy(stats, player_id)
    return vstats
    

def get_player_weapon_stats_table(stats, player_id):
    '''return string with HTML formatted table with weapon statistics'''
    html_weapon_stats_table = ''
    # 'weapon type' => tuple(accuracy, efficiency, kills) - all strings
    zero_tuple = ('N/A', 'N/A', '0')
    weapons_order = ['G', 'MG', 'SG', 'GL', 'RL', 'LG', 'RG', 'PG', 'BFG']
    weapon_stats = {}
    for w in weapons_order:
        weapon_stats[w] = zero_tuple
                     
    conn = stats.db_get()
    c = conn.execute('''select match_player_weapon_stats.type, 
                        sum(match_player_weapon_stats.hits), 
                        sum(match_player_weapon_stats.shots), 
                        sum(match_player_weapon_stats.kills) 
                        from 
                        match_player_weapon_stats 
                        where 
                        match_player_weapon_stats.match_player_stats_id 
                        in 
                        (select id from match_player_stats where alias_id = ?)
                        group by 
                        match_player_weapon_stats.type''', (player_id, ))
    for row in c:
        weapon_type = row[0]
        hits = float(row[1])
        shots = float(row[2])
        kills = row[3]
        accuracy = 'N/A'
        efficiency = 'N/A'
        if shots > 0.0:
            accuracy = '%d%%' % (100.0 * hits / shots)
            efficiency = '%.02f' % ((float(kills) / shots))
        weapon_stats[weapon_type] = (accuracy, efficiency, kills)
    
    tmpl = string.Template(stats.template_get('player-stats-weapon-stats-table-entry'))
    for w_key in weapons_order:
        html_weapon_stats_table += tmpl.substitute(weapon_img = stats.image_path_get(w_key),
                                                   weapon_kills = weapon_stats[w_key][2],
                                                   weapon_accuracy = weapon_stats[w_key][0],
                                                   weapon_efficiency = weapon_stats[w_key][1])
    return html_weapon_stats_table

def get_player_item_stats_table(stats, player_id):
    '''return string with HTML formatted table with items statistics'''
    html_items_stats_table = ''
    # 'item type' => tuple(pickups, time) - all strings
    zero_tuple = ('0', 'N/A')
    items_order = ['RA', 'YA', 'GA', 'Quad', 'BattleSuit', 'Invis', 'Regen', 'Flight']
    items_stats = {}
    for i in items_order:
        items_stats[i] = zero_tuple
                     
    conn = stats.db_get()
    c = conn.execute('''select match_player_items_stats.type, 
                        sum(match_player_items_stats.pickups), 
                        sum(match_player_items_stats.time)
                        from 
                        match_player_items_stats 
                        where  
                        match_player_items_stats.match_player_stats_id 
                        in  
                        (select id from match_player_stats where alias_id = ?)
                        group by 
                        match_player_items_stats.type;''', (player_id, ))
    for row in c:
        item_type = row[0]
        pickups = row[1]
        hold_time = 'N/A'
        if row[2] != 0:
            hold_time = '%ds' % (int(row[2]))
        items_stats[item_type] = (pickups, hold_time)
    
    tmpl = string.Template(stats.template_get('player-stats-item-stats-table-entry'))
    for i_key in items_order:
        html_items_stats_table += tmpl.substitute(item_img = stats.image_path_get(i_key),
                                                   item_pickups = items_stats[i_key][0],
                                                   item_time = items_stats[i_key][1])
    return html_items_stats_table

def get_alias_checkbox_list(stats):
    html_alias_checkbox_list = ''
    tmpl = string.Template(stats.template_get('modify-player-alias-checkbox'))
    conn = stats.db_get()
    c = conn.execute('select id, alias from player_aliases')
    for row in c:
        alias_id = row[0]
        alias_name = row[1]
        html_alias_checkbox_list += tmpl.substitute(alias_id = alias_id,
                                                     alias_name = alias_name)
    return html_alias_checkbox_list
        
def output_main_page(stats):
    html_aliases_table = get_aliases_table(stats)
    html_players_table = get_players_table(stats)
    html_matches_table = get_recent_matches_table(stats)
    tmpl = string.Template(stats.template_get('main-page'))
    print tmpl.substitute(script_name = SCRIPT_NAME,
                          aliases_list = html_aliases_table,
                          players_list = html_players_table,
                          matches_list = html_matches_table)

def output_show_add_player_page(stats):
    tmpl = string.Template(stats.template_get('add-player'))
    print tmpl.substitute(script_name = SCRIPT_NAME)
                          
def output_show_modify_player_page(stats, cgi_fs):
    player_id = cgi_fs.getvalue('player_id', None)
    if not player_id:
        raise StatsError('player_id not set')
    html_alias_checkbox_list = get_alias_checkbox_list(stats)
    tmpl = string.Template(stats.template_get('modify-player'))
    print tmpl.substitute(script_name = SCRIPT_NAME,
                          aliases_list = html_alias_checkbox_list,
                          player_id = player_id)

def output_match_stats_page(stats, cgi_fs):
    print '<h3>match stats</h3>'

def output_player_stats_page(stats, cgi_fs):
    '''output player stats page
    if player_id == 0 then there the player is not defined do not show anything'''
    player_id = cgi_fs.getvalue('player_id', None)
    if not player_id:
        raise StatsError('player_id not set')
    html_weapon_stats_table = get_player_weapon_stats_table(stats, player_id)
    html_item_stats_table = get_player_item_stats_table(stats, player_id)
    html_player_name = get_player_name(stats, player_id)
    player_vital_stats = get_player_vital_stats(stats, player_id)
    tmpl = string.Template(stats.template_get('player-stats'))
    print tmpl.substitute(player_name = html_player_name,
                          weapon_stats = html_weapon_stats_table,
                          item_stats = html_item_stats_table,
                          player_matches = player_vital_stats['matches'],
                          player_wins = player_vital_stats['wins'],
                          player_losses = player_vital_stats['losses'],
                          player_kills = player_vital_stats['kills'],
                          player_deaths = player_vital_stats['deaths'],
                          player_suicides = player_vital_stats['suicides'],
                          player_hits = player_vital_stats['hits'],
                          player_shots = player_vital_stats['shots'],
                          player_accuracy = player_vital_stats['accuracy'])


# read configuration
SCRIPT_NAME = os.environ['SCRIPT_NAME']
#cgi_get_request = cgi.parse()
cgi_fs = cgi.FieldStorage()
stats_handler = Stats()

print 'Content-Type: text/html'
print
print 'query string: %s' % (os.environ['REQUEST_URI'])
#print '%s' % (repr(cgi_get_request))
print '<p>fs: %s</p>' % (repr(cgi_fs))
#cgi.print_environ()
req = cgi_fs.getvalue('req', 'none')
print '<p>%s</p>' % (req)
if req == 'show-match-stats':
    output_match_stats_page(stats_handler, cgi_fs)
elif req == 'show-player-stats':
    output_player_stats_page(stats_handler, cgi_fs)
elif req == 'show-add-player':
    output_show_add_player_page(stats_handler)
elif req == 'show-modify-player':
    output_show_modify_player_page(stats_handler, cgi_fs)
elif req == 'add-player':
    pass
elif req == 'modify-player':
    pass
else:
    output_main_page(stats_handler)


