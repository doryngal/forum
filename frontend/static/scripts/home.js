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

document.addEventListener("DOMContentLoaded", () => {
  const cards = Array.from(document.querySelectorAll(".cards"));
  const tagFilters = document.querySelectorAll(".tag-filter");
  const sortFilters = document.querySelectorAll(".sort-filter");

  let selectedTags = new Set();
  let sortBy = "new";

  function renderCards() {
    let filtered = [...cards];

    if (selectedTags.size > 0) {
      filtered = filtered.filter(card => {
        const tags = card.dataset.tags.split(",");
        return Array.from(selectedTags).every(tag => tags.includes(tag));
      });
    }

    filtered.sort((a, b) => {
      const dateA = parseInt(a.dataset.date);
      const dateB = parseInt(b.dataset.date);
      const likesA = parseInt(a.dataset.likes);
      const likesB = parseInt(b.dataset.likes);
      const commentsA = parseInt(a.dataset.comments);
      const commentsB = parseInt(b.dataset.comments);

      if (sortBy === "new") return dateB - dateA;
      if (sortBy === "old") return dateA - dateB;
      if (sortBy === "hot") return likesB - likesA;
      if (sortBy === "rising") return commentsB - commentsA;
      return 0;
    });

    cards.forEach(card => card.style.display = "none");
    filtered.forEach(card => card.style.display = "block");
  }

  function updateTagStyles(tagName, isSelected) {
    document.querySelectorAll(".tag-filter").forEach(tagEl => {
      if (tagEl.textContent.trim() === tagName) {
        tagEl.classList.toggle("selected", isSelected);
      }
    });
  }

  tagFilters.forEach(tag => {
    tag.addEventListener("click", e => {
      e.preventDefault();
      const tagName = tag.textContent.trim();

      const isSelected = selectedTags.has(tagName);
      if (isSelected) {
        selectedTags.delete(tagName);
      } else {
        selectedTags.add(tagName);
      }

      updateTagStyles(tagName, !isSelected);
      renderCards();
    });
  });

  sortFilters.forEach(sort => {
    sort.addEventListener("click", e => {
      e.preventDefault();
      sortBy = sort.dataset.sort;
      renderCards();
    });
  });

  renderCards();
});

// Add JavaScript for handling likes/dislikes
document.querySelectorAll('.btn-like').forEach(button => {
  button.addEventListener('click', async function() {
    const postId = this.dataset.postId;
    try {
      const response = await fetch(`/api/posts/${postId}/like`, {
        method: 'POST',
        credentials: 'include'
      });

      if (response.ok) {
        const data = await response.json();
        this.querySelector('.like-count').textContent = data.likes;
        // You might also want to update the dislike count if user had disliked before
      }
    } catch (error) {
      console.error('Error:', error);
    }
  });
});

// Similar for dislike button
document.querySelectorAll('.btn-dislike').forEach(button => {
  button.addEventListener('click', async function() {
    const postId = this.dataset.postId;
    try {
      const response = await fetch(`/api/posts/${postId}/dislike`, {
        method: 'POST',
        credentials: 'include'
      });

      if (response.ok) {
        const data = await response.json();
        this.querySelector('.dislike-count').textContent = data.dislikes;
      }
    } catch (error) {
      console.error('Error:', error);
    }
  });
});