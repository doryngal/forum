const toggleMode = document.querySelector("#toggle-icon");
const lightMode = `
<i class="ri-moon-line"></i>
<span class="sidebar-item-title">Dark Mode</span>`;
const darkMode = `
<i class="ri-sun-line"></i>
<span class="sidebar-item-title">Light Mode</span>`;
toggleMode.innerHTML = lightMode;

// DARK MODE:
function switchTheme() {
  if (toggleMode.classList.contains("dark")) {
    toggleMode.classList.remove("dark");
    toggleMode.innerHTML = "";
    toggleMode.innerHTML = lightMode;

    document.documentElement.removeAttribute("data-theme");

    // SAVE USER PREFERENCES
    localStorage.removeItem("toggleMode");
  } else {
    toggleMode.classList.add("dark");
    toggleMode.innerHTML = "";
    toggleMode.innerHTML = darkMode;

    document.documentElement.setAttribute("data-theme", "dark");

    // SAVE USER PREFERENCES
    localStorage.setItem("toggleMode", "dark");
  }
  smoothSidebar();
}

toggleMode.addEventListener("click", switchTheme);

// TOGGLE SIDEBAR:
const sidebar = document.getElementById("sidebar");
const header = document.getElementById("header");
const toggleSide = document.getElementById("toggleside");
const mobileSearchContainer = document.getElementById(
  "mobile-search-container"
);

const visible = document.querySelector(".visible");
const mobileClose = document.getElementById("close-sidebar");
const mobileSearch = document.getElementById("mobile-search");

const closeSidebar = document.getElementById("close-sidebar");
closeSidebar.addEventListener("click", () => {
  sidebar.classList.toggle("active");
  header.classList.toggle("sided");
  smoothSidebar();
});

toggleSide.addEventListener("click", () => {
  sidebar.classList.toggle("active");
  header.classList.toggle("sided");
  smoothSidebar();
});

mobileSearchContainer.addEventListener("click", () => {
  if (!sidebar.classList.contains("active")) {
    sidebar.classList.toggle("active");
    header.classList.toggle("sided");
    smoothSidebar();
  }
});

function smoothSidebar() {
  const sidebarListTitles = document.querySelectorAll(".sidebar-item-title");
  if (sidebar.classList.contains("active")) {
    toggleSide.innerHTML = `<i class="ri-menu-2-line"></i>`;
    if (window.innerWidth <= 380) {
      mobileClose.style.display = "flex";
    }
    mobileSearch.classList.add("sided");
    setTimeout(() => {
      visible.style.opacity = 1;
      visible.innerText = "kumi";
      mobileClose.style.opacity = 1;
    }, 300);
    sidebarListTitles.forEach((el) => {
      el.style.display = "inline";
      setTimeout(() => {
        el.style.opacity = 1;
      }, 400);
    });
  } else {
    toggleSide.innerHTML = `<i class="ri-menu-3-line"></i>`;
    visible.style.opacity = 0;
    mobileClose.style.display = "none";
    setTimeout(() => {
      mobileSearch.classList.remove("sided");
    }, 750);
    sidebarListTitles.forEach((el) => {
      el.style.opacity = 0;
      setTimeout(() => {
        el.style.display = "none";
        visible.innerText = "";
      }, 200);
    });
  }
}

// SAVE THEME PREFERENCES:
if (toggleMode) {
  let isDarkMode =
    localStorage.getItem("toggleMode") !== null &&
    localStorage.getItem("toggleMode") === "dark";

  if (isDarkMode) {
    switchTheme();
  }
}

if (document.documentElement.hasAttribute("data-theme")) {
  fullLoader.style.backgroundColor = "#000000";
  const loaderImg = document.getElementById("loader-img");
  loaderImg.style.filter = "invert()";
}

// NOTIFICAtIONS:
const notifications = document.getElementById("notifications");
notifications.style.opacity = "0";
notifications.style.display = "none";
const notificationIcon = document.getElementById("notification-icon");

notificationIcon.addEventListener("click", () => {
  if (notifications.style.display === "none") {
    notifications.style.display = "flex";
    setTimeout(() => {
      notifications.style.opacity = "1";
    }, 50);
  } else {
    notifications.style.opacity = "0";
    setTimeout(() => {
      notifications.style.display = "none";
    }, 200);
  }
});
