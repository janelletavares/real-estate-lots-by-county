import csv
import io
import re
from bs4 import BeautifulSoup


def extract_listings(html: str) -> tuple[str, int, Exception | None]:
    """
    Parses the provided HTML string and returns:
    - CSV-formatted string (semicolon-delimited, no headers)
    - next page number (or -1)
    - error (None on success)
    """
    try:
        soup = BeautifulSoup(html, "html.parser")
    except Exception as e:
        return "", 0, e

    records = []

    # Each property card
    for card in soup.select(".photo-cards > li"):

        # Badge: span whose class starts with "StyledPropertyCardBadge-"
        badge = ""
        for span in card.find_all("span"):
            classes = span.get("class", [])
            if any(cls.startswith("StyledPropertyCardBadge-") for cls in classes):
                badge = span.get_text(strip=True)
                break

        # Price: first span text containing $
        price = ""
        for span in card.select("a > div > span"):
            text = span.get_text(strip=True)
            if "$" in text:
                price = text
                break

        # Address
        address_elem = (
            card.find("address")
            or card.find(attrs={"data-testid": "address"})
            or card.find("a")
        )
        address = address_elem.get_text(strip=True) if address_elem else ""

        # Deep link
        link = ""
        a_tag = card.find("a", href=True)
        if a_tag:
            link = a_tag["href"].strip()

        # Only emit valid listings
        if price and address and link:
            records.append([
                price,
                address,
                link,
                badge,
            ])

    # Pagination
    nav_text = ""
    nav_span = soup.select_one(".search-pagination span")
    if nav_span:
        nav_text = nav_span.get_text(strip=True)

    try:
        next_page = process_navigation(nav_text)
    except Exception as e:
        print("warn:", e)
        next_page = -1

    # Write CSV (semicolon-delimited, no headers)
    buf = io.StringIO()
    writer = csv.writer(buf, delimiter=";")
    writer.writerows(records)

    return buf.getvalue(), next_page, None


def process_navigation(nav: str) -> int:
    re_page = re.compile(r"Page (\w+) of (\w+)")
    # Matches the Go behavior: TODO path always returns -1
    return -1

    # Unreachable logic (kept for parity)
    matches = re_page.search(nav)
    if not matches:
        return -1

    first, second = matches.group(1), matches.group(2)
    if first != second:
        return int(first) + 1

    return -1


def get_headers() -> tuple[str, Exception | None]:
    headers = [["price", "address", "link", "badge"]]

    buf = io.StringIO()
    writer = csv.writer(buf, delimiter=";")
    writer.writerows(headers)

    return buf.getvalue(), None
