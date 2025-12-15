 new Swiper(".mySwiper", {
  loop: true,
  autoplay: { delay: 2400 },
  pagination: { el: ".swiper-pagination", clickable: true },
});


const btn = document.getElementById("profileDropdownBtn");
const menu = document.getElementById("profileDropdown");

btn.addEventListener("click", () => {
  menu.classList.toggle("hidden");
});

// Klik di luar dropdown untuk menutup
document.addEventListener("click", (e) => {
  if (!btn.contains(e.target) && !menu.contains(e.target)) {
    menu.classList.add("hidden");
  }
});


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




// Custom Popup System
function showPopup(title, message, type = 'info') {
  // Remove existing popup if any
  const existingPopup = document.getElementById('customPopup');
  if (existingPopup) {
    existingPopup.remove();
  }

  // Create popup elements
  const popup = document.createElement('div');
  popup.id = 'customPopup';
  popup.className = 'fixed inset-0 bg-black bg-opacity-50 z-[100] flex items-center justify-center';
  const popupContent = document.createElement('div');
  popupContent.className = 'bg-white rounded-2xl shadow-2xl p-6 w-11/12 max-w-md transform transition-all popup-fade-in';

  // Set icon based on type
  let icon = '';
  let bgColor = '';
  let iconColor = '';

  switch (type) {
    case 'success':
      icon = `
                    <svg class="w-12 h-12 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                `;
      bgColor = 'from-green-50 to-green-100';
      iconColor = 'text-green-500';
      break;
    case 'warning':
      icon = `
                    <svg class="w-12 h-12 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                    </svg>
                `;
      bgColor = 'from-yellow-50 to-yellow-100';
      iconColor = 'text-yellow-500';
      break;
    case 'error':
      icon = `
                    <svg class="w-12 h-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                `;
      bgColor = 'from-red-50 to-red-100';
      iconColor = 'text-red-500';
      break;
    default:
      icon = `
                    <svg class="w-12 h-12 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                `;
      bgColor = 'from-blue-50 to-blue-100';
      iconColor = 'text-blue-500';
  }

  popupContent.innerHTML = `
            <div class="text-center">
                <div class="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-gradient-to-br ${bgColor} mb-4">
                    ${icon}
                </div>
                <h3 class="text-xl font-bold text-gray-900 mb-2">${title}</h3>
                <p class="text-gray-600 mb-6 leading-relaxed">${message}</p>
                <button class="px-6 py-3 bg-blue-500 text-white font-semibold rounded-lg hover:bg-blue-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-blue-200">
                    OK
                </button>
            </div>
        `;
  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  const okButton = popup.querySelector('button');
  okButton.addEventListener('click', () => {
    popupContent.classList.remove('popup-fade-in');
    popupContent.classList.add('popup-fade-out');
    setTimeout(() => popup.remove(), 300);
  });

  // Auto close after 5 seconds for non-error messages
  if (type !== 'error') {
    setTimeout(() => {
      if (popup.parentNode) {
        popupContent.classList.remove('popup-fade-in');
        popupContent.classList.add('popup-fade-out');
        setTimeout(() => popup.remove(), 300);
      }
    }, 5000);
  }
}

