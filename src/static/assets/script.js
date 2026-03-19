// ===== UTILITY: JWT TOKEN =====
// Helper terpusat untuk mengambil token dari localStorage.
// Semua fetch ke protected endpoint memanggil ini.
function getAuthToken() {
  return localStorage.getItem('userToken') || '';
}

// Header default untuk request yang butuh autentikasi
function authHeaders(extra = {}) {
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${getAuthToken()}`,
    ...extra
  };
}

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
const refreshBtn  = document.getElementById('refreshBtn');
const refreshIcon = document.getElementById('refreshIcon');

if (refreshBtn) {
  refreshBtn.addEventListener('click', function () {
    refreshIcon.style.transform  = 'rotate(0deg)';
    refreshIcon.style.transition = 'transform 0.6s linear';
    setTimeout(() => { refreshIcon.style.transform = 'rotate(360deg)'; }, 10);
    setTimeout(() => { window.location.reload(); }, 300);
  });
}

// ===== PROFILE DROPDOWN =====
// FIX: Semua akses elemen navbar dibungkus null-check.
// Script ini di-load di login.html dan register.html juga,
// di mana elemen-elemen navbar tidak ada — tanpa null-check
// baris `btn.addEventListener(...)` akan crash dan menghentikan
// seluruh eksekusi script.
const btn  = document.getElementById('profileDropdownBtn');
const menu = document.getElementById('profileDropdown');

if (btn && menu) {
  // FIX: hanya JS click yang mengontrol dropdown (group-hover CSS dihapus dari index.html)
  btn.addEventListener('click', (e) => {
    e.stopPropagation(); // cegah document click langsung menutup lagi
    menu.classList.toggle('hidden');
  });

  // Logout button — ambil tombol kedua (index 1) karena tombol pertama adalah Settings
  const logoutBtn = menu.querySelectorAll('button')[1];
  if (logoutBtn) {
    logoutBtn.addEventListener('click', function () {
      localStorage.removeItem('isGuestUser');
      localStorage.removeItem('guestLoginTime');
      localStorage.removeItem('isAdmin');
      localStorage.removeItem('userToken');
      window.location.href = 'login.html';
    });
  }

  // Tutup dropdown saat klik di luar
  document.addEventListener('click', (e) => {
    if (!btn.contains(e.target) && !menu.contains(e.target)) {
      menu.classList.add('hidden');
    }
  });
}

// ===== MOBILE MENU =====
const mobileMenuBtn = document.getElementById('mobileMenuBtn');
const navMenu       = document.getElementById('navMenu');

if (mobileMenuBtn && navMenu) {
  mobileMenuBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    // CSS sudah handle show/hide via media query + .mobile-menu-active
    // JS hanya toggle class — tidak perlu cek innerWidth
    navMenu.classList.toggle('mobile-menu-active');
    // Update aria untuk aksesibilitas
    const isOpen = navMenu.classList.contains('mobile-menu-active');
    mobileMenuBtn.setAttribute('aria-expanded', isOpen);
  });

  // Tutup menu saat klik di luar area navbar
  document.addEventListener('click', (e) => {
    if (!mobileMenuBtn.contains(e.target) && !navMenu.contains(e.target)) {
      navMenu.classList.remove('mobile-menu-active');
      mobileMenuBtn.setAttribute('aria-expanded', 'false');
    }
  });

  // Tutup menu saat salah satu link/button di dalamnya diklik
  navMenu.querySelectorAll('a, button').forEach(el => {
    el.addEventListener('click', () => {
      navMenu.classList.remove('mobile-menu-active');
      mobileMenuBtn.setAttribute('aria-expanded', 'false');
    });
  });

  // Reset saat resize ke desktop agar state tidak stuck
  window.addEventListener('resize', () => {
    if (window.innerWidth >= 768) {
      navMenu.classList.remove('mobile-menu-active');
      mobileMenuBtn.setAttribute('aria-expanded', 'false');
    }
  });
}

// ===== GSAP ANIMATION =====
// FIX: Cek keberadaan GSAP sebelum memanggilnya.
// Jika library tidak ter-load (network error, CDN down),
// script tidak crash dan halaman tetap fungsional.
window.addEventListener('load', () => {
  if (typeof gsap !== 'undefined') {
    gsap.utils.toArray('.fade-slide').forEach((el, i) => {
      gsap.to(el, { opacity: 1, y: 0, duration: 0.6, delay: i * 0.15, ease: 'power2.out' });
    });
  }
});

// ===== CUSTOM POPUP SYSTEM =====
function showPopup(title, message, type = 'info') {
  const existingPopup = document.getElementById('customPopup');
  if (existingPopup) existingPopup.remove();

  const popup = document.createElement('div');
  popup.id        = 'customPopup';
  popup.className = 'fixed inset-0 bg-black bg-opacity-50 z-[100] flex items-center justify-center';

  const popupContent = document.createElement('div');
  popupContent.className = 'bg-white rounded-2xl shadow-2xl p-6 w-11/12 max-w-md transform transition-all popup-fade-in';

  const configs = {
    success: { icon: `<svg class="w-12 h-12 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`, bg: 'from-green-50 to-green-100' },
    warning: { icon: `<svg class="w-12 h-12 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path></svg>`, bg: 'from-yellow-50 to-yellow-100' },
    error:   { icon: `<svg class="w-12 h-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`, bg: 'from-red-50 to-red-100' },
    info:    { icon: `<svg class="w-12 h-12 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`, bg: 'from-blue-50 to-blue-100' },
  };
  const { icon, bg } = configs[type] || configs.info;

  popupContent.innerHTML = `
    <div class="text-center">
      <div class="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-gradient-to-br ${bg} mb-4">${icon}</div>
      <h3 class="text-xl font-bold text-gray-900 mb-2">${title}</h3>
      <p class="text-gray-600 mb-6 leading-relaxed">${message}</p>
      <button class="px-6 py-3 bg-blue-500 text-white font-semibold rounded-lg hover:bg-blue-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-blue-200">OK</button>
    </div>`;

  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  const closePopup = () => {
    popupContent.classList.replace('popup-fade-in', 'popup-fade-out');
    setTimeout(() => popup.remove(), 300);
  };

  popup.querySelector('button').addEventListener('click', closePopup);

  if (type !== 'error') {
    setTimeout(() => { if (popup.parentNode) closePopup(); }, 5000);
  }
}

// ===== CONFIRMATION POPUP =====
function showConfirmationPopup(title, message, onConfirm) {
  const existingPopup = document.getElementById('customConfirmationPopup');
  if (existingPopup) existingPopup.remove();

  const popup = document.createElement('div');
  popup.id        = 'customConfirmationPopup';
  popup.className = 'fixed inset-0 bg-black bg-opacity-50 z-[100] flex items-center justify-center';

  const popupContent = document.createElement('div');
  popupContent.className = 'bg-white rounded-2xl shadow-2xl p-6 w-11/12 max-w-md transform transition-all popup-fade-in';
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
    popupContent.classList.replace('popup-fade-in', 'popup-fade-out');
    setTimeout(() => popup.remove(), 300);
  };

  document.getElementById('confirmBtn').addEventListener('click', () => { onConfirm(); closeThis(); });
  document.getElementById('cancelBtn').addEventListener('click', closeThis);
}

