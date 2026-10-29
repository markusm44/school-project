<script>
  import { navigate } from "svelte-routing";
  import Fa from "svelte-fa";
  import { faMoon } from "@fortawesome/free-solid-svg-icons";
  import { handleFileInput } from "./util/image";
  import { BackendHost } from "../main";
  import { postData } from "./util/postAndGet";
  export let toggleDarkMode;
  export let login;

  let userEmail = "";
  let userPassword = "";
  let userFirstName = "";
  let userLastName = "";
  let userDateOfBirth = "";
  let userAvatar = null;
  let userNickname = "";
  let userAboutMe = "";
  let dataUrl = null;
  let imageData = {
    base64String: "",
    name: "",
  };

  async function handleFileInputChange(event) {
    const result = await handleFileInput(event);
    if (result.dataUrl && result.imageInfo.name) {
      imageData = {
        base64String: result.dataUrl,
        name: result.imageInfo.name,
      };
    } else {
      imageData = {
        base64String: "",
        name: "",
      };
    }
  }

  const loggedInBool = () => {
    navigate("/login", { replace: false });
  };

  const tryRegister = () => {
    if (
      userEmail == "" ||
      userFirstName == "" ||
      userLastName == "" ||
      userPassword == "" ||
      userDateOfBirth == ""
    ) {
      alert(
        "Register failed: Email, First Name, Last Name, Password & Date of Birth are required!",
      );
      return;
    }

    if (userEmail.search("@") == -1) {
      alert("Register failed: Incorrect email entered");
      return;
    }

    let data = {
      nickname: userNickname,
      email: userEmail,
      firstName: userFirstName,
      lastName: userLastName,
      passwordHash: userPassword,
      image: imageData,
      dateOfBirth: userDateOfBirth,
      aboutMe: userAboutMe,
    };
    postData(`${BackendHost}/register`, data)
      .then(async (value) => {
        if (value.createAccount) {
          let dataForlogin = {
            email: data.email,
            password: data.passwordHash,
          };

          await login(dataForlogin, true);
          //go to main feed

          userEmail = "";
          userPassword = "";
          userFirstName = "";
          userLastName = "";
          userDateOfBirth = "";
          userAvatar = null;
          userNickname = "";
          userAboutMe = "";
          dataUrl = null;
          imageData = {
            base64String: "",
            name: "",
          };

          navigate("/feed", { replace: false });
        } else {
          alert(value.msg);
        }
      })
      .catch((e) => {
        console.log(e);
      });
  };
</script>

<div class="dark-mode-container">
  <button class="toggle-dark-mode" on:click={toggleDarkMode}>
    <Fa icon={faMoon} />
  </button>
</div>
<div class="register-container">
  <div class="input-container">
    <label class="input-label" for="email">Email</label>
    <input
      type="email"
      id="email"
      autocomplete="off"
      class="input"
      bind:value={userEmail}
    />
  </div>
  <div class="input-container">
    <label class="input-label" for="password">Password</label>
    <input
      type="password"
      id="password"
      class="input"
      bind:value={userPassword}
    />
  </div>
  <div class="input-container">
    <label class="input-label" for="firstName">First Name</label>
    <input
      type="text"
      id="firstName"
      autocomplete="off"
      class="input"
      bind:value={userFirstName}
    />
  </div>
  <div class="input-container">
    <label class="input-label" for="lastName">Last Name</label>
    <input
      type="text"
      id="lastName"
      autocomplete="off"
      class="input"
      bind:value={userLastName}
    />
  </div>
  <div class="input-container">
    <label class="input-label" for="dateOfBirth">Date Of Birth</label>
    <input
      type="date"
      id="dateOfBirth"
      class="input"
      bind:value={userDateOfBirth}
    />
  </div>
  <div class="input-container">
    <label class="input-label" for="avatarIcon"
      >Avatar Icon Url (optional)</label
    >
    <input
      type="file"
      accept="image/png, image/jpg, image/jpeg, image/gif"
      id="avatarIcon"
      class="input"
      on:change={(e) => handleFileInputChange(e)}
    />
  </div>
  <div class="input-container">
    <label class="input-label" for="nickname">Nickname (optional)</label>
    <input type="text" id="nickname" class="input" bind:value={userNickname} />
  </div>
  <div class="input-container">
    <label class="input-label" for="aboutMe">About Me (optional)</label>
    <input type="text" id="aboutMe" class="input" bind:value={userAboutMe} />
  </div>
  <div class="button-container">
    <button class="form-button" on:click={tryRegister}>Register</button>
  </div>
  <div class="form-footer-container">
    <p class="form-footer">
      Already have an account? <button class="form-link" on:click={loggedInBool}
        >Login</button
      >
    </p>
  </div>
</div>

<style>
  .register-container {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);

    display: flex;
    flex-direction: column;
    align-items: center;

    padding: 1rem;
    width: 26rem;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 0px 15px;
  }
</style>