// Custom Confirmation Popup System
function showConfirmationPopup(title, message, onConfirm) {
  // Remove existing popup if any
  const existingPopup = document.getElementById('customConfirmationPopup');
  if (existingPopup) {
    existingPopup.remove();
  }

  // Create popup elements
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
                    <button id="confirmBtn" class="px-6 py-3 bg-green-500 text-white font-semibold rounded-lg hover:bg-green-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-green-200">
                        Ya
                    </button>
                    <button id="cancelBtn" class="px-6 py-3 bg-red-500 text-white font-semibold rounded-lg hover:bg-red-600 transform transition-all hover:scale-105 focus:outline-none focus:ring-4 focus:ring-red-200">
                        Tidak
                    </button>
                </div>
            </div>
        `;
  popup.appendChild(popupContent);
  document.body.appendChild(popup);

  document.getElementById('confirmBtn').addEventListener('click', function() {
    onConfirm();
    popupContent.classList.remove('popup-fade-in');
    popupContent.classList.add('popup-fade-out');
    setTimeout(() => popup.remove(), 300);
  });

  document.getElementById('cancelBtn').addEventListener('click', function() {
    popupContent.classList.remove('popup-fade-in');
    popupContent.classList.add('popup-fade-out');
    setTimeout(() => popup.remove(), 300);
  });
}

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
  setTimeout(() => {
    popupElement.classList.add('hidden');
  }, 300);
}


// Function to update the Quick Summary title based on the current date
function updateQuickSummaryTitle() {
  const quickSummaryTitleElement = document.getElementById('quickSummaryTitle');
  if (quickSummaryTitleElement) {
    const currentDate = new Date();
    if (currentDate.getDate() === 1) { // Check if it's the 1st of the month
      quickSummaryTitleElement.textContent = 'Quick Summary Bulan Ini';
    } else {
      quickSummaryTitleElement.textContent = 'Quick Summary'; // Default text for other days
    }
  }
}

// Member Status Management
document.addEventListener('DOMContentLoaded', async function () {
  
  let members = [];
  
  // DOM elements for DOM
  const memberStatusPopup = document.getElementById('memberStatusPopup');
  const memberList = document.getElementById('memberList');
  const closePopupBtn = document.getElementById('closePopup');
  const statusContainers = document.querySelectorAll('.status-container');
  const workOrdersTableBody = document.getElementById('workOrdersTableBody');

  // Take order popup elements
  const takeOrderPopup = document.getElementById('takeOrderPopup');
  const closeTakeOrderPopupBtn = document.getElementById('closeTakeOrderPopup');
  const cancelTakeOrderBtn = document.getElementById('cancelTakeOrderBtn');
  const confirmTakeOrderBtn = document.getElementById('confirmTakeOrderBtn');
  const openSelectHelperOperatorModalBtn = document.getElementById('openSelectHelperOperatorModalBtn');

  // Create order popup elements
  const createOrderPopup = document.getElementById('createOrderPopup');
  const closeCreateOrderPopupBtn = document.getElementById('closeCreateOrderPopup');
  const cancelCreateOrderBtn = document.getElementById('cancelCreateOrderBtn');
  const createOrderForm = document.getElementById('createOrderForm');
  const createOrderBtn = document.getElementById('createOrderBtn');
  const orderRequesterInput = document.getElementById('orderRequester');

  // New elements for specific location
  const orderLocationSelect = document.getElementById('orderLocation');
  const specificLocationContainer = document.getElementById('specificLocationContainer');
  const specificLocationInput = document.getElementById('specificLocation');

  // Search elements
  const memberSearchInput = document.getElementById('memberSearchInput');
  const searchDropdown = document.getElementById('searchDropdown');
  const searchResults = document.getElementById('searchResults');

  // Status filter tabs
  const statusFilterTabs = document.querySelectorAll('.status-filter-tab');
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

  let workOrders = [];

  // Current order being processed
  let currentOrder = null;
  let additionalOperators = [];
  let selectedOperators = [];

  // Comprehensive safety checklist items for all locations
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
    // Default safety checklist 
    'default': [
      { id: 'def1', text: 'Gunakan pelindung mata (goggles)', required: false },
      { id: 'def2', text: 'Gunakan Sarung Tangan', required: false },
      { id: 'def3', text: 'Pastikan area kerja aman', required: true },
      { id: 'def4', text: 'Matikan listrik sebelum bekerja', required: true }
    ]
  };

  // Fungsi baru untuk mengambil dan me-render semua data dari API
  async function fetchAndRenderWorkOrders() {
    try {
      const response = await fetch('/api/workorders');
      if (!response.ok) {
        throw new Error('Gagal mengambil data work order dari server');
      }
      // Perbarui variabel global workOrders dengan data dari server
      workOrders = await response.json(); 
      if (!Array.isArray(workOrders)) {
        // Jika server mengembalikan null atau bukan array, inisialisasi sebagai array kosong
        workOrders = [];
      }
    } catch (error) {
      console.error("Error fetching work orders:", error);
      showPopup('Error', 'Gagal memuat data work order dari server.', 'error');
      workOrders = []; // Pastikan workOrders adalah array kosong jika fetch gagal
    }
    // Render ulang tabel dan perbarui hitungan ringkasan
    populateWorkOrdersTable();
    updateSummaryCounts();
  }

  // Fungsi untuk me-refresh semua data dari API
  function refreshAllDataFromAPI() {
    fetchAndRenderWorkOrders();
    fetchMembers();
  }

  await fetchMembers();
  await fetchAndRenderWorkOrders(); // Panggilan awal untuk memuat data saat halaman dibuka

 // Biar aldi ga kepanggil terus wkwkw
  let currentUser = null;
  async function fetchCurrentUser() {
    try {
      const response = await fetch('/api/currentuser');
      if (!response.ok) {
        throw new Error('Gagal mengambil data user saat ini dari server');
      }
      currentUser = await response.json();
    } catch (error) {
      console.error("Error fetching current user:", error);
      showPopup('Error', 'Gagal memuat data user saat ini dari server.', 'error');
      currentUser = null; // Pastikan currentUser null jika fetch gagal
    }
  }

  // Initialize member images on page load
  initializeMemberImages();

  // Populate work orders table
  populateWorkOrdersTable();

  // Update summary counts
  updateSummaryCounts();

  // Open popup when any status button is clicked
  statusContainers.forEach(container => {
    container.addEventListener('click', function (e) {
      // Prevent opening popup if clicking on a member image
      if (!e.target.closest('.member-images') && !e.target.closest('.more-members')) {
        // Set filter to clicked status
        const status = this.dataset.status;
        currentStatusFilter = status;

        // Update filter tabs
        updateFilterTabs(status);

        openMemberStatusPopup(status);
      }
    });
  });

  // Close popup when close button is clicked
  closePopupBtn.addEventListener('click', function () {
    hideAnimatedPopup(memberStatusPopup);
  });


  // Take order popup event listeners
  closeTakeOrderPopupBtn.addEventListener('click', function () {
    hideAnimatedPopup(takeOrderPopup);
    resetTakeOrderForm();
  });

  cancelTakeOrderBtn.addEventListener('click', function () {
    hideAnimatedPopup(takeOrderPopup);
    resetTakeOrderForm();
  });


  // Confirm take order
  confirmTakeOrderBtn.addEventListener('click', function () {
    confirmTakeOrder();
  });

  // Add more operators button
  openSelectHelperOperatorModalBtn.addEventListener('click', openSelectHelperOperatorModal);

  // Helper Operator Select popup elements
  const selectHelperOperatorModal = document.getElementById('selectHelperOperatorModal');
  const closeSelectHelperOperatorModalBtn = document.getElementById('closeSelectHelperOperatorModalBtn');
  const availableStandbyOperatorsList = document.getElementById('availableStandbyOperatorsList');

  // Close helper operator select popup
  closeSelectHelperOperatorModalBtn.addEventListener('click', function () {
    hideAnimatedPopup(selectHelperOperatorModal);
  });

  // Function to open the select helper operator modal
  function openSelectHelperOperatorModal() {
    showAnimatedPopup(selectHelperOperatorModal);
    populateAvailableStandbyOperators();
  }

  // Function to populate available standby operators in the modal
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
          <button class="add-helper-operator-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center" data-member-id="${member.id}" title="Tambahkan sebagai operator bantuan">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </button>
        `;
      availableStandbyOperatorsList.appendChild(memberDiv);
    });

    // Add event listeners for the new add helper operator buttons
    document.querySelectorAll('.add-helper-operator-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const memberId = parseInt(this.dataset.memberId);
        addHelperOperator(memberId);
      });
    });
  }

  // Function to add a helper operator
  function addHelperOperator(memberId) {
    if (!additionalOperators.includes(memberId)) {
      additionalOperators.push(memberId);
      showPopup('Operator Ditambahkan', `${members.find(m => m.id === memberId).name} ditambahkan sebagai operator bantuan.`, 'success');
      // Optionally, update the display in the takeOrderPopup to show the added operator
      populateStandbyOperatorsInTakeOrderPopup(); // A new function to update the list
      hideAnimatedPopup(selectHelperOperatorModal); // Close the selection modal
    } else {
      showPopup('Peringatan', 'Operator ini sudah ditambahkan!', 'warning');
    }
  }

  // Function to populate standby operators (including selected additional operators) in the takeOrderPopup
  function populateStandbyOperatorsInTakeOrderPopup() {
    const standbyOperatorsListDiv = document.getElementById('standbyOperatorsList');
    standbyOperatorsListDiv.innerHTML = '';

    // Combine initially selected (checkbox) operators with additional operators
    const allSelectedOperatorIds = [...selectedOperators, ...additionalOperators];
    const uniqueSelectedOperatorIds = [...new Set(allSelectedOperatorIds)]; // Ensure no duplicates

    if (uniqueSelectedOperatorIds.length === 0) {
      standbyOperatorsListDiv.innerHTML = '<p class="text-gray-500 text-center py-4">Tidak ada operator bantuan yang dipilih.</p>';
      return;
    }

    uniqueSelectedOperatorIds.forEach(memberId => {
      const member = members.find(m => m.id === memberId);
      if (member) {
        const operatorDiv = document.createElement('div');
        operatorDiv.className = 'flex items-center gap-3 p-2 bg-blue-50 rounded-lg shadow-sm';
        operatorDiv.innerHTML = `
            <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full">
            <div class="flex-1">
              <div class="font-medium">${member.name}</div>
              <div class="text-xs text-blue-700">Operator Bantuan</div>
            </div>
            <button class="remove-helper-operator-btn text-red-500 hover:text-red-700 transition-colors" data-member-id="${member.id}" title="Hapus operator bantuan">
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          `;
        standbyOperatorsListDiv.appendChild(operatorDiv);
      }
    });

    // Add event listeners for remove helper operator buttons
    document.querySelectorAll('.remove-helper-operator-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const memberIdToRemove = parseInt(this.dataset.memberId);
        removeHelperOperator(memberIdToRemove);
      });
    });
  }

  // Function to remove a helper operator
  function removeHelperOperator(memberId) {
    additionalOperators = additionalOperators.filter(id => id !== memberId);

    showPopup('Operator Dihapus', `${members.find(m => m.id === memberId).name} dihapus dari operator bantuan.`, 'info');
    populateStandbyOperatorsInTakeOrderPopup(); // Refresh the list
  }

  // Create order popup event listeners
  createOrderBtn.addEventListener('click', function () {
    showAnimatedPopup(createOrderPopup);
  });

  closeCreateOrderPopupBtn.addEventListener('click', function () {
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    // Hide specific location field when closing
    specificLocationContainer.classList.add('hidden');
  });

  cancelCreateOrderBtn.addEventListener('click', function () {
    hideAnimatedPopup(createOrderPopup);
    createOrderForm.reset();
    // Hide specific location field when canceling
    specificLocationContainer.classList.add('hidden');
  });




  // Event listener for location dropdown to show/hide specific location input
  orderLocationSelect.addEventListener('change', function () {
    const selectedLocation = this.value;

    // Define specific location prompts for each main location
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
      'Ruang Guru': 'Contoh: Ruang Horenso,',
      'Ruang Yayasan': 'Contoh: Ruang Ketua Yayasan',
      'default': 'Contoh: Nomor ruang, lantai, atau area spesifik'
    };

    if (selectedLocation) {
      // Show the specific location input
      specificLocationContainer.classList.remove('hidden');

      // Set the placeholder based on the selected location
      const prompt = locationPrompts[selectedLocation] || locationPrompts['default'];
      specificLocationInput.placeholder = prompt;
    } else {
      // Hide the specific location input if no location is selected
      specificLocationContainer.classList.add('hidden');
    }
  });

  // Handle form submission
  createOrderForm.addEventListener('submit', function (e) {
    e.preventDefault();

    const priority = document.getElementById('orderPriority').value;
    const requesterName = document.getElementById('orderRequester').value;
    const location = document.getElementById('orderLocation').value;
    const specificLocation = document.getElementById('specificLocation').value;
    const device = document.getElementById('orderDevice').value;
    const problem = document.getElementById('orderProblem').value;

    // Dapatkan waktu saat ini dalam format yang dibutuhkan backend untuk display dan sort
    const now = new Date();
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    const currentTimeDisplay = `${hours}:${minutes}`; // Contoh: "23:13"
    const currentTimeSort = `${hours}:${minutes}:00`; // Contoh: "23:13:00"

    // lokasi
    const finalLocation = specificLocation ? `${location} - ${specificLocation}` : location;

    // Data send to Go lang
    const payload = {
      priority: priority,
      time_display: currentTimeDisplay, // Time display
      time_sort: currentTimeSort,      // Time format
      requester: requesterName,
      location: finalLocation,
      device: device,
      problem: problem,
      working_hours: '0 menit',
      status: 'pending',
      executors: [], // Sesuai skema, kosong saat membuat
      safety_checklist: [] // Sesuai skema, kosong saat membuat
    };

    // Mengirim data ke API backend Go Lang
    fetch('/api/workorders', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
    .then(response => {
      if (!response.ok) {
        // Tangani error dari server (misalnya 400 atau 500)
        throw new Error('Gagal menyimpan order. Status: ' + response.status);
      }
      return response.json(); // Server harus mengembalikan ID baru
    })
    .then(data => {
      // Setelah sukses
      const newOrderId = data.id; // Asumsi backend mengembalikan objek { id: X }
        
      // Logika refresh, reset form, dan notifikasi
      hideAnimatedPopup(createOrderPopup);
      createOrderForm.reset();
      specificLocationContainer.classList.add('hidden');

      // Panggil fungsi refresh untuk memuat ulang semua data dari server
      refreshAllDataFromAPI();

      showPopup('Work Order Berhasil Dibuat!', `Work Order #${newOrderId} telah berhasil dibuat dan disimpan di database.`, 'success');
    })
    .catch(error => {
      console.error('Error saat membuat order:', error);
      showPopup('Error', 'Terjadi kesalahan saat menghubungi server.', 'error');
    });
  });

  // Status filter tab click events
  statusFilterTabs.forEach(tab => {
    tab.addEventListener('click', function () {
      const status = this.dataset.statusFilter;
      currentStatusFilter = status;

      // Update the filter tabs UI
      updateFilterTabs(status);

      // Update the member list
      populateMemberList(status);
    });
  });

  // Search functionality
  memberSearchInput.addEventListener('focus', function () {
    showSearchDropdown();
    populateSearchResults();
  });

  // Search functionality
  memberSearchInput.addEventListener('focus', function () {
    showSearchDropdown();
    populateSearchResults();
  });

  memberSearchInput.addEventListener('input', function () {
    const searchTerm = this.value.toLowerCase();
    populateSearchResults(searchTerm);
  });

  // Close search dropdown when clicking outside
  document.addEventListener('click', function (e) {
    if (!e.target.closest('.relative')) {
      searchDropdown.classList.add('hidden');
    }
  });

  // Fungsi asinkron untuk mengambil data dari Go API
  async function fetchMembers() {
    try {
      // Panggil endpoint API Go Anda (sesuaikan URL jika perlu)
      const response = await fetch('http://localhost:8080/api/members');

      if (!response.ok) {
        throw new Error('Network response was not ok: ' + response.statusText);
      }

      // Konversi respons JSON menjadi array objek JS
      members = await response.json();
      console.log("Data members berhasil di-fetch dari Go API:", members);

    } catch (error) {
      console.error("Error fetching members:", error);
      // Tampilkan pesan error di UI jika perlu
      if (memberList) {
        memberList.innerHTML = '<div class="text-center py-4 text-red-500">Failed to load member data.</div>';
      }
    }
  }

  // Function to populate work orders table
  function populateWorkOrdersTable() {
    workOrdersTableBody.innerHTML = '';

    // Sort work orders by priority (high first)
    const sortedWorkOrders = [...workOrders].sort((a, b) => {
      const priorityOrder = { high: 0, medium: 1, low: 2 };
      return priorityOrder[a.priority] - priorityOrder[b.priority];
    });

    sortedWorkOrders.forEach(order => {
      const row = document.createElement('tr');

      // Add high priority class for blinking effect, only if not completed
      if (order.priority === 'high' && order.status !== 'completed') {
        row.classList.add('high-priority');
      }

      // Get priority badge HTML
      let priorityBadge = '';
      switch (order.priority) {
        case 'high':
          priorityBadge = '<span class="priority-badge priority-high">High Priority</span>';
          break;
        case 'medium':
          priorityBadge = '<span class="priority-badge priority-medium">Medium</span>';
          break;
        case 'low':
          priorityBadge = '<span class="priority-badge priority-low">Low</span>';
          break;
      }

      // Get status badge HTML
      let statusBadge = '';
      switch (order.status) {
        case 'pending':
          statusBadge = '<span class="status-badge status-pending">Pending</span>';
          break;
        case 'progress':
          statusBadge = '<span class="status-badge status-progress">On Progress</span>';
          break;
        case 'completed':
          statusBadge = '<span class="status-badge status-completed">Completed</span>';
          break;
      }

      // Get requester information - handle both ID and name cases
      let requesterName = 'Unknown';
      if (typeof order.requester === 'number') {
        const requester = requesters.find(r => r.id === order.requester);
        requesterName = requester ? requester.name : 'Unknown';
      } else if (typeof order.requester === 'string') {
        requesterName = order.requester;
      }

      // Get executors HTML - only show members with "onjob" status
      // If order is completed, show all assigned executors (regardless of current status)
      let executorsHtml = '<div class="flex -space-x-2">';
      if (order.executors && order.executors.length > 0) {
        order.executors.forEach(executorId => {
          const member = members.find(m => m.id === executorId);
          if (member) {
            executorsHtml += `<img src="/static/public/${member.avatar}" alt="${member.name}" title="${member.name}" class="member-avatar-small">`;
          }
        });
      }
      executorsHtml += '</div>';


      // Action Buttons
      let actionButtons = '<div class="flex items-center gap-2">';
      if (order.status === 'pending') {
        // Check if order already has workers assigned
        if (order.executors && order.executors.length > 0) {
          // Order has workers, show "Add Worker" button
          actionButtons += `<button class="add-worker-btn bg-green-500 text-white rounded-full p-1 hover:bg-green-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Tambah worker">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
              </button>`;
        } else {
          // No workers assigned, show "Take Order" button
          actionButtons += `<button class="take-order-btn bg-blue-500 text-white rounded-full p-1 hover:bg-blue-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Ambil order ini">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
              </button>`;
        }
      } else if (order.status === 'progress') {
        actionButtons += `<button class="done-btn bg-green-500 text-white rounded-full p-1 hover:bg-green-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Tandai sebagai selesai">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
              </button>`;
      }
      // Always allow deletion, but place it consistently
      actionButtons += `<button class="delete-btn bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-all h-7 w-7 flex items-center justify-center" data-order-id="${order.id}" title="Hapus order">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
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


    // Add event listeners for take order buttons (empty executor slots become take order buttons)
    document.querySelectorAll('.take-order-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const orderId = parseInt(this.dataset.orderId);
        openTakeOrderPopup(orderId);
      });
    });

    // Add event listeners for add worker buttons (for pending orders with existing workers)
    document.querySelectorAll('.add-worker-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const orderId = parseInt(this.dataset.orderId);
        openAddWorkerPopup(orderId);
      });
    });

    // Add event listeners for delete buttons
    document.querySelectorAll('.delete-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const orderId = parseInt(this.dataset.orderId);
        deleteOrder(orderId);
      });
    });

    // Add event listeners for done buttons
    document.querySelectorAll('.done-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const orderId = parseInt(this.dataset.orderId);
        markOrderDone(orderId);
      });
    });
  }




  function openTakeOrderPopup(orderId) {
    const order = workOrders.find(o => o.id === orderId);
    if (!order) return;

    if (!currentUser) {
      showPopup('Error', 'Tidak dapat mengambil order, data pengguna tidak ditemukan.', 'error');
      return;
    }

    // Check if current user is already assigned to this order
    if (order.executors.includes(currentUser.id)) {
      showPopup('Peringatan', 'Anda sudah terdaftar sebagai pelaksana untuk order ini!', 'warning');
      return;
    }

    currentOrder = order;
    additionalOperators = [];

    // Populate order details
    document.getElementById('popupOrderId').textContent = order.id;
    document.getElementById('popupPriority').textContent = order.priority.charAt(0).toUpperCase() + order.priority.slice(1);
    document.getElementById('popupLocation').textContent = order.location;
    document.getElementById('popupDevice').textContent = order.device;
    document.getElementById('popupProblem').textContent = order.problem;

    // Populate standby operators
    populateStandbyOperatorsInTakeOrderPopup();

    // Populate safety checklist
    populateSafetyChecklist(order.location);

    // Show popup
    showAnimatedPopup(takeOrderPopup);
  }

  // Function to open the add worker popup for pending orders with existing workers
  function openAddWorkerPopup(orderId) {
    const order = workOrders.find(o => o.id === orderId);
    if (!order) return;

    // Check if order is still pending and has workers
    if (order.status !== 'pending') {
      showPopup('Error', 'Hanya bisa menambahkan worker ke order yang masih pending!', 'error');
      return;
    }

    currentOrder = order;
    additionalOperators = [];

    // Populate order details
    document.getElementById('popupOrderId').textContent = order.id;
    document.getElementById('popupPriority').textContent = order.priority.charAt(0).toUpperCase() + order.priority.slice(1);
    document.getElementById('popupLocation').textContent = order.location;
    document.getElementById('popupDevice').textContent = order.device;
    document.getElementById('popupProblem').textContent = order.problem;

    // Populate available workers (exclude existing executors)
    populateAvailableWorkersForAddWorker(order.executors);

    // Show the helper operator modal directly
    showAnimatedPopup(selectHelperOperatorModal);
  }

  // Function to populate available workers (excluding existing executors)
  function populateAvailableWorkersForAddWorker(existingExecutorIds) {
    availableStandbyOperatorsList.innerHTML = '';

    // Filter out members who are already assigned to this order and only show standby members
    const availableMembers = members.filter(m => 
      m.status === 'standby' && !existingExecutorIds.includes(m.id)
    );

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
          <button class="add-helper-operator-btn bg-green-500 text-white rounded-full p-2 hover:bg-green-600 transition-colors h-8 w-8 flex items-center justify-center" data-member-id="${member.id}" title="Tambahkan sebagai worker">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </button>
        `;
      availableStandbyOperatorsList.appendChild(memberDiv);
    });

    // Add event listeners for the add helper operator buttons
    document.querySelectorAll('.add-helper-operator-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const memberId = parseInt(this.dataset.memberId);
        addWorkerToOrder(memberId);
      });
    });
  }


  // Function to add worker to order
  function addWorkerToOrder(memberId) {
    if (!currentOrder) return;

    // Check if member is already assigned to this order
    if (currentOrder.executors.includes(memberId)) {
      showPopup('Peringatan', 'Worker ini sudah terdaftar untuk order ini!', 'warning');
      return;
    }

    // Add the member to the order's executor list
    currentOrder.executors.push(memberId);

    // Update member status to "onjob"
    updateMemberStatus(memberId, 'onjob');

    // Update the work orders array
    const orderIndex = workOrders.findIndex(o => o.id === currentOrder.id);
    if (orderIndex !== -1) {
      workOrders[orderIndex] = currentOrder;
    }

    // Refresh the table
    populateWorkOrdersTable();

    // Update summary counts
    updateSummaryCounts();

    // Close the modal
    hideAnimatedPopup(selectHelperOperatorModal);

    // Show success message
    showPopup('Worker Ditambahkan', `${members.find(m => m.id === memberId).name} berhasil ditambahkan sebagai worker untuk order #${currentOrder.id}.`, 'success');
  }



  // Function to populate safety checklist
  function populateSafetyChecklist(location) {
    const safetyChecklistDiv = document.getElementById('safetyChecklist');
    safetyChecklistDiv.innerHTML = '';

    // Extract the main location from the location string (remove specific location details)
    let mainLocation = location;
    if (location.includes(' - ')) {
      mainLocation = location.split(' - ')[0];
    }

    // Get safety checklist items for the location, or use default if not found
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



  // Function to confirm take order (TRANSACTION API CALL)
  function confirmTakeOrder() {
      console.log("DEBUG: currentUser:", currentUser);
      console.log("DEBUG: members array:", members);
      if (!currentOrder) return;

      // 1. Ambil data yang dibutuhkan untuk dikirim ke Backend
      const allAssignedOperatorIds = [...new Set([currentUser.id, ...additionalOperators])];
      
      // ... (Logika validasi allAssignedOperatorIds dan requiredCheckboxes tetap sama) ...
      // ...
      
      // Collect safety checklist data (mengambil ID item yang diceklis)
      const safetyChecklist = [];
      document.querySelectorAll('#safetyChecklist input').forEach(checkbox => {
          if (checkbox.checked) {
              safetyChecklist.push(checkbox.id);
          }
      });

      // 2. Buat Payload untuk API
      const payload = {
          order_id: currentOrder.id,
          executors: allAssignedOperatorIds, // Array of operator IDs (int/string)
          safety_checklist_items: safetyChecklist, // Array of checklist IDs (string)
          status: 'progress' // Langsung set status di backend
      };

      // 3. Panggil API: POST /api/workorders/take
      fetch(`/api/workorders/${currentOrder.id}/take`, {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json',
          },
          body: JSON.stringify(payload),
      })
      .then(response => {
          if (!response.ok) {
              throw new Error('Gagal mengambil order. Status: ' + response.status);
          }
          // DEBUG: Log raw text response untuk melihat apa yang dikirim server
          return response.text().then(text => {
            console.log("DEBUG: Raw response from /take:", text);
            // Hanya parse jika text tidak kosong, jika kosong anggap sebagai objek kosong
            return text ? JSON.parse(text) : {};
          });
      })
      .then(data => {
          // 4. Setelah Sukses
          // Tampilkan popup sukses DULU, sebelum me-reset state
          showPopup('Order Berhasil Diambil!', `Berhasil mengambil order #${data.order_id}!`, 'success');
          
          // Baru jalankan sisanya
          refreshAllDataFromAPI(); 
          hideAnimatedPopup(takeOrderPopup);
          resetTakeOrderForm();
      })
      .catch(error => {
          console.error('Error saat konfirmasi ambil order:', error);
          showPopup('Error', 'Terjadi kesalahan saat menyimpan perubahan ke database.', 'error');
      });
  }

  // Function to reset take order form
  function resetTakeOrderForm() {
    currentOrder = null;
    additionalOperators = [];

    // Reset checkboxes in the take order popup
    const standbyOperatorsListDiv = document.getElementById('standbyOperatorsList');
    if (standbyOperatorsListDiv) {
      standbyOperatorsListDiv.innerHTML = '';
    }

    document.querySelectorAll('#safetyChecklist input').forEach(checkbox => {
      checkbox.checked = false;
    });
  }

  // Function to mark an order as done (UPDATE API CALL)
  function markOrderDone(orderId) {
      // ... (Logika pengambilan completionTime tetap sama) ...
      const now = new Date();
      const hours = String(now.getHours()).padStart(2, '0');
      const minutes = String(now.getMinutes()).padStart(2, '0');
      const completionTime = `${hours}:${minutes}`; // Waktu display

      // 1. Buat Payload
      const payload = {
          status: 'completed',
          completed_at_display: completionTime // Kirim waktu selesai ke backend
      };

      // 2. Panggil API: PATCH /api/workorders/{orderId}/complete
      fetch(`/api/workorders/${orderId}/complete`, {
          method: 'PATCH', // Atau PUT
          headers: {
              'Content-Type': 'application/json',
          },
          body: JSON.stringify(payload),
      })
      .then(response => {
          if (!response.ok) {
              throw new Error('Gagal menyelesaikan order. Status: ' + response.status);
          }
          return response.json();
      })
      .then(data => {
          // 3. Setelah Sukses
          // Tidak perlu lagi manipulasi localStorage/local array workOrders/members
          
          // Refresh seluruh data dari API
          refreshAllDataFromAPI();

          // Show success message
          showPopup('Order Selesai!', `Order #${orderId} berhasil ditandai selesai!\nWaktu selesai: ${completionTime}`, 'success');
      })
      .catch(error => {
          console.error('Error saat menyelesaikan order:', error);
          showPopup('Error', 'Terjadi kesalahan saat memperbarui status order.', 'error');
      });
  }

  // Function to delete an order (DELETE API CALL)
  function deleteOrder(orderId) {
      showConfirmationPopup(
          'Konfirmasi Hapus Order',
          `Apakah Anda yakin ingin menghapus order #${orderId}?`,
          function() {
              // 1. Panggil API: DELETE /api/workorders/{orderId}
              fetch(`/api/workorders/${orderId}`, {
                  method: 'DELETE',
                  headers: {
                      'Content-Type': 'application/json',
                  },
              })
              .then(response => {
                  if (!response.ok) {
                      throw new Error('Gagal menghapus order. Status: ' + response.status);
                  }
                  return response.json();
              })
              .then(data => {
                  // 2. Setelah Sukses
                  // Tidak perlu lagi manipulasi localStorage/local array
                  
                  // Refresh seluruh data dari API
                  refreshAllDataFromAPI();

                  // Show success message
                  showPopup('Order Dihapus!', `Order #${orderId} telah berhasil dihapus dari database.`, 'success');
              })
              .catch(error => {
                  console.error('Error saat menghapus order:', error);
                  showPopup('Error', 'Terjadi kesalahan saat menghapus order.', 'error');
              });
          }
      );
  }

  // Function to update summary counts
  function updateSummaryCounts() {
    const totalOrders = workOrders.length;
    const pendingOrders = workOrders.filter(o => o.status === 'pending').length;
    const progressOrders = workOrders.filter(o => o.status === 'progress').length;
    const completedOrders = workOrders.filter(o => o.status === 'completed').length; // Explicitly count completed

    document.getElementById('totalOrdersCount').textContent = totalOrders;
    document.getElementById('pendingOrdersCount').textContent = pendingOrders;
    document.getElementById('progressOrdersCount').textContent = progressOrders;
    document.getElementById('completedOrdersCount').textContent = completedOrders;

    // If Kaizen popup is open, refresh its evaluation
    const kPopupExisting = document.getElementById('kaizenPopup');
    if (kPopupExisting && !kPopupExisting.classList.contains('hidden')) {
      // renderKaizenEvaluation is declared later; call if available
      if (typeof renderKaizenEvaluation === 'function') renderKaizenEvaluation();
    }
  }

  // Function to update filter tabs UI
  function updateFilterTabs(activeStatus) {
    statusFilterTabs.forEach(tab => {
      if (tab.dataset.statusFilter === activeStatus) {
        tab.classList.remove('bg-gray-200', 'text-gray-700');
        tab.classList.add('bg-blue-500', 'text-white');
      } else {
        tab.classList.remove('bg-blue-500', 'text-white');
        tab.classList.add('bg-gray-200', 'text-gray-700');
      }
    });
  }

  // Function to show search dropdown
  function showSearchDropdown() {
    searchDropdown.classList.remove('hidden');
  }

  // Function to hide search dropdown
  function hideSearchDropdown() {
    searchDropdown.classList.add('hidden');
  }

  // Function to populate search results
  function populateSearchResults(searchTerm = '') {
    searchResults.innerHTML = '';

    const filteredMembers = members.filter(member =>
      member.name.toLowerCase().includes(searchTerm)
    );

    if (filteredMembers.length === 0) {
      searchResults.innerHTML = '<div class="text-center py-4 text-gray-500">No members found</div>';
      return;
    }

    filteredMembers.forEach(member => {
      const memberItem = document.createElement('div');
      memberItem.className = 'flex items-center gap-3 px-4 py-2 hover:bg-gray-100 cursor-pointer transition-colors';

      // Get status color based on member status
      let statusColor = 'bg-gray-500';
      let statusText = 'Unknown';

      switch (member.status) {
        case 'standby':
          statusColor = 'bg-green-500';
          statusText = 'Stand By';
          break;
        case 'onjob':
          statusColor = 'bg-blue-500';
          statusText = 'On Job';
          break;
        case 'support':
          statusColor = 'bg-yellow-400';
          statusText = 'Support';
          break;
        case 'nextshift':
          statusColor = 'bg-purple-500';
          statusText = 'Next Shift';
          break;
        case 'offduty':
          statusColor = 'bg-black-500';
          statusText = 'Off Duty';
          break;
      }

      memberItem.innerHTML = `
          <img src="/static/public/${member.avatar}" alt="${member.name}" class="w-10 h-10 rounded-full object-cover">
          <div class="flex-1">
            <div class="font-medium text-gray-800">${member.name}</div>
            <div class="flex items-center gap-2 text-xs text-gray-600">
              <span class="w-2 h-2 rounded-full ${statusColor}"></span>
              <span>${statusText}</span>
            </div>
          </div>
        `;

      memberItem.addEventListener('click', function () {
        selectMember(member);
      });

      searchResults.appendChild(memberItem);
    });
  }

  // Function to handle member selection
  function selectMember(member) {
    memberSearchInput.value = member.name;
    hideSearchDropdown();

    // Highlight the selected member in their status container
    highlightMember(member.id);
  }

  // Function to highlight a selected member
  function highlightMember(memberId) {
    // Remove any existing highlights
    document.querySelectorAll('.member-highlight').forEach(el => {
      el.classList.remove('ring-2', 'ring-blue-500', 'ring-offset-2');
    });

    // Add highlight to the selected member
    const memberImg = document.querySelector(`img[data-member-id="${memberId}"]`);
    if (memberImg) {
      memberImg.classList.add('member-highlight', 'ring-2', 'ring-blue-500', 'ring-offset-2');

      // Remove highlight after 3 seconds
      setTimeout(() => {
        memberImg.classList.remove('ring-2', 'ring-blue-500', 'ring-offset-2');
      }, 3000);
    }
  }

  // Function to initialize member images on page load
  function initializeMemberImages() {
    members.forEach(member => {
      const statusContainer = document.getElementById(`status-${member.status}`);
      if (statusContainer) {
        const memberImagesContainer = statusContainer.querySelector('.member-images');
        const memberImg = document.createElement('img');
        let avatar = `/static/public/${member.avatar}`
        memberImg.src = avatar;
        memberImg.alt = member.name;
        memberImg.className = 'w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm';
        memberImg.dataset.memberId = member.id;
        memberImagesContainer.appendChild(memberImg);
      }
    });

    // Update display for all status containers
    statusContainers.forEach(container => {
      updateMemberDisplay(container);
    });
  }

  // Function to open member status popup
  async function openMemberStatusPopup(statusFilter = 'all') {
    // Pastikan data sudah di-fetch sebelum mengisi list
    if (members.length === 0) {
      await fetchMembers();
    }
    showAnimatedPopup(memberStatusPopup);
    populateMemberList(statusFilter);
  }

  // Function to populate member list in popup
  function populateMemberList(statusFilter = 'all') {
    memberList.innerHTML = '';

    let filteredMembers = members;

    // Filter by status if not "all"
    if (statusFilter !== 'all') {
      filteredMembers = members.filter(member => member.status === statusFilter);
    }

    if (filteredMembers.length === 0) {
      memberList.innerHTML = '<div class="text-center py-4 text-gray-500">No members found for this status</div>';
      return;
    }

    filteredMembers.forEach(member => {
      const memberItem = document.createElement('div');
      memberItem.className = 'flex items-center justify-between p-4 bg-gray-50 rounded-lg';
      let avatar = `/static/public/${member.avatar}`

      memberItem.innerHTML = `
          <div class="flex items-center gap-3">
            <img src="${avatar}" alt="${member.name}" class="w-12 h-12 rounded-full">
            <span class="font-medium">${member.name}</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600">Status:</span>
            <select class="status-select px-3 py-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" data-member-id="${member.id}">
              <option value="standby" ${member.status === 'standby' ? 'selected' : ''}>Stand By</option>
              <option value="onjob" ${member.status === 'onjob' ? 'selected' : ''}>On Job</option>
              <option value="support" ${member.status === 'support' ? 'selected' : ''}>Support</option>
              <option value="nextshift" ${member.status === 'nextshift' ? 'selected' : ''}>Next Shift</option>
              <option value="offduty" ${member.status === 'offduty' ? 'selected' : ''}>Off Duty</option>
            </select>
          </div>
        `;

      memberList.appendChild(memberItem);
    });

    // Add event listeners to status selects
    document.querySelectorAll('.status-select').forEach(select => {
      select.addEventListener('change', function () {
        const memberId = parseInt(this.dataset.memberId);
        const newStatus = this.value;
        updateMemberStatus(memberId, newStatus);
      });
    });
  }

  // Function to update member status
  function updateMemberStatus(memberId, newStatus) {
    // Find the member in the array
    const memberIndex = members.findIndex(m => m.id === memberId);
    if (memberIndex !== -1) {
      const oldStatus = members[memberIndex].status;
      members[memberIndex].status = newStatus;

      // If changing from "onjob" to something else, remove from all work orders
      if (oldStatus === 'onjob' && newStatus !== 'onjob') {
        removeFromAllWorkOrders(memberId);
      }

      // Update the UI
      updateStatusUI(members[memberIndex], oldStatus, newStatus);

      // If the popup is open and we're filtering by status, refresh the list
      if (!memberStatusPopup.classList.contains('hidden') && currentStatusFilter !== 'all') {
        populateMemberList(currentStatusFilter);
      }

      // Refresh the work orders table to reflect changes
      populateWorkOrdersTable();
    }
  }

  // Function to remove a member from all work orders
  function removeFromAllWorkOrders(memberId) {
    workOrders.forEach(order => {
      const executorIndex = order.executors.indexOf(memberId);
      if (executorIndex !== -1) {
        order.executors.splice(executorIndex, 1);

        // If no executors left, change status to pending
        if (order.executors.length === 0 && order.status === 'progress') {
          order.status = 'pending';
        }
      }
    });

    // Update summary counts
    updateSummaryCounts();
  }

  // Function to update the UI when member status changes
  function updateStatusUI(member, oldStatus, newStatus) {
    // Remove member from old status container
    const oldStatusContainer = document.getElementById(`status-${oldStatus}`);
    if (oldStatusContainer) {
      const memberImg = oldStatusContainer.querySelector(`img[data-member-id="${member.id}"]`);
      if (memberImg) {
        memberImg.remove();
        // Update display for old status container
        updateMemberDisplay(oldStatusContainer);
      }
    }

    // Add member to new status container
    const newStatusContainer = document.getElementById(`status-${newStatus}`);
    if (newStatusContainer) {
      const memberImagesContainer = newStatusContainer.querySelector('.member-images');
      const memberImg = document.createElement('img');
      let avatar = `/static/public/${member.avatar}`
      memberImg.src = avatar;
      memberImg.alt = member.name;
      memberImg.className = 'w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm';
      memberImg.dataset.memberId = member.id;
      memberImagesContainer.appendChild(memberImg);

      // Update display for new status container
      updateMemberDisplay(newStatusContainer);
    }
  }

  // Function to update member display (show only 3 images and +N if more)
  function updateMemberDisplay(container) {
    const memberImages = container.querySelectorAll('.member-images img');
    const moreMembersDiv = container.querySelector('.more-members');

    // If there are more than 3 members, show only the first 3 and +N
    if (memberImages.length > 3) {
      // Hide all images beyond the first 3
      for (let i = 3; i < memberImages.length; i++) {
        memberImages[i].style.display = 'none';
      }

      // Show the +N indicator
      moreMembersDiv.classList.remove('hidden');
      moreMembersDiv.textContent = `+${memberImages.length - 3}`;
    } else {
      // Show all images
      memberImages.forEach(img => {
        img.style.display = 'block';
      });

      // Hide the +N indicator
      moreMembersDiv.classList.add('hidden');
    }
  }

});