// ===== POPUP ANIMATION HELPERS =====
function showAnimatedPopup(popupElement) {
  const content = popupElement.firstElementChild;
  popupElement.classList.remove('hidden');
  content.classList.remove('popup-fade-out');
  content.classList.add('popup-fade-in');
}

function hideAnimatedPopup(popupElement) {
  const content = popupElement.firstElementChild;
  content.classList.remove('popup-fade-in');
  content.classList.add('popup-fade-out');
  setTimeout(() => popupElement.classList.add('hidden'), 300);
}

// ===== QUICK SUMMARY TITLE =====
function updateQuickSummaryTitle() {
  const el = document.getElementById('quickSummaryTitle');
  if (el) {
    el.textContent = new Date().getDate() === 1 ? 'Quick Summary Bulan Ini' : 'Quick Summary';
  }
}

// ===== LOADING STATE =====
function setButtonLoading(btn, isLoading) {
  if (!btn) return;
  btn.disabled      = isLoading;
  btn.style.opacity = isLoading ? '0.6' : '1';
  btn.style.cursor  = isLoading ? 'not-allowed' : '';
}

// ===== MAIN DASHBOARD LOGIC =====
document.addEventListener('DOMContentLoaded', async function () {

  // Script ini di-load di login.html & register.html juga.
  // Elemen-elemen dashboard tidak akan ada di sana.
  // FIX: cek keberadaan elemen kunci sebelum menjalankan logika dashboard.
  const workOrdersTableBody = document.getElementById('workOrdersTableBody');
  if (!workOrdersTableBody) return; // Bukan halaman dashboard, stop di sini

  updateQuickSummaryTitle();

  let members            = [];
  let workOrders         = [];
  let currentOrder       = null;
  let additionalOperators = [];

  const isGuestUser = localStorage.getItem('isGuestUser') === 'true';

  if (isGuestUser) {
    setTimeout(() => {
      const createOrderPopup = document.getElementById('createOrderPopup');
      if (createOrderPopup) createOrderPopup.classList.remove('hidden');
    }, 500);
  }

  function checkGuestRestriction(action = 'action') {
    if (isGuestUser) {
      showPopup('Access Denied', `Guests can only create work orders. ${action} is not allowed.`, 'warning');
      return true;
    }
    return false;
  }

  // ===== DOM REFERENCES =====
  const memberStatusPopup              = document.getElementById('memberStatusPopup');
  const memberList                     = document.getElementById('memberList');
  const closePopupBtn                  = document.getElementById('closePopup');
  const statusContainers               = document.querySelectorAll('.status-container');
  const takeOrderPopup                 = document.getElementById('takeOrderPopup');
  const closeTakeOrderPopupBtn         = document.getElementById('closeTakeOrderPopup');
  const cancelTakeOrderBtn             = document.getElementById('cancelTakeOrderBtn');
  const confirmTakeOrderBtn            = document.getElementById('confirmTakeOrderBtn');
  const openSelectHelperOperatorModalBtn = document.getElementById('openSelectHelperOperatorModalBtn');
  const createOrderPopup               = document.getElementById('createOrderPopup');
  const closeCreateOrderPopupBtn       = document.getElementById('closeCreateOrderPopup');
  const cancelCreateOrderBtn           = document.getElementById('cancelCreateOrderBtn');
  const exitGuestBtn                   = document.getElementById('exitGuestBtn');
  const createOrderForm                = document.getElementById('createOrderForm');
  const createOrderBtn                 = document.getElementById('createOrderBtn');
  const orderLocationSelect            = document.getElementById('orderLocation');
  const specificLocationContainer      = document.getElementById('specificLocationContainer');
  const specificLocationInput          = document.getElementById('specificLocation');
  const memberSearchInput              = document.getElementById('memberSearchInput');
  const searchDropdown                 = document.getElementById('searchDropdown');
  const searchResults                  = document.getElementById('searchResults');
  const statusFilterTabs               = document.querySelectorAll('.status-filter-tab');
  const selectHelperOperatorModal      = document.getElementById('selectHelperOperatorModal');
  const closeSelectHelperOperatorModalBtn = document.getElementById('closeSelectHelperOperatorModalBtn');
  const availableStandbyOperatorsList  = document.getElementById('availableStandbyOperatorsList');

  let currentStatusFilter = 'all';

  // ===== SAFETY CHECKLIST DATA =====
  const safetyChecklistItems = {
    'CCTV':        [{ id: 'cctv1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'cctv2', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'cctv3', text: 'Pastikan area kerja aman', required: true }, { id: 'cctv4', text: 'Gunakan Sarung Tangan', required: true }],
    'WiFi':        [{ id: 'wifi1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'wifi2', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'wifi3', text: 'Pastikan area kerja aman', required: true }, { id: 'wifi4', text: 'Gunakan Sarung Tangan', required: true }],
    'Gedung A':    [{ id: 'ga1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'ga2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'ga3', text: 'Pastikan area kerja aman', required: true }, { id: 'ga4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'ga5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung B':    [{ id: 'gb1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'gb2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'gb3', text: 'Pastikan area kerja aman', required: true }, { id: 'gb4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'gb5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung B Baru': [{ id: 'gbb1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'gbb2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'gbb3', text: 'Pastikan area kerja aman', required: true }, { id: 'gbb4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'gbb5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung C':    [{ id: 'gc1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'gc2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'gc3', text: 'Pastikan area kerja aman', required: true }, { id: 'gc4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'gc5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung D':    [{ id: 'gd1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'gd2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'gd3', text: 'Pastikan area kerja aman', required: true }, { id: 'gd4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'gd5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung E':    [{ id: 'ge1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'ge2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'ge3', text: 'Pastikan area kerja aman', required: true }, { id: 'ge4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'ge5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung F':    [{ id: 'gf1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'gf2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'gf3', text: 'Pastikan area kerja aman', required: true }, { id: 'gf4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'gf5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung G':    [{ id: 'gg1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'gg2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'gg3', text: 'Pastikan area kerja aman', required: true }, { id: 'gg4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'gg5', text: 'Gunakan sepatu safety', required: true }],
    'Gedung TKI':  [{ id: 'gt1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'gt2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'gt3', text: 'Pastikan area kerja aman', required: true }, { id: 'gt4', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'gt5', text: 'Gunakan sepatu safety', required: true }],
    'Ruang Guru':  [{ id: 'rg1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'rg2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'rg3', text: 'Pastikan area kerja aman', required: false }, { id: 'rg4', text: 'Matikan listrik sebelum bekerja', required: true }],
    'Ruang Yayasan': [{ id: 'ry1', text: 'Pastikan sirkulasi udara baik', required: false }, { id: 'ry2', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'ry3', text: 'Matikan listrik sebelum bekerja', required: true }, { id: 'ry4', text: 'Gunakan Sarung Tangan', required: false }, { id: 'ry5', text: 'Pastikan area kerja aman', required: true }],
    'default':     [{ id: 'def1', text: 'Gunakan pelindung mata (goggles)', required: false }, { id: 'def2', text: 'Gunakan Sarung Tangan', required: false }, { id: 'def3', text: 'Pastikan area kerja aman', required: true }, { id: 'def4', text: 'Matikan listrik sebelum bekerja', required: true }],
  };

  // ===== API FUNCTIONS =====

  async function fetchMembers() {
    try {
      const r = await fetch('/api/members');
      if (!r.ok) throw new Error(r.statusText);
      const json = await r.json();
      // FIX: handle format response baru { code, message, data: [...] }
      members = Array.isArray(json) ? json : (Array.isArray(json.data) ? json.data : []);
    } catch (err) {
      console.error('Error fetching members:', err);
      if (memberList) memberList.innerHTML = '<div class="text-center py-4 text-red-500">Failed to load member data.</div>';
      members = [];
    }
  }

  async function fetchAndRenderWorkOrders() {
    try {
      const r = await fetch('/api/workorders');
      if (!r.ok) throw new Error('Gagal mengambil data work order dari server');
      const json = await r.json();
      // FIX: handle format response baru { code, message, data: [...] }
      workOrders = Array.isArray(json) ? json : (Array.isArray(json.data) ? json.data : []);
    } catch (err) {
      console.error('Error fetching work orders:', err);
      showPopup('Error', 'Gagal memuat data work order dari server.', 'error');
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

  // FIX: Kirim Authorization header ke semua protected endpoint
  async function apiUpdateMemberStatus(memberId, newStatus) {
    const r = await fetch(`/api/members/${memberId}/status`, {
      method: 'PATCH',
      headers: authHeaders(),
      body: JSON.stringify({ status: newStatus })
    });
    if (!r.ok) throw new Error('Gagal memperbarui status member. Status: ' + r.status);
    const json = await r.json();
    return unwrapData(json);
  }

  // ===== INIT =====
  await fetchMembers();
  await fetchAndRenderWorkOrders();
  initializeMemberImages();
  updateSummaryCounts();

  if (isGuestUser && exitGuestBtn) exitGuestBtn.classList.remove('hidden');

  // ===== STATUS CONTAINER CLICK =====
  statusContainers.forEach(container => {
    container.addEventListener('click', function (e) {
      if (!e.target.closest('.member-images') && !e.target.closest('.more-members')) {
        if (checkGuestRestriction('Viewing/managing member status')) return;
        const status = this.dataset.status;
        currentStatusFilter = status;
        updateFilterTabs(status);
        openMemberStatusPopup(status);
      }
    });
  });

  // ===== MEMBER STATUS POPUP =====
  closePopupBtn.addEventListener('click', () => hideAnimatedPopup(memberStatusPopup));

  statusFilterTabs.forEach(tab => {
    tab.addEventListener('click', function () {
      const status = this.dataset.statusFilter;
      currentStatusFilter = status;
      updateFilterTabs(status);
      populateMemberList(status);
    });
  });

  // ===== TAKE ORDER POPUP =====
  closeTakeOrderPopupBtn.addEventListener('click', () => { hideAnimatedPopup(takeOrderPopup); resetTakeOrderForm(); });
  cancelTakeOrderBtn.addEventListener('click',     () => { hideAnimatedPopup(takeOrderPopup); resetTakeOrderForm(); });
  confirmTakeOrderBtn.addEventListener('click',    () => confirmTakeOrder());
  openSelectHelperOperatorModalBtn.addEventListener('click', openSelectHelperOperatorModal);

  // ===== HELPER OPERATOR MODAL =====
  closeSelectHelperOperatorModalBtn.addEventListener('click', () => hideAnimatedPopup(selectHelperOperatorModal));

  function openSelectHelperOperatorModal() {
    showAnimatedPopup(selectHelperOperatorModal);
    populateAvailableStandbyOperators();
  }

  function populateAvailableStandbyOperators() {
    availableStandbyOperatorsList.innerHTML = '';
    const standbyMembers = members.filter(m => m.status === 'standby');

    if (standbyMembers.length === 0) {
      availableStandbyOperatorsList.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada operator standby tersedia</p>';
      return;
    }

    standbyMembers.forEach(member => {
      const div = document.createElement('div');
      div.className = 'flex items-center justify-between p-2 bg-gray-50 rounded-lg';
      div.innerHTML = `
        <div class="flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full">
          <span class="font-medium">${member.name}</span>
        </div>
        <button class="add-helper-operator-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center"
          data-member-id="${member.id}" title="Tambahkan sebagai operator bantuan">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        </button>`;
      availableStandbyOperatorsList.appendChild(div);
    });

    document.querySelectorAll('.add-helper-operator-btn').forEach(b => {
      b.addEventListener('click', function () { addHelperOperator(parseInt(this.dataset.memberId)); });
    });
  }

  function addHelperOperator(memberId) {
    if (additionalOperators.includes(memberId)) {
      showPopup('Peringatan', 'Operator ini sudah ditambahkan!', 'warning');
      return;
    }
    additionalOperators.push(memberId);
    const member = members.find(m => m.id === memberId);
    showPopup('Operator Ditambahkan', `${member.name} ditambahkan sebagai operator bantuan.`, 'success');
    populateStandbyOperatorsInTakeOrderPopup();
    hideAnimatedPopup(selectHelperOperatorModal);
  }

  function populateStandbyOperatorsInTakeOrderPopup() {
    const listDiv = document.getElementById('standbyOperatorsList');
    listDiv.innerHTML = '';

    if (additionalOperators.length === 0) {
      listDiv.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada operator bantuan yang dipilih.</p>';
      return;
    }

    additionalOperators.forEach(memberId => {
      const member = members.find(m => m.id === memberId);
      if (!member) return;
      const div = document.createElement('div');
      div.className = 'flex items-center gap-3 p-2 bg-blue-50 rounded-lg shadow-sm';
      div.innerHTML = `
        <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full">
        <div class="flex-1">
          <div class="font-medium">${member.name}</div>
          <div class="text-xs text-blue-700">Operator Bantuan</div>
        </div>
        <button class="remove-helper-operator-btn text-red-500 hover:text-red-700 transition-colors"
          data-member-id="${member.id}" title="Hapus operator bantuan">
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>`;
      listDiv.appendChild(div);
    });

    document.querySelectorAll('.remove-helper-operator-btn').forEach(b => {
      b.addEventListener('click', function () { removeHelperOperator(parseInt(this.dataset.memberId)); });
    });
  }

  function removeHelperOperator(memberId) {
    additionalOperators = additionalOperators.filter(id => id !== memberId);
    const member = members.find(m => m.id === memberId);
    showPopup('Operator Dihapus', `${member.name} dihapus dari operator bantuan.`, 'info');
    populateStandbyOperatorsInTakeOrderPopup();
  }

  // ===== CREATE ORDER POPUP =====
  createOrderBtn.addEventListener('click', () => showAnimatedPopup(createOrderPopup));

  closeCreateOrderPopupBtn.addEventListener('click', () => {
    if (isGuestUser) { showPopup('Guest Restriction', 'You must create at least one work order before closing.', 'warning'); return; }
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    specificLocationContainer.classList.add('hidden');
  });

  cancelCreateOrderBtn.addEventListener('click', () => {
    if (isGuestUser) { showPopup('Guest Restriction', 'You must create at least one work order before closing.', 'warning'); return; }
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    specificLocationContainer.classList.add('hidden');
  });

  if (exitGuestBtn) {
    exitGuestBtn.addEventListener('click', function () {
      localStorage.removeItem('isGuestUser');
      localStorage.removeItem('guestLoginTime');
      hideAnimatedPopup(createOrderPopup);
      createOrderForm.reset();
      specificLocationContainer.classList.add('hidden');
      exitGuestBtn.classList.add('hidden');
      showPopup('Guest Mode Ended', 'You have exited guest mode. You can now view the dashboard.', 'success');
    });
  }

  // ===== LOCATION DROPDOWN =====
  orderLocationSelect.addEventListener('change', function () {
    const locationPrompts = {
      'Gedung A': 'Contoh: Lantai 2, Ruang Kelas', 'Gedung B': 'Contoh: Lantai 1, Lorong Kelas',
      'Gedung B Baru': 'Contoh: Lantai 3, Lorong Kelas', 'Gedung C': 'Contoh: Lantai 1, Lorong Kelas',
      'Gedung D': 'Contoh: Lantai 2, Ruang PPDB', 'Gedung E': 'Contoh: Lantai 1, Bengkel',
      'Gedung F': 'Contoh: Lantai 1, Ruang Kelas', 'Gedung G': 'Contoh: Pintu Masuk Workshop',
      'Gedung TKI': 'Contoh: Lantai 1', 'Ruang Guru': 'Contoh: Ruang Horenso',
      'Ruang Yayasan': 'Contoh: Ruang Ketua Yayasan', 'default': 'Contoh: Nomor ruang, lantai, atau area spesifik'
    };
    if (this.value) {
      specificLocationContainer.classList.remove('hidden');
      specificLocationInput.placeholder = locationPrompts[this.value] || locationPrompts['default'];
    } else {
      specificLocationContainer.classList.add('hidden');
    }
  });

  // ===== CREATE ORDER SUBMIT =====
  createOrderForm.addEventListener('submit', function (e) {
    e.preventDefault();
    const submitBtn = this.querySelector('[type="submit"]') || createOrderBtn;
    if (submitBtn.disabled) return;
    setButtonLoading(submitBtn, true);

    const now = new Date();
    const hh  = String(now.getHours()).padStart(2, '0');
    const mm  = String(now.getMinutes()).padStart(2, '0');
    const loc = document.getElementById('orderLocation').value;
    const spc = document.getElementById('specificLocation').value;

    const payload = {
      priority:      document.getElementById('orderPriority').value,
      time_display:  `${hh}:${mm}`,
      time_sort:     `${hh}:${mm}:00`,
      requester:     document.getElementById('orderRequester').value,
      location:      spc ? `${loc} - ${spc}` : loc,
      device:        document.getElementById('orderDevice').value,
      problem:       document.getElementById('orderProblem').value,
      working_hours: '0 menit',
      status:        'pending',
      executors:     [],
      safety_checklist: []
    };

    // FIX: Kirim Authorization header untuk endpoint yang dilindungi
    fetch('/api/workorders', {
      method:  'POST',
      headers: authHeaders(),
      body:    JSON.stringify(payload)
    })
    .then(r => {
      if (!r.ok) throw new Error('Gagal menyimpan order. Status: ' + r.status);
      return r.json();
    })
    .then(json => {
      // FIX: unwrap format response baru { code, message, data: { id } }
      const data = unwrapData(json);
      hideAnimatedPopup(createOrderPopup);
      createOrderForm.reset();
      specificLocationContainer.classList.add('hidden');
      refreshAllDataFromAPI();
      showPopup('Work Order Berhasil Dibuat!', `Work Order #${data.id} telah berhasil dibuat dan disimpan.`, 'success');
    })
    .catch(err => {
      console.error('Error saat membuat order:', err);
      showPopup('Error', 'Terjadi kesalahan saat menghubungi server.', 'error');
    })
    .finally(() => setButtonLoading(submitBtn, false));
  });

  // ===== SEARCH =====
  memberSearchInput.addEventListener('focus', function () {
    searchDropdown.classList.remove('hidden');
    populateSearchResults();
  });

  memberSearchInput.addEventListener('input', function () {
    populateSearchResults(this.value.toLowerCase());
  });

  document.addEventListener('click', function (e) {
    if (!e.target.closest('.relative')) searchDropdown.classList.add('hidden');
  });

  function populateSearchResults(searchTerm = '') {
    searchResults.innerHTML = '';
    const filtered = members.filter(m => m.name.toLowerCase().includes(searchTerm));

    if (filtered.length === 0) {
      searchResults.innerHTML = '<div class="text-center py-4 text-gray-500 text-sm">Member tidak ditemukan</div>';
      return;
    }

    const statusMap = {
      standby:   { color: 'bg-green-500',  text: 'Stand By'  },
      onjob:     { color: 'bg-blue-500',   text: 'On Job'    },
      support:   { color: 'bg-yellow-400', text: 'Support'   },
      nextshift: { color: 'bg-purple-500', text: 'Next Shift' },
      offduty:   { color: 'bg-gray-500',   text: 'Off Duty'  }
    };

    filtered.forEach(member => {
      const s    = statusMap[member.status] || { color: 'bg-gray-500', text: 'Unknown' };
      const item = document.createElement('div');
      item.className = 'flex items-center gap-3 px-4 py-2.5 hover:bg-blue-50 cursor-pointer transition-colors';
      item.innerHTML = `
        <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full object-cover flex-shrink-0">
        <div class="flex-1 min-w-0">
          <div class="font-medium text-gray-800 text-sm truncate">${member.name}</div>
          <div class="flex items-center gap-1.5 text-xs text-gray-500">
            <span class="w-2 h-2 rounded-full flex-shrink-0 ${s.color}"></span>
            <span>${s.text}</span>
            ${member.role ? `<span class="text-gray-300">·</span><span class="text-gray-400">${member.role}</span>` : ''}
          </div>
        </div>
        <svg class="w-4 h-4 text-gray-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
        </svg>`;
      item.addEventListener('click', () => {
        memberSearchInput.value = '';
        searchDropdown.classList.add('hidden');
        showOperatorProfile(member);
      });
      searchResults.appendChild(item);
    });
  }

  function showOperatorProfile(member) {
    const statusMap = {
      standby:   { text: 'Stand By',  bg: 'bg-green-100',  color: 'text-green-700',  dot: '#22c55e' },
      onjob:     { text: 'On Job',    bg: 'bg-blue-100',   color: 'text-blue-700',   dot: '#3b82f6' },
      support:   { text: 'Support',   bg: 'bg-yellow-100', color: 'text-yellow-700', dot: '#eab308' },
      nextshift: { text: 'Next Shift',bg: 'bg-purple-100', color: 'text-purple-700', dot: '#a855f7' },
      offduty:   { text: 'Off Duty',  bg: 'bg-gray-100',   color: 'text-gray-600',   dot: '#9ca3af' },
    };
    const headerColors = {
      standby:   'linear-gradient(135deg,#14532d,#166534)',
      onjob:     'linear-gradient(135deg,#1e3a8a,#1d4ed8)',
      support:   'linear-gradient(135deg,#713f12,#92400e)',
      nextshift: 'linear-gradient(135deg,#4c1d95,#6d28d9)',
      offduty:   'linear-gradient(135deg,#1e293b,#374151)',
    };

    const s = statusMap[member.status] || statusMap.offduty;
    const hc = headerColors[member.status] || headerColors.offduty;

    document.getElementById('opModalHeader').style.background = hc;
    document.getElementById('opModalAvatar').src = '/static/public/' + (member.avatar || 'default-avatar.png');
    document.getElementById('opModalAvatar').alt = member.name;
    document.getElementById('opModalName').textContent = member.name;
    document.getElementById('opModalRole').textContent = member.role || '—';
    document.getElementById('opModalRoleDetail').textContent = member.role || '—';

    const statusEl = document.getElementById('opModalStatus');
    statusEl.textContent = s.text;
    statusEl.className = `text-sm font-semibold px-3 py-1 rounded-full ${s.bg} ${s.color}`;

    document.getElementById('operatorProfileModal').classList.remove('hidden');
  }

  // ===== MEMBER IMAGES =====
  function initializeMemberImages() {
    statusContainers.forEach(c => {
      const container = c.querySelector('.member-images');
      if (container) container.innerHTML = '';
    });

    members.forEach(member => {
      const statusContainer = document.getElementById(`status-${member.status}`);
      if (statusContainer) {
        const container = statusContainer.querySelector('.member-images');
        const img = document.createElement('img');
        img.src            = `/static/public/${member.avatar}`;
        img.alt            = member.name;
        img.className      = 'w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm';
        img.dataset.memberId = member.id;
        container.appendChild(img);
      }
    });

    statusContainers.forEach(c => updateMemberDisplay(c));
  }

  async function openMemberStatusPopup(statusFilter = 'all') {
    if (members.length === 0) await fetchMembers();
    showAnimatedPopup(memberStatusPopup);
    populateMemberList(statusFilter);
  }

  function populateMemberList(statusFilter = 'all') {
    memberList.innerHTML = '';
    const filtered = statusFilter === 'all' ? members : members.filter(m => m.status === statusFilter);

    if (filtered.length === 0) {
      memberList.innerHTML = '<div class="text-center py-4 text-gray-500">No members found for this status</div>';
      return;
    }

    filtered.forEach(member => {
      const item = document.createElement('div');
      item.className = 'flex items-center justify-between p-4 bg-gray-50 rounded-lg';
      item.innerHTML = `
        <div class="flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-12 h-12 rounded-full">
          <span class="font-medium">${member.name}</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-600">Status:</span>
          <select class="status-select px-3 py-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" data-member-id="${member.id}">
            <option value="standby"   ${member.status === 'standby'   ? 'selected' : ''}>Stand By</option>
            <option value="onjob"     ${member.status === 'onjob'     ? 'selected' : ''}>On Job</option>
            <option value="support"   ${member.status === 'support'   ? 'selected' : ''}>Support</option>
            <option value="nextshift" ${member.status === 'nextshift' ? 'selected' : ''}>Next Shift</option>
            <option value="offduty"   ${member.status === 'offduty'   ? 'selected' : ''}>Off Duty</option>
          </select>
        </div>`;
      memberList.appendChild(item);
    });

    document.querySelectorAll('.status-select').forEach(select => {
      select.addEventListener('change', function () {
        updateMemberStatus(parseInt(this.dataset.memberId), this.value, this);
      });
    });
  }

  function updateMemberStatus(memberId, newStatus, selectEl) {
    const idx = members.findIndex(m => m.id === memberId);
    if (idx === -1) return;

    const oldStatus = members[idx].status;
    if (selectEl) setButtonLoading(selectEl, true);

    apiUpdateMemberStatus(memberId, newStatus)
      .then(() => {
        members[idx].status = newStatus;
        if (oldStatus === 'onjob' && newStatus !== 'onjob') syncRemoveMemberFromOrders(memberId);
        updateStatusUI(members[idx], oldStatus, newStatus);
        if (!memberStatusPopup.classList.contains('hidden') && currentStatusFilter !== 'all') {
          populateMemberList(currentStatusFilter);
        }
        populateWorkOrdersTable();
        updateSummaryCounts();
        showPopup('Status Diperbarui', `Status ${members[idx].name} berhasil diubah.`, 'success');
      })
      .catch(err => {
        console.error('Error updating member status:', err);
        if (selectEl) selectEl.value = oldStatus;
        showPopup('Error', 'Gagal memperbarui status member ke server.', 'error');
      })
      .finally(() => { if (selectEl) setButtonLoading(selectEl, false); });
  }

  // FIX: Kirim Authorization header ke semua PATCH workorder
  function syncRemoveMemberFromOrders(memberId) {
    workOrders
      .filter(o => o.executors && o.executors.includes(memberId))
      .forEach(order => {
        const newExecutors = order.executors.filter(id => id !== memberId);
        const newStatus    = newExecutors.length === 0 && order.status === 'progress' ? 'pending' : order.status;
        fetch(`/api/workorders/${order.id}`, {
          method:  'PATCH',
          headers: authHeaders(),
          body:    JSON.stringify({ executors: newExecutors, status: newStatus })
        })
        .then(r => { if (!r.ok) throw new Error('Gagal sinkronisasi executor order #' + order.id); })
        .catch(err => console.error(err));
      });
  }

  function updateStatusUI(member, oldStatus, newStatus) {
    const oldC = document.getElementById(`status-${oldStatus}`);
    if (oldC) {
      const img = oldC.querySelector(`img[data-member-id="${member.id}"]`);
      if (img) { img.remove(); updateMemberDisplay(oldC); }
    }
    const newC = document.getElementById(`status-${newStatus}`);
    if (newC) {
      const container = newC.querySelector('.member-images');
      const img = document.createElement('img');
      img.src              = `/static/public/${member.avatar}`;
      img.alt              = member.name;
      img.className        = 'w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm';
      img.dataset.memberId = member.id;
      container.appendChild(img);
      updateMemberDisplay(newC);
    }
  }

  function updateMemberDisplay(container) {
    const imgs   = container.querySelectorAll('.member-images img');
    const moreEl = container.querySelector('.more-members');
    if (imgs.length > 3) {
      for (let i = 3; i < imgs.length; i++) imgs[i].style.display = 'none';
      moreEl.classList.remove('hidden');
      moreEl.textContent = `+${imgs.length - 3}`;
    } else {
      imgs.forEach(img => img.style.display = 'block');
      moreEl.classList.add('hidden');
    }
  }

  // ===== WORK ORDERS TABLE =====
  function populateWorkOrdersTable() {
    workOrdersTableBody.innerHTML = '';

    const sorted = [...workOrders].sort((a, b) => {
      return ({ high: 0, medium: 1, low: 2 }[a.priority] ?? 3) - ({ high: 0, medium: 1, low: 2 }[b.priority] ?? 3);
    });

    sorted.forEach(order => {
      const row = document.createElement('tr');
      if (order.priority === 'high' && order.status !== 'completed') row.classList.add('high-priority');

      const priorityBadge = `<span class="priority-badge priority-${order.priority}">${
        { high: 'High Priority', medium: 'Medium', low: 'Low' }[order.priority] || order.priority
      }</span>`;

      const statusBadge = `<span class="status-badge status-${order.status}">${
        { pending: 'Pending', progress: 'On Progress', completed: 'Completed' }[order.status] || order.status
      }</span>`;

      // Requester — backend sekarang selalu string, tapi handle juga fallback number
      const requesterName = typeof order.requester === 'string' ? order.requester : (order.requester || 'Unknown');

      let executorsHtml = '<div class="flex -space-x-2">';
      (order.executors || []).forEach(execId => {
        const member = members.find(m => m.id === execId || m.id === parseInt(execId));
        if (member) executorsHtml += `<img src="/static/public/${member.avatar}" alt="${member.name}" title="${member.name}" class="member-avatar-small">`;
      });
      executorsHtml += '</div>';

      let actionButtons = '<div class="flex items-center gap-2">';
      if (order.status === 'pending') {
        if ((order.executors || []).length > 0) {
          actionButtons += `<button class="add-worker-btn bg-green-500 text-white rounded-full p-1 hover:bg-green-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Tambah worker">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/></svg>
          </button>`;
        } else {
          actionButtons += `<button class="take-order-btn bg-blue-500 text-white rounded-full p-1 hover:bg-blue-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Ambil order ini">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/></svg>
          </button>`;
        }
      } else if (order.status === 'progress') {
        actionButtons += `<button class="done-btn bg-green-500 text-white rounded-full p-1 hover:bg-green-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Tandai sebagai selesai">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
        </button>`;
      }
      actionButtons += `<button class="delete-btn bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Hapus order">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
      </button>`;
      actionButtons += '</div>';

      row.innerHTML = `
        <td class="py-3 px-2 text-sm">${order.id}</td>
        <td class="py-3 px-2 text-sm">${priorityBadge}</td>
        <td class="py-3 px-2 text-sm">${order.time || '-'}</td>
        <td class="py-3 px-2 text-sm">${requesterName}</td>
        <td class="py-3 px-2 text-sm">${order.location || '-'}</td>
        <td class="py-3 px-2 text-sm">${order.device || '-'}</td>
        <td class="py-3 px-2 text-sm">${order.problem || '-'}</td>
        <td class="py-3 px-2 text-sm">${executorsHtml}</td>
        <td class="py-3 px-2 text-sm" id="wh-${order.id}">${renderWorkingHours(order)}</td>
        <td class="py-3 px-2 text-sm">${statusBadge}</td>
        <td class="py-3 px-2 text-sm">${actionButtons}</td>`;
      workOrdersTableBody.appendChild(row);
    });

    document.querySelectorAll('.take-order-btn').forEach(b => {
      b.addEventListener('click', function () {
        if (checkGuestRestriction('Taking orders')) return;
        openTakeOrderPopup(parseInt(this.dataset.orderId));
      });
    });
    document.querySelectorAll('.add-worker-btn').forEach(b => {
      b.addEventListener('click', function () {
        if (checkGuestRestriction('Adding workers')) return;
        openAddWorkerPopup(parseInt(this.dataset.orderId));
      });
    });
    document.querySelectorAll('.delete-btn').forEach(b => {
      b.addEventListener('click', function () {
        if (checkGuestRestriction('Deleting orders')) return;
        deleteOrder(parseInt(this.dataset.orderId));
      });
    });
    document.querySelectorAll('.done-btn').forEach(b => {
      b.addEventListener('click', function () {
        if (checkGuestRestriction('Completing orders')) return;
        markOrderDone(parseInt(this.dataset.orderId));
      });
    });

    // Render mobile card view dengan data dan members yang sama
    if (typeof window.renderMobileCards === 'function') {
      window.renderMobileCards(sorted, members);
    }
  }

  // ===== TAKE ORDER =====
  function openTakeOrderPopup(orderId) {
    const order = workOrders.find(o => o.id === orderId);
    if (!order) return;

    const standbyMembers = members.filter(m => m.status === 'standby');
    if (standbyMembers.length === 0) {
      showPopup('Peringatan', 'Tidak ada operator standby yang tersedia untuk mengambil order ini.', 'warning');
      return;
    }

    currentOrder         = order;
    additionalOperators  = [];

    document.getElementById('popupOrderId').textContent  = order.id;
    document.getElementById('popupPriority').textContent = order.priority.charAt(0).toUpperCase() + order.priority.slice(1);
    document.getElementById('popupLocation').textContent = order.location;
    document.getElementById('popupDevice').textContent   = order.device;
    document.getElementById('popupProblem').textContent  = order.problem;

    populateStandbyOperatorsInTakeOrderPopup();
    populateSafetyChecklist(order.location);
    showAnimatedPopup(takeOrderPopup);
  }

  function openAddWorkerPopup(orderId) {
    const order = workOrders.find(o => o.id === orderId);
    if (!order) return;
    if (order.status !== 'pending') {
      showPopup('Error', 'Hanya bisa menambahkan worker ke order yang masih pending!', 'error');
      return;
    }

    currentOrder        = order;
    additionalOperators = [];

    document.getElementById('popupOrderId').textContent  = order.id;
    document.getElementById('popupPriority').textContent = order.priority.charAt(0).toUpperCase() + order.priority.slice(1);
    document.getElementById('popupLocation').textContent = order.location;
    document.getElementById('popupDevice').textContent   = order.device;
    document.getElementById('popupProblem').textContent  = order.problem;

    populateAvailableWorkersForAddWorker(order.executors);
    showAnimatedPopup(selectHelperOperatorModal);
  }

  function populateAvailableWorkersForAddWorker(existingExecutorIds) {
    availableStandbyOperatorsList.innerHTML = '';
    const normalized = (existingExecutorIds || []).map(id => parseInt(id));
    const available  = members.filter(m => m.status === 'standby' && !normalized.includes(m.id));

    if (available.length === 0) {
      availableStandbyOperatorsList.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada worker standby yang tersedia</p>';
      return;
    }

    available.forEach(member => {
      const div = document.createElement('div');
      div.className = 'flex items-center justify-between p-2 bg-gray-50 rounded-lg';
      div.innerHTML = `
        <div class="flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full">
          <span class="font-medium">${member.name}</span>
        </div>
        <button class="add-worker-direct-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center"
          data-member-id="${member.id}" title="Tambahkan sebagai worker">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        </button>`;
      availableStandbyOperatorsList.appendChild(div);
    });

    document.querySelectorAll('.add-worker-direct-btn').forEach(b => {
      b.addEventListener('click', function () { addWorkerToOrder(parseInt(this.dataset.memberId)); });
    });
  }

  // FIX: Kirim Authorization header
  function addWorkerToOrder(memberId) {
    if (!currentOrder) return;

    const normalized = (currentOrder.executors || []).map(id => parseInt(id));
    if (normalized.includes(memberId)) {
      showPopup('Peringatan', 'Worker ini sudah terdaftar untuk order ini!', 'warning');
      return;
    }

    const newExecutors = [...normalized, memberId];

    fetch(`/api/workorders/${currentOrder.id}`, {
      method:  'PATCH',
      headers: authHeaders(),
      body:    JSON.stringify({ executors: newExecutors })
    })
    .then(r => {
      if (!r.ok) throw new Error('Gagal menambahkan worker. Status: ' + r.status);
      return r.json();
    })
    .then(() => {
      const member = members.find(m => m.id === memberId);
      hideAnimatedPopup(selectHelperOperatorModal);
      refreshAllDataFromAPI();
      showPopup('Worker Ditambahkan', `${member.name} berhasil ditambahkan ke order #${currentOrder.id}.`, 'success');
      currentOrder = null;
    })
    .catch(err => {
      console.error('Error saat menambahkan worker:', err);
      showPopup('Error', 'Terjadi kesalahan saat menambahkan worker.', 'error');
    });
  }

  function populateSafetyChecklist(location) {
    const div = document.getElementById('safetyChecklist');
    div.innerHTML = '';
    const mainLoc = location.includes(' - ') ? location.split(' - ')[0] : location;
    const items   = safetyChecklistItems[mainLoc] || safetyChecklistItems['default'];

    if (items.length === 0) {
      div.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada checklist safety untuk lokasi ini</p>';
      return;
    }

    items.forEach(item => {
      const itemDiv  = document.createElement('div');
      itemDiv.className = 'flex items-center gap-3';
      const checkbox = document.createElement('input');
      checkbox.type  = 'checkbox';
      checkbox.className = 'custom-checkbox';
      checkbox.id    = item.id;
      checkbox.dataset.required = item.required;
      const label    = document.createElement('label');
      label.htmlFor  = item.id;
      label.className = 'flex-1 cursor-pointer';
      label.innerHTML = `${item.text} ${item.required ? '<span class="text-red-500">*</span>' : ''}`;
      itemDiv.appendChild(checkbox);
      itemDiv.appendChild(label);
      div.appendChild(itemDiv);
    });
  }

  // ===== CONFIRM TAKE ORDER =====
  // FIX: Kirim Authorization header
  function confirmTakeOrder() {
    if (!currentOrder) return;
    if (additionalOperators.length === 0) {
      showPopup('Peringatan', 'Pilih minimal satu operator untuk mengambil order ini.', 'warning');
      return;
    }

    setButtonLoading(confirmTakeOrderBtn, true);

    const safetyChecklist = [];
    document.querySelectorAll('#safetyChecklist input').forEach(cb => {
      if (cb.checked) safetyChecklist.push(cb.id);
    });

    const payload = {
      order_id:              currentOrder.id,
      executors:             additionalOperators,
      safety_checklist_items: safetyChecklist,
      status:                'progress'
    };

    fetch(`/api/workorders/${currentOrder.id}/take`, {
      method:  'POST',
      headers: authHeaders(),
      body:    JSON.stringify(payload)
    })
    .then(r => {
      if (!r.ok) throw new Error('Gagal mengambil order. Status: ' + r.status);
      return r.text().then(text => text ? JSON.parse(text) : {});
    })
    .then(() => {
      // Simpan waktu mulai pengerjaan di localStorage untuk live counter
      const timerKey = `order_timer_${currentOrder.id}`;
      localStorage.setItem(timerKey, Date.now().toString());
      showPopup('Order Berhasil Diambil!', `Berhasil mengambil order #${currentOrder.id}!`, 'success');
      hideAnimatedPopup(takeOrderPopup);
      resetTakeOrderForm();
      refreshAllDataFromAPI();
    })
    .catch(err => {
      console.error('Error saat konfirmasi ambil order:', err);
      showPopup('Error', 'Terjadi kesalahan saat menyimpan perubahan ke database.', 'error');
    })
    .finally(() => setButtonLoading(confirmTakeOrderBtn, false));
  }

  function resetTakeOrderForm() {
    currentOrder        = null;
    additionalOperators = [];
    const listDiv = document.getElementById('standbyOperatorsList');
    if (listDiv) listDiv.innerHTML = '';
    document.querySelectorAll('#safetyChecklist input').forEach(cb => cb.checked = false);
  }

  // ===== MARK ORDER DONE =====
  // FIX: Kirim Authorization header
  function markOrderDone(orderId) {
    const doneBtn = document.querySelector(`.done-btn[data-order-id="${orderId}"]`);
    setButtonLoading(doneBtn, true);

    const now = new Date();
    const hh  = String(now.getHours()).padStart(2, '0');
    const mm  = String(now.getMinutes()).padStart(2, '0');
    const completionTime = `${hh}:${mm}`;

    fetch(`/api/workorders/${orderId}/complete`, {
      method:  'PATCH',
      headers: authHeaders(),
      body:    JSON.stringify({ status: 'completed', completed_at_display: completionTime })
    })
    .then(r => {
      if (!r.ok) throw new Error('Gagal menyelesaikan order. Status: ' + r.status);
      return r.text().then(text => text ? JSON.parse(text) : {});
    })
    .then(() => {
      // Simpan waktu terakhir stopwatch sebelum dihapus, untuk ditampilkan saat completed
      const finalStartedAt = localStorage.getItem(`order_timer_${orderId}`);
      if (finalStartedAt) {
        const finalElapsed = Math.floor((Date.now() - parseInt(finalStartedAt)) / 1000);
        localStorage.setItem(`order_timer_final_${orderId}`, formatElapsed(finalElapsed));
        localStorage.removeItem(`order_timer_${orderId}`);
      }
      refreshAllDataFromAPI();
      showPopup('Order Selesai!', `Order #${orderId} berhasil ditandai selesai!\nWaktu selesai: ${completionTime}`, 'success');
    })
    .catch(err => {
      console.error('Error saat menyelesaikan order:', err);
      showPopup('Error', 'Terjadi kesalahan saat memperbarui status order.', 'error');
      setButtonLoading(doneBtn, false);
    });
  }

  // ===== DELETE ORDER =====
  // FIX: Kirim Authorization header
  function deleteOrder(orderId) {
    showConfirmationPopup(
      'Konfirmasi Hapus Order',
      `Apakah Anda yakin ingin menghapus order #${orderId}?`,
      () => {
        fetch(`/api/workorders/${orderId}`, {
          method:  'DELETE',
          headers: authHeaders()
        })
        .then(r => {
          if (!r.ok) throw new Error('Gagal menghapus order. Status: ' + r.status);
          return r.text().then(text => text ? JSON.parse(text) : {});
        })
        .then(() => {
          refreshAllDataFromAPI();
          showPopup('Order Dihapus!', `Order #${orderId} telah berhasil dihapus dari database.`, 'success');
        })
        .catch(err => {
          console.error('Error saat menghapus order:', err);
          showPopup('Error', 'Terjadi kesalahan saat menghapus order.', 'error');
        });
      }
    );
  }


  // ===== LIVE WORKING HOURS COUNTER =====

  // Format detik → HH.MM.SS (stopwatch style)
  function formatElapsed(totalSeconds) {
    const h = Math.floor(totalSeconds / 3600);
    const m = Math.floor((totalSeconds % 3600) / 60);
    const s = totalSeconds % 60;
    return `${String(h).padStart(2,'0')}.${String(m).padStart(2,'0')}.${String(s).padStart(2,'0')}`;
  }

  // Render initial value untuk kolom working hours
  function renderWorkingHours(order) {
    if (order.status === 'completed') {
      // Tampilkan waktu yang tersimpan (dari localStorage sebelum dihapus)
      const savedTime = localStorage.getItem(`order_timer_final_${order.id}`);
      if (savedTime) return savedTime;
      if (order.workingHours != null) return order.workingHours + ' mnt';
      return '-';
    }
    if (order.status === 'progress') {
      const startedAt = localStorage.getItem(`order_timer_${order.id}`);
      if (startedAt) {
        const elapsed = Math.floor((Date.now() - parseInt(startedAt)) / 1000);
        return formatElapsed(elapsed);
      }
      // Order progress tapi tidak ada timer lokal (diambil di device lain)
      return '00.00.00';
    }
    return '-';
  }

  // Tick setiap detik — update stopwatch untuk semua order yang sedang progress
  setInterval(() => {
    document.querySelectorAll('[id^="wh-"]').forEach(cell => {
      const orderId   = cell.id.slice(3); // hapus "wh-"
      const startedAt = localStorage.getItem(`order_timer_${orderId}`);
      if (!startedAt) return;
      const elapsed = Math.floor((Date.now() - parseInt(startedAt)) / 1000);
      cell.textContent = formatElapsed(elapsed);
    });
  }, 1000);

  // ===== SUMMARY COUNTS =====
  function updateSummaryCounts() {
    document.getElementById('totalOrdersCount').textContent     = workOrders.length;
    document.getElementById('pendingOrdersCount').textContent   = workOrders.filter(o => o.status === 'pending').length;
    document.getElementById('progressOrdersCount').textContent  = workOrders.filter(o => o.status === 'progress').length;
    document.getElementById('completedOrdersCount').textContent = workOrders.filter(o => o.status === 'completed').length;
  }

  // ===== OPERATOR PROFILE MODAL (dari search bar) =====
  const operatorProfileModal = document.getElementById('operatorProfileModal');
  if (operatorProfileModal) {
    document.getElementById('closeOperatorProfileModal').addEventListener('click', () => {
      operatorProfileModal.classList.add('hidden');
    });
    operatorProfileModal.addEventListener('click', (e) => {
      if (e.target === operatorProfileModal) operatorProfileModal.classList.add('hidden');
    });
  }

  // ===== PRIORITY RADIO SYNC =====
  // Sync visual priority selector → hidden <select id="orderPriority">
  document.querySelectorAll('#prioritySelector input[type="radio"]').forEach(radio => {
    radio.addEventListener('change', function() {
      const sel = document.getElementById('orderPriority');
      if (sel) sel.value = this.value;
    });
  });
  // Set default: medium checked on load
  const defaultPriorityRadio = document.querySelector('#prioritySelector input[value="medium"]');
  if (defaultPriorityRadio) defaultPriorityRadio.checked = true;

  // ===== FILTER TABS =====
  function updateFilterTabs(activeStatus) {
    statusFilterTabs.forEach(tab => {
      const isActive = tab.dataset.statusFilter === activeStatus;
      tab.classList.toggle('bg-blue-500',  isActive);
      tab.classList.toggle('text-white',   isActive);
      tab.classList.toggle('bg-gray-200',  !isActive);
      tab.classList.toggle('text-gray-700', !isActive);
    });
  }

  // Expose fungsi-fungsi yang dibutuhkan oleh renderMobileCards di index.html
  // Semua fungsi ini ada di dalam closure DOMContentLoaded sehingga tidak
  // accessible dari luar tanpa di-assign ke window secara eksplisit.
  window.openTakeOrderPopup    = openTakeOrderPopup;
  window.openAddWorkerPopup    = openAddWorkerPopup;
  window.markOrderDone         = markOrderDone;
  window.deleteOrder           = deleteOrder;
  window.checkGuestRestriction = checkGuestRestriction;

}); // END DOMContentLoaded

// ===== LOGIN / REGISTER PAGE LOGIC =====
// FIX: Logika login/register dipindah sepenuhnya ke inline script di login.html & register.html.
// Blok ini sekarang hanya menangani hal yang TIDAK bisa diletakkan di inline script:
// yaitu elemen yang mungkin ada di KEDUA halaman (tidak ada saat ini).
// Duplikasi event listener yang sebelumnya ada di sini sudah dihapus agar tidak
// konflik dengan script inline di login.html dan register.html.