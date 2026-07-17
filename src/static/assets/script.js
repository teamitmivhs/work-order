// Identity data is safe for display; the JWT itself stays in an HttpOnly cookie.
function getCurrentUserClaims() {
  try {
    return JSON.parse(localStorage.getItem("currentUser") || "{}");
  } catch {
    return {};
  }
}

function getSavedTheme() {
  return (
    localStorage.getItem("woTheme") || localStorage.getItem("theme") || "light"
  );
}

function updateThemeControls(theme) {
  document.querySelectorAll("[data-theme-toggle]").forEach((btn) => {
    btn.setAttribute(
      "aria-label",
      theme === "dark" ? "Switch to light theme" : "Switch to dark theme",
    );
    btn.setAttribute(
      "title",
      theme === "dark" ? "Switch to light theme" : "Switch to dark theme",
    );
  });
}

function applySavedTheme() {
  const theme = getSavedTheme();
  document.documentElement.dataset.theme = theme;
  document.body?.setAttribute("data-theme", theme);
  updateThemeControls(theme);
}

function toggleTheme() {
  const next = getSavedTheme() === "dark" ? "light" : "dark";
  localStorage.setItem("woTheme", next);
  localStorage.setItem("theme", next);
  applySavedTheme();
}

applySavedTheme();
document.addEventListener("DOMContentLoaded", () => {
  applySavedTheme();
  document.querySelectorAll("[data-theme-toggle]").forEach((btn) => {
    btn.addEventListener("click", toggleTheme);
    const icon = btn.querySelector("[data-theme-icon]");
    if (icon)
      icon.textContent =
        (localStorage.getItem("woTheme") || "light") === "dark" ? "LT" : "DK";
  });
});

// Header default untuk request yang butuh autentikasi
function authHeaders(extra = {}) {
  return {
    "Content-Type": "application/json",
    ...extra,
  };
}

function isCurrentUserAdmin() {
  return ["Admin", "Guru"].includes(getCurrentUserClaims().role);
}

function getCurrentUserId() {
  return Number(getCurrentUserClaims().id);
}

function isCurrentUserAssigned(order) {
  const userID = getCurrentUserId();
  return (
    Array.isArray(order.executors) &&
    order.executors.map(Number).includes(userID)
  );
}

window.isCurrentUserAdmin = isCurrentUserAdmin;
window.isCurrentUserAssigned = isCurrentUserAssigned;

// ===== UTILITY: Unwrap API Response =====
// Go backend (error.go yang sudah diperbaiki) sekarang return format:
//   success → { code, message, data: <payload> }
//   error   → { code, message, details? }
// Fungsi ini mengekstrak .data jika ada, fallback ke object itu sendiri
// agar kode yang mengakses field seperti .id tetap bekerja.
function unwrapData(json) {
  if (json && json.data !== undefined) return json.data;
  return json;
}

// ===== REFRESH PAGE =====
const refreshBtn = document.getElementById("refreshBtn");
const refreshIcon = document.getElementById("refreshIcon");

if (refreshBtn) {
  refreshBtn.addEventListener("click", function () {
    refreshIcon.style.transform = "rotate(0deg)";
    refreshIcon.style.transition = "transform 0.6s linear";
    setTimeout(() => {
      refreshIcon.style.transform = "rotate(360deg)";
    }, 10);
    setTimeout(() => {
      window.location.reload();
    }, 300);
  });
}

// ===== PROFILE DROPDOWN =====
// FIX: Semua akses elemen navbar dibungkus null-check.
// Script ini di-load di login.html dan register.html juga,
// di mana elemen-elemen navbar tidak ada — tanpa null-check
// baris `btn.addEventListener(...)` akan crash dan menghentikan
// seluruh eksekusi script.
const btn = document.getElementById("profileDropdownBtn");
const menu = document.getElementById("profileDropdown");
const profileDrawerBackdrop = document.getElementById("profileDrawerBackdrop");

if (btn && menu) {
  function setProfileMenu(open) {
    menu.classList.toggle("hidden", !open);
    profileDrawerBackdrop?.classList.toggle("hidden", !open);
    btn.setAttribute("aria-expanded", open);
  }

  // FIX: hanya JS click yang mengontrol dropdown (group-hover CSS dihapus dari index.html)
  btn.addEventListener("click", (e) => {
    e.stopPropagation(); // cegah document click langsung menutup lagi
    setProfileMenu(menu.classList.contains("hidden"));
  });

  // Logout button — cari berdasarkan ID, bukan index yang rapuh
  const logoutBtn = document.getElementById("logoutBtn");
  if (logoutBtn) {
    logoutBtn.addEventListener("click", async function () {
      await fetch("/api/logout", { method: "POST" }).catch(() => {});
      localStorage.removeItem("isGuestUser");
      localStorage.removeItem("guestLoginTime");
      localStorage.removeItem("isAdmin");
      localStorage.removeItem("userToken");
      localStorage.removeItem("currentUser");
      window.location.href = "/login";
    });
  }

  // Tutup dropdown saat klik di luar
  document.addEventListener("click", (e) => {
    if (!btn.contains(e.target) && !menu.contains(e.target)) {
      setProfileMenu(false);
    }
  });

  profileDrawerBackdrop?.addEventListener("click", (e) => {
    e.stopPropagation();
    setProfileMenu(false);
  });
}

// ===== MOBILE MENU =====
const mobileMenuBtn = document.getElementById("mobileMenuBtn");
const navMenu = document.getElementById("navMenu");

if (mobileMenuBtn && navMenu) {
  function setMobileMenu(open) {
    navMenu.classList.toggle("mobile-menu-active", open);
    mobileMenuBtn.setAttribute("aria-expanded", String(open));
  }

  mobileMenuBtn.addEventListener("click", (e) => {
    e.stopPropagation();
    setMobileMenu(!navMenu.classList.contains("mobile-menu-active"));
  });

  // Tutup menu saat klik di luar area navbar
  document.addEventListener("click", (e) => {
    if (!mobileMenuBtn.contains(e.target) && !navMenu.contains(e.target)) {
      setMobileMenu(false);
    }
  });

  // Tutup menu saat salah satu link/button di dalamnya diklik
  navMenu.querySelectorAll("a, button").forEach((el) => {
    el.addEventListener("click", () => {
      setMobileMenu(false);
    });
  });

  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") setMobileMenu(false);
  });

  // Reset saat resize ke desktop agar state tidak stuck
  window.addEventListener("resize", () => {
    if (window.innerWidth >= 768) {
      setMobileMenu(false);
    }
  });
}

document.querySelectorAll("#navMenu a[href]").forEach((link) => {
  const currentPath = window.location.pathname.replace(/\/$/, "") || "/";
  const linkPath = new URL(link.href, window.location.origin).pathname.replace(
    /\/$/,
    "",
  );
  if (linkPath === currentPath) {
    link.classList.add("active");
    link.setAttribute("aria-current", "page");
  }
});

// ===== CUSTOM POPUP SYSTEM =====
function showPopup(title, message, type = "info") {
  const existingPopup = document.getElementById("customPopup");
  if (existingPopup) existingPopup.remove();

  const popup = document.createElement("div");
  popup.id = "customPopup";
  popup.className =
    "fixed inset-0 bg-black bg-opacity-50 z-[100] flex items-center justify-center";

  const popupContent = document.createElement("div");
  popupContent.className =
    "bg-white rounded-2xl shadow-2xl p-6 w-11/12 max-w-md transform transition-all popup-fade-in";

  const configs = {
    success: {
      icon: `<svg class="w-12 h-12 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`,
      bg: "from-green-50 to-green-100",
    },
    warning: {
      icon: `<svg class="w-12 h-12 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path></svg>`,
      bg: "from-yellow-50 to-yellow-100",
    },
    error: {
      icon: `<svg class="w-12 h-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`,
      bg: "from-red-50 to-red-100",
    },
    info: {
      icon: `<svg class="w-12 h-12 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`,
      bg: "from-blue-50 to-blue-100",
    },
  };
  const { icon, bg } = configs[type] || configs.info;

  popupContent.innerHTML = `
    <div class="text-center">
      <div class="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-gradient-to-br ${bg} mb-4">${icon}</div>
      <h3 class="text-xl font-bold text-gray-900 mb-2">${title}</h3>
      <p class="text-gray-600 mb-6 leading-relaxed">${message}</p>
      <div class="popup-actions">
        <button class="px-6 py-3 bg-blue-500 text-white font-semibold rounded-lg hover:bg-blue-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-blue-200">OK</button>
      </div>
    </div>`;

  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  const closePopup = () => {
    popupContent.classList.replace("popup-fade-in", "popup-fade-out");
    setTimeout(() => popup.remove(), 300);
  };

  popup.querySelector("button").addEventListener("click", closePopup);

  if (type !== "error") {
    setTimeout(() => {
      if (popup.parentNode && popup.dataset.persistent !== "true") closePopup();
    }, 5000);
  }
}

function showActionPopup(title, message, actionText, onAction, type = "info") {
  showPopup(title, message, type);
  const popup = document.getElementById("customPopup");
  const buttonRow = popup?.querySelector(".popup-actions");
  const okBtn = buttonRow?.querySelector("button");
  if (!popup || !buttonRow || !okBtn) return;
  popup.dataset.persistent = "true";

  okBtn.textContent = "Nanti";
  okBtn.className =
    "px-5 py-3 border border-gray-300 text-gray-700 font-semibold rounded-lg hover:bg-gray-100";

  const actionBtn = document.createElement("button");
  actionBtn.type = "button";
  actionBtn.textContent = actionText;
  actionBtn.className =
    "px-5 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-200";

  buttonRow.className = "flex flex-col sm:flex-row justify-center gap-2";
  buttonRow.appendChild(actionBtn);
  actionBtn.addEventListener("click", async () => {
    actionBtn.disabled = true;
    actionBtn.textContent = "Mengaktifkan...";
    try {
      await onAction();
      popup.remove();
    } catch (err) {
      console.warn("Action popup failed:", err);
      actionBtn.disabled = false;
      actionBtn.textContent = actionText;
    }
  });
}

// ===== CONFIRMATION POPUP =====
function showConfirmationPopup(title, message, onConfirm) {
  const existingPopup = document.getElementById("customConfirmationPopup");
  if (existingPopup) existingPopup.remove();

  const popup = document.createElement("div");
  popup.id = "customConfirmationPopup";
  popup.className =
    "fixed inset-0 bg-black bg-opacity-50 z-[100] flex items-center justify-center";

  const popupContent = document.createElement("div");
  popupContent.className =
    "bg-white rounded-2xl shadow-2xl p-6 w-11/12 max-w-md transform transition-all popup-fade-in";
  popupContent.innerHTML = `
    <div class="text-center">
      <div class="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-gradient-to-br from-yellow-50 to-yellow-100 mb-4">
        <svg class="w-12 h-12 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
        </svg>
      </div>
      <h3 class="text-xl font-bold text-gray-900 mb-2">${title}</h3>
      <p class="text-gray-600 mb-6 leading-relaxed">${message}</p>
      <div class="flex justify-center gap-4">
        <button id="confirmBtn" class="px-6 py-3 bg-green-500 text-white font-semibold rounded-lg hover:bg-green-600 transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-green-200">Ya</button>
        <button id="cancelBtn"  class="px-6 py-3 bg-red-500  text-white font-semibold rounded-lg hover:bg-red-600  transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-red-200">Tidak</button>
      </div>
    </div>`;

  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  const closeThis = () => {
    popupContent.classList.replace("popup-fade-in", "popup-fade-out");
    setTimeout(() => popup.remove(), 300);
  };

  document.getElementById("confirmBtn").addEventListener("click", () => {
    onConfirm();
    closeThis();
  });
  document.getElementById("cancelBtn").addEventListener("click", closeThis);
}

