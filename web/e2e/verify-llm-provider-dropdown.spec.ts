import { test, expect } from '@playwright/test'

test('Verify LLM Provider dropdown in Inventory Settings', async ({ page }) => {
  const uniqueId = Date.now()
  const email = `user${uniqueId}@example.com`
  const password = 'password123'
  const name = 'Test User'

  // 1. Sign up
  await page.goto('/signup')
  await page.fill('#name', name)
  await page.fill('#email', email)
  await page.fill('#password', password)
  await page.click('button[type="submit"]')

  // 2. Handle inventory selection if redirected
  // Depending on redirect speed, we might land on /select-inventory
  await page.waitForURL(/\/select-inventory/)

  // Create a new inventory
  await page.fill('input[placeholder="e.g. Home, Office"]', 'Test Household')
  await page.click('button:has-text("Create Household")')

  // 3. Navigate to Inventory page
  await page.waitForURL('/') // Dashboard
  await page.goto('/inventory')

  // 4. Open Settings
  const settingsButton = page.getByTitle('Settings')
  await expect(settingsButton).toBeVisible()
  await settingsButton.click()

  // 5. Verify LLM Provider dropdown
  const providerSelect = page.locator('select#llm_provider') // Expecting a select element
  await expect(providerSelect).toBeVisible()

  // Verify options
  const options = await providerSelect.locator('option').allTextContents()
  // We expect "Select a provider", "Gemini", "OpenAI"
  expect(options).toContain('Gemini')
  expect(options).toContain('OpenAI')

  // Verify values
  const optionValues = await providerSelect
    .locator('option')
    .evaluateAll((opts) => opts.map((o) => (o as HTMLOptionElement).value))
  expect(optionValues).toContain('gemini')
  expect(optionValues).toContain('openai')

  // Verify default disabled option
  const defaultOption = providerSelect.locator('option[disabled]')
  await expect(defaultOption).toHaveText('Select a provider')
})
