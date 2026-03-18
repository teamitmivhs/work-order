// ===== WELCOME BANNER ANIMATION =====
// Slide up animation menggunakan CSS keyframes - NO DEPENDENCIES

// CSS animation sudah ditangani di index.html
// Tidak perlu JavaScript animation untuk welcome banner

// ===== REFRESH PAGE FUNCTIONALITY =====
const refreshBtn = document.getElementById('refreshBtn');
const refreshIcon = document.getElementById('refreshIcon');

if (refreshBtn) {
  refreshBtn.addEventListener('click', function() {
    refreshIcon.style.transform = 'rotate(0deg)';
    refreshIcon.style.transition = 'transform 0.6s linear';
    setTimeout(() => {
      refreshIcon.style.transform = 'rotate(360deg)';
    }, 10);
    setTimeout(() => {
      window.location.reload();
    }, 300);
  });
}

const btn = document.getElementById("profileDropdownBtn");
const menu = document.getElementById("profileDropdown");

btn.addEventListener("click", () => {
  menu.classList.toggle("hidden");
});

// Logout button handler
const logoutBtn = menu.querySelector('button');
if (logoutBtn) {
  logoutBtn.addEventListener('click', function() {
    localStorage.removeItem('isGuestUser');
    localStorage.removeItem('guestLoginTime');
    localStorage.removeItem('isAdmin');
    localStorage.removeItem('userToken');
    window.location.href = 'login.html';
  });
}

// Klik di luar dropdown untuk menutup
document.addEventListener("click", (e) => {
  if (!btn.contains(e.target) && !menu.contains(e.target)) {
    menu.classList.add("hidden");
  }
});

// Mobile menu functionality
const mobileMenuBtn = document.getElementById('mobileMenuBtn');
const navMenu = document.getElementById('navMenu');

if (mobileMenuBtn && navMenu) {
  mobileMenuBtn.addEventListener('click', () => {
    navMenu.classList.toggle('hidden');
    navMenu.classList.toggle('mobile-menu-active');
  });

  document.addEventListener('click', (e) => {
    if (!mobileMenuBtn.contains(e.target) && !navMenu.contains(e.target)) {
      navMenu.classList.add('hidden');
      navMenu.classList.remove('mobile-menu-active');
    }
  });
}

// GSAP Fade + Slide Animation
window.addEventListener("load", () => {
  gsap.utils.toArray(".fade-slide").forEach((el, i) => {
    gsap.to(el, {
      opacity: 1,
      y: 0,
      duration: 0.6,
      delay: i * 0.15,
      ease: "power2.out",
    });
  });
});

// ===== CUSTOM POPUP SYSTEM =====
function showPopup(title, message, type = 'info') {
  const existingPopup = document.getElementById('customPopup');
  if (existingPopup) existingPopup.remove();

  const popup = document.createElement('div');
  popup.id = 'customPopup';
  popup.className = 'fixed inset-0 bg-black bg-opacity-50 z-[100] flex items-center justify-center';

  const popupContent = document.createElement('div');
  popupContent.className = 'bg-white rounded-2xl shadow-2xl p-6 w-11/12 max-w-md transform transition-all popup-fade-in';

  let icon = '';
  let bgColor = '';

  switch (type) {
    case 'success':
      icon = `<svg class="w-12 h-12 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`;
      bgColor = 'from-green-50 to-green-100';
      break;
    case 'warning':
      icon = `<svg class="w-12 h-12 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path></svg>`;
      bgColor = 'from-yellow-50 to-yellow-100';
      break;
    case 'error':
      icon = `<svg class="w-12 h-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`;
      bgColor = 'from-red-50 to-red-100';
      break;
    default:
      icon = `<svg class="w-12 h-12 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`;
      bgColor = 'from-blue-50 to-blue-100';
  }

  popupContent.innerHTML = `
    <div class="text-center">
      <div class="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-gradient-to-br ${bgColor} mb-4">${icon}</div>
      <h3 class="text-xl font-bold text-gray-900 mb-2">${title}</h3>
      <p class="text-gray-600 mb-6 leading-relaxed">${message}</p>
      <button class="px-6 py-3 bg-blue-500 text-white font-semibold rounded-lg hover:bg-blue-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-blue-200">OK</button>
    </div>
  `;
  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  popup.querySelector('button').addEventListener('click', () => {
    popupContent.classList.replace('popup-fade-in', 'popup-fade-out');
    setTimeout(() => popup.remove(), 300);
  });

  if (type !== 'error') {
    setTimeout(() => {
      if (popup.parentNode) {
        popupContent.classList.replace('popup-fade-in', 'popup-fade-out');
        setTimeout(() => popup.remove(), 300);
      }
    }, 5000);
  }
}

