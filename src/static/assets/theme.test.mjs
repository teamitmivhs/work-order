import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import vm from "node:vm";

class FakeElement {
    constructor() {
        this.attributes = new Map();
        this.dataset = {};
        this.listeners = new Map();
        this.themeLabel = null;
    }

    setAttribute(name, value) {
        this.attributes.set(name, String(value));
    }

    getAttribute(name) {
        return this.attributes.get(name) || null;
    }

    querySelector(selector) {
        return selector === "[data-theme-label]" ? this.themeLabel : null;
    }

    addEventListener(type, listener) {
        this.listeners.set(type, listener);
    }

    click() {
        this.listeners.get("click")?.({ currentTarget: this });
    }
}

const html = new FakeElement();
const body = new FakeElement();
const toggle = new FakeElement();
toggle.themeLabel = { textContent: "" };

const documentListeners = new Map();
const windowListeners = new Map();
const values = new Map();

const document = {
    body,
    documentElement: html,
    readyState: "loading",
    addEventListener(type, listener) {
        documentListeners.set(type, listener);
    },
    querySelectorAll(selector) {
        return selector === "[data-theme-toggle]" ? [toggle] : [];
    },
};

const localStorage = {
    getItem(key) {
        return values.get(key) || null;
    },
    setItem(key, value) {
        values.set(key, String(value));
    },
};

const window = {
    addEventListener(type, listener) {
        windowListeners.set(type, listener);
    },
};

const source = readFileSync(new URL("./theme.js", import.meta.url), "utf8");
vm.runInNewContext(source, { document, localStorage, window });

assert.equal(html.getAttribute("data-theme"), "light");

document.readyState = "complete";
documentListeners.get("DOMContentLoaded")?.();
assert.equal(toggle.getAttribute("aria-pressed"), "false");
assert.equal(toggle.themeLabel.textContent, "Mode Gelap");

toggle.click();
assert.equal(html.getAttribute("data-theme"), "dark");
assert.equal(body.getAttribute("data-theme"), "dark");
assert.equal(values.get("woTheme"), "dark");
assert.equal(values.get("theme"), "dark");
assert.equal(toggle.getAttribute("aria-pressed"), "true");
assert.equal(toggle.themeLabel.textContent, "Mode Terang");

values.set("woTheme", "light");
windowListeners.get("storage")?.({ key: "woTheme" });
assert.equal(html.getAttribute("data-theme"), "light");
assert.equal(body.getAttribute("data-theme"), "light");
assert.equal(toggle.getAttribute("aria-pressed"), "false");

console.log("theme controller behavior: ok");
