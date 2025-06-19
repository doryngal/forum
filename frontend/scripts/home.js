// TOGGLE SORTING:
const sortingBy = document.getElementById("sortingBy");
const sortContainer = document.getElementById("sort-container");
sortContainer.style.display = "none";
sortingBy.addEventListener("click", () => {
  if (sortContainer.style.display == "none") {
    sortContainer.style.display = "block";
  } else {
    sortContainer.style.display = "none";
  }
});

// SORTING CHOICE:
const sortingList = document.querySelectorAll("#sort-container-by");

sortingList.forEach((el, i) => {
  el.addEventListener("click", () => {
    sortingList.forEach((els) => {
      els.classList.remove("selected-li");
    });
    el.classList.add("selected-li");

    // NEW:
    switch (i) {
      // Latest
      case 0:
        sortingBy.innerHTML = `Latest Topics: <i class="ri-arrow-drop-down-line"></i>`;
        break;
      // Oldest
      case 1:
        sortingBy.innerHTML = `Oldest Topics: <i class="ri-arrow-drop-down-line"></i>`;
        break;
      // Hottest
      case 2:
        sortingBy.innerHTML = `Hottest Topics: <i class="ri-arrow-drop-down-line"></i>`;
        break;
      // Rising
      case 3:
        sortingBy.innerHTML = `Rising Topics: <i class="ri-arrow-drop-down-line"></i>`;
        break;
    }
  });
});

// LIKE/DISLIKE:
const btnLike = document.querySelectorAll("#btn-like");
const btnDislike = document.querySelectorAll("#btn-dislike");

btnLike.forEach((btn, i) => {
  btn.addEventListener("click", () => {
    if (btnDislike[i].classList.contains("btn-black")) {
      btnDislike[i].classList.remove("btn-black");
    }
    btn.classList.toggle("btn-black");
  });
});

btnDislike.forEach((btn, i) => {
  btn.addEventListener("click", () => {
    if (btnLike[i].classList.contains("btn-black")) {
      btnLike[i].classList.remove("btn-black");
    }
    btn.classList.toggle("btn-black");
  });
});
