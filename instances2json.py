#!/usr/bin/python3
import json
import sys

from bs4 import BeautifulSoup
from urllib3.util import Retry
import requests
from requests.adapters import HTTPAdapter


headers = {
    'User-Agent': 'Mozilla/5.0 MozhiInstanceFetcher/1.0 (+codeberg.org/aryak/mozhi)'
}

session = requests.Session()
session.headers.update(headers)
retries = Retry(total=5,
                connect=3,       # ConnectTimeoutError
                read=False,      # ReadTimeoutError or ProtocolError
                redirect=False,  # obvi, any redirections
                status=2,        # Status codes by server
                backoff_factor=1,
                backoff_max=30,  # Just to be sure that script don't go sleep for a minute
                respect_retry_after_header=False)
http_adapter = HTTPAdapter(max_retries=retries)
session.mount('https://', http_adapter)
session.mount('http://', http_adapter)


print("Getting HTML")
# Get the HTML from the page
r = session.get('https://codeberg.org/aryak/mozhi/src/branch/master/README.md') # XXX: relied on Codeberg

# Parse the HTML
soup = BeautifulSoup(r.text, 'html.parser')

print("Scraping started")

# Get table after Instances header
instances_h2 = soup.find("h2", string="Instances")
try:
    table = instances_h2.find_next_sibling("table")
except AttributeError:
    print("Instances header not found")
    sys.exit(-1)

# Get all rows and columns. Skip the first row because it's the header
rows = table.find_all('tr')[1:]

def get_net_type(url: str):
    url = url.strip("/")
    if url.endswith(".onion"):
        return "onion"
    elif url.endswith(".i2p"):
        return "i2p"
    elif url.endswith(".loki"):
        return "lokinet"
    return "link"

theJson = []

for row in rows:

    link = row.find_all('td')[0].find('a')['href']
    cloudflare = row.find_all('td')[1].text
    country = row.find_all('td')[2].text
    host = row.find_all('td')[3].text

    print("Scraping " + row.find_all('td')[0].find('a')['href'] + ' instance...')
    isCloudflare = cloudflare == "Yes"

    try:
        if get_net_type(url=link) == "link":
            r = session.get(link + '/', headers=headers)
        else:
            print(f"Non-clearnet mirror [{row.find_all('td')[0].find('a').get_text()}]. Skipping check")
        if r.status_code != 200:
            print("Error while fetching " + link + '/. We got a ' + str(r.status_code) + ' status code. Skipping...')
            continue
    except:
        print("Error while fetching " + link + '/. Skipping...')
        continue

    theJson.append({
        'country': country,
        get_net_type(url=link): link,
        'cloudflare': isCloudflare,
        'host': host,
    })


print("Scraping finished. Saving JSON...")

# save JSON
with open('instances.json', 'w') as outfile:
    json.dump(theJson, outfile, indent=4)
    print("File saved as instances.json")
