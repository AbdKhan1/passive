const puppeteer = require('puppeteer');
const fetch = require('node-fetch');

const fetchUsername = async () => {
    try {
        const response = await fetch('http://localhost:8080/get-username');
        if (!response.ok) {
            throw new Error('HTTP Error! Status: ' + response.status);
        }
        const username = await response.text();
        return username;
    } catch (error) {
        console.error('Error fetching username:', error);
        return '';
    }
};

const sendHtml = async (html) => {
    try {
        const response = await fetch('http://localhost:8080/receive-html', {
            method: 'POST',
            body: html,
        });
        if (!response.ok) {
            throw new Error('HTTP Error sending HTML. Status: ' + response.status);
        }
        console.log('HTML sent successfully.');
    } catch (error) {
        console.error('Error sending HTML:', error);
    }
};

(async () => {
    try {
        const browser = await puppeteer.launch();
        const page = await browser.newPage();
        const username = await fetchUsername();

        await page.goto('https://whatsmyname.app/?q=' + username);
        await page.waitForSelector('.card.text-white.bg-success.h-200');

        const dynamicHTML = await page.evaluate(() => {
            return document.querySelector('.card.text-white.bg-success.h-200').innerHTML;
        });

        await browser.close();

        await sendHtml(dynamicHTML);
    } catch (error) {
        console.error('An error occurred:', error);
    }
})();
