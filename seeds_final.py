import pandas as pd
import psycopg2

conn = psycopg2.connect("postgres://khoibrbbxcrewx:ce8b6a7426e124d3cde178a327ba490722821e82720eee2f5206a19c9c76fa69@ec2-54-83-203-198.compute-1.amazonaws.com:5432/d5p8bq6iroraf4")
cur = conn.cursor()

df = pd.read_csv('SciCoding_Set1.csv')

for _, row in df.iterrows():
    for user in [1, 2]:
        cur.execute("insert into Modified \
        (article_id, text, lm, wave, subj_id, user_id, coded) values \
        (%s, %s, %s, %s, %s, %s, false)", \
        (row['articlid'], row['response'], row['l/m'], int(row['wave']), row['user'], user))

conn.commit()
cur.close()
conn.close()
