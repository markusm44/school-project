<script>
  import Navbar from "./components/Navbar.svelte";
  import { navigate } from "svelte-routing";
  export let toggleDarkMode;
  export let navToNotifications;

  setTimeout(() => {
    const path = window.location.pathname;
    const newPath = path.split("/");
    navigate(`/${newPath[2]}`, { replace: false });
  }, 300);
</script>

<Navbar {toggleDarkMode} {navToNotifications} />
<main />
