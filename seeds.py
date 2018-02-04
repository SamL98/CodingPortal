import csv
import psycopg2

conn = psycopg2.connect("dbname=cp user=samlerner")
cur = conn.cursor()

user = 1
i = 0
with open('sample.csv') as f:
    reader = csv.reader(f, delimiter=',')
    for row in reader:
        if i == 0:
            i += 1
            continue
        if i >= 37:
            user = 2 
        cur.execute("insert into Modified \
        (article_id, text, lm, wave, subj_id, user_id) values \
        (%s, %s, %s, %s, %s, %s)", \
        (row[2], row[3], row[1], int(row[4]), row[0], user))
        i += 1

conn.commit()
cur.close()
conn.close()
