import os
import time
import json
from playwright.sync_api import sync_playwright

BASE_URL = "http://localhost:5174"

def run():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context()
        page = context.new_page()

        # Login (using last created user from previous runs might be hard, so create new user)
        # Actually, let's just use the flow again.

        email = f"screenshot_{int(time.time())}@example.com"
        password = "password123"

        page.goto(f"{BASE_URL}/signup")
        page.fill("#name", "Screenshot User")
        page.fill("#email", email)
        page.fill("#password", password)
        page.click("button[type='submit']")
        page.wait_for_url(f"{BASE_URL}/")

        # Create Inventory
        page.goto(f"{BASE_URL}/select-inventory")
        page.fill("#name", "Screenshot Inv")
        page.click("button:has-text('Create Household')")
        page.wait_for_url(f"{BASE_URL}/")

        # Get Inventory ID from API? Or just proceed.
        # We need to create a canonical product "Milk".
        # We need token.
        storage = page.evaluate("localStorage.getItem('auth-storage')")
        token = json.loads(storage)['state']['token']

        # Get inventory ID
        resp = page.request.get(f"http://localhost:8080/inventories", headers={"Authorization": f"Bearer {token}"})
        inv_id = resp.json()[0]['id']

        # Create Canonical Product
        page.request.post(f"http://localhost:8080/inventories/{inv_id}/canonical-products",
                          data={"name": "Milk"},
                          headers={"Authorization": f"Bearer {token}", "Content-Type": "application/json"})

        # Create List
        page.goto(f"{BASE_URL}/shopping-lists")
        page.click("button:has-text('New List')")
        page.fill("#name", "My List")
        page.click("button:has-text('Create')")
        time.sleep(1)
        page.click("text=My List")

        # Add Item
        page.click("button:has-text('Add Item')")
        page.fill("input[placeholder='Search products...']", "Milk")
        time.sleep(1)
        page.click("text=Milk")
        time.sleep(0.5)
        # Click Add Item in dialog
        page.click("div[role='dialog'] button:has-text('Add Item')") # Or dialog
        time.sleep(1)

        # Mark as shopped
        page.reload()
        time.sleep(1)
        page.click("input[type='checkbox']")

        # Complete Shopping
        page.click("button:has-text('Complete Shopping')")
        page.wait_for_selector("select")

        # Select Generic
        page.locator("select").nth(1).select_option(value="__CREATE_GENERIC__")
        time.sleep(1)

        # Select Default Variant
        page.locator("select").nth(2).select_option(value="__CREATE_DEFAULT_VARIANT__")
        time.sleep(1)

        # Take screenshot
        page.screenshot(path="verification/generic_transaction.png")
        print("Screenshot taken: verification/generic_transaction.png")

        browser.close()

if __name__ == "__main__":
    run()
