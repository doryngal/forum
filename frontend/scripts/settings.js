// CHANGE THEME
const themeToggle = document.querySelectorAll("#theme-toggle");
const themeCheck = document.querySelectorAll("#theme-check");

themesP1 = [
  "var(--black)",
  "var(--primary)",
  "var(--secondary)",
  "#ae2fc5",
  "#e58816",
  "#213bca",
];

themesP2 = [
  "var(--semi-black)",
  "#3c6408",
  "var(--paragraph)",
  "#821697",
  "#a26417",
  "var(--secondary)",
];

// CHANGE NAV & SIDE
const navLogo = document.querySelector(".nav-logo-container");

themeToggle.forEach((el, i) => {
  el.addEventListener("click", () => {
    themeCheck.forEach((els) => {
      els.innerHTML = ``;
    });
    navLogo.style.backgroundColor = themesP1[i];
    header.style.backgroundColor = themesP1[i];
    sidebar.style.backgroundColor = themesP2[i];
    themeCheck[i].innerHTML = `<i class="ri-check-line"></i>`;
    themeToggle.forEach((element) => {
      element.classList.remove("active");
    });
    el.classList.add("active");
  });
});

// DELETE ACCOUNT
