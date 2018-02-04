import csv
import psycopg2

conn = psycopg2.connect("postgres://khoibrbbxcrewx:ce8b6a7426e124d3cde178a327ba490722821e82720eee2f5206a19c9c76fa69@ec2-54-83-203-198.compute-1.amazonaws.com:5432/d5p8bq6iroraf4")
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
