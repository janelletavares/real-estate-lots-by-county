# Redfin scraper

This tool gathers specific real estate stats per zip code and then organizes by county into a csv file.

## Usage

### Pre-requisites
* python3
* potentially virtualenv

### Running for the first time

You can download this source code from Github, which is easier than doing a git clone.
Save this to a convenient location on the computer you would like to run the script: [code](https://github.com/janelletavares/real-estate-lots-by-county/archive/refs/heads/main.zip)

It's easiest to create a virtual environment so that the system python doesn't conflict with the python in this workspace.

These instructions came from a Linux system.  Your mileage may vary.

To set up the environment and install the dependencies, do this at the top of this repo:
```
python3 -m venv redfin
source redfin/bin/activate
pip install -r requirements.txt
```

Next, drag the desired states that you want to rerun from the done/ dir to the not_done dir.  Each StateName directory contains county files.  Those county files are a JSON list of zip codes.  *This file structure must remain for the script to work.*  Do not drag county files out of their respective state directory.

Lastly, to start the script:
```
./run.sh
```

When the script is done checking each zip code, you can file the data by county in the file `output/counts_by_county.csv`. The extension is csv, but the delimiter used was a `;` as some values contain commas, so it would not be the best delimiter.

### Running in parallel
There is no automatic mechanism for increasing the concurrency in this project.  Several attempts were abandoned.  There is one known way to run through one or more states faster; you can start the script in two or more copies of this project manually.  Just make sure you do not pull the same counties into the not_done directory or you will be duplicating work and wasting time. If you do split up a single state, make sure that the JSON county files are always under a directory with the appropriate state name.
