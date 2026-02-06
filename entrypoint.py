import time

from redfin import extract_land_listings
from selenium import webdriver
import json
import pprint
import csv
import os

def fetch_and_write(name, zipcodes):
    n = name.replace(" ", "_")
    for_sale = []
    sold = []

    try:

        for zipcode in zipcodes:
            print("{}: {}".format(name, zipcode))
            z = zipcode.zfill(5)

            driver = webdriver.Chrome()
            for_sale_url = "https://www.redfin.com/zipcode/{}/filter/property-type=land".format(z)
            driver.get(for_sale_url)
            land_for_sale = extract_land_listings(driver)
            for_sale += land_for_sale
            print(f"Found {len(land_for_sale)} for sale land listings")
            driver.quit()

            #time.sleep(30)
            driver = webdriver.Chrome()
            sold_url = "https://www.redfin.com/zipcode/{}/filter/property-type=land,include=sold-1yr".format(z)
            driver.get(sold_url)
            land_sold = extract_land_listings(driver)
            sold += land_sold
            print(f"Found {len(land_sold)} sold land listings")
            driver.quit()

            for_sale_in_range = [x for x in for_sale if int(x["price"]) >= 10000 and int(x["price"]) <= 300000]
            sold_in_range = [x for x in sold if int(x["price"]) >= 10000 and int(x["price"]) <= 300000]

    except Exception as e:
        print(e)
        return None
        #exit(1)


    filename = "output/{}_for_sale.json".format(n)
    p = pprint.pformat(for_sale, compact=True).replace("'", '"')
    with open(filename, 'w') as fp:
        fp.write(p)

    filename = "output/{}_sold.json".format(n)
    p = pprint.pformat(sold, compact=True).replace("'", '"')
    with open(filename, 'w') as fp:
        fp.write(p)

    return {'all_sold_count': len(sold), 'all_for_sale_count': len(for_sale), 'in_range_sold_count': len(sold_in_range), 'in_range_for_sale_count': len(for_sale_in_range)}

def main():
    f = open(os.path.join("input", "one.json"),'r')
    data = json.load(f)
    f.close()

    with open('output/counts_by_county.csv', 'w', newline='') as csvfile:
        fieldnames = ['state', 'county', 'all_sold_count', 'all_for_sale_count', 'in_range_sold_count',
                          'in_range_for_sale_count']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter =';')
        writer.writeheader()

        for idx, state in enumerate(data):
            counties = data[state]
            for county in counties:
                zipcodes = counties[county]
                # one file for all zipcodes
                d = fetch_and_write(county, zipcodes)
                if d != None:
                    d["state"] = state
                    d["county"] = county
                    writer.writerow(d)


if __name__ == "__main__":
    main()