// ===== CUSTOM CONFIRMATION POPUP =====
function showConfirmationPopup(title, message, onConfirm) {
  const existingPopup = document.getElementById('customConfirmationPopup');
  if (existingPopup) existingPopup.remove();

  const popup = document.createElement('div');
  popup.id = 'customConfirmationPopup';
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
        <button id="confirmBtn" class="px-6 py-3 bg-green-500 text-white font-semibold rounded-lg hover:bg-green-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-green-200">Ya</button>
        <button id="cancelBtn" class="px-6 py-3 bg-red-500 text-white font-semibold rounded-lg hover:bg-red-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-red-200">Tidak</button>
      </div>
    </div>
  `;
  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  const closePopup = () => {
    popupContent.classList.replace('popup-fade-in', 'popup-fade-out');
    setTimeout(() => popup.remove(), 300);
  };

  document.getElementById('confirmBtn').addEventListener('click', () => {
    onConfirm();
    closePopup();
  });

  document.getElementById('cancelBtn').addEventListener('click', closePopup);
}

// ===== POPUP ANIMATION HELPERS =====
function showAnimatedPopup(popupElement) {
  const popupContent = popupElement.firstElementChild;
  popupElement.classList.remove('hidden');
  popupContent.classList.remove('popup-fade-out');
  popupContent.classList.add('popup-fade-in');
}

function hideAnimatedPopup(popupElement) {
  const popupContent = popupElement.firstElementChild;
  popupContent.classList.remove('popup-fade-in');
  popupContent.classList.add('popup-fade-out');
  setTimeout(() => popupElement.classList.add('hidden'), 300);
}

// ===== QUICK SUMMARY TITLE =====
function updateQuickSummaryTitle() {
  const el = document.getElementById('quickSummaryTitle');
  if (el) {
    el.textContent = new Date().getDate() === 1 ? 'Quick Summary Bulan Ini' : 'Quick Summary';
  }
}

// ===== LOADING STATE HELPER =====
// Mencegah double-click saat request sedang berjalan
function setButtonLoading(btn, isLoading) {
  if (!btn) return;
  btn.disabled = isLoading;
  btn.style.opacity = isLoading ? '0.6' : '1';
  btn.style.cursor = isLoading ? 'not-allowed' : '';
}

// ===== MAIN DASHBOARD LOGIC =====
document.addEventListener('DOMContentLoaded', async function () {

  // Panggil updateQuickSummaryTitle — FIX: sebelumnya tidak pernah dipanggil
  updateQuickSummaryTitle();

  let members = [];
  let workOrders = [];
  let currentOrder = null;
  let additionalOperators = []; // Operator tambahan yang dipilih lewat modal
  // FIX: selectedOperators dihapus karena tidak ada mekanisme pengisian checkbox-nya;
  //      semua operator dipilih melalui additionalOperators saja.

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
  const memberStatusPopup         = document.getElementById('memberStatusPopup');
  const memberList                = document.getElementById('memberList');
  const closePopupBtn             = document.getElementById('closePopup');
  const statusContainers          = document.querySelectorAll('.status-container');
  const workOrdersTableBody       = document.getElementById('workOrdersTableBody');
  const takeOrderPopup            = document.getElementById('takeOrderPopup');
  const closeTakeOrderPopupBtn    = document.getElementById('closeTakeOrderPopup');
  const cancelTakeOrderBtn        = document.getElementById('cancelTakeOrderBtn');
  const confirmTakeOrderBtn       = document.getElementById('confirmTakeOrderBtn');
  const openSelectHelperOperatorModalBtn = document.getElementById('openSelectHelperOperatorModalBtn');
  const createOrderPopup          = document.getElementById('createOrderPopup');
  const closeCreateOrderPopupBtn  = document.getElementById('closeCreateOrderPopup');
  const cancelCreateOrderBtn      = document.getElementById('cancelCreateOrderBtn');
  const exitGuestBtn              = document.getElementById('exitGuestBtn');
  const createOrderForm           = document.getElementById('createOrderForm');
  const createOrderBtn            = document.getElementById('createOrderBtn');
  const orderLocationSelect       = document.getElementById('orderLocation');
  const specificLocationContainer = document.getElementById('specificLocationContainer');
  const specificLocationInput     = document.getElementById('specificLocation');
  const memberSearchInput         = document.getElementById('memberSearchInput');
  const searchDropdown            = document.getElementById('searchDropdown');
  const searchResults             = document.getElementById('searchResults');
  const statusFilterTabs          = document.querySelectorAll('.status-filter-tab');
  const selectHelperOperatorModal = document.getElementById('selectHelperOperatorModal');
  const closeSelectHelperOperatorModalBtn = document.getElementById('closeSelectHelperOperatorModalBtn');
  const availableStandbyOperatorsList     = document.getElementById('availableStandbyOperatorsList');

  let currentStatusFilter = 'all';

  // Sample requester data
  const requesters = [
    { id: 101, name: 'Michael Scott', department: 'Management' },
    { id: 102, name: 'Dwight Schrute', department: 'Sales' },
    { id: 103, name: 'Jim Halpert', department: 'Sales' },
    { id: 104, name: 'Pam Beesly', department: 'Reception' },
    { id: 105, name: 'Oscar Martinez', department: 'Accounting' },
    { id: 106, name: 'Angela Martin', department: 'Accounting' },
    { id: 107, name: 'Kevin Malone', department: 'Accounting' },
    { id: 108, name: 'Stanley Hudson', department: 'Sales' },
    { id: 109, name: 'Phyllis Vance', department: 'Sales' },
    { id: 110, name: 'Meredith Palmer', department: 'Supplier Relations' }
  ];

  // ===== SAFETY CHECKLIST =====
  const safetyChecklistItems = {
    'CCTV': [
      { id: 'cctv1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'cctv2', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'cctv3', text: 'Pastikan area kerja aman', required: true },
      { id: 'cctv4', text: 'Gunakan Sarung Tangan', required: true }
    ],
    'WiFi': [
      { id: 'wifi1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'wifi2', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'wifi3', text: 'Pastikan area kerja aman', required: true },
      { id: 'wifi4', text: 'Gunakan Sarung Tangan', required: true }
    ],
    'Gedung A': [
      { id: 'ga1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'ga2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'ga3', text: 'Pastikan area kerja aman', required: true },
      { id: 'ga4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'ga5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung B': [
      { id: 'gb1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'gb2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'gb3', text: 'Pastikan area kerja aman', required: true },
      { id: 'gb4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'gb5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung B Baru': [
      { id: 'gbb1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'gbb2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'gbb3', text: 'Pastikan area kerja aman', required: true },
      { id: 'gbb4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'gbb5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung C': [
      { id: 'gc1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'gc2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'gc3', text: 'Pastikan area kerja aman', required: true },
      { id: 'gc4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'gc5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung D': [
      { id: 'gd1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'gd2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'gd3', text: 'Pastikan area kerja aman', required: true },
      { id: 'gd4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'gd5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung E': [
      { id: 'ge1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'ge2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'ge3', text: 'Pastikan area kerja aman', required: true },
      { id: 'ge4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'ge5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung F': [
      { id: 'gf1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'gf2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'gf3', text: 'Pastikan area kerja aman', required: true },
      { id: 'gf4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'gf5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung G': [
      { id: 'gg1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'gg2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'gg3', text: 'Pastikan area kerja aman', required: true },
      { id: 'gg4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'gg5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Gedung TKI': [
      { id: 'gt1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'gt2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'gt3', text: 'Pastikan area kerja aman', required: true },
      { id: 'gt4', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'gt5', text: 'Gunakan sepatu safety', required: true }
    ],
    'Ruang Guru': [
      { id: 'rg1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'rg2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'rg3', text: 'Pastikan area kerja aman', required: false },
      { id: 'rg4', text: 'Matikan listrik sebelum bekerja', required: true }
    ],
    'Ruang Yayasan': [
      { id: 'ry1', text: 'Pastikan sirkulasi udara baik', required: false },
      { id: 'ry2', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'ry3', text: 'Matikan listrik sebelum bekerja', required: true },
      { id: 'ry4', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'ry5', text: 'Pastikan area kerja aman', required: true }
    ],
    'default': [
      { id: 'def1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'def2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'def3', text: 'Pastikan area kerja aman', required: true },
      { id: 'def4', text: 'Matikan listrik sebelum bekerja', required: true }
    ]
  };

  // ===== API FUNCTIONS =====

  async function fetchMembers() {
    try {
      const response = await fetch('/api/members');
      if (!response.ok) throw new Error('Network response was not ok: ' + response.statusText);
      members = await response.json();
      if (!Array.isArray(members)) members = [];
      console.log("Data members berhasil di-fetch dari Go API:", members);
    } catch (error) {
      console.error("Error fetching members:", error);
      if (memberList) memberList.innerHTML = '<div class="text-center py-4 text-red-500">Failed to load member data.</div>';
      members = [];
    }
  }

  async function fetchAndRenderWorkOrders() {
    try {
      const response = await fetch('/api/workorders');
      if (!response.ok) throw new Error('Gagal mengambil data work order dari server');
      workOrders = await response.json();
      if (!Array.isArray(workOrders)) workOrders = [];
    } catch (error) {
      console.error("Error fetching work orders:", error);
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

  // FIX: updateMemberStatus sekarang memanggil API ke backend
  async function apiUpdateMemberStatus(memberId, newStatus) {
    const response = await fetch(`/api/members/${memberId}/status`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status: newStatus })
    });
    if (!response.ok) throw new Error('Gagal memperbarui status member. Status: ' + response.status);
    return response.json();
  }

  // ===== INIT =====
  await fetchMembers();
  await fetchAndRenderWorkOrders();
  initializeMemberImages();
  updateSummaryCounts();

  if (isGuestUser && exitGuestBtn) {
    exitGuestBtn.classList.remove('hidden');
  }

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
  closeTakeOrderPopupBtn.addEventListener('click', () => {
    hideAnimatedPopup(takeOrderPopup);
    resetTakeOrderForm();
  });

  cancelTakeOrderBtn.addEventListener('click', () => {
    hideAnimatedPopup(takeOrderPopup);
    resetTakeOrderForm();
  });

  confirmTakeOrderBtn.addEventListener('click', () => confirmTakeOrder());

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
      const memberDiv = document.createElement('div');
      memberDiv.className = 'flex items-center justify-between p-2 bg-gray-50 rounded-lg';
      memberDiv.innerHTML = `
        <div class="flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full">
          <span class="font-medium">${member.name}</span>
        </div>
        <button class="add-helper-operator-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center"
          data-member-id="${member.id}" title="Tambahkan sebagai operator bantuan">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        </button>
      `;
      availableStandbyOperatorsList.appendChild(memberDiv);
    });

    document.querySelectorAll('.add-helper-operator-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        addHelperOperator(parseInt(this.dataset.memberId));
      });
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

  // FIX: Hanya pakai additionalOperators; selectedOperators dihapus
  function populateStandbyOperatorsInTakeOrderPopup() {
    const standbyOperatorsListDiv = document.getElementById('standbyOperatorsList');
    standbyOperatorsListDiv.innerHTML = '';

    if (additionalOperators.length === 0) {
      standbyOperatorsListDiv.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada operator bantuan yang dipilih.</p>';
      return;
    }

    additionalOperators.forEach(memberId => {
      const member = members.find(m => m.id === memberId);
      if (!member) return;

      const operatorDiv = document.createElement('div');
      operatorDiv.className = 'flex items-center gap-3 p-2 bg-blue-50 rounded-lg shadow-sm';
      operatorDiv.innerHTML = `
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
        </button>
      `;
      standbyOperatorsListDiv.appendChild(operatorDiv);
    });

    document.querySelectorAll('.remove-helper-operator-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        removeHelperOperator(parseInt(this.dataset.memberId));
      });
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
    if (isGuestUser) {
      showPopup('Guest Restriction', 'You must create at least one work order before closing.', 'warning');
      return;
    }
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    specificLocationContainer.classList.add('hidden');
  });

  cancelCreateOrderBtn.addEventListener('click', () => {
    if (isGuestUser) {
      showPopup('Guest Restriction', 'You must create at least one work order before closing.', 'warning');
      return;
    }
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    specificLocationContainer.classList.add('hidden');
  });

  // FIX: exitGuestBtn hanya punya SATU event listener (sebelumnya duplikat dua listener berbeda)
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
      'Gedung A': 'Contoh: Lantai 2, Ruang Kelas',
      'Gedung B': 'Contoh: Lantai 1, Lorong Kelas',
      'Gedung B Baru': 'Contoh: Lantai 3, Lorong Kelas',
      'Gedung C': 'Contoh: Lantai 1, Lorong Kelas',
      'Gedung D': 'Contoh: Lantai 2, Ruang PPDB',
      'Gedung E': 'Contoh: Lantai 1, Bengkel',
      'Gedung F': 'Contoh: Lantai 1, Ruang Kelas',
      'Gedung G': 'Contoh: Pintu Masuk Workshop',
      'Gedung TKI': 'Contoh: Lantai 1',
      'Ruang Guru': 'Contoh: Ruang Horenso',
      'Ruang Yayasan': 'Contoh: Ruang Ketua Yayasan',
      'default': 'Contoh: Nomor ruang, lantai, atau area spesifik'
    };

    if (this.value) {
      specificLocationContainer.classList.remove('hidden');
      specificLocationInput.placeholder = locationPrompts[this.value] || locationPrompts['default'];
    } else {
      specificLocationContainer.classList.add('hidden');
    }
  });

  // ===== CREATE ORDER FORM SUBMIT =====
  createOrderForm.addEventListener('submit', function (e) {
    e.preventDefault();

    const submitBtn = this.querySelector('[type="submit"]') || createOrderBtn;
    if (submitBtn.disabled) return; // FIX: Cegah double submit
    setButtonLoading(submitBtn, true);

    const priority         = document.getElementById('orderPriority').value;
    const requesterName    = document.getElementById('orderRequester').value;
    const location         = document.getElementById('orderLocation').value;
    const specificLocation = document.getElementById('specificLocation').value;
    const device           = document.getElementById('orderDevice').value;
    const problem          = document.getElementById('orderProblem').value;

    const now = new Date();
    const hh = String(now.getHours()).padStart(2, '0');
    const mm = String(now.getMinutes()).padStart(2, '0');

    const payload = {
      priority,
      time_display: `${hh}:${mm}`,
      time_sort: `${hh}:${mm}:00`,
      requester: requesterName,
      location: specificLocation ? `${location} - ${specificLocation}` : location,
      device,
      problem,
      working_hours: '0 menit',
      status: 'pending',
      executors: [],
      safety_checklist: []
    };

    fetch('/api/workorders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(response => {
      if (!response.ok) throw new Error('Gagal menyimpan order. Status: ' + response.status);
      return response.json();
    })
    .then(data => {
      hideAnimatedPopup(createOrderPopup);
      createOrderForm.reset();
      specificLocationContainer.classList.add('hidden');
      refreshAllDataFromAPI();
      showPopup('Work Order Berhasil Dibuat!', `Work Order #${data.id} telah berhasil dibuat dan disimpan di database.`, 'success');
    })
    .catch(error => {
      console.error('Error saat membuat order:', error);
      showPopup('Error', 'Terjadi kesalahan saat menghubungi server.', 'error');
    })
    .finally(() => {
      setButtonLoading(submitBtn, false); // FIX: Selalu re-enable tombol setelah selesai
    });
  });

  // ===== SEARCH =====
  // FIX: Duplikat event listener 'focus' dihapus, hanya satu yang tersisa
  memberSearchInput.addEventListener('focus', function () {
    searchDropdown.classList.remove('hidden');
    populateSearchResults();
  });

  memberSearchInput.addEventListener('input', function () {
    populateSearchResults(this.value.toLowerCase());
  });

  document.addEventListener('click', function (e) {
    if (!e.target.closest('.relative')) {
      searchDropdown.classList.add('hidden');
    }
  });

  function populateSearchResults(searchTerm = '') {
    searchResults.innerHTML = '';
    const filteredMembers = members.filter(m => m.name.toLowerCase().includes(searchTerm));

    if (filteredMembers.length === 0) {
      searchResults.innerHTML = '<div class="text-center py-4 text-gray-500">No members found</div>';
      return;
    }

    const statusMap = {
      standby:   { color: 'bg-green-500',  text: 'Stand By' },
      onjob:     { color: 'bg-blue-500',   text: 'On Job' },
      support:   { color: 'bg-yellow-400', text: 'Support' },
      nextshift: { color: 'bg-purple-500', text: 'Next Shift' },
      offduty:   { color: 'bg-gray-500',   text: 'Off Duty' }
    };

    filteredMembers.forEach(member => {
      const s = statusMap[member.status] || { color: 'bg-gray-500', text: 'Unknown' };
      const item = document.createElement('div');
      item.className = 'flex items-center gap-3 px-4 py-2 hover:bg-gray-100 cursor-pointer transition-colors';
      item.innerHTML = `
        <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full object-cover">
        <div class="flex-1">
          <div class="font-medium text-gray-800">${member.name}</div>
          <div class="flex items-center gap-2 text-xs text-gray-600">
            <span class="w-2 h-2 rounded-full ${s.color}"></span>
            <span>${s.text}</span>
          </div>
        </div>
      `;
      item.addEventListener('click', () => {
        memberSearchInput.value = member.name;
        searchDropdown.classList.add('hidden');
        highlightMember(member.id);
      });
      searchResults.appendChild(item);
    });
  }

  function highlightMember(memberId) {
    document.querySelectorAll('.member-highlight').forEach(el => {
      el.classList.remove('ring-2', 'ring-blue-500', 'ring-offset-2');
    });
    const img = document.querySelector(`img[data-member-id="${memberId}"]`);
    if (img) {
      img.classList.add('member-highlight', 'ring-2', 'ring-blue-500', 'ring-offset-2');
      setTimeout(() => img.classList.remove('ring-2', 'ring-blue-500', 'ring-offset-2'), 3000);
    }
  }

  // ===== MEMBER IMAGES =====
  function initializeMemberImages() {
    // Bersihkan dulu semua container sebelum mengisi ulang
    statusContainers.forEach(container => {
      const memberImagesContainer = container.querySelector('.member-images');
      if (memberImagesContainer) memberImagesContainer.innerHTML = '';
    });

    members.forEach(member => {
      const statusContainer = document.getElementById(`status-${member.status}`);
      if (statusContainer) {
        const memberImagesContainer = statusContainer.querySelector('.member-images');
        const img = document.createElement('img');
        img.src = `/static/public/${member.avatar}`;
        img.alt = member.name;
        img.className = 'w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm';
        img.dataset.memberId = member.id;
        memberImagesContainer.appendChild(img);
      }
    });

    statusContainers.forEach(container => updateMemberDisplay(container));
  }

  async function openMemberStatusPopup(statusFilter = 'all') {
    if (members.length === 0) await fetchMembers();
    showAnimatedPopup(memberStatusPopup);
    populateMemberList(statusFilter);
  }

  function populateMemberList(statusFilter = 'all') {
    memberList.innerHTML = '';
    const filteredMembers = statusFilter === 'all' ? members : members.filter(m => m.status === statusFilter);

    if (filteredMembers.length === 0) {
      memberList.innerHTML = '<div class="text-center py-4 text-gray-500">No members found for this status</div>';
      return;
    }

    filteredMembers.forEach(member => {
      const item = document.createElement('div');
      item.className = 'flex items-center justify-between p-4 bg-gray-50 rounded-lg';
      item.innerHTML = `
        <div class="flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-12 h-12 rounded-full">
          <span class="font-medium">${member.name}</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-600">Status:</span>
          <select class="status-select px-3 py-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            data-member-id="${member.id}">
            <option value="standby"   ${member.status === 'standby'   ? 'selected' : ''}>Stand By</option>
            <option value="onjob"     ${member.status === 'onjob'     ? 'selected' : ''}>On Job</option>
            <option value="support"   ${member.status === 'support'   ? 'selected' : ''}>Support</option>
            <option value="nextshift" ${member.status === 'nextshift' ? 'selected' : ''}>Next Shift</option>
            <option value="offduty"   ${member.status === 'offduty'   ? 'selected' : ''}>Off Duty</option>
          </select>
        </div>
      `;
      memberList.appendChild(item);
    });

    document.querySelectorAll('.status-select').forEach(select => {
      select.addEventListener('change', function () {
        updateMemberStatus(parseInt(this.dataset.memberId), this.value, this);
      });
    });
  }

  // FIX: updateMemberStatus sekarang memanggil API dan disable select selama request
  function updateMemberStatus(memberId, newStatus, selectEl) {
    const memberIndex = members.findIndex(m => m.id === memberId);
    if (memberIndex === -1) return;

    const oldStatus = members[memberIndex].status;
    if (selectEl) setButtonLoading(selectEl, true);

    apiUpdateMemberStatus(memberId, newStatus)
      .then(() => {
        members[memberIndex].status = newStatus;

        // Jika member di-off dari onjob, hapus dari semua work order via API
        if (oldStatus === 'onjob' && newStatus !== 'onjob') {
          syncRemoveMemberFromOrders(memberId);
        }

        updateStatusUI(members[memberIndex], oldStatus, newStatus);

        if (!memberStatusPopup.classList.contains('hidden') && currentStatusFilter !== 'all') {
          populateMemberList(currentStatusFilter);
        }

        populateWorkOrdersTable();
        updateSummaryCounts();
        showPopup('Status Diperbarui', `Status ${members[memberIndex].name} berhasil diubah.`, 'success');
      })
      .catch(error => {
        console.error('Error updating member status:', error);
        // FIX: Rollback select ke nilai lama jika API gagal
        if (selectEl) selectEl.value = oldStatus;
        showPopup('Error', 'Gagal memperbarui status member ke server.', 'error');
      })
      .finally(() => {
        if (selectEl) setButtonLoading(selectEl, false);
      });
  }

  // FIX: Hapus member dari order via API jika status berubah dari onjob
  function syncRemoveMemberFromOrders(memberId) {
    const affectedOrders = workOrders.filter(o => o.executors && o.executors.includes(memberId));
    affectedOrders.forEach(order => {
      const newExecutors = order.executors.filter(id => id !== memberId);
      const newStatus = newExecutors.length === 0 && order.status === 'progress' ? 'pending' : order.status;

      fetch(`/api/workorders/${order.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ executors: newExecutors, status: newStatus })
      })
      .then(res => { if (!res.ok) throw new Error('Gagal sinkronisasi executor order #' + order.id); })
      .catch(err => console.error(err));
    });
  }

  function updateStatusUI(member, oldStatus, newStatus) {
    const oldContainer = document.getElementById(`status-${oldStatus}`);
    if (oldContainer) {
      const img = oldContainer.querySelector(`img[data-member-id="${member.id}"]`);
      if (img) { img.remove(); updateMemberDisplay(oldContainer); }
    }

    const newContainer = document.getElementById(`status-${newStatus}`);
    if (newContainer) {
      const memberImagesContainer = newContainer.querySelector('.member-images');
      const img = document.createElement('img');
      img.src = `/static/public/${member.avatar}`;
      img.alt = member.name;
      img.className = 'w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm';
      img.dataset.memberId = member.id;
      memberImagesContainer.appendChild(img);
      updateMemberDisplay(newContainer);
    }
  }

  function updateMemberDisplay(container) {
    const memberImages = container.querySelectorAll('.member-images img');
    const moreMembersDiv = container.querySelector('.more-members');

    if (memberImages.length > 3) {
      for (let i = 3; i < memberImages.length; i++) memberImages[i].style.display = 'none';
      moreMembersDiv.classList.remove('hidden');
      moreMembersDiv.textContent = `+${memberImages.length - 3}`;
    } else {
      memberImages.forEach(img => img.style.display = 'block');
      moreMembersDiv.classList.add('hidden');
    }
  }

  // ===== WORK ORDERS TABLE =====
  function populateWorkOrdersTable() {
    workOrdersTableBody.innerHTML = '';

    const sortedWorkOrders = [...workOrders].sort((a, b) => {
      const priorityOrder = { high: 0, medium: 1, low: 2 };
      return priorityOrder[a.priority] - priorityOrder[b.priority];
    });

    sortedWorkOrders.forEach(order => {
      const row = document.createElement('tr');
      if (order.priority === 'high' && order.status !== 'completed') row.classList.add('high-priority');

      const priorityLabels = { high: 'priority-high', medium: 'priority-medium', low: 'priority-low' };
      const priorityText   = { high: 'High Priority', medium: 'Medium', low: 'Low' };
      const statusLabels   = { pending: 'status-pending', progress: 'status-progress', completed: 'status-completed' };
      const statusText     = { pending: 'Pending', progress: 'On Progress', completed: 'Completed' };

      const priorityBadge = `<span class="priority-badge ${priorityLabels[order.priority]}">${priorityText[order.priority]}</span>`;
      const statusBadge   = `<span class="status-badge ${statusLabels[order.status]}">${statusText[order.status]}</span>`;

      let requesterName = 'Unknown';
      if (typeof order.requester === 'number') {
        requesterName = requesters.find(r => r.id === order.requester)?.name || 'Unknown';
      } else if (typeof order.requester === 'string') {
        requesterName = order.requester;
      }

      let executorsHtml = '<div class="flex -space-x-2">';
      if (order.executors?.length > 0) {
        order.executors.forEach(executorId => {
          // FIX: normalisasi tipe ID (string vs number dari backend)
          const member = members.find(m => m.id === executorId || m.id === parseInt(executorId));
          if (member) {
            executorsHtml += `<img src="/static/public/${member.avatar}" alt="${member.name}" title="${member.name}" class="member-avatar-small">`;
          }
        });
      }
      executorsHtml += '</div>';

      let actionButtons = '<div class="flex items-center gap-2">';
      if (order.status === 'pending') {
        if (order.executors?.length > 0) {
          actionButtons += `<button class="add-worker-btn bg-green-500 text-white rounded-full p-1 hover:bg-green-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Tambah worker">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
          </button>`;
        } else {
          actionButtons += `<button class="take-order-btn bg-blue-500 text-white rounded-full p-1 hover:bg-blue-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Ambil order ini">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
          </button>`;
        }
      } else if (order.status === 'progress') {
        actionButtons += `<button class="done-btn bg-green-500 text-white rounded-full p-1 hover:bg-green-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Tandai sebagai selesai">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </button>`;
      }
      actionButtons += `<button class="delete-btn bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Hapus order">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
        </svg>
      </button>`;
      actionButtons += '</div>';

      row.innerHTML = `
        <td class="py-3 px-2 text-sm">${order.id}</td>
        <td class="py-3 px-2 text-sm">${priorityBadge}</td>
        <td class="py-3 px-2 text-sm">${order.time}</td>
        <td class="py-3 px-2 text-sm">${requesterName}</td>
        <td class="py-3 px-2 text-sm">${order.location}</td>
        <td class="py-3 px-2 text-sm">${order.device}</td>
        <td class="py-3 px-2 text-sm">${order.problem}</td>
        <td class="py-3 px-2 text-sm">${executorsHtml}</td>
        <td class="py-3 px-2 text-sm">${order.workingHours || '-'}</td>
        <td class="py-3 px-2 text-sm">${statusBadge}</td>
        <td class="py-3 px-2 text-sm">${actionButtons}</td>
      `;
      workOrdersTableBody.appendChild(row);
    });

    document.querySelectorAll('.take-order-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        if (checkGuestRestriction('Taking orders')) return;
        openTakeOrderPopup(parseInt(this.dataset.orderId));
      });
    });

    document.querySelectorAll('.add-worker-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        if (checkGuestRestriction('Adding workers')) return;
        openAddWorkerPopup(parseInt(this.dataset.orderId));
      });
    });

    document.querySelectorAll('.delete-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        if (checkGuestRestriction('Deleting orders')) return;
        deleteOrder(parseInt(this.dataset.orderId));
      });
    });

    document.querySelectorAll('.done-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        if (checkGuestRestriction('Completing orders')) return;
        markOrderDone(parseInt(this.dataset.orderId));
      });
    });
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

    currentOrder = order;
    additionalOperators = [];

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

    currentOrder = order;
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
    // FIX: normalisasi tipe ID saat filter
    const normalizedExisting = (existingExecutorIds || []).map(id => parseInt(id));
    const availableMembers = members.filter(m => m.status === 'standby' && !normalizedExisting.includes(m.id));

    if (availableMembers.length === 0) {
      availableStandbyOperatorsList.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada worker standby yang tersedia</p>';
      return;
    }

    availableMembers.forEach(member => {
      const memberDiv = document.createElement('div');
      memberDiv.className = 'flex items-center justify-between p-2 bg-gray-50 rounded-lg';
      memberDiv.innerHTML = `
        <div class="flex items-center gap-3">
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full">
          <span class="font-medium">${member.name}</span>
        </div>
        <button class="add-worker-direct-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center"
          data-member-id="${member.id}" title="Tambahkan sebagai worker">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        </button>
      `;
      availableStandbyOperatorsList.appendChild(memberDiv);
    });

    // FIX: Gunakan class berbeda (.add-worker-direct-btn) agar tidak bentrok dengan
    //      handler .add-helper-operator-btn di populateAvailableStandbyOperators
    document.querySelectorAll('.add-worker-direct-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        addWorkerToOrder(parseInt(this.dataset.memberId));
      });
    });
  }

  // FIX: addWorkerToOrder sekarang memanggil API
  function addWorkerToOrder(memberId) {
    if (!currentOrder) return;

    const normalizedExecutors = (currentOrder.executors || []).map(id => parseInt(id));
    if (normalizedExecutors.includes(memberId)) {
      showPopup('Peringatan', 'Worker ini sudah terdaftar untuk order ini!', 'warning');
      return;
    }

    const newExecutors = [...normalizedExecutors, memberId];

    fetch(`/api/workorders/${currentOrder.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ executors: newExecutors })
    })
    .then(res => {
      if (!res.ok) throw new Error('Gagal menambahkan worker. Status: ' + res.status);
      return res.json();
    })
    .then(() => {
      const member = members.find(m => m.id === memberId);
      hideAnimatedPopup(selectHelperOperatorModal);
      refreshAllDataFromAPI();
      showPopup('Worker Ditambahkan', `${member.name} berhasil ditambahkan ke order #${currentOrder.id}.`, 'success');
      currentOrder = null;
    })
    .catch(error => {
      console.error('Error saat menambahkan worker:', error);
      showPopup('Error', 'Terjadi kesalahan saat menambahkan worker.', 'error');
    });
  }

  function populateSafetyChecklist(location) {
    const safetyChecklistDiv = document.getElementById('safetyChecklist');
    safetyChecklistDiv.innerHTML = '';

    const mainLocation = location.includes(' - ') ? location.split(' - ')[0] : location;
    const items = safetyChecklistItems[mainLocation] || safetyChecklistItems['default'];

    if (items.length === 0) {
      safetyChecklistDiv.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada checklist safety untuk lokasi ini</p>';
      return;
    }

    items.forEach(item => {
      const itemDiv = document.createElement('div');
      itemDiv.className = 'flex items-center gap-3';

      const checkbox = document.createElement('input');
      checkbox.type = 'checkbox';
      checkbox.className = 'custom-checkbox';
      checkbox.id = item.id;
      checkbox.dataset.required = item.required;

      const label = document.createElement('label');
      label.htmlFor = item.id;
      label.className = 'flex-1 cursor-pointer';
      label.innerHTML = `${item.text} ${item.required ? '<span class="text-red-500">*</span>' : ''}`;

      itemDiv.appendChild(checkbox);
      itemDiv.appendChild(label);
      safetyChecklistDiv.appendChild(itemDiv);
    });
  }

  // ===== CONFIRM TAKE ORDER =====
  function confirmTakeOrder() {
    if (!currentOrder) return;

    if (additionalOperators.length === 0) {
      showPopup('Peringatan', 'Pilih minimal satu operator untuk mengambil order ini.', 'warning');
      return;
    }

    // FIX: Disable tombol saat request berlangsung untuk cegah double click
    setButtonLoading(confirmTakeOrderBtn, true);

    const safetyChecklist = [];
    document.querySelectorAll('#safetyChecklist input').forEach(checkbox => {
      if (checkbox.checked) safetyChecklist.push(checkbox.id);
    });

    const payload = {
      order_id: currentOrder.id,
      executors: additionalOperators,
      safety_checklist_items: safetyChecklist,
      status: 'progress'
    };

    fetch(`/api/workorders/${currentOrder.id}/take`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(response => {
      if (!response.ok) throw new Error('Gagal mengambil order. Status: ' + response.status);
      return response.text().then(text => {
        console.log("DEBUG: Raw response from /take:", text);
        return text ? JSON.parse(text) : {};
      });
    })
    .then(() => {
      showPopup('Order Berhasil Diambil!', `Berhasil mengambil order #${currentOrder.id}!`, 'success');
      hideAnimatedPopup(takeOrderPopup);
      resetTakeOrderForm();
      refreshAllDataFromAPI();
    })
    .catch(error => {
      console.error('Error saat konfirmasi ambil order:', error);
      showPopup('Error', 'Terjadi kesalahan saat menyimpan perubahan ke database.', 'error');
    })
    .finally(() => {
      setButtonLoading(confirmTakeOrderBtn, false); // FIX: Re-enable tombol setelah selesai
    });
  }

  function resetTakeOrderForm() {
    currentOrder = null;
    additionalOperators = [];
    const standbyOperatorsListDiv = document.getElementById('standbyOperatorsList');
    if (standbyOperatorsListDiv) standbyOperatorsListDiv.innerHTML = '';
    document.querySelectorAll('#safetyChecklist input').forEach(cb => cb.checked = false);
  }

  // ===== MARK ORDER DONE =====
  function markOrderDone(orderId) {
    const btn = document.querySelector(`.done-btn[data-order-id="${orderId}"]`);
    setButtonLoading(btn, true); // FIX: Disable saat request berlangsung

    const now = new Date();
    const hh = String(now.getHours()).padStart(2, '0');
    const mm = String(now.getMinutes()).padStart(2, '0');
    const completionTime = `${hh}:${mm}`;

    const payload = {
      status: 'completed',
      completed_at_display: completionTime
    };

    fetch(`/api/workorders/${orderId}/complete`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    // FIX: Gunakan .text() lalu parse manual, aman jika response body kosong
    .then(response => {
      if (!response.ok) throw new Error('Gagal menyelesaikan order. Status: ' + response.status);
      return response.text().then(text => text ? JSON.parse(text) : {});
    })
    .then(() => {
      refreshAllDataFromAPI();
      showPopup('Order Selesai!', `Order #${orderId} berhasil ditandai selesai!\nWaktu selesai: ${completionTime}`, 'success');
    })
    .catch(error => {
      console.error('Error saat menyelesaikan order:', error);
      showPopup('Error', 'Terjadi kesalahan saat memperbarui status order.', 'error');
      setButtonLoading(btn, false); // FIX: Re-enable jika gagal (jika sukses, tabel di-refresh otomatis)
    });
  }

  // ===== DELETE ORDER =====
  function deleteOrder(orderId) {
    showConfirmationPopup(
      'Konfirmasi Hapus Order',
      `Apakah Anda yakin ingin menghapus order #${orderId}?`,
      () => {
        fetch(`/api/workorders/${orderId}`, {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json' }
        })
        // FIX: Gunakan .text() lalu parse manual, aman jika response body kosong
        .then(response => {
          if (!response.ok) throw new Error('Gagal menghapus order. Status: ' + response.status);
          return response.text().then(text => text ? JSON.parse(text) : {});
        })
        .then(() => {
          refreshAllDataFromAPI();
          showPopup('Order Dihapus!', `Order #${orderId} telah berhasil dihapus dari database.`, 'success');
        })
        .catch(error => {
          console.error('Error saat menghapus order:', error);
          showPopup('Error', 'Terjadi kesalahan saat menghapus order.', 'error');
        });
      }
    );
  }

  // ===== SUMMARY COUNTS =====
  function updateSummaryCounts() {
    document.getElementById('totalOrdersCount').textContent     = workOrders.length;
    document.getElementById('pendingOrdersCount').textContent   = workOrders.filter(o => o.status === 'pending').length;
    document.getElementById('progressOrdersCount').textContent  = workOrders.filter(o => o.status === 'progress').length;
    document.getElementById('completedOrdersCount').textContent = workOrders.filter(o => o.status === 'completed').length;

    const kPopup = document.getElementById('kaizenPopup');
    if (kPopup && !kPopup.classList.contains('hidden') && typeof renderKaizenEvaluation === 'function') {
      renderKaizenEvaluation();
    }
  }

  // ===== FILTER TABS =====
  function updateFilterTabs(activeStatus) {
    statusFilterTabs.forEach(tab => {
      const isActive = tab.dataset.statusFilter === activeStatus;
      tab.classList.toggle('bg-blue-500', isActive);
      tab.classList.toggle('text-white', isActive);
      tab.classList.toggle('bg-gray-200', !isActive);
      tab.classList.toggle('text-gray-700', !isActive);
    });
  }

}); // END DOMContentLoaded

