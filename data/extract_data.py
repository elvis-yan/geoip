import csv
import json


def extract_data():
    table = {}
    with open('../GeoLite2-City-CSV_20180703/GeoLite2-City-Locations-en.csv') as f:
        rd = csv.DictReader(f)
        # print('-----------> ', next(rd))
        for row in rd:
            table[row['geoname_id']] = (row['city_name'], row['country_iso_code'])
    
    data = []
    with open('../GeoLite2-City-CSV_20180703/GeoLite2-City-Blocks-IPv4.csv') as f:
        rd = csv.DictReader(f)
        # next(rd)
        for row in rd:
            # if row['geoname_id'] == '':
            data.append( (row['network'], table.get(row['geoname_id'], ('--','--'))) )
    print('#network:', len(data))

    with open('data.csv', 'w') as f:
        wt = csv.DictWriter(f, ['network', 'city_name', 'country_iso_code'])
        wt.writeheader()
        for (nw, (city, country)) in data:
            wt.writerow({'network': nw, 'city_name': city, 'country_iso_code': country})


if __name__ == '__main__':
    extract_data()