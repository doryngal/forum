const btnBack = document.getElementById("btn-back");
btnBack.addEventListener("click", () => {
  history.back();
});

document.addEventListener("DOMContentLoaded", () => {
  const buttons = document.querySelectorAll(".info-profile-box");
  const sections = document.querySelectorAll(".posts");

  buttons.forEach(button => {
    button.addEventListener("click", () => {
      // Удалить active со всех кнопок и секций
      buttons.forEach(b => b.classList.remove("active"));
      sections.forEach(s => s.classList.remove("active"));

      // Добавить active к выбранной кнопке и секции
      button.classList.add("active");
      const targetId = button.getAttribute("data-target");
      const targetSection = document.getElementById(targetId);
      if (targetSection) {
        targetSection.classList.add("active");
      }
    });
  });

  // По умолчанию показать первую секцию
  const firstSection = document.getElementById("posts");
  if (firstSection) {
    firstSection.classList.add("active");
  }
});