// ===== LOGIN PAGE LOGIC =====
document.addEventListener('DOMContentLoaded', function() {
  const modal        = document.getElementById("myModal");
  const span         = document.getElementsByClassName("close")[0];
  const modalMessage = document.getElementById("modal-message");

  function showModal(message) {
    if (modal && modalMessage) {
      modalMessage.textContent = message;
      modal.style.display = "block";
    }
  }

  if (span) {
    span.onclick = () => { if (modal) modal.style.display = "none"; };
  }

  window.onclick = (event) => {
    if (event.target === modal) modal.style.display = "none";
  };

  // Guest login
  const guestLoginBtn = document.getElementById('guestLoginBtn');
  if (guestLoginBtn) {
    guestLoginBtn.addEventListener('click', function() {
      localStorage.setItem('isGuestUser', 'true');
      localStorage.setItem('guestLoginTime', new Date().toISOString());
      showModal('Welcome Guest! You can only create work orders.');
      setTimeout(() => { window.location.href = 'index.html'; }, 1500);
    });
  }

  // Login form
  const loginForm = document.querySelector('.login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', function(event) {
      event.preventDefault();
      const submitBtn = loginForm.querySelector('[type="submit"]');
      setButtonLoading(submitBtn, true);

      const name     = event.target.username.value;
      const password = event.target.password.value;

      fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, password })
      })
      .then(res => res.json())
      .then(data => {
        if (data.message) {
          localStorage.removeItem('isGuestUser');
          localStorage.removeItem('guestLoginTime');
          localStorage.setItem('isAdmin', 'true');
          showModal('Login successful!');
          setTimeout(() => { window.location.href = 'index.html'; }, 1000);
        } else {
          showModal(data.error || 'Login failed');
          setButtonLoading(submitBtn, false);
        }
      })
      .catch(error => {
        console.error('Error:', error);
        showModal('An error occurred. Please try again.');
        setButtonLoading(submitBtn, false);
      });
    });
  }

  // Register form
  const registerForm = document.querySelector('.register-form');
  if (registerForm) {
    registerForm.addEventListener('submit', function(event) {
      event.preventDefault();
      const submitBtn = registerForm.querySelector('[type="submit"]');
      setButtonLoading(submitBtn, true);

      const name     = event.target.username.value;
      const password = event.target.password.value;

      fetch('/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, password })
      })
      .then(res => res.json())
      .then(data => {
        if (data.message) {
          showModal('Registration successful! Please log in.');
          setTimeout(() => { window.location.href = 'login.html'; }, 2000);
        } else {
          showModal(data.error || 'Registration failed');
          setButtonLoading(submitBtn, false);
        }
      })
      .catch(error => {
        console.error('Error:', error);
        showModal('An error occurred. Please try again.');
        setButtonLoading(submitBtn, false);
      });
    });
  }
});