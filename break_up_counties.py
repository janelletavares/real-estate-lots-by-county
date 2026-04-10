import os
import json
from pathlib import Path


def main():
    f = open(os.path.join("input", "zips_by_county.json"),'r')
    data = json.load(f)
    f.close()

    for idx, state in enumerate(data):
        counties = data[state]
        s = state.replace(" ", "_")
        Path("not_done/"+s).mkdir(parents=True, exist_ok=True)
        for county in counties:
            fn = "not_done/{}/{}.json".format(s, county.replace(" ", "_"))
            with open(fn, 'w', encoding="utf-8") as f:
                json.dump(counties[county], f, indent=2)

if __name__ == "__main__":
    main()
