#!/usr/bin/python3
import requests
import json
from bs4 import BeautifulSoup

print("Getting HTML")

headers = {
    'User-Agent': 'Mozilla/5.0 MozhiInstanceFetcher/1.0 (+codeberg.org/aryak/mozhi)'
}

# Get the HTML from the page
r = requests.get('https://codeberg.org/aryak/mozhi', headers=headers)

# Parse the HTML
soup = BeautifulSoup(r.text, 'html.parser')

print("Scraping started")

# Get tables
tables = soup.find_all('table')

# Get table with header 'Master Branch'
table = tables[1]

# Get all rows and columns. Skip the first row because it's the header
rows = table.find_all('tr')[1:]

theJson = []

for row in rows:

    link = row.find_all('td')[0].find('a')['href']
    cloudflare = row.find_all('td')[1].text
    country = row.find_all('td')[2].text
    host = row.find_all('td')[3].text

    print("Scraping " + row.find_all('td')[0].find('a')['href'] + ' instance...')
    if cloudflare == 'Yes':
        isCloudflare = True
    else:
        isCloudflare = False

    try:
        r = requests.get(link + '/', headers=headers)
        if r.status_code != 200:
            print("Error while fetching " + link + '/. We got a ' + str(r.status_code) + ' status code. Skipping...')
            continue
    except:
        print("Error while fetching " + link + '/. Skipping...')
        continue

    theJson.append({
        'country': country,
        'link': link,
        'cloudflare': isCloudflare,
        'host': host,
    })


print("Scraping finished. Saving JSON...")

# save JSON
with open('instances.json', 'w') as outfile:
    json.dump(theJson, outfile, indent=4)
    print("File saved as instances.json")
