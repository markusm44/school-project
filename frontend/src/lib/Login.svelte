<script>
  import { navigate } from "svelte-routing";
  import Fa from "svelte-fa";
  import { faMoon } from "@fortawesome/free-solid-svg-icons";
  import { BackendHost } from "../main";
  import { getData } from "./util/postAndGet";
  export let toggleDarkMode;
  export let login;
  export let cookieLogin;

  logInBasedOnCookie();

  function logInBasedOnCookie() {
    getData(`${BackendHost}/cookieLogin`).then((data) => {
      if (data.bool) {
        cookieLogin(data.resp);
      } else {
        return;
      }
    });
  }

  let userEmail = "";
  let userPassword = "";

  const loggedInBool = () => {
    navigate("/register", { replace: false });
  };

  const tryLogin = () => {
    let data = {
      email: userEmail,
      password: userPassword,
    };
    login(data, false);
  };
</script>

<div class="dark-mode-container">
  <button class="toggle-dark-mode" on:click={toggleDarkMode}>
    <Fa icon={faMoon} />
  </button>
</div>
<div class="login-container">
  <div class="input-container">
    <label class="input-label" for="email">Email</label>
    <input
      type="email"
      id="email"
      autocomplete="off"
      required
      class="input"
      bind:value={userEmail}
    />
  </div>
  <div class="input-container">
    <label class="input-label" for="password">Password</label>
    <input
      type="password"
      id="password"
      required
      class="input"
      bind:value={userPassword}
    />
  </div>
  <div class="button-container">
    <button class="form-button" on:click={tryLogin}>Login</button>
  </div>
  <div class="form-footer-container">
    <p class="form-footer">
      Don't have an account? <button class="form-link" on:click={loggedInBool}
        >Register</button
      >
    </p>
  </div>
</div>

<style>
  .login-container {
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
