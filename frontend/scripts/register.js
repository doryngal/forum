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
  let profilePicture = `<img src="https://source.boringavatars.com/beam/40/${username}?square&colors=7dc81b,8ACD4F,26E4B9,0051ff,47139A">
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
