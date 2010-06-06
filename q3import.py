#!/usr/bin/env python

import os
import os.path
import sys
import sqlite3
import re
import time
import string
import hashlib
import xml.dom.minidom as xdm

STATS_DB_NAME = 'stat.db'
SCHEMA_SCRIPT = 'schema.sql'

def get_attr_typed(node_w_attr, attr_name, type_cast = str):
    return type_cast(node_w_attr.attributes[attr_name].value)

class Q3StatImportFilesList:
    '''wrapper to build a list of files to import'''
    def __init__(self):
        self.__xml_file_pattern = re.compile('^[0-9]{2}_[0-9]{2}_[0-9]{2}\.xml$')
        self.__import_files = []
    def path_walk_clbk(self, arg, dir_name, file_names):
        # filter only files which match the naming pattern
        matching_files = filter(self.__xml_file_pattern.match, file_names)
        # prepend dir path to file names
        import_files = map(lambda x: dir_name + os.sep + x, matching_files) 
        self.__import_files += import_files
    def get_import_files(self):
        return self.__import_files
    def add_import_file(self, path):
        self.__import_files.append(path)
    
class Q3FileListImporter:
    '''a wrapper for importing files list'''
    def __init__(self, db_conn, q3_files_list):
        '''db_conn - sqlite database connection
        q3_files_list - instance of Q3StatImportFilesList'''
        self.db = db_conn
        self.q3ifl = q3_files_list

    def import_files(self):
        print 'starting file import'
        for f in self.q3ifl.get_import_files():
            self.import_single_file(f)

    def import_single_file(self, filename):
        print 'importing file %s' % (filename)
        imp = Q3FileImporter(self.db, filename)
        imp.import_file()

