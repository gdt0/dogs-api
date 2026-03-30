const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  await page.setViewportSize({ width: 1280, height: 900 });
  
  try {
    await page.goto('https://lee-dogs-api.loca.lt/', { waitUntil: 'networkidle', timeout: 30000 });
    await page.waitForTimeout(2000); // Wait for images
    await page.screenshot({ path: 'screenshot.png', fullPage: false });
    console.log('Screenshot saved to screenshot.png');
  } catch (e) {
    console.error('Error:', e.message);
  }
  
  await browser.close();
})();