document.addEventListener('DOMContentLoaded', function() {
    // Modal functionality
    const modal = document.getElementById("myModal");
    const span = document.getElementsByClassName("close")[0];
    const modalMessage = document.getElementById("modal-message");

    function showModal(message) {
        if (modal && modalMessage) {
            modalMessage.textContent = message;
            modal.style.display = "block";
        }
    }

    if(span) {
        span.onclick = function() {
            if(modal) {
                modal.style.display = "none";
            }
        }
    }

    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    }

    // Login form submission
    const loginForm = document.querySelector('.login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', function(event) {
            event.preventDefault();
            const name = event.target.username.value;
            const password = event.target.password.value;
            
            fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name, password })
            })
            .then(response => response.json())
            .then(data => {
                if (data.message) {
                    showModal('Login successful!');
                    setTimeout(() => {
                        window.location.href = 'index.html';
                    }, 1000);
                } else {
                    showModal(data.error || 'Login failed');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                showModal('An error occurred. Please try again.');
            });
        });
    }

    // Registration form submission
    const registerForm = document.querySelector('.register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', function(event) {
            event.preventDefault();
            const name = event.target.username.value;
            const password = event.target.password.value;
            
            fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name, password })
            })
            .then(response => response.json())
            .then(data => {
                if (data.message) {
                    showModal('Registration successful! Please log in.');
                    setTimeout(() => {
                        window.location.href = 'login.html';
                    }, 2000);
                } else {
                    showModal(data.error || 'Registration failed');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                showModal('An error occurred. Please try again.');
            });
        });
    }
});