class Q3FileImporter:
    '''actual file importer
    will handle reading xml and adding proper entries into DB'''
    def __init__(self, db_conn, file_name):
        self.db = db_conn
        self.filename = file_name

    def get_digest(self):
        f = open(self.filename)
        digest = hashlib.md5(f.read())
        self.digest = digest.hexdigest()

    def import_file(self):
        '''start file import'''
        self.get_digest()
        self.document = xdm.parse(self.filename)
        match_id = self.import_match_info()
        self.import_players_info(match_id)
        print 'commit DB'
        self.db.commit()

    def import_match_info(self):
        '''add match info
        return match ID'''
        match_node = self.document.getElementsByTagName('match')[0]
        # map name should always be lowercase
        match_map = string.lower(get_attr_typed(match_node, 'map'))
        # game type alwasy uppercase
        match_type = string.upper(get_attr_typed(match_node, 'type'))
        # date format is year/month/day hour:minute:second
        match_date = int(time.mktime(
                time.strptime(get_attr_typed(match_node, 'datetime'), 
                              '%Y/%m/%d %H:%M:%S')
                ))
        print 'map: %s type %s date: %d file digest: %s' % (match_map, match_type, match_date, self.digest)
        # dump everything into DB
        c = self.db.execute('insert into matches(date, type, map, source_file_digest) values (?, ?, ?, ?)',
                        (match_date, match_type, match_map, self.digest,))
        return c.lastrowid

    
    def import_players_info(self, match_id):
        players = self.document.getElementsByTagName('player')
        for p in players:
           player_stats_id = self.import_player_stats_info(p, match_id)
           self.import_player_weapons_info(p, player_stats_id)
           self.import_player_items_info(p, player_stats_id)

    def get_player_and_alias_id(self, player_name):
        '''find player ID and alias ID for given player/nick name
        return tuple (player ID, alias ID)'''
        player_id = 0
        alias_id = 0
        c = self.db.execute('select id, player_id from player_aliases where alias = ?',
                            (player_name,))
        row = c.fetchone()
        if not row:
            #entry does not exist
            print 'player %s does not exist yet, adding' % (player_name)
            player_id = self.add_player_name(player_name)
            alias_id = self.add_alias_for_player_id(player_name, player_id)
        else:
            player_id = row['player_id']
            alias_id = row['id']

        return (player_id, alias_id)

    def add_player_name(self, player_name):
        '''add player name to players table
        return player ID'''
        c = self.db.execute('insert into players(name) values(?)', (player_name,))
        return c.lastrowid # player name entry ID
    
    def add_alias_for_player_id(self, alias_name, player_id):
        '''add alias for given player_id
        return alias_id'''
        c = self.db.execute('insert into player_aliases(player_id, alias) values(?, ?)',
                            (player_id, alias_name,))
        return c.lastrowid # alias entry ID
        
        
    def import_player_stats_info(self, player_element, match_id):
        '''import player statistics
        return player_stats_id which is used as key in weapon and items statistics'''
        player_stats_id = 0
        # this will be filled in the for loop
        player_stats = {
            'score' : 0,
            'kills' : 0,
            'deaths' : 0,
            'suicides' : 0,
            'net' : 0,
            'damage_given' : 0,
            'damage_taken' : 0,
            'health_total' : 0,
            'armor_total' : 0 
            }
        # translate names from xml to column names
        # if entry is missing here, then lowercase name is
        # used
        translate_stats = {
            'DamageGiven' : 'damage_given',
            'DamageTaken' : 'damage_taken',
            'HealthTotal' : 'health_total',
            'ArmorTotal' : 'armor_total'
            }

        player_name = get_attr_typed(player_element, 'name')
        print 'player:', player_name
        stats = player_element.getElementsByTagName('stat')
        
        # first fill up the statistics dictionary
        for s in stats:
            name = get_attr_typed(s, 'name')
            value = get_attr_typed(s, 'value', int)
            print 's name: %s value: %s' % (name, value)
            # first try to translate the stat name to get the key into player_stats
            stat_key = translate_stats.get(name)
            if not stat_key:
                # translation failed, try lowercase name
                stat_key = string.lower(name)
            if player_stats.has_key(stat_key):
                # the key was found, update the value
                player_stats[stat_key] = value
            else:
                print 'WARNING: player %s has unknown stat %s' % (player_name, name)

        # now insert the entry into database
        print player_stats
        (player_id, alias_id) = self.get_player_and_alias_id(player_name)
        
        # first value is ID, set to NULL - DB fills that automatically
        c = self.db.execute('insert into match_player_stats values(NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)',
                            (match_id, alias_id,
                             player_stats['score'], player_stats['kills'],
                             player_stats['deaths'], player_stats['suicides'],
                             player_stats['net'], player_stats['damage_given'],
                             player_stats['damage_taken'], player_stats['health_total'],
                             player_stats['armor_total'],))
        # dump everything into DB
        return c.lastrowid # stats entry ID

    def import_player_weapons_info(self, player_element, player_stats_id):
        '''import list of weapon statistics'''
        weapons = player_element.getElementsByTagName('weapon')
        for w in weapons:
            name = string.upper(get_attr_typed(w, 'name'))
            hits = get_attr_typed(w, 'hits', int)
            shots = get_attr_typed(w, 'shots', int)
            kills = get_attr_typed(w, 'kills', int)
            print 'w name: %s hits: %d shots: %d kills: %d' % (name, hits, shots, kills)
            # dump everything into DB
            c = self.db.execute('insert into match_player_weapon_stats(match_player_stats_id, type, hits, shots, kills) values (?, ?, ?, ?, ?)',
                                (player_stats_id, name, hits, shots, kills, ))

    def import_player_items_info(self, player_element, player_stats_id):
        '''import list of items statistics'''
        items = player_element.getElementsByTagName('item')
        for i in items:
            name = get_attr_typed(i, 'name')
            pickups = get_attr_typed(i, 'pickups', int)
            keep_time = 0
            try: 
                keep_time = get_attr_typed(i, 'time', int)
            except:
                pass
            print 'i name: %s pickups: %d keep_time: %d' % (name, pickups, keep_time)
            # dump everything into DB
            c = self.db.execute('insert into match_player_items_stats(match_player_stats_id, type, pickups, time) values (?, ?, ?, ?)',
                                (player_stats_id, name, pickups, keep_time,))

def init_db_schema(db):
    '''initialize DB schema, read the init script from file'''
    schemaf = open(SCHEMA_SCRIPT)
    db.executescript(schemaf.read())

def main():
    q3ifl = Q3StatImportFilesList()
    if len(sys.argv) == 1:
        stats_path = os.path.expanduser('~/.q3a/cpma/stats')
        if not os.path.exists(stats_path):
            print 'Q3 CPMA stats dir (%s) not found' % (stats_path)
            sys.exit(1)
        
        print 'stats path: %s' % (stats_path)
        os.path.walk(stats_path, q3ifl.path_walk_clbk, None)
    else:
        # shortcut to import just one file
        q3ifl.add_import_file(sys.argv[1])

    for f in q3ifl.get_import_files():
        print f
   
    create_schema = False
    if not os.path.exists(STATS_DB_NAME):
        create_schema = True

    db = sqlite3.connect(STATS_DB_NAME)
    db.row_factory = sqlite3.Row
    if create_schema:
        print 'DB is fresh, creating schema'
        init_db_schema(db)
    
    q3i = Q3FileListImporter(db, q3ifl)
    q3i.import_files()

if __name__ == '__main__':
    main()

