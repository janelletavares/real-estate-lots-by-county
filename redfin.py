from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import Select
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from selenium.webdriver.common.keys import Keys
import time
import pprint


def extract_land_listings(driver, timeout=5):
    wait = WebDriverWait(driver, timeout)
    listings = []

    # 1. Search for provided string
    #search_box = wait.until(
    #    EC.element_to_be_clickable((By.XPATH, "//*[@id='search-box-input']"))
    #)
    #search_box.clear()
    #search_box.send_keys(search_text)
    #search_box.send_keys(Keys.ENTER)

    # print("looking for filters button")
    # # 2. Open Filters menu
    # filters_button = wait.until(
    #     EC.element_to_be_clickable((By.XPATH, "//*[@id='WideSidepaneHeader--container']/div/form/div[5]/button"))
    # )
    # filters_button.click()
    #
    # # 3. Select "Sold"
    # sold_toggle = wait.until(
    #     EC.element_to_be_clickable((By.XPATH, "//*[@id='filterContent']/div/div[1]/div/div/div/div[3]/span"))
    # )
    # sold_toggle.click()
    #
    # # 4. Select Home Type → Land
    # home_type_button = wait.until(
    #     EC.element_to_be_clickable((By.XPATH, "//*[@id='filterContent']/div/div[5]/div[2]/div/div[4]"))
    # )
    # home_type_button.click()
    #
    # #land_checkbox = wait.until(
    # #    EC.element_to_be_clickable((By.XPATH, "//label[.//span[text()='Land']]//input"))
    # #)
    # #if not land_checkbox.is_selected():
    # #    land_checkbox.click()
    #
    # # Close Home Type dropdown (click outside or re-click button)
    # home_type_button.click()
    #
    # # 5. Sold Within → Last 1 year
    # sold_within_dropdown = wait.until(
    #     EC.element_to_be_clickable((By.XPATH, "//*[@id='filterContent']/div/div[6]/div[2]/div/div/div/label/div/select"))
    # )
    # sold_within_dropdown.click()
    #
    # last_year_option = wait.until(
    #     EC.element_to_be_clickable((By.XPATH, "//span[text()='Last 1 year']"))
    # )
    # last_year_option.click()
    #
    # # 6. Click "See homes"
    # see_homes_button = wait.until(
    #     EC.element_to_be_clickable((By.XPATH, "//*[@id='searchForm']/form/div[2]/div/div/div/button/span"))
    # )
    # see_homes_button.click()

    while True:
        # Wait for listings to render
        try:
            wait.until(
                EC.presence_of_all_elements_located(
                    # //*[@id="results-display"]/div[4]/div/div/div/div[1]/div/div
                    (By.XPATH, "//*[@id='results-display']/div[4]/div/div/div/div[1]/div/div")
                )
            )
        except Exception as e:
            print("no listings found: {}".format(e))
            return listings

        only_results = driver.find_element(By.XPATH, "//*[@id='results-display']/div[4]/div/div[1]")

        cards = only_results.find_elements(By.CLASS_NAME, "bp-Homecard") # "bp-InteractiveHomecard" bp-Homecard  MapHomecardWrapper
        #print("{} cards\n".format(len(cards)))

        for i, card in enumerate(cards):
            #print(i)
            try:
                # #NearbyMapHomeCard_0 > div > div > div.bp-Homecard__Content.bp-Homecard__Content--custom.cleanAnchor > a
                # //*[@id="NearbyMapHomeCard_0"]/div/div/div[3]/a
                link_el = card.find_element(By.XPATH, "//div/div/div[3]/a") #"bp-Homecard__Address")
                url = link_el.get_attribute("href")
                address = link_el.text.strip()
            except Exception as e:
                print("cannot find url or address {}".format(str(e)))
                continue

            try:
                price = card.find_element(
                    By.CLASS_NAME, "bp-Homecard__Price--value"
                ).text.replace("$", "").replace(",", "")
            except:
                price = None

            try:
                date = card.find_element(
                    By.CLASS_NAME, "bp-Homecard__Sash"
                ).text
            except:
                date = None
            if date == "ABOUT THIS HOME":
                date = None

            try:
                other = card.find_element(
                    By.CLASS_NAME, "KeyFactsExtension"
                ).text
            except:
                other = None

            if price != None:
                listings.append({
                    "address": address,
                    "price": price,
                    "date": date,
                    "url": url,
                    "other": other
                })
            del address, price, date, url


        # Attempt to go to next page
        try:
            next_button = wait.until(
                EC.element_to_be_clickable(
                    (By.XPATH, "//*[@id='results-display']/div[4]/div/div[3]/div/button[2]")
                )
            )
            #print("found the next page button")

            # If disabled, we are done
            if "disabled" in next_button.get_attribute("class").lower():
                break

            #driver.execute_script("arguments[0].scrollIntoView(true);", next_button)
            #time.sleep(1)
            next_button.click()

            # Allow React to re-render
            time.sleep(2)

        except Exception:
            # print("no more pages")
            #return listings
            break

    #pprint.pprint(listings)
    return listings
