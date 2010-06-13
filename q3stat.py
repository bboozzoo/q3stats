#!/usr/bin/env python

import cgi
import cgitb
import ConfigParser
import sqlite3
import string
import os
import Cookie

# to see backtraces
cgitb.enable()

# config
CONFIG_FILE='/etc/q3stat.conf'
# there are read from config
DB_PATH=''
IMG_PATH=''
SKEL_PATH=''
skeleton_files = { 'players-list' : 'main.xhtml',
                   'player-stats' : 'player-stats.xhtml' }

# read configuration
cp = ConfigParser.SafeConfigParser()
cp.read(CONFIG_FILE)
DB_PATH = cp.get('database', 'path')
IMG_PATH = cp.get('general', 'img-path')
SKEL_PATH = cp.get('general', 'skel-path')


def list_players(field_storage):
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    subst_str = ''
    c = conn.execute('select alias from player_aliases order by alias asc')
    for row in c:
        s = '<option>%s</option>\n' % (row['alias'])
        subst_str += s
        
    template_string = open(SKEL_PATH + os.sep + skeleton_files['players-list']).read()
    tmpl = string.Template(template_string)
    print tmpl.substitute(players_list = subst_str)
    
def show_stats(player_name):
    template_string = open(SKEL_PATH + os.sep + skeleton_files['player-stats']).read()
    tmpl = string.Template(template_string)
    print tmpl.substitute(player_name = player_name)

def show_match_stats():
    pass

fs = cgi.FieldStorage()

print 'Content-Type: text/html'
print


player_name = fs.getvalue('player-name', None)

if not player_name:
    list_players(fs)
else:
    show_stats(player_name)


#c = conn.execute('select count(*), map from matches group by map order by count(map) desc')
#for row in c:
#    print row
#cgi.test()
