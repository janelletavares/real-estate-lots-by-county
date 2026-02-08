from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait, Select
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from selenium.webdriver.firefox.service import Service as FirefoxService
from selenium.webdriver.firefox.options import Options as FirefoxOptions
import os


def scrape_sold_lots_html(
    url: str,
    output_dir: str,
    timeout: int = 20,
):
    """
    Loads a listings page, applies filters:
      - Status: Sold
      - Home Type: Lots/Land only
      - More: Sold in last 12 months
    Then saves HTML for every results page.
    """

    os.makedirs(output_dir, exist_ok=True)

    #firefox_options = FirefoxOptions()
    #firefox_options.add_argument("--width=1920")
    #firefox_options.add_argument("--height=1080")
    #firefox_options.binary_location = '/usr/bin/firefox'
    #driver.execute_script("arguments[0].click();", element)

    #driver = webdriver.Firefox(service=Service(GeckoDriverManager().install()))
    #service = FirefoxService(executable_path='bin/geckodriver')

    #driver = webdriver.Firefox(
    #    service=service,
    #    options=firefox_options
    #)

    driver = webdriver.Chrome()
    wait = WebDriverWait(driver, timeout)

    try:
        driver.get(url)

        # -------------------------
        # STATUS MENU → SOLD
        # -------------------------
        status_button = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//button[.//text()[contains(., 'Status')]]")
            )
        )
        status_button.click()

        sold_radio = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//label[.//text()[contains(., 'Sold')]]//input[@type='radio']")
            )
        )
        driver.execute_script("arguments[0].click();", sold_radio)

        apply_button = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//button[.//text()[contains(., 'Apply')]]")
            )
        )
        apply_button.click()

        # -------------------------
        # HOME TYPE → LOTS / LAND ONLY
        # -------------------------
        home_type_button = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//button[.//text()[contains(., 'Home type')]]")
            )
        )
        home_type_button.click()

        # Uncheck all checkboxes first
        checkboxes = driver.find_elements(By.XPATH, "//input[@type='checkbox']")
        for cb in checkboxes:
            if cb.is_selected():
                driver.execute_script("arguments[0].click();", cb)

        lots_checkbox = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//label[.//text()[contains(., 'Lot')]]//input[@type='checkbox']")
            )
        )
        if not lots_checkbox.is_selected():
            driver.execute_script("arguments[0].click();", lots_checkbox)

        apply_button = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//button[.//text()[contains(., 'Apply')]]")
            )
        )
        apply_button.click()

        # -------------------------
        # MORE → SOLD IN LAST → 12 MONTHS
        # -------------------------
        more_button = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//button[.//text()[contains(., 'More')]]")
            )
        )
        more_button.click()

        sold_in_last_select = wait.until(
            EC.presence_of_element_located(
                (By.XPATH, "//label[.//text()[contains(., 'Sold in last')]]//select")
            )
        )

        Select(sold_in_last_select).select_by_visible_text("12 months")

        apply_button = wait.until(
            EC.element_to_be_clickable(
                (By.XPATH, "//button[.//text()[contains(., 'Apply')]]")
            )
        )
        apply_button.click()

        # -------------------------
        # PAGINATION + HTML SAVE
        # -------------------------
        page_num = 1

        while True:
            wait.until(EC.presence_of_element_located((By.TAG_NAME, "body")))

            html_path = os.path.join(output_dir, f"results_page_{page_num}.html")
            with open(html_path, "w", encoding="utf-8") as f:
                f.write(driver.page_source)

            print(f"Saved {html_path}")

            # Try to click "Next"
            try:
                next_button = wait.until(
                    EC.element_to_be_clickable(
                        (By.XPATH, "//a[.//text()[contains(., 'Next')]]")
                    )
                )
                driver.execute_script("arguments[0].click();", next_button)
                page_num += 1

            except (TimeoutException, NoSuchElementException):
                # No more pages
                break

    finally:
        driver.quit()
