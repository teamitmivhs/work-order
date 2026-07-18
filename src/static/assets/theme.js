(function () {
    "use strict";

    var STORAGE_KEY = "woTheme";
    var LEGACY_STORAGE_KEY = "theme";
    var DEFAULT_THEME = "light";

    function normalizeTheme(value) {
        return value === "dark" ? "dark" : DEFAULT_THEME;
    }

    function getSavedTheme() {
        try {
            return normalizeTheme(
                localStorage.getItem(STORAGE_KEY) ||
                    localStorage.getItem(LEGACY_STORAGE_KEY),
            );
        } catch (_error) {
            return DEFAULT_THEME;
        }
    }

    function updateThemeControls(theme) {
        var isDark = theme === "dark";
        var actionLabel = isDark
            ? "Aktifkan mode terang"
            : "Aktifkan mode gelap";

        document.querySelectorAll("[data-theme-toggle]").forEach(function (control) {
            control.setAttribute("aria-pressed", String(isDark));
            control.setAttribute("aria-label", actionLabel);
            control.setAttribute("title", actionLabel);

            var label = control.querySelector("[data-theme-label]");
            if (label) {
                label.textContent = isDark ? "Mode Terang" : "Mode Gelap";
            }
        });
    }

    function applyTheme(theme) {
        var nextTheme = normalizeTheme(theme || getSavedTheme());
        document.documentElement.setAttribute("data-theme", nextTheme);

        if (document.body) {
            document.body.setAttribute("data-theme", nextTheme);
        }

        updateThemeControls(nextTheme);
    }

    function saveTheme(theme) {
        try {
            localStorage.setItem(STORAGE_KEY, theme);
            localStorage.setItem(LEGACY_STORAGE_KEY, theme);
        } catch (_error) {
            // The current page can still switch theme when storage is unavailable.
        }
    }

    function toggleTheme() {
        var currentTheme = normalizeTheme(
            document.documentElement.getAttribute("data-theme") ||
                getSavedTheme(),
        );
        var nextTheme = currentTheme === "dark" ? "light" : "dark";

        saveTheme(nextTheme);
        applyTheme(nextTheme);
    }

    function bindThemeControls() {
        document.querySelectorAll("[data-theme-toggle]").forEach(function (control) {
            if (control.dataset.themeBound === "true") {
                return;
            }

            control.dataset.themeBound = "true";
            control.addEventListener("click", toggleTheme);
        });
    }

    function initializeTheme() {
        applyTheme(getSavedTheme());
        bindThemeControls();
    }

    applyTheme(getSavedTheme());

    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", initializeTheme, {
            once: true,
        });
    } else {
        initializeTheme();
    }

    window.addEventListener("storage", function (event) {
        if (
            event.key === STORAGE_KEY ||
            event.key === LEGACY_STORAGE_KEY ||
            event.key === null
        ) {
            applyTheme(getSavedTheme());
        }
    });
})();