function showTextInputPopup(title, message, onSubmit) {
  const existingPopup = document.getElementById("customTextInputPopup");
  if (existingPopup) existingPopup.remove();

  const popup = document.createElement("div");
  popup.id = "customTextInputPopup";
  popup.className =
    "fixed inset-0 bg-black bg-opacity-50 z-[101] flex items-center justify-center";

  const popupContent = document.createElement("div");
  popupContent.className =
    "popup-inner bg-white rounded-2xl shadow-2xl p-6 w-11/12 max-w-md transform transition-all popup-fade-in";
  popupContent.innerHTML = `
    <div>
      <h3 class="popup-title text-xl font-bold text-gray-900 mb-2"></h3>
      <p class="popup-message text-gray-600 mb-4 leading-relaxed"></p>
      <textarea id="textInputPopupValue" class="w-full min-h-[120px] rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm text-gray-900 outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-100" maxlength="1000" placeholder="Tulis catatan untuk requester..."></textarea>
      <p class="text-xs text-gray-500 mt-2">Catatan wajib diisi dan akan terlihat oleh requester saat tracking.</p>
      <div class="flex justify-end gap-3 mt-5">
        <button id="textInputCancelBtn" type="button" class="px-5 py-2.5 border border-gray-300 text-gray-700 font-semibold rounded-lg hover:bg-gray-100 focus:outline-none focus:ring-4 focus:ring-gray-200">Batal</button>
        <button id="textInputSubmitBtn" type="button" class="px-5 py-2.5 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-200">Kirim</button>
      </div>
    </div>`;

  popupContent.querySelector(".popup-title").textContent = title;
  popupContent.querySelector(".popup-message").textContent = message;
  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  const input = document.getElementById("textInputPopupValue");
  const submitBtn = document.getElementById("textInputSubmitBtn");
  const cancelBtn = document.getElementById("textInputCancelBtn");

  const closeThis = () => {
    popupContent.classList.replace("popup-fade-in", "popup-fade-out");
    setTimeout(() => popup.remove(), 300);
  };

  submitBtn.addEventListener("click", () => {
    const value = input.value.trim();
    if (!value) {
      input.focus();
      input.classList.add("border-red-400", "focus:border-red-500", "focus:ring-red-100");
      return;
    }
    onSubmit(value);
    closeThis();
  });
  cancelBtn.addEventListener("click", closeThis);
  input.focus();
}

// ===== POPUP ANIMATION HELPERS =====
function showAnimatedPopup(popupElement) {
  const content = popupElement.firstElementChild;
  popupElement.classList.remove("hidden");
  content.classList.remove("popup-fade-out");
  content.classList.add("popup-fade-in");
}

function hideAnimatedPopup(popupElement) {
  const content = popupElement.firstElementChild;
  content.classList.remove("popup-fade-in");
  content.classList.add("popup-fade-out");
  setTimeout(() => popupElement.classList.add("hidden"), 300);
}

// ===== QUICK SUMMARY TITLE =====
function updateQuickSummaryTitle() {
  const el = document.getElementById("quickSummaryTitle");
  if (el) {
    el.textContent =
      new Date().getDate() === 1 ? "Quick Summary Bulan Ini" : "Quick Summary";
  }
}

// ===== LOADING STATE =====
function setButtonLoading(btn, isLoading) {
  if (!btn) return;
  btn.disabled = isLoading;
  btn.style.opacity = isLoading ? "0.6" : "1";
  btn.style.cursor = isLoading ? "not-allowed" : "";
}

