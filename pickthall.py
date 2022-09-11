from time import sleep
import requests
from bs4 import BeautifulSoup

quran = []

with open('surahs.csv') as f:
    lines = f.readlines()

    for line in lines:
        words = line.strip().split(',')
        surah_id = int(words[0])
        surah_sz = int(words[1])
        quran.append((surah_id, surah_sz))

added = set()

with open('pickthall.csv') as f:
    lines = f.readlines()

    for line in lines:
        words = line.strip().split(',')
        surah_id = int(words[0])
        verse_id = int(words[1])
        added.add((surah_id, verse_id))

csv = open('pickthall.csv', 'a')

for surah_id, surah_sz in quran:
    for verse_id in range(1, surah_sz + 1):
        if (surah_id, verse_id) in added:
            continue

        r = requests.get('https://corpus.quran.com/translation.jsp?chapter={}&verse={}'.format(surah_id, verse_id))

        soup = BeautifulSoup(r.text, 'html.parser')
        text = soup.body.find(text='Pickthall').parent.parent.text.lstrip('Pickthall:').strip()
        row = '{},{},"{}"'.format(surah_id, verse_id, text)

        csv.write('{}\n'.format(row))

        print('added: surah={}, verse={}'.format(surah_id, verse_id))
