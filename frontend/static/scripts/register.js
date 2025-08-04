const usernameInput = document.getElementById("username");
const passwordInput = document.getElementById("password");
const repeatPasswordInput = document.getElementById("repeat-password");

const horizontalRule = document.querySelectorAll(".horizontal-rule");
// function generateProfilePicture() {
//   const profilePictureContainer = document.getElementById(
//     ".profile-picture-container"
//   );

//   profilePictureContainer.appendChild(profilePicture);
// }

usernameInput.addEventListener("change", (e) => {
  let username = e.target.value;
  let profilePicture = `<img src="./static/assets/images/user.png" alt="{{.AuthorUsername}}>
  <h5>Your unique avatar</h5>`;
  const profilePictureContainer = document.querySelector(
    ".profile-picture-container"
  );
  profilePictureContainer.innerHTML = profilePicture;
});

passwordInput.addEventListener("change", (e) => {
  console.log(e.currentTarget.value);
});

repeatPasswordInput.addEventListener("change", (e) => {
  console.log(e.currentTarget.value);
});