// ===== MAIN DASHBOARD LOGIC =====
document.addEventListener("DOMContentLoaded", async function () {
  // Script ini di-load di login.html & register.html juga.
  // Elemen-elemen dashboard tidak akan ada di sana.
  // FIX: cek keberadaan elemen kunci sebelum menjalankan logika dashboard.
  const workOrdersTableBody = document.getElementById("workOrdersTableBody");
  if (!workOrdersTableBody) return; // Bukan halaman dashboard, stop di sini

  updateQuickSummaryTitle();

  let members = [];
  let workOrders = [];
  let currentOrder = null;
  let additionalOperators = [];
  let pendingDocumentationOrderId = null;
  let pendingCompleteOrderId = null;
  let documentationInput = null;
  let knownOrderIds = new Set();
  let notificationReady = false;
  let orderNotificationRegistration = null;
  let notificationPermissionRequestInFlight = false;

  const isGuestUser = localStorage.getItem("isGuestUser") === "true";
  if (!localStorage.getItem("currentUser")) {
    window.location.replace("/login");
    return;
  }
  if (isGuestUser) {
    window.location.replace("/guest");
    return;
  }

  function checkGuestRestriction(action = "action") {
    if (isGuestUser) {
      showPopup(
        "Access Denied",
        `Guests can only create work orders. ${action} is not allowed.`,
        "warning",
      );
      return true;
    }
    return false;
  }

  function isAvailableWorker(member) {
    return (
      member.status === "standby" ||
      (member.role === "Guru" && member.status !== "onjob")
    );
  }

  function isVisibleStatusMember(member) {
    return member.role !== "Guru" || member.status === "onjob";
  }

  // ===== DOM REFERENCES =====
  const memberStatusPopup = document.getElementById("memberStatusPopup");
  const memberList = document.getElementById("memberList");
  const closePopupBtn = document.getElementById("closePopup");
  const statusContainers = document.querySelectorAll(".status-container");
  const takeOrderPopup = document.getElementById("takeOrderPopup");
  const closeTakeOrderPopupBtn = document.getElementById("closeTakeOrderPopup");
  const cancelTakeOrderBtn = document.getElementById("cancelTakeOrderBtn");
  const confirmTakeOrderBtn = document.getElementById("confirmTakeOrderBtn");
  const openSelectHelperOperatorModalBtn = document.getElementById(
    "openSelectHelperOperatorModalBtn",
  );
  const createOrderPopup = document.getElementById("createOrderPopup");
  const closeCreateOrderPopupBtn = document.getElementById(
    "closeCreateOrderPopup",
  );
  const cancelCreateOrderBtn = document.getElementById("cancelCreateOrderBtn");
  const exitGuestBtn = document.getElementById("exitGuestBtn");
  const createOrderForm = document.getElementById("createOrderForm");
  const createOrderBtn = document.getElementById("createOrderBtn");
  const orderLocationSelect = document.getElementById("orderLocation");
  const specificLocationContainer = document.getElementById(
    "specificLocationContainer",
  );
  const specificLocationInput = document.getElementById("specificLocation");
  const memberSearchInput = document.getElementById("memberSearchInput");
  const searchDropdown = document.getElementById("searchDropdown");
  const searchResults = document.getElementById("searchResults");
  const statusFilterTabs = document.querySelectorAll(".status-filter-tab");
  const workOrderStatusTabs = document.querySelectorAll(
    ".work-order-status-tab",
  );
  const workOrderSearchInput = document.getElementById("workOrderSearchInput");
  const workOrderDateFilter = document.getElementById("workOrderDateFilter");
  const clearWorkOrderFiltersBtn = document.getElementById(
    "clearWorkOrderFilters",
  );
  const heroPrimaryActionBtns = document.querySelectorAll(
    "[data-hero-primary-action]",
  );
  const selectHelperOperatorModal = document.getElementById(
    "selectHelperOperatorModal",
  );
  const closeSelectHelperOperatorModalBtn = document.getElementById(
    "closeSelectHelperOperatorModalBtn",
  );
  const availableStandbyOperatorsList = document.getElementById(
    "availableStandbyOperatorsList",
  );
  const takeOrderOperatorsTitle = document.getElementById(
    "takeOrderOperatorsTitle",
  );
  const takeOrderHelperBtnText = document.getElementById(
    "takeOrderHelperBtnText",
  );
  const takeOrderSafetySection = document.querySelector(".take-order-safety");

  let currentStatusFilter = "all";
  let activeWorkOrderStatus = "pending";
  let workOrderSearchQuery = "";
  let workOrderDateQuery = "";

  function scrollToPendingWorkOrders() {
    setActiveWorkOrderStatus("pending");
    document
      .querySelector(".workorders-panel")
      ?.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  function setupRoleBasedOrderActions() {
    const admin = isCurrentUserAdmin();
    if (!admin && createOrderBtn) createOrderBtn.classList.add("hidden");

    const createMarkup = `
      <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
        <path d="M12 4.5a.75.75 0 01.75.75v6.75h6.75a.75.75 0 010 1.5h-6.75v6.75a.75.75 0 01-1.5 0v-6.75H4.5a.75.75 0 010-1.5h6.75V5.25a.75.75 0 01.75-.75z"/>
      </svg>
      Create New Order`;
    const takeMarkup = `
      <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v8m-4-4h8m4 0a8 8 0 11-16 0 8 8 0 0116 0z"/>
      </svg>
      Take Work Order`;

    heroPrimaryActionBtns.forEach((btn) => {
      btn.innerHTML = admin ? createMarkup : takeMarkup;
      btn.setAttribute(
        "aria-label",
        admin ? "Create new work order" : "Take pending work order",
      );
      btn.addEventListener("click", () => {
        if (admin) showAnimatedPopup(createOrderPopup);
        else scrollToPendingWorkOrders();
      });
    });
  }

  setupRoleBasedOrderActions();

  function ensureDocumentationInput() {
    if (documentationInput) return documentationInput;
    documentationInput = document.createElement("input");
    documentationInput.type = "file";
    documentationInput.setAttribute("accept", "image/*");
    documentationInput.setAttribute("capture", "environment");
    documentationInput.style.position = "fixed";
    documentationInput.style.left = "-9999px";
    documentationInput.style.top = "0";
    documentationInput.style.width = "1px";
    documentationInput.style.height = "1px";
    documentationInput.style.opacity = "0";
    documentationInput.addEventListener("change", async function () {
      const file = this.files && this.files[0];
      const orderId = pendingDocumentationOrderId;
      this.value = "";
      if (!file || !orderId) return;

      try {
        await uploadDocumentationPhoto(orderId, file);
        await refreshAllDataFromAPI();
        const shouldComplete = pendingCompleteOrderId === orderId;
        pendingDocumentationOrderId = null;
        pendingCompleteOrderId = null;
        if (shouldComplete) {
          showPopup(
            "Foto Tersimpan",
            "Foto bukti berhasil disimpan. Menyelesaikan work order...",
            "success",
          );
          markOrderDone(orderId, true);
        } else {
          showPopup(
            "Foto Tersimpan",
            "Foto bukti berhasil ditambahkan ke work order.",
            "success",
          );
        }
      } catch (err) {
        pendingDocumentationOrderId = null;
        pendingCompleteOrderId = null;
        console.error("Error upload documentation photo:", err);
        showPopup(
          "Upload Gagal",
          err.message || "Gagal menyimpan foto dokumentasi.",
          "error",
        );
      }
    });
    document.body.appendChild(documentationInput);
    return documentationInput;
  }

  function loadImageFile(file) {
    return new Promise((resolve, reject) => {
      const url = URL.createObjectURL(file);
      const img = new Image();
      img.onload = () => {
        URL.revokeObjectURL(url);
        resolve(img);
      };
      img.onerror = () => {
        URL.revokeObjectURL(url);
        reject(new Error("Gagal membaca foto."));
      };
      img.src = url;
    });
  }

  async function compressImageForUpload(file, options = {}) {
    const maxDimension = options.maxDimension || 1600;
    const quality = options.quality || 0.82;
    const imageTypes = ["image/jpeg", "image/jpg", "image/png", "image/webp"];
    if (!imageTypes.includes(String(file.type).toLowerCase())) return file;
    if (file.size <= 2 * 1024 * 1024) return file;

    try {
      const img = await loadImageFile(file);
      const longestSide = Math.max(
        img.naturalWidth || img.width,
        img.naturalHeight || img.height,
      );
      const scale = longestSide > maxDimension ? maxDimension / longestSide : 1;
      const width = Math.max(
        1,
        Math.round((img.naturalWidth || img.width) * scale),
      );
      const height = Math.max(
        1,
        Math.round((img.naturalHeight || img.height) * scale),
      );
      const canvas = document.createElement("canvas");
      canvas.width = width;
      canvas.height = height;
      const ctx = canvas.getContext("2d");
      ctx.drawImage(img, 0, 0, width, height);
      const blob = await new Promise((resolve) =>
        canvas.toBlob(resolve, "image/jpeg", quality),
      );
      if (!blob) return file;
      const baseName = file.name.replace(/\.[^.]+$/, "") || "documentation";
      return new File([blob], `${baseName}.jpg`, {
        type: "image/jpeg",
        lastModified: Date.now(),
      });
    } catch (err) {
      console.warn("Image compression failed, uploading original file:", err);
      return file;
    }
  }

  async function uploadDocumentationPhoto(orderId, file) {
    const photoFile = await compressImageForUpload(file);
    const formData = new FormData();
    formData.append("photo", photoFile, photoFile.name || "documentation.jpg");
    const r = await fetch(`/api/workorders/${orderId}/documentation`, {
      method: "POST",
      body: formData,
    });
    if (!r.ok) {
      const err = await r.json().catch(() => ({}));
      if (r.status === 413) {
        throw new Error(
          "Foto masih terlalu besar untuk server. Coba pilih foto lain atau reload nginx/backend setelah update limit upload.",
        );
      }
      throw new Error(
        err.message || "Gagal upload foto dokumentasi. Status: " + r.status,
      );
    }
    const json = await r.json();
    return unwrapData(json);
  }

  function requestDocumentationPhoto(orderId, autoComplete = false) {
    const order = workOrders.find((o) => o.id === orderId);
    if (!order) return;
    if (!isCurrentUserAdmin() && !isCurrentUserAssigned(order)) {
      showPopup(
        "Access Denied",
        "Hanya operator yang ditugaskan atau admin yang bisa menambahkan foto bukti.",
        "warning",
      );
      return;
    }
    pendingDocumentationOrderId = orderId;
    pendingCompleteOrderId = autoComplete ? orderId : null;
    ensureDocumentationInput().click();
  }
  window.requestDocumentationPhoto = requestDocumentationPhoto;

  function elapsedFromServerStartedAt(order) {
    if (order && typeof order.progressSeconds === "number") {
      const fetchedAt = order._fetchedAt || Date.now();
      return Math.max(
        0,
        order.progressSeconds + Math.floor((Date.now() - fetchedAt) / 1000),
      );
    }
    if (!order || !order.startedAt || order.status !== "progress") return null;
    const parts = String(order.startedAt).split(":").map(Number);
    if (parts.length < 2 || parts.some(Number.isNaN)) return null;
    const started = new Date();
    started.setHours(parts[0] || 0, parts[1] || 0, parts[2] || 0, 0);
    let elapsed = Math.floor((Date.now() - started.getTime()) / 1000);
    if (elapsed < 0) elapsed += 24 * 60 * 60;
    return Math.max(0, elapsed);
  }

  function isTouchMobileLike() {
    return (
      window.matchMedia &&
      window.matchMedia("(hover: none) and (pointer: coarse)").matches
    );
  }

  async function registerOrderNotificationWorker() {
    if (!("serviceWorker" in navigator)) return null;
    if (orderNotificationRegistration) return orderNotificationRegistration;
    try {
      orderNotificationRegistration = await navigator.serviceWorker.register(
        "/static/assets/workorder-sw.js",
        {
          scope: "/static/assets/",
        },
      );
      return orderNotificationRegistration;
    } catch (err) {
      console.warn("Service worker notification registration failed:", err);
      return null;
    }
  }

  function urlBase64ToUint8Array(base64String) {
    const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
    const base64 = (base64String + padding)
      .replace(/-/g, "+")
      .replace(/_/g, "/");
    const rawData = atob(base64);
    const outputArray = new Uint8Array(rawData.length);
    for (let i = 0; i < rawData.length; i += 1) {
      outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
  }

  async function syncPushSubscription() {
    if (!("Notification" in window) || Notification.permission !== "granted")
      return;
    if (!window.isSecureContext) return;
    const registration = await registerOrderNotificationWorker();
    if (!registration || !("pushManager" in registration)) return;

    const keyResponse = await fetch("/api/notifications/vapid-public-key", {
      headers: authHeaders(),
    });
    if (!keyResponse.ok) return;
    const keyJson = await keyResponse.json();
    const publicKey = unwrapData(keyJson).publicKey;
    if (!publicKey) return;

    const savedKey = localStorage.getItem("workOrderPushVapidPublicKey");
    let subscription = await registration.pushManager.getSubscription();
    if (subscription && savedKey && savedKey !== publicKey) {
      await subscription.unsubscribe();
      subscription = null;
    }
    if (!subscription) {
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(publicKey),
      });
    }
    localStorage.setItem("workOrderPushVapidPublicKey", publicKey);

    await fetch("/api/notifications/subscribe", {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify(subscription),
    });
  }

  async function requestOrderNotificationPermission() {
    if (!("Notification" in window)) return "unsupported";
    if (Notification.permission !== "default") return Notification.permission;
    if (notificationPermissionRequestInFlight) return;
    notificationPermissionRequestInFlight = true;
    try {
      const permission = await Notification.requestPermission();
      if (permission === "granted") await syncPushSubscription();
      return permission;
    } catch (err) {
      console.warn("Notification permission request failed:", err);
      return "error";
    } finally {
      notificationPermissionRequestInFlight = false;
    }
  }

  function showNotificationSetupPrompt() {
    if (!("Notification" in window)) {
      showPopup(
        "Notifikasi Tidak Didukung",
        "Browser ini tidak mendukung notifikasi work order.",
        "warning",
      );
      return;
    }
    if (!window.isSecureContext) {
      showPopup(
        "Butuh HTTPS",
        "Push notification hanya aktif di HTTPS atau localhost. Buka dashboard lewat HTTPS agar Brave bisa subscribe.",
        "warning",
      );
      return;
    }
    if (!("serviceWorker" in navigator)) {
      showPopup(
        "Push Tidak Didukung",
        "Service worker tidak tersedia di browser ini. Notifikasi hanya bisa diterima saat dashboard terbuka.",
        "warning",
      );
      return;
    }
    if (Notification.permission === "granted") {
      syncPushSubscription().catch((err) =>
        console.warn("Push subscription sync failed:", err),
      );
      return;
    }
    if (Notification.permission === "denied") {
      showPopup(
        "Notifikasi Diblokir",
        "Buka pengaturan site di browser, izinkan Notifications untuk alamat ini, lalu reload dashboard.",
        "warning",
      );
      return;
    }

    registerOrderNotificationWorker().then((registration) => {
      if (!registration || !("pushManager" in registration)) {
        showPopup(
          "Push Tidak Aktif di Brave",
          "Aktifkan push messaging di pengaturan Brave, lalu reload dashboard dan klik Aktifkan Notifikasi lagi.",
          "warning",
        );
        return;
      }

      showActionPopup(
        "Aktifkan Notifikasi Work Order",
        "Izinkan notifikasi supaya work order baru tetap muncul di ponsel meski tab Work Order sedang ditutup.",
        "Aktifkan Notifikasi",
        async () => {
          const permission = await requestOrderNotificationPermission();
          if (permission === "granted") {
            showPopup(
              "Notifikasi Aktif",
              "Browser ini sudah subscribe notifikasi work order.",
              "success",
            );
          } else {
            showPopup(
              "Notifikasi Belum Aktif",
              "Permission belum diberikan. Aktifkan dari pengaturan site browser untuk menerima push.",
              "warning",
            );
          }
        },
        "info",
      );
    });
  }

  async function showWorkOrderBrowserNotification(title, message, order) {
    if (!("Notification" in window) || Notification.permission !== "granted") {
      if ("vibrate" in navigator) navigator.vibrate([180, 80, 180]);
      return;
    }

    const options = {
      body: message,
      icon: "/static/public/itlogo.png",
      badge: "/static/public/itlogo.png",
      tag: `work-order-${order && order.id ? order.id : "incoming"}`,
      renotify: true,
      data: {
        url: "/",
        workOrderId: order && order.id ? order.id : null,
      },
    };

    try {
      const registration = await registerOrderNotificationWorker();
      if (registration && typeof registration.showNotification === "function") {
        await registration.showNotification(title, options);
      } else {
        new Notification(title, options);
      }
      if ("vibrate" in navigator) navigator.vibrate([120, 60, 120]);
    } catch (err) {
      console.warn("Browser notification failed:", err);
      try {
        new Notification(title, options);
      } catch (e) {
        /* ignore */
      }
    }
  }

  function primeOrderNotifications() {
    knownOrderIds = new Set(workOrders.map((order) => Number(order.id)));
    notificationReady = true;
    registerOrderNotificationWorker();

    if ("Notification" in window && Notification.permission === "default") {
      setTimeout(showNotificationSetupPrompt, 600);
    } else if ("Notification" in window && Notification.permission === "granted") {
      syncPushSubscription().catch((err) =>
        console.warn("Push subscription sync failed:", err),
      );
    } else if ("Notification" in window && Notification.permission === "denied") {
      setTimeout(showNotificationSetupPrompt, 600);
    }
  }

  function notifyIncomingOrders(latestOrders) {
    if (!notificationReady) return;
    const incoming = latestOrders.filter(
      (order) => !knownOrderIds.has(Number(order.id)),
    );
    if (incoming.length === 0) return;
    incoming.forEach((order) => knownOrderIds.add(Number(order.id)));
    const newest = incoming[0];
    const priorityLabel =
      {
        low: "Rendah",
        medium: "Sedang",
        high: "Tinggi",
        urgent: "Urgent",
      }[String(newest.priority || "").toLowerCase()] || "";
    const title = `${incoming.length === 1 ? "Work Order Baru" : `${incoming.length} Work Order Baru`}${priorityLabel ? ` • Prioritas ${priorityLabel}` : ""}`;
    const message = [
      `Perangkat: ${newest.device || "Belum diketahui"}`,
      `Lokasi: ${newest.location || "Belum diisi"}`,
      `Kendala: ${newest.problem || "Belum ada deskripsi"}`,
    ].join("\n");
    showPopup(title, message, "info");
    showWorkOrderBrowserNotification(title, message, newest);
  }

  async function pollIncomingOrders() {
    try {
      const r = await fetch("/api/workorders", {
        headers: authHeaders(),
        cache: "no-store",
      });
      if (!r.ok) return;
      const json = await r.json();
      const latestOrders = normalizeWorkOrdersFromApi(json);
      notifyIncomingOrders(latestOrders);
      workOrders = latestOrders;
      populateWorkOrdersTable();
      updateSummaryCounts();
    } catch (err) {
      console.error("Order notification poll failed:", err);
    }
  }

  let realtimeRefreshTimer;

  function scheduleRealtimeRefresh() {
    clearTimeout(realtimeRefreshTimer);
    realtimeRefreshTimer = setTimeout(async () => {
      await fetchMembers();
      initializeMemberImages();
      await pollIncomingOrders();
    }, 100);
  }

  function connectRealtimeUpdates() {
    if (!("EventSource" in window)) return;

    const stream = new EventSource("/api/events");
    stream.onmessage = scheduleRealtimeRefresh;
    window.addEventListener("pagehide", () => stream.close(), { once: true });
  }

  // ===== SAFETY CHECKLIST DATA =====
  const safetyChecklistItems = {
    CCTV: [
      {
        id: "cctv1",
        text: "Gunakan pelindung mata (goggles)",
        required: false,
      },
      { id: "cctv2", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "cctv3", text: "Pastikan area kerja aman", required: true },
      { id: "cctv4", text: "Gunakan Sarung Tangan", required: true },
    ],
    WiFi: [
      {
        id: "wifi1",
        text: "Gunakan pelindung mata (goggles)",
        required: false,
      },
      { id: "wifi2", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "wifi3", text: "Pastikan area kerja aman", required: true },
      { id: "wifi4", text: "Gunakan Sarung Tangan", required: true },
    ],
    "Gedung A": [
      { id: "ga1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "ga2", text: "Gunakan Sarung Tangan", required: false },
      { id: "ga3", text: "Pastikan area kerja aman", required: true },
      { id: "ga4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "ga5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung B": [
      { id: "gb1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "gb2", text: "Gunakan Sarung Tangan", required: false },
      { id: "gb3", text: "Pastikan area kerja aman", required: true },
      { id: "gb4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "gb5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung B Baru": [
      { id: "gbb1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "gbb2", text: "Gunakan Sarung Tangan", required: false },
      { id: "gbb3", text: "Pastikan area kerja aman", required: true },
      { id: "gbb4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "gbb5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung C": [
      { id: "gc1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "gc2", text: "Gunakan Sarung Tangan", required: false },
      { id: "gc3", text: "Pastikan area kerja aman", required: true },
      { id: "gc4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "gc5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung D": [
      { id: "gd1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "gd2", text: "Gunakan Sarung Tangan", required: false },
      { id: "gd3", text: "Pastikan area kerja aman", required: true },
      { id: "gd4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "gd5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung E": [
      { id: "ge1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "ge2", text: "Gunakan Sarung Tangan", required: false },
      { id: "ge3", text: "Pastikan area kerja aman", required: true },
      { id: "ge4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "ge5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung F": [
      { id: "gf1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "gf2", text: "Gunakan Sarung Tangan", required: false },
      { id: "gf3", text: "Pastikan area kerja aman", required: true },
      { id: "gf4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "gf5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung G": [
      { id: "gg1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "gg2", text: "Gunakan Sarung Tangan", required: false },
      { id: "gg3", text: "Pastikan area kerja aman", required: true },
      { id: "gg4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "gg5", text: "Gunakan sepatu safety", required: true },
    ],
    "Gedung TKI": [
      { id: "gt1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "gt2", text: "Gunakan Sarung Tangan", required: false },
      { id: "gt3", text: "Pastikan area kerja aman", required: true },
      { id: "gt4", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "gt5", text: "Gunakan sepatu safety", required: true },
    ],
    "Ruang Guru": [
      { id: "rg1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "rg2", text: "Gunakan Sarung Tangan", required: false },
      { id: "rg3", text: "Pastikan area kerja aman", required: false },
      { id: "rg4", text: "Matikan listrik sebelum bekerja", required: true },
    ],
    "Ruang Yayasan": [
      { id: "ry1", text: "Pastikan sirkulasi udara baik", required: false },
      { id: "ry2", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "ry3", text: "Matikan listrik sebelum bekerja", required: true },
      { id: "ry4", text: "Gunakan Sarung Tangan", required: false },
      { id: "ry5", text: "Pastikan area kerja aman", required: true },
    ],
    default: [
      { id: "def1", text: "Gunakan pelindung mata (goggles)", required: false },
      { id: "def2", text: "Gunakan Sarung Tangan", required: false },
      { id: "def3", text: "Pastikan area kerja aman", required: true },
      { id: "def4", text: "Matikan listrik sebelum bekerja", required: true },
    ],
  };

  // ===== API FUNCTIONS =====

  function normalizeWorkOrdersFromApi(json) {
    const data = Array.isArray(json)
      ? json
      : Array.isArray(json.data)
        ? json.data
        : [];
    const fetchedAt = Date.now();
    return data.map((order) => {
      if (order?.status !== "progress") {
        localStorage.removeItem(`order_timer_${order.id}`);
      }
      return { ...order, _fetchedAt: fetchedAt };
    });
  }

  async function fetchMembers() {
    try {
      const r = await fetch("/api/members", {
        headers: authHeaders(),
        cache: "no-store",
      });
      if (!r.ok) throw new Error(r.statusText);
      const json = await r.json();
      // FIX: handle format response baru { code, message, data: [...] }
      members = Array.isArray(json)
        ? json
        : Array.isArray(json.data)
          ? json.data
          : [];
    } catch (err) {
      console.error("Error fetching members:", err);
      if (memberList)
        memberList.innerHTML =
          '<div class="text-center py-4 text-red-500">Failed to load member data.</div>';
      members = [];
    }
  }

  async function fetchAndRenderWorkOrders() {
    try {
      const r = await fetch("/api/workorders", {
        headers: authHeaders(),
        cache: "no-store",
      });
      if (!r.ok) throw new Error("Gagal mengambil data work order dari server");
      const json = await r.json();
      // FIX: handle format response baru { code, message, data: [...] }
      workOrders = normalizeWorkOrdersFromApi(json);
    } catch (err) {
      console.error("Error fetching work orders:", err);
      showPopup("Error", "Gagal memuat data work order dari server.", "error");
      workOrders = [];
    }
    populateWorkOrdersTable();
    updateSummaryCounts();
  }

  async function refreshAllDataFromAPI() {
    await fetchMembers();
    initializeMemberImages();
    await fetchAndRenderWorkOrders();
  }

  // FIX: Kirim Authorization header ke semua protected endpoint.
  // /api/status hanya untuk status user sendiri; update member lain harus lewat
  // /api/members/:id/status agar ID target tidak diabaikan oleh backend.
  async function apiUpdateMemberStatus(memberId, newStatus) {
    const currentUserId = Number(getCurrentUserClaims().id);
    const targetMemberId = Number(memberId);
    const isSelfUpdate = targetMemberId === currentUserId;
    const endpoint = isSelfUpdate
      ? "/api/status"
      : `/api/members/${targetMemberId}/status`;
    const method = isSelfUpdate ? "POST" : "PATCH";
    const r = await fetch(endpoint, {
      method,
      headers: authHeaders(),
      body: JSON.stringify({ status: newStatus }),
    });
    if (!r.ok) {
      const err = await r.json().catch(() => ({}));
      throw new Error(
        err.message || "Gagal memperbarui status member. Status: " + r.status,
      );
    }
    const json = await r.json();
    return unwrapData(json);
  }

  // ===== INIT =====
  if (isGuestUser) {
    // Guest: tidak load work orders dari DB (data milik operator lain)
    // Hanya load members untuk status bar, lalu buka popup create order
    await fetchMembers();
    initializeMemberImages();
    updateSummaryCounts();
    setTimeout(() => {
      if (createOrderPopup) showAnimatedPopup(createOrderPopup);
      if (exitGuestBtn) exitGuestBtn.classList.remove("hidden");
    }, 300);
  } else {
    await fetchMembers();
    await fetchAndRenderWorkOrders();
    initializeMemberImages();
    updateSummaryCounts();
    primeOrderNotifications();
    connectRealtimeUpdates();
    setInterval(scheduleRealtimeRefresh, 60000);
    document.addEventListener("visibilitychange", () => {
      if (document.visibilityState === "visible") scheduleRealtimeRefresh();
    });
  }

  workOrderStatusTabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      setActiveWorkOrderStatus(tab.dataset.workOrderStatus || "pending");
    });
  });
  if (workOrderSearchInput) {
    workOrderSearchInput.addEventListener("input", () => {
      workOrderSearchQuery = workOrderSearchInput.value;
      populateWorkOrdersTable();
    });
  }
  if (workOrderDateFilter) {
    workOrderDateFilter.addEventListener("change", () => {
      workOrderDateQuery = workOrderDateFilter.value;
      populateWorkOrdersTable();
    });
  }
  if (clearWorkOrderFiltersBtn) {
    clearWorkOrderFiltersBtn.addEventListener("click", () => {
      workOrderSearchQuery = "";
      workOrderDateQuery = "";
      if (workOrderSearchInput) workOrderSearchInput.value = "";
      if (workOrderDateFilter) workOrderDateFilter.value = "";
      populateWorkOrdersTable();
    });
  }

  // ===== STATUS CONTAINER CLICK =====
  statusContainers.forEach((container) => {
    container.addEventListener("click", function (e) {
      if (
        !e.target.closest(".member-images") &&
        !e.target.closest(".more-members")
      ) {
        if (checkGuestRestriction("Viewing member status")) return;
        const status = this.dataset.status;
        currentStatusFilter = status;
        updateFilterTabs(status);
        openMemberStatusPopup(status);
      }
    });
  });

  // ===== MEMBER STATUS POPUP =====
  closePopupBtn.addEventListener("click", () =>
    hideAnimatedPopup(memberStatusPopup),
  );

  statusFilterTabs.forEach((tab) => {
    tab.addEventListener("click", function () {
      const status = this.dataset.statusFilter;
      currentStatusFilter = status;
      updateFilterTabs(status);
      populateMemberList(status);
    });
  });

  // ===== TAKE ORDER POPUP =====
  closeTakeOrderPopupBtn.addEventListener("click", () => {
    hideAnimatedPopup(takeOrderPopup);
    resetTakeOrderForm();
  });
  cancelTakeOrderBtn.addEventListener("click", () => {
    hideAnimatedPopup(takeOrderPopup);
    resetTakeOrderForm();
  });
  confirmTakeOrderBtn.addEventListener("click", () => confirmTakeOrder());
  openSelectHelperOperatorModalBtn.addEventListener(
    "click",
    openSelectHelperOperatorModal,
  );

  // ===== HELPER OPERATOR MODAL =====
  closeSelectHelperOperatorModalBtn.addEventListener("click", () =>
    hideAnimatedPopup(selectHelperOperatorModal),
  );

  function openSelectHelperOperatorModal() {
    if (!isCurrentUserAdmin()) return;
    showAnimatedPopup(selectHelperOperatorModal);
    populateAvailableStandbyOperators();
  }

  function populateAvailableStandbyOperators() {
    availableStandbyOperatorsList.innerHTML = "";
    const standbyMembers = members.filter(isAvailableWorker);

    if (standbyMembers.length === 0) {
      availableStandbyOperatorsList.innerHTML =
        '<p class="text-gray-500 text-center py-4">Tidak ada pelaksana tersedia</p>';
      return;
    }

    standbyMembers.forEach((member) => {
      const div = document.createElement("div");
      div.className =
        "helper-operator-item flex items-center justify-between gap-3 p-2 bg-gray-50 rounded-lg";
      div.innerHTML = `
        <div class="helper-operator-person flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full object-cover">
          <span class="helper-operator-name font-medium">${member.name}</span>
        </div>
        <button class="add-helper-operator-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center"
          data-member-id="${member.id}" title="Tambahkan sebagai operator bantuan">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        </button>`;
      availableStandbyOperatorsList.appendChild(div);
    });

    document.querySelectorAll(".add-helper-operator-btn").forEach((b) => {
      b.addEventListener("click", function () {
        addHelperOperator(parseInt(this.dataset.memberId));
      });
    });
  }

  function addHelperOperator(memberId) {
    if (additionalOperators.includes(memberId)) {
      showPopup("Peringatan", "Operator ini sudah ditambahkan!", "warning");
      return;
    }
    additionalOperators.push(memberId);
    const member = members.find((m) => m.id === memberId);
    showPopup(
      "Operator Ditambahkan",
      `${member.name} ditambahkan sebagai operator bantuan.`,
      "success",
    );
    populateStandbyOperatorsInTakeOrderPopup();
    hideAnimatedPopup(selectHelperOperatorModal);
  }

  function populateStandbyOperatorsInTakeOrderPopup() {
    const listDiv = document.getElementById("standbyOperatorsList");
    listDiv.innerHTML = "";

    if (additionalOperators.length === 0) {
      listDiv.innerHTML =
        '<p class="text-gray-500 text-center py-4">Tidak ada operator bantuan yang dipilih.</p>';
      return;
    }

    additionalOperators.forEach((memberId) => {
      const member = members.find((m) => m.id === memberId);
      if (!member) return;
      const div = document.createElement("div");
      div.className =
        "selected-helper-item flex items-center gap-3 p-2 rounded-lg shadow-sm";
      div.innerHTML = `
        <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full object-cover">
        <div class="selected-helper-copy flex-1 min-w-0">
          <div class="selected-helper-name font-medium">${member.name}</div>
          <div class="selected-helper-role text-xs">Operator Bantuan</div>
        </div>
        <button class="remove-helper-operator-btn text-red-500 hover:text-red-700 transition-colors"
          data-member-id="${member.id}" title="Hapus operator bantuan">
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>`;
      listDiv.appendChild(div);
    });

    document.querySelectorAll(".remove-helper-operator-btn").forEach((b) => {
      b.addEventListener("click", function () {
        removeHelperOperator(parseInt(this.dataset.memberId));
      });
    });
  }

  function removeHelperOperator(memberId) {
    additionalOperators = additionalOperators.filter((id) => id !== memberId);
    const member = members.find((m) => m.id === memberId);
    showPopup(
      "Operator Dihapus",
      `${member.name} dihapus dari operator bantuan.`,
      "info",
    );
    populateStandbyOperatorsInTakeOrderPopup();
  }

  // ===== CREATE ORDER POPUP =====
  createOrderBtn.addEventListener("click", () => {
    if (!isCurrentUserAdmin()) {
      scrollToPendingWorkOrders();
      return;
    }
    showAnimatedPopup(createOrderPopup);
  });

  closeCreateOrderPopupBtn.addEventListener("click", () => {
    if (isGuestUser) {
      showPopup(
        "Guest Restriction",
        "You must create at least one work order before closing.",
        "warning",
      );
      return;
    }
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    specificLocationContainer.classList.add("hidden");
  });

  cancelCreateOrderBtn.addEventListener("click", () => {
    if (isGuestUser) {
      showPopup(
        "Guest Restriction",
        "You must create at least one work order before closing.",
        "warning",
      );
      return;
    }
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    specificLocationContainer.classList.add("hidden");
  });

  if (exitGuestBtn) {
    exitGuestBtn.addEventListener("click", async function () {
      await fetch("/api/logout", { method: "POST" }).catch(() => {});
      localStorage.removeItem("isGuestUser");
      localStorage.removeItem("guestLoginTime");
      localStorage.removeItem("userToken");
      localStorage.removeItem("currentUser");
      localStorage.removeItem("isAdmin");
      // Reload agar closure isGuestUser ter-refresh dan tabel dimuat ulang bersih
      window.location.reload();
    });
  }

  // ===== LOCATION DROPDOWN =====
  orderLocationSelect.addEventListener("change", function () {
    const locationPrompts = {
      "Gedung A": "Contoh: Lantai 2, Ruang Kelas",
      "Gedung B": "Contoh: Lantai 1, Lorong Kelas",
      "Gedung B Baru": "Contoh: Lantai 3, Lorong Kelas",
      "Gedung C": "Contoh: Lantai 1, Lorong Kelas",
      "Gedung D": "Contoh: Lantai 2, Ruang PPDB",
      "Gedung E": "Contoh: Lantai 1, Bengkel",
      "Gedung F": "Contoh: Lantai 1, Ruang Kelas",
      "Gedung G": "Contoh: Pintu Masuk Workshop",
      "Gedung TKI": "Contoh: Lantai 1",
      "Ruang Guru": "Contoh: Ruang Horenso",
      "Ruang Yayasan": "Contoh: Ruang Ketua Yayasan",
      default: "Contoh: Nomor ruang, lantai, atau area spesifik",
    };
    if (this.value) {
      specificLocationContainer.classList.remove("hidden");
      specificLocationInput.placeholder =
        locationPrompts[this.value] || locationPrompts["default"];
    } else {
      specificLocationContainer.classList.add("hidden");
    }
  });

  // ===== CREATE ORDER SUBMIT =====
  createOrderForm.addEventListener("submit", function (e) {
    e.preventDefault();
    if (!isCurrentUserAdmin()) {
      showPopup(
        "Akses Ditolak",
        "Hanya admin yang bisa membuat work order.",
        "warning",
      );
      hideAnimatedPopup(createOrderPopup);
      scrollToPendingWorkOrders();
      return;
    }
    const submitBtn = this.querySelector('[type="submit"]') || createOrderBtn;
    if (submitBtn.disabled) return;
    setButtonLoading(submitBtn, true);

    const now = new Date();
    const hh = String(now.getHours()).padStart(2, "0");
    const mm = String(now.getMinutes()).padStart(2, "0");
    const loc = document.getElementById("orderLocation").value;
    const spc = document.getElementById("specificLocation").value;

    const payload = {
      priority: document.getElementById("orderPriority").value,
      time_display: `${hh}:${mm}`,
      time_sort: `${hh}:${mm}:00`,
      requester: document.getElementById("orderRequester").value,
      location: spc ? `${loc} - ${spc}` : loc,
      device: document.getElementById("orderDevice").value,
      problem: document.getElementById("orderProblem").value,
      working_hours: "0 menit",
      status: "pending",
      executors: [],
      safety_checklist: [],
    };

    // FIX: Kirim Authorization header untuk endpoint yang dilindungi
    fetch("/api/workorders", {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify(payload),
    })
      .then((r) => {
        if (!r.ok)
          throw new Error("Gagal menyimpan order. Status: " + r.status);
        return r.json();
      })
      .then((json) => {
        const data = unwrapData(json);
        hideAnimatedPopup(createOrderPopup);
        createOrderForm.reset();
        specificLocationContainer.classList.add("hidden");
        // Guest: refresh untuk tampilkan order yang baru dibuat (hanya milik guest)
        // Non-guest: refresh semua data seperti biasa
        refreshAllDataFromAPI();
        showPopup(
          "Work Order Berhasil Dibuat!",
          `Work order ${data.trackingCode || ""} berhasil dibuat dan disimpan.`,
          "success",
        );
      })
      .catch((err) => {
        console.error("Error saat membuat order:", err);
        showPopup(
          "Error",
          "Terjadi kesalahan saat menghubungi server.",
          "error",
        );
      })
      .finally(() => setButtonLoading(submitBtn, false));
  });

  // ===== SEARCH =====
  memberSearchInput.addEventListener("focus", function () {
    searchDropdown.classList.remove("hidden");
    populateSearchResults();
  });

  memberSearchInput.addEventListener("input", function () {
    populateSearchResults(this.value.toLowerCase());
  });

  document.addEventListener("click", function (e) {
    if (!e.target.closest(".relative")) searchDropdown.classList.add("hidden");
  });

  function populateSearchResults(searchTerm = "") {
    searchResults.innerHTML = "";
    const filtered = members.filter((m) =>
      m.name.toLowerCase().includes(searchTerm),
    );

    if (filtered.length === 0) {
      searchResults.innerHTML =
        '<div class="text-center py-4 text-gray-500 text-sm">Member tidak ditemukan</div>';
      return;
    }

    const statusMap = {
      standby: { color: "bg-green-500", text: "Stand By" },
      onjob: { color: "bg-blue-500", text: "On Job" },
      nextshift: { color: "bg-purple-500", text: "Next Shift" },
      offduty: { color: "bg-gray-500", text: "Off Duty" },
    };

    filtered.forEach((member) => {
      const s = statusMap[member.status] || {
        color: "bg-gray-500",
        text: "Unknown",
      };
      const item = document.createElement("div");
      item.className =
        "flex items-center gap-3 px-4 py-2.5 hover:bg-blue-50 cursor-pointer transition-colors";
      item.innerHTML = `
        <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full object-cover flex-shrink-0">
        <div class="flex-1 min-w-0">
          <div class="font-medium text-gray-800 text-sm truncate">${member.name}</div>
          <div class="flex items-center gap-1.5 text-xs text-gray-500">
            <span class="w-2 h-2 rounded-full flex-shrink-0 ${s.color}"></span>
            <span>${s.text}</span>
            ${member.division ? `<span class="text-gray-300">·</span><span class="text-gray-400">${member.division}</span>` : ""}
          </div>
        </div>
        <svg class="w-4 h-4 text-gray-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
        </svg>`;
      item.addEventListener("click", () => {
        memberSearchInput.value = "";
        searchDropdown.classList.add("hidden");
        showOperatorProfile(member);
      });
      searchResults.appendChild(item);
    });
  }

  function showOperatorProfile(member) {
    const statusMap = {
      standby: {
        text: "Stand By",
        bg: "bg-green-100",
        color: "text-green-700",
        dot: "#22c55e",
      },
      onjob: {
        text: "On Job",
        bg: "bg-blue-100",
        color: "text-blue-700",
        dot: "#3b82f6",
      },
      nextshift: {
        text: "Next Shift",
        bg: "bg-purple-100",
        color: "text-purple-700",
        dot: "#a855f7",
      },
      offduty: {
        text: "Off Duty",
        bg: "bg-gray-100",
        color: "text-gray-600",
        dot: "#9ca3af",
      },
    };
    const headerColors = {
      standby: "linear-gradient(135deg,#14532d,#166534)",
      onjob: "linear-gradient(135deg,#1e3a8a,#1d4ed8)",
      nextshift: "linear-gradient(135deg,#4c1d95,#6d28d9)",
      offduty: "linear-gradient(135deg,#1e293b,#374151)",
    };

    const s = statusMap[member.status] || statusMap.offduty;
    const hc = headerColors[member.status] || headerColors.offduty;

    document.getElementById("opModalHeader").style.background = hc;
    document.getElementById("opModalAvatar").src =
      "/static/public/" + (member.avatar || "default-avatar.png");
    document.getElementById("opModalAvatar").alt = member.name;
    document.getElementById("opModalName").textContent = member.name;
    document.getElementById("opModalRole").textContent = member.division || "—";
    document.getElementById("opModalRoleDetail").textContent =
      member.division || "—";

    const statusEl = document.getElementById("opModalStatus");
    statusEl.textContent = s.text;
    statusEl.className = `text-sm font-semibold px-3 py-1 rounded-full ${s.bg} ${s.color}`;

    document.getElementById("operatorProfileModal").classList.remove("hidden");
  }

  // ===== MEMBER IMAGES =====
  function initializeMemberImages() {
    statusContainers.forEach((c) => {
      const container = c.querySelector(".member-images");
      if (container) container.innerHTML = "";
    });

    members.filter(isVisibleStatusMember).forEach((member) => {
      const statusContainer = document.getElementById(
        `status-${member.status}`,
      );
      if (statusContainer) {
        const container = statusContainer.querySelector(".member-images");
        const img = document.createElement("img");
        img.src = `/static/public/${member.avatar || "default-avatar.png"}`;
        img.alt = member.name;
        img.onerror = () => {
          img.onerror = null;
          img.src = "/static/public/default-avatar.png";
        };
        img.className =
          "w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm";
        img.dataset.memberId = member.id;
        container.appendChild(img);
      }
    });

    statusContainers.forEach((c) => updateMemberDisplay(c));
  }

  async function openMemberStatusPopup(statusFilter = "all") {
    if (members.length === 0) await fetchMembers();
    showAnimatedPopup(memberStatusPopup);
    populateMemberList(statusFilter);
  }

  function populateMemberList(statusFilter = "all") {
    memberList.innerHTML = "";
    const visibleMembers = members.filter(isVisibleStatusMember);
    const filtered =
      statusFilter === "all"
        ? visibleMembers
        : visibleMembers.filter((m) => m.status === statusFilter);
    const statusMap = {
      standby: {
        text: "Stand By",
        badge: "member-status-badge--standby",
      },
      onjob: {
        text: "On Job",
        badge: "member-status-badge--onjob",
      },
      nextshift: {
        text: "Next Shift",
        badge: "member-status-badge--nextshift",
      },
      offduty: {
        text: "Off Duty",
        badge: "member-status-badge--offduty",
      },
    };

    if (filtered.length === 0) {
      memberList.innerHTML =
        '<div class="text-center py-4 text-gray-500">No members found for this status</div>';
      return;
    }

    filtered.forEach((member) => {
      const status = statusMap[member.status] || statusMap.offduty;
      const item = document.createElement("div");
      item.className = "member-status-item p-4 bg-gray-50 rounded-lg";
      item.innerHTML = `
        <div class="member-status-person flex items-center gap-3">
          <img src="/static/public/${member.avatar || "default-avatar.png"}" alt="${member.name}" class="w-12 h-12 rounded-full object-cover" onerror="this.onerror=null;this.src='/static/public/default-avatar.png';">
          <span class="member-status-copy">
            <span class="member-status-name font-medium">${member.name}</span>
            <span class="member-status-division">${member.division || "Divisi belum diisi"}</span>
          </span>
        </div>
        <span class="member-status-badge ${status.badge}">${status.text}</span>`;
      memberList.appendChild(item);
    });
  }

  function updateMemberStatus(memberId, newStatus, selectEl) {
    const idx = members.findIndex((m) => m.id === memberId);
    if (idx === -1) return;

    const oldStatus = members[idx].status;
    if (selectEl) setButtonLoading(selectEl, true);

    apiUpdateMemberStatus(memberId, newStatus)
      .then(() => {
        members[idx].status = newStatus;
        if (oldStatus === "onjob" && newStatus !== "onjob")
          syncRemoveMemberFromOrders(memberId);
        updateStatusUI(members[idx], oldStatus, newStatus);
        if (
          !memberStatusPopup.classList.contains("hidden") &&
          currentStatusFilter !== "all"
        ) {
          populateMemberList(currentStatusFilter);
        }
        populateWorkOrdersTable();
        updateSummaryCounts();
        showPopup(
          "Status Diperbarui",
          `Status ${members[idx].name} berhasil diubah.`,
          "success",
        );
      })
      .catch((err) => {
        console.error("Error updating member status:", err);
        if (selectEl) selectEl.value = oldStatus;
        showPopup(
          "Error",
          "Gagal memperbarui status member ke server.",
          "error",
        );
      })
      .finally(() => {
        if (selectEl) setButtonLoading(selectEl, false);
      });
  }

  // FIX: Kirim Authorization header ke semua PATCH workorder
  function syncRemoveMemberFromOrders(memberId) {
    workOrders
      .filter((o) => o.executors && o.executors.includes(memberId))
      .forEach((order) => {
        const newExecutors = order.executors.filter((id) => id !== memberId);
        const newStatus =
          newExecutors.length === 0 && order.status === "progress"
            ? "pending"
            : order.status;
        fetch(`/api/workorders/${order.id}`, {
          method: "PATCH",
          headers: authHeaders(),
          body: JSON.stringify({ executors: newExecutors, status: newStatus }),
        })
          .then((r) => {
            if (!r.ok)
              throw new Error("Gagal sinkronisasi executor order #" + order.id);
          })
          .catch((err) => console.error(err));
      });
  }

  function updateStatusUI(member, oldStatus, newStatus) {
    const oldC = document.getElementById(`status-${oldStatus}`);
    if (oldC) {
      const img = oldC.querySelector(`img[data-member-id="${member.id}"]`);
      if (img) {
        img.remove();
        updateMemberDisplay(oldC);
      }
    }
    const newC = document.getElementById(`status-${newStatus}`);
    if (newC) {
      const container = newC.querySelector(".member-images");
      const img = document.createElement("img");
      img.src = `/static/public/${member.avatar || "default-avatar.png"}`;
      img.alt = member.name;
      img.onerror = () => {
        img.onerror = null;
        img.src = "/static/public/default-avatar.png";
      };
      img.className =
        "w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm";
      img.dataset.memberId = member.id;
      container.appendChild(img);
      updateMemberDisplay(newC);
    }
  }

  function updateMemberDisplay(container) {
    const imgs = container.querySelectorAll(".member-images img");
    const moreEl = container.querySelector(".more-members");
    if (imgs.length > 3) {
      for (let i = 3; i < imgs.length; i++) imgs[i].style.display = "none";
      moreEl.classList.remove("hidden");
      moreEl.textContent = `+${imgs.length - 3}`;
    } else {
      imgs.forEach((img) => (img.style.display = "block"));
      moreEl.classList.add("hidden");
    }
  }

  // ===== WORK ORDERS TABLE =====
  function requesterName(order) {
    return typeof order.requester === "string"
      ? order.requester
      : order.requester || "Unknown";
  }

  function orderDateKey(order) {
    const raw = order.createdAt || order.completedAt || "";
    return String(raw).slice(0, 10);
  }

  function matchesWorkOrderFilters(order) {
    const query = workOrderSearchQuery.trim().toLowerCase();
    const matchesText =
      !query ||
      [
        order.id,
        requesterName(order),
        order.device,
        order.location,
        order.problem,
        order.priority,
        order.trackingCode,
      ]
        .filter(Boolean)
        .some((value) => String(value).toLowerCase().includes(query));
    const matchesDate = !workOrderDateQuery || orderDateKey(order) === workOrderDateQuery;
    return matchesText && matchesDate;
  }

  function populateWorkOrdersTable() {
    workOrdersTableBody.innerHTML = "";
    updateWorkOrderCategoryTabs();

    const sorted = [...workOrders]
      .filter((order) => order.status === activeWorkOrderStatus)
      .filter(matchesWorkOrderFilters)
      .sort((a, b) => {
        return (
          ({ high: 0, medium: 1, low: 2 }[a.priority] ?? 3) -
          ({ high: 0, medium: 1, low: 2 }[b.priority] ?? 3)
        );
      });

    if (sorted.length === 0) {
      workOrdersTableBody.innerHTML = `<tr><td colspan="11" class="py-8 px-2 text-center text-sm text-gray-400">Tidak ada work order ${workOrderStatusText(activeWorkOrderStatus)}.</td></tr>`;
      if (typeof window.renderMobileCards === "function") {
        window.renderMobileCards(sorted, members, activeWorkOrderStatus);
      }
      return;
    }

    sorted.forEach((order, index) => {
      const row = document.createElement("tr");
      if (order.priority === "high" && order.status !== "completed")
        row.classList.add("high-priority");

      const priorityBadge = `<span class="priority-badge priority-${order.priority}">${
        { high: "High Priority", medium: "Medium", low: "Low" }[
          order.priority
        ] || order.priority
      }</span>`;

      const statusBadge = `<span class="status-badge status-${order.status}">${
        {
          pending: "Pending",
          progress: "On Progress",
          completed: "Completed",
          rejected: "Rejected",
        }[order.status] || order.status
      }</span>`;

      // Requester — backend sekarang selalu string, tapi handle juga fallback number
      const requester = requesterName(order);

      const executorIds = (order.executors || []).map(Number);
      let executorsHtml = `<button type="button" class="executors-profile-btn flex -space-x-2" data-member-ids="${executorIds.join(",")}" title="Lihat pelaksana">`;
      executorIds.forEach((execId) => {
        const member = members.find(
          (m) => m.id === execId || m.id === parseInt(execId),
        );
        if (member)
          executorsHtml += `<img src="/static/public/${member.avatar}" alt="${member.name}" class="member-avatar-small">`;
      });
      executorsHtml += "</button>";

      const canActOnProgressOrder =
        isCurrentUserAdmin() || isCurrentUserAssigned(order);
      let actionButtons = '<div class="flex items-center gap-2">';
      if (order.status === "pending") {
        if (isCurrentUserAdmin()) {
          actionButtons += `<button class="take-order-btn wo-action-btn wo-action-take" data-order-id="${order.id}" title="Approve dan mulai work order">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
          </button>`;
          actionButtons += `<button class="reject-order-btn wo-action-btn wo-action-delete" data-order-id="${order.id}" title="Reject work order">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>`;
          if ((order.executors || []).length > 0) {
            actionButtons += `<button class="add-worker-btn wo-action-btn wo-action-add" data-order-id="${order.id}" title="Tambah worker">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/></svg>
          </button>`;
          }
        } else if (!isCurrentUserAssigned(order)) {
          actionButtons += `<button class="take-order-btn wo-action-btn wo-action-take" data-order-id="${order.id}" title="Ambil order ini">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/></svg>
          </button>`;
        }
      } else if (order.status === "progress") {
        if (canActOnProgressOrder) {
          actionButtons += `<button class="documentation-btn wo-action-btn wo-action-photo" data-order-id="${order.id}" title="${order.documentationPhoto ? "Ganti foto bukti" : "Tambah foto bukti"}">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8h3l2-3h8l2 3h3v11H3V8Zm9 8a4 4 0 100-8 4 4 0 000 8Z"/></svg>
          </button>`;
          actionButtons += `<button class="done-btn wo-action-btn wo-action-done" data-order-id="${order.id}" title="Tandai sebagai selesai">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
          </button>`;
        }
      }
      if (isCurrentUserAdmin() && order.status !== "progress") {
        actionButtons += `<button class="delete-btn wo-action-btn wo-action-delete" data-order-id="${order.id}" title="Hapus order">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
      </button>`;
      }
      actionButtons += "</div>";

      row.innerHTML = `
        <td class="py-3 px-2 text-sm">${index + 1}</td>
        <td class="py-3 px-2 text-sm">${priorityBadge}</td>
        <td class="py-3 px-2 text-sm">${order.time || "-"}</td>
        <td class="py-3 px-2 text-sm">${requester}</td>
        <td class="py-3 px-2 text-sm">${order.location || "-"}</td>
        <td class="py-3 px-2 text-sm">${order.device || "-"}</td>
        <td class="py-3 px-2 text-sm">${order.problem || "-"}</td>
        <td class="py-3 px-2 text-sm">${executorsHtml}</td>
        <td class="py-3 px-2 text-sm" id="wh-${order.id}">${renderWorkingHours(order)}</td>
        <td class="py-3 px-2 text-sm">${statusBadge}</td>
        <td class="py-3 px-2 text-sm">${actionButtons}</td>`;
      workOrdersTableBody.appendChild(row);
    });

    document.querySelectorAll(".take-order-btn").forEach((b) => {
      b.addEventListener("click", function () {
        if (checkGuestRestriction("Taking orders")) return;
        openTakeOrderPopup(parseInt(this.dataset.orderId));
      });
    });
    document.querySelectorAll(".add-worker-btn").forEach((b) => {
      b.addEventListener("click", function () {
        if (checkGuestRestriction("Adding workers")) return;
        openAddWorkerPopup(parseInt(this.dataset.orderId));
      });
    });
    document.querySelectorAll(".delete-btn").forEach((b) => {
      b.addEventListener("click", function () {
        if (checkGuestRestriction("Deleting orders")) return;
        deleteOrder(parseInt(this.dataset.orderId));
      });
    });
    document.querySelectorAll(".reject-order-btn").forEach((b) => {
      b.addEventListener("click", function () {
        if (checkGuestRestriction("Rejecting orders")) return;
        rejectOrder(parseInt(this.dataset.orderId));
      });
    });
    document.querySelectorAll(".documentation-btn").forEach((b) => {
      b.addEventListener("click", function () {
        if (checkGuestRestriction("Uploading documentation")) return;
        requestDocumentationPhoto(parseInt(this.dataset.orderId), false);
      });
    });
    document.querySelectorAll(".done-btn").forEach((b) => {
      b.addEventListener("click", function () {
        if (checkGuestRestriction("Completing orders")) return;
        markOrderDone(parseInt(this.dataset.orderId));
      });
    });

    // Render mobile card view dengan data dan members yang sama
    if (typeof window.renderMobileCards === "function") {
      window.renderMobileCards(sorted, members, activeWorkOrderStatus);
    }
  }

  document.addEventListener("click", (e) => {
    const btn = e.target.closest(".executors-profile-btn");
    if (!btn) return;
    const memberIds = (btn.dataset.memberIds || "")
      .split(",")
      .map((id) => parseInt(id))
      .filter(Boolean);
    showExecutorsModal(memberIds);
  });

  function showExecutorsModal(memberIds) {
    const modal = document.getElementById("executorsModal");
    const list = document.getElementById("executorsModalList");
    if (!modal || !list) return;

    const selectedMembers = memberIds
      .map((id) => members.find((m) => Number(m.id) === id))
      .filter(Boolean);

    list.innerHTML =
      selectedMembers
        .map(
          (member) => `
          <div class="executor-list-item">
            <img src="/static/public/${member.avatar || "default-avatar.png"}" alt="${member.name}" class="executor-list-avatar">
            <div class="executor-list-copy">
              <p class="executor-list-name">${member.name}</p>
              <p class="executor-list-division">${member.division || "Divisi belum diisi"}</p>
            </div>
          </div>`,
        )
        .join("") || '<p class="text-sm text-gray-500">Belum ada pelaksana.</p>';

    modal.classList.remove("hidden");
  }

  // ===== TAKE ORDER =====
  function openTakeOrderPopup(orderId) {
    const order = workOrders.find((o) => o.id === orderId);
    if (!order) return;

    if (!isCurrentUserAdmin() && isCurrentUserAssigned(order)) {
      showPopup(
        "Menunggu Approval",
        "Kamu sudah masuk sebagai operator. Tunggu Admin atau Guru menyetujui untuk mulai progress work order.",
        "info",
      );
      return;
    }

    const standbyMembers = members.filter(isAvailableWorker);
    if (isCurrentUserAdmin() && standbyMembers.length === 0) {
      showPopup(
        "Peringatan",
        "Tidak ada pelaksana yang tersedia untuk mengambil order ini.",
        "warning",
      );
      return;
    }

    currentOrder = order;
    additionalOperators = isCurrentUserAdmin()
      ? (order.executors || []).map(Number).filter(Boolean)
      : [getCurrentUserId()];

    document.getElementById("popupOrderId").textContent =
      order.trackingCode || "-";
    document.getElementById("popupPriority").textContent =
      order.priority.charAt(0).toUpperCase() + order.priority.slice(1);
    document.getElementById("popupLocation").textContent = order.location;
    document.getElementById("popupDevice").textContent = order.device;
    document.getElementById("popupProblem").textContent = order.problem;
    if (takeOrderOperatorsTitle)
      takeOrderOperatorsTitle.textContent = isCurrentUserAdmin()
        ? "Pilih Operator"
        : "Ambil Work Order";
    if (takeOrderHelperBtnText)
      takeOrderHelperBtnText.textContent = isCurrentUserAdmin()
        ? "Tambah Operator"
        : "Ambil Work Order";
    if (openSelectHelperOperatorModalBtn)
      openSelectHelperOperatorModalBtn.classList.toggle(
        "hidden",
        !isCurrentUserAdmin(),
      );
    if (takeOrderSafetySection)
      takeOrderSafetySection.classList.toggle("hidden", !isCurrentUserAdmin());

    if (isCurrentUserAdmin()) {
      populateStandbyOperatorsInTakeOrderPopup();
    } else {
      const currentUser = members.find(
        (m) => Number(m.id) === getCurrentUserId(),
      );
      const listDiv = document.getElementById("standbyOperatorsList");
      listDiv.innerHTML = currentUser
        ? `
        <div class="selected-helper-item flex items-center gap-3 p-2 rounded-lg shadow-sm">
          <img src="/static/public/${currentUser.avatar || "default-avatar.png"}" alt="${currentUser.name}" class="w-10 h-10 rounded-full object-cover">
          <div class="selected-helper-copy flex-1 min-w-0">
            <div class="selected-helper-name font-medium">${currentUser.name}</div>
            <div class="selected-helper-role text-xs">Menunggu approval admin</div>
          </div>
        </div>`
        : '<p class="text-gray-500 text-center py-4">Akun operator tidak ditemukan.</p>';
    }
    populateSafetyChecklist(order.location);
    showAnimatedPopup(takeOrderPopup);
  }

  function openAddWorkerPopup(orderId) {
    const order = workOrders.find((o) => o.id === orderId);
    if (!order) return;
    if (order.status !== "pending") {
      showPopup(
        "Error",
        "Hanya bisa menambahkan worker ke order yang masih pending!",
        "error",
      );
      return;
    }

    currentOrder = order;
    additionalOperators = [];

    document.getElementById("popupOrderId").textContent =
      order.trackingCode || "-";
    document.getElementById("popupPriority").textContent =
      order.priority.charAt(0).toUpperCase() + order.priority.slice(1);
    document.getElementById("popupLocation").textContent = order.location;
    document.getElementById("popupDevice").textContent = order.device;
    document.getElementById("popupProblem").textContent = order.problem;

    populateAvailableWorkersForAddWorker(order.executors);
    showAnimatedPopup(selectHelperOperatorModal);
  }

  function populateAvailableWorkersForAddWorker(existingExecutorIds) {
    availableStandbyOperatorsList.innerHTML = "";
    const normalized = (existingExecutorIds || []).map((id) => parseInt(id));
    const available = members.filter(
      (m) => isAvailableWorker(m) && !normalized.includes(Number(m.id)),
    );

    if (available.length === 0) {
      availableStandbyOperatorsList.innerHTML =
        '<p class="text-gray-500 text-center py-4">Tidak ada pelaksana yang tersedia</p>';
      return;
    }

    available.forEach((member) => {
      const div = document.createElement("div");
      div.className =
        "helper-operator-item flex items-center justify-between gap-3 p-2 bg-gray-50 rounded-lg";
      div.innerHTML = `
        <div class="helper-operator-person flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full object-cover">
          <span class="helper-operator-name font-medium">${member.name}</span>
        </div>
        <button class="add-worker-direct-btn add-helper-operator-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center"
          data-member-id="${member.id}" title="Tambahkan sebagai worker">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        </button>`;
      availableStandbyOperatorsList.appendChild(div);
    });

    document.querySelectorAll(".add-worker-direct-btn").forEach((b) => {
      b.addEventListener("click", function () {
        addWorkerToOrder(parseInt(this.dataset.memberId));
      });
    });
  }

  // FIX: Kirim Authorization header
  function addWorkerToOrder(memberId) {
    if (!currentOrder) return;

    const normalized = (currentOrder.executors || []).map((id) => parseInt(id));
    if (normalized.includes(memberId)) {
      showPopup(
        "Peringatan",
        "Worker ini sudah terdaftar untuk order ini!",
        "warning",
      );
      return;
    }

    const newExecutors = [...normalized, memberId];

    fetch(`/api/workorders/${currentOrder.id}`, {
      method: "PATCH",
      headers: authHeaders(),
      body: JSON.stringify({ executors: newExecutors }),
    })
      .then((r) => {
        if (!r.ok)
          throw new Error("Gagal menambahkan worker. Status: " + r.status);
        return r.json();
      })
      .then(() => {
        const member = members.find((m) => m.id === memberId);
        hideAnimatedPopup(selectHelperOperatorModal);
        refreshAllDataFromAPI();
        showPopup(
          "Worker Ditambahkan",
          `${member.name} berhasil ditambahkan ke order #${currentOrder.id}.`,
          "success",
        );
        currentOrder = null;
      })
      .catch((err) => {
        console.error("Error saat menambahkan worker:", err);
        showPopup(
          "Error",
          "Terjadi kesalahan saat menambahkan worker.",
          "error",
        );
      });
  }

  function populateSafetyChecklist(location) {
    const div = document.getElementById("safetyChecklist");
    div.innerHTML = "";
    const mainLoc = location.includes(" - ")
      ? location.split(" - ")[0]
      : location;
    const items =
      safetyChecklistItems[mainLoc] || safetyChecklistItems["default"];

    if (items.length === 0) {
      div.innerHTML =
        '<p class="text-gray-500 text-center py-4">Tidak ada checklist safety untuk lokasi ini</p>';
      return;
    }

    items.forEach((item) => {
      const itemDiv = document.createElement("div");
      itemDiv.className = "safety-check-item flex items-center gap-3";
      const checkbox = document.createElement("input");
      checkbox.type = "checkbox";
      checkbox.className = "custom-checkbox";
      checkbox.id = item.id;
      checkbox.dataset.required = item.required;
      const label = document.createElement("label");
      label.htmlFor = item.id;
      label.className = "flex-1 cursor-pointer";
      label.innerHTML = `${item.text} ${item.required ? '<span class="text-red-500">*</span>' : ""}`;
      itemDiv.appendChild(checkbox);
      itemDiv.appendChild(label);
      div.appendChild(itemDiv);
    });
  }

  // ===== CONFIRM TAKE ORDER =====
  // FIX: Kirim Authorization header
  function confirmTakeOrder() {
    if (!currentOrder) return;
    if (isCurrentUserAdmin() && additionalOperators.length === 0) {
      showPopup(
        "Peringatan",
        "Pilih minimal satu operator sebelum mulai progress.",
        "warning",
      );
      return;
    }

    setButtonLoading(confirmTakeOrderBtn, true);

    const safetyChecklist = [];
    document.querySelectorAll("#safetyChecklist input").forEach((cb) => {
      if (cb.checked) safetyChecklist.push(cb.id);
    });

    const payload = isCurrentUserAdmin()
      ? {
          order_id: currentOrder.id,
          executors: additionalOperators,
          safety_checklist_items: safetyChecklist,
          status: "progress",
        }
      : {
          order_id: currentOrder.id,
          executors: [getCurrentUserId()],
          safety_checklist_items: [],
          status: "pending",
        };

    fetch(`/api/workorders/${currentOrder.id}/take`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify(payload),
    })
      .then((r) => {
        if (!r.ok)
          throw new Error("Gagal mengambil order. Status: " + r.status);
        return r.text().then((text) => (text ? JSON.parse(text) : {}));
      })
      .then(() => {
        if (isCurrentUserAdmin()) {
          const timerKey = `order_timer_${currentOrder.id}`;
          localStorage.setItem(timerKey, Date.now().toString());
        }
        showPopup(
          isCurrentUserAdmin()
            ? "Work Order Dimulai"
            : "Work Order Diambil",
          isCurrentUserAdmin()
            ? `Work order ${currentOrder.trackingCode || ""} mulai progress.`
            : `Kamu sudah masuk ke work order ${currentOrder.trackingCode || ""}. Tunggu Admin atau Guru menyetujui untuk mulai progress.`,
          "success",
        );
        hideAnimatedPopup(takeOrderPopup);
        resetTakeOrderForm();
        setActiveWorkOrderStatus(isCurrentUserAdmin() ? "progress" : "pending");
        refreshAllDataFromAPI();
      })
      .catch((err) => {
        console.error("Error saat konfirmasi ambil order:", err);
        showPopup(
          "Error",
          "Terjadi kesalahan saat menyimpan perubahan ke database.",
          "error",
        );
      })
      .finally(() => setButtonLoading(confirmTakeOrderBtn, false));
  }

  function resetTakeOrderForm() {
    currentOrder = null;
    additionalOperators = [];
    const listDiv = document.getElementById("standbyOperatorsList");
    if (listDiv) listDiv.innerHTML = "";
    document
      .querySelectorAll("#safetyChecklist input")
      .forEach((cb) => (cb.checked = false));
  }

  // ===== MARK ORDER DONE =====
  // FIX: Kirim Authorization header
  function markOrderDone(orderId, skipPhotoCheck = false) {
    const order = workOrders.find((o) => o.id === orderId);
    if (!order) return;
    if (!isCurrentUserAdmin() && !isCurrentUserAssigned(order)) {
      showPopup(
        "Access Denied",
        "Hanya operator yang ditugaskan atau admin yang bisa menyelesaikan work order ini.",
        "warning",
      );
      return;
    }
    if (!skipPhotoCheck && !order.documentationPhoto) {
      showPopup(
        "Foto Bukti Diperlukan",
        "Ambil foto dokumentasi terlebih dahulu. Setelah foto tersimpan, order akan otomatis diselesaikan.",
        "warning",
      );
      requestDocumentationPhoto(orderId, true);
      return;
    }

    const doneBtn = document.querySelector(
      `.done-btn[data-order-id="${orderId}"], .mc-done-btn[data-order-id="${orderId}"]`,
    );
    setButtonLoading(doneBtn, true);

    const now = new Date();
    const hh = String(now.getHours()).padStart(2, "0");
    const mm = String(now.getMinutes()).padStart(2, "0");
    const completionTime = `${hh}:${mm}`;

    fetch(`/api/workorders/${orderId}/complete`, {
      method: "PATCH",
      headers: authHeaders(),
      body: JSON.stringify({
        status: "completed",
        completed_at_display: completionTime,
      }),
    })
      .then((r) => {
        if (!r.ok)
          throw new Error("Gagal menyelesaikan order. Status: " + r.status);
        return r.text().then((text) => (text ? JSON.parse(text) : {}));
      })
      .then(() => {
        // Simpan waktu terakhir stopwatch sebelum dihapus, untuk ditampilkan saat completed
        const finalStartedAt = localStorage.getItem(`order_timer_${orderId}`);
        if (finalStartedAt) {
          const finalElapsed = Math.floor(
            (Date.now() - parseInt(finalStartedAt)) / 1000,
          );
          localStorage.setItem(
            `order_timer_final_${orderId}`,
            formatElapsed(finalElapsed),
          );
          localStorage.removeItem(`order_timer_${orderId}`);
        }
        setActiveWorkOrderStatus("completed");
        refreshAllDataFromAPI();
        showPopup(
          "Order Selesai!",
          `Work order ${order.trackingCode || ""} berhasil ditandai selesai!\nWaktu selesai: ${completionTime}`,
          "success",
        );
      })
      .catch((err) => {
        console.error("Error saat menyelesaikan order:", err);
        showPopup(
          "Error",
          "Terjadi kesalahan saat memperbarui status order.",
          "error",
        );
        setButtonLoading(doneBtn, false);
      });
  }

  // ===== DELETE ORDER =====
  // FIX: Kirim Authorization header
  function deleteOrder(orderId) {
    const order = workOrders.find((item) => item.id === orderId);
    const reference = order?.trackingCode || "ini";
    if (order?.status === "progress") {
      showPopup(
        "Order Sedang Dikerjakan",
        "Work order On Progress harus diselesaikan sebelum dapat dihapus.",
        "error",
      );
      return;
    }
    showConfirmationPopup(
      "Konfirmasi Hapus Order",
      `Apakah Anda yakin ingin menghapus work order ${reference}?`,
      () => {
        fetch(`/api/workorders/${orderId}`, {
          method: "DELETE",
          headers: authHeaders(),
        })
          .then((r) => {
            if (!r.ok)
              throw new Error("Gagal menghapus order. Status: " + r.status);
            return r.text().then((text) => (text ? JSON.parse(text) : {}));
          })
          .then(() => {
            refreshAllDataFromAPI();
            showPopup(
              "Order Dihapus!",
              `Work order ${reference} telah berhasil dihapus dari database.`,
              "success",
            );
          })
          .catch((err) => {
            console.error("Error saat menghapus order:", err);
            showPopup(
              "Error",
              "Terjadi kesalahan saat menghapus order.",
              "error",
            );
          });
      },
    );
  }

  function rejectOrder(orderId) {
    showConfirmationPopup(
      "Konfirmasi Reject Order",
      `Apakah benar ingin reject work order #${orderId}?`,
      () => {
        showTextInputPopup(
          "Berikan catatan untuk requester",
          `Catatan ini akan tampil di tracking work order #${orderId}.`,
          (reason) => {
            fetch(`/api/workorders/${orderId}/reject`, {
              method: "PATCH",
              headers: authHeaders(),
              body: JSON.stringify({ reason }),
            })
              .then((r) => {
                if (!r.ok)
                  throw new Error("Gagal reject order. Status: " + r.status);
                return r.text().then((text) => (text ? JSON.parse(text) : {}));
              })
              .then(() => {
                setActiveWorkOrderStatus("rejected");
                refreshAllDataFromAPI();
                showPopup(
                  "Order Ditolak",
                  `Work order #${orderId} sudah direject dan catatan bisa dilihat requester.`,
                  "success",
                );
              })
              .catch((err) => {
                console.error("Error saat reject order:", err);
                showPopup(
                  "Error",
                  "Terjadi kesalahan saat reject order.",
                  "error",
                );
              });
          },
        );
      },
    );
  }
  window.rejectOrder = rejectOrder;

  // ===== LIVE WORKING HOURS COUNTER =====

  // Format detik → HH.MM.SS (stopwatch style)
  function formatElapsed(totalSeconds) {
    const h = Math.floor(totalSeconds / 3600);
    const m = Math.floor((totalSeconds % 3600) / 60);
    const s = totalSeconds % 60;
    return `${String(h).padStart(2, "0")}.${String(m).padStart(2, "0")}.${String(s).padStart(2, "0")}`;
  }

  function formatWorkingMinutes(minutes) {
    const value = Number(minutes);
    if (!Number.isFinite(value)) return "-";
    if (value < 60) return `${value} mnt`;
    const h = Math.floor(value / 60);
    const m = value % 60;
    return m ? `${h} jam ${m} mnt` : `${h} jam`;
  }

  // Render initial value untuk kolom working hours
  function renderWorkingHours(order) {
    if (order.status === "completed") {
      localStorage.removeItem(`order_timer_${order.id}`);
      if (order.workingHours != null)
        return formatWorkingMinutes(order.workingHours);
      // Fallback hanya kalau backend belum punya durasi, misalnya Rust timer gagal stop.
      const savedTime = localStorage.getItem(`order_timer_final_${order.id}`);
      if (savedTime) return savedTime;
      return "-";
    }
    if (order.status === "progress") {
      const serverElapsed = elapsedFromServerStartedAt(order);
      if (serverElapsed != null) return formatElapsed(serverElapsed);
      const startedAt = localStorage.getItem(`order_timer_${order.id}`);
      if (startedAt) {
        const elapsed = Math.floor((Date.now() - parseInt(startedAt)) / 1000);
        return formatElapsed(elapsed);
      }
      // Order progress tapi tidak ada timer lokal (diambil di device lain)
      return "00.00.00";
    }
    return "-";
  }

  // Tick setiap detik — update stopwatch hanya untuk order yang masih progress.
  // Kalau order diselesaikan dari device lain, status server menjadi completed dan
  // timer localStorage desktop harus dihentikan agar tidak terus jalan.
  setInterval(() => {
    document.querySelectorAll('[id^="wh-"]').forEach((cell) => {
      const orderId = cell.id.slice(3); // hapus "wh-"
      const order = workOrders.find((o) => String(o.id) === String(orderId));
      if (!order) return;
      if (order.status !== "progress") {
        localStorage.removeItem(`order_timer_${orderId}`);
        cell.textContent = renderWorkingHours(order);
        return;
      }
      const serverElapsed = elapsedFromServerStartedAt(order);
      if (serverElapsed != null) {
        cell.textContent = formatElapsed(serverElapsed);
        return;
      }
      const startedAt = localStorage.getItem(`order_timer_${orderId}`);
      if (!startedAt) return;
      const elapsed = Math.floor((Date.now() - parseInt(startedAt)) / 1000);
      cell.textContent = formatElapsed(elapsed);
    });
  }, 1000);

  // ===== SUMMARY COUNTS =====
  function workOrderStatusText(status) {
    return (
      {
        pending: "pending",
        progress: "on progress",
        completed: "completed",
        rejected: "rejected",
      }[status] || status
    );
  }

  function setText(id, value) {
    const el = document.getElementById(id);
    if (el) el.textContent = value;
  }

  function updateWorkOrderCategoryTabs() {
    const counts = {
      pending: workOrders.filter((o) => o.status === "pending").length,
      progress: workOrders.filter((o) => o.status === "progress").length,
      completed: workOrders.filter((o) => o.status === "completed").length,
      rejected: workOrders.filter((o) => o.status === "rejected").length,
    };

    setText("categoryPendingCount", counts.pending);
    setText("categoryProgressCount", counts.progress);
    setText("categoryCompletedCount", counts.completed);
    setText("categoryRejectedCount", counts.rejected);

    workOrderStatusTabs.forEach((tab) => {
      const isActive = tab.dataset.workOrderStatus === activeWorkOrderStatus;
      tab.classList.toggle("active", isActive);
      tab.setAttribute("aria-pressed", String(isActive));
    });
  }

  function setActiveWorkOrderStatus(status) {
    activeWorkOrderStatus = status || "pending";
    populateWorkOrdersTable();
  }

  function updateSummaryCounts() {
    document.getElementById("totalOrdersCount").textContent = workOrders.length;
    document.getElementById("pendingOrdersCount").textContent =
      workOrders.filter((o) => o.status === "pending").length;
    document.getElementById("progressOrdersCount").textContent =
      workOrders.filter((o) => o.status === "progress").length;
    document.getElementById("completedOrdersCount").textContent =
      workOrders.filter((o) => o.status === "completed").length;
    updateWorkOrderCategoryTabs();
    renderHeroOperations();
  }

  function renderHeroOperations() {
    const list = document.getElementById("heroWorkOrderList");
    if (!list) return;

    const activeOrders = workOrders.filter((order) =>
      ["pending", "progress"].includes(order.status),
    );
    const progressOrders = activeOrders.filter(
      (order) => order.status === "progress",
    );
    const highPriorityOrders = activeOrders.filter(
      (order) => order.priority === "high",
    );
    const trackableOrders = workOrders.filter(
      (order) => order.status !== "rejected",
    );
    const completionRate = trackableOrders.length
      ? Math.round(
          (trackableOrders.filter((order) => order.status === "completed")
            .length /
            trackableOrders.length) *
            100,
        )
      : 0;
    const onDutyStaff = members.filter(
      (member) =>
        isVisibleStatusMember(member) &&
        ["standby", "onjob"].includes(member.status),
    ).length;

    setText("heroActiveOrders", activeOrders.length);
    setText("heroOnDutyStaff", onDutyStaff);
    setText("heroCompletionRate", `${completionRate}%`);
    setText("heroPreviewActive", activeOrders.length);
    setText("heroPreviewProgress", progressOrders.length);
    setText("heroPreviewHigh", highPriorityOrders.length);
    setText("heroPreviewCompletion", `${completionRate}%`);

    const progress = document.getElementById("heroCompletionProgress");
    const progressBar = document.getElementById("heroCompletionBar");
    progress?.setAttribute("aria-valuenow", String(completionRate));
    if (progressBar) progressBar.style.width = `${completionRate}%`;

    list.replaceChildren();
    const previewOrders = [...activeOrders]
      .sort(
        (a, b) =>
          ({ high: 0, medium: 1, low: 2 }[a.priority] ?? 3) -
            ({ high: 0, medium: 1, low: 2 }[b.priority] ?? 3) ||
          Number(b.id) - Number(a.id),
      )
      .slice(0, 3);

    if (previewOrders.length === 0) {
      const empty = document.createElement("p");
      empty.className = "operations-empty";
      empty.textContent = "Tidak ada work order aktif saat ini.";
      list.appendChild(empty);
      return;
    }

    const priorityLabels = {
      high: "Tinggi",
      medium: "Sedang",
      low: "Rendah",
    };
    const statusLabels = {
      pending: "Menunggu",
      progress: "Dikerjakan",
    };

    previewOrders.forEach((order) => {
      const row = document.createElement("article");
      row.className = "operations-row";

      const main = document.createElement("div");
      main.className = "operations-row-main";
      const top = document.createElement("div");
      top.className = "operations-row-top";
      const title = document.createElement("p");
      title.className = "operations-row-title";
      title.textContent = `${order.device || "Perangkat IT"} · ${order.problem || `Work order #${order.id}`}`;
      const priority = document.createElement("span");
      const priorityKey = priorityLabels[order.priority]
        ? order.priority
        : "low";
      priority.className = `operations-priority operations-priority-${priorityKey}`;
      priority.textContent = priorityLabels[priorityKey];
      priority.setAttribute(
        "aria-label",
        `Prioritas ${priorityLabels[priorityKey]}`,
      );
      top.append(title, priority);

      const meta = document.createElement("div");
      meta.className = "operations-row-meta";
      const location = document.createElement("span");
      location.textContent = order.location || "Lokasi belum diisi";
      meta.appendChild(location);

      const assignees = document.createElement("span");
      assignees.className = "operations-assignees";
      (order.executors || []).slice(0, 2).forEach((executorId) => {
        const member = members.find(
          (item) => Number(item.id) === Number(executorId),
        );
        if (!member) return;
        const photo = document.createElement("img");
        photo.className = "operations-assignee";
        const memberName = String(member.name || "Staf IT");
        photo.src = `/static/public/${member.avatar || "default-avatar.png"}`;
        photo.alt = memberName;
        photo.title = memberName;
        photo.onerror = () => {
          photo.onerror = null;
          photo.src = "/static/public/default-avatar.png";
        };
        assignees.appendChild(photo);
      });
      main.append(top, meta);

      const statusKey = statusLabels[order.status] ? order.status : "pending";
      const status = document.createElement("span");
      status.className = `operations-status operations-status-${statusKey}`;
      status.textContent = statusLabels[statusKey];
      row.append(main, assignees, status);
      list.appendChild(row);
    });
  }

  // ===== OPERATOR PROFILE MODAL (dari search bar) =====
  const operatorProfileModal = document.getElementById("operatorProfileModal");
  if (operatorProfileModal) {
    document
      .getElementById("closeOperatorProfileModal")
      .addEventListener("click", () => {
        operatorProfileModal.classList.add("hidden");
      });
    operatorProfileModal.addEventListener("click", (e) => {
      if (e.target === operatorProfileModal)
        operatorProfileModal.classList.add("hidden");
    });
  }

  const executorsModal = document.getElementById("executorsModal");
  if (executorsModal) {
    document
      .getElementById("closeExecutorsModal")
      ?.addEventListener("click", () => executorsModal.classList.add("hidden"));
    executorsModal.addEventListener("click", (e) => {
      if (e.target === executorsModal) executorsModal.classList.add("hidden");
    });
  }

  // ===== PRIORITY RADIO SYNC =====
  // Sync visual priority selector → hidden <select id="orderPriority">
  document
    .querySelectorAll('#prioritySelector input[type="radio"]')
    .forEach((radio) => {
      radio.addEventListener("change", function () {
        const sel = document.getElementById("orderPriority");
        if (sel) sel.value = this.value;
      });
    });
  // Set default: medium checked on load
  const defaultPriorityRadio = document.querySelector(
    '#prioritySelector input[value="medium"]',
  );
  if (defaultPriorityRadio) defaultPriorityRadio.checked = true;

  // ===== FILTER TABS =====
  function updateFilterTabs(activeStatus) {
    statusFilterTabs.forEach((tab) => {
      const isActive = tab.dataset.statusFilter === activeStatus;
      tab.classList.toggle("bg-blue-500", isActive);
      tab.classList.toggle("text-white", isActive);
      tab.classList.toggle("bg-gray-200", !isActive);
      tab.classList.toggle("text-gray-700", !isActive);
    });
  }

  // Expose fungsi-fungsi yang dibutuhkan oleh renderMobileCards di index.html
  // Semua fungsi ini ada di dalam closure DOMContentLoaded sehingga tidak
  // accessible dari luar tanpa di-assign ke window secara eksplisit.
  window.openTakeOrderPopup = openTakeOrderPopup;
  window.openAddWorkerPopup = openAddWorkerPopup;
  window.markOrderDone = markOrderDone;
  window.deleteOrder = deleteOrder;
  window.checkGuestRestriction = checkGuestRestriction;
  window.refreshDashboardData = refreshAllDataFromAPI;
}); // END DOMContentLoaded

// ===== LOGIN / REGISTER PAGE LOGIC =====
// FIX: Logika login/register dipindah sepenuhnya ke inline script di login.html & register.html.
// Blok ini sekarang hanya menangani hal yang TIDAK bisa diletakkan di inline script:
// yaitu elemen yang mungkin ada di KEDUA halaman (tidak ada saat ini).
// Duplikasi event listener yang sebelumnya ada di sini sudah dihapus agar tidak
// konflik dengan script inline di login.html dan register.html.
