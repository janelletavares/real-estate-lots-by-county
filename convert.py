import csv
import json
from collections import defaultdict

INPUT_CSV = "input/zips_by_county.csv"
OUTPUT_JSON = "input/zips_by_county.json"


def csv_to_county_zip_json(input_csv, output_json):
    # Nested dictionary:
    # { state: { county: set(zipcodes) } }
    data = defaultdict(lambda: defaultdict(set))

    with open(input_csv, newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f)

        for row in reader:
            state = row["County State"].strip()
            county = row["County"].strip()
            zipcode = row["ZIP"].strip()

            data[state][county].add(zipcode)

    # Convert sets to sorted lists for JSON serialization
    output = {
        state: {
            county: sorted(zips)
            for county, zips in counties.items()
        }
        for state, counties in data.items()
    }

    with open(output_json, "w", encoding="utf-8") as f:
        json.dump(output, f, indent=2)

    print(f"Wrote {output_json}")


if __name__ == "__main__":
    csv_to_county_zip_json(INPUT_CSV, OUTPUT_JSON)
