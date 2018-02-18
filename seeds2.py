import csv
import psycopg2

conn = psycopg2.connect("postgres://khoibrbbxcrewx:ce8b6a7426e124d3cde178a327ba490722821e82720eee2f5206a19c9c76fa69@ec2-54-83-203-198.compute-1.amazonaws.com:5432/d5p8bq6iroraf4")
cur = conn.cursor()

user = 0
with open('originals.csv') as f:
    reader = csv.reader(f, delimiter=',')
    for row in reader:
        cur.execute("insert into Originals \
        (article_id, text, user_id) values \
        (%s, %s, %s)", \
        (row[0], row[1], user))

conn.commit()
cur.close()
conn.close()
