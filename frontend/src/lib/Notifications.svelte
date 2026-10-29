<script>
  import { navigate } from "svelte-routing";
  import Navbar from "./components/Navbar.svelte";
  import Contacts from "./components/Contacts.svelte";
  import Notification from "./components/Notification.svelte";
  import { user } from "./user/user";
  export let toggleDarkMode;
  export let navToNotifications;
  export let joinGroup;
  export let deleteNotificationFromDB;
  export let getMessageHistory;
  export let removeFromChatListNotificationArray;
  export let navToGroup;

  $user.recipientID = null;

  const navToGroupList = () => {
    navigate("/groups", { replace: true });
  };
</script>

<Navbar {toggleDarkMode} {navToNotifications} />
<main class="main-page-container">
  <div class="main-left-sidebar">
    <div class="main-left-sidebar-content">
      <div class="groups-page-button">
        <button class="sidebar-button" on:click={navToGroupList}
          >Look for Groups...</button
        >
      </div>
      <div class="sidebar-list">
        {#each $user.userJoinedGroups as joinedGroup, index (index)}
          <button
            class="sidebar-button"
            on:click={() => navToGroup(joinedGroup.groupID, joinedGroup.name)}
            >{joinedGroup.name}</button
          >
        {/each}
      </div>
    </div>
  </div>
  <div class="posts-list-container">
    <div class="posts-list-content">
      <div class="posts-list">
        {#each $user.allNotificationArray as notification, index (index)}
          <Notification
            notif={notification}
            {joinGroup}
            {deleteNotificationFromDB}
          />
        {/each}
      </div>
    </div>
  </div>
  <Contacts {getMessageHistory} {removeFromChatListNotificationArray} />
</main>

<style>
  .main-page-container {
    position: absolute;
    top: 3rem;
    width: 100%;
  }

  .main-left-sidebar {
    position: fixed;
    width: 16rem;
    height: 100%;
  }

  .main-left-sidebar-content {
    padding-top: 0.5rem;
    height: 92%;
  }

  .groups-page-button {
    padding-bottom: 0.5rem;

    display: flex;
    justify-content: center;

    border-bottom: 1px solid #5d5d5d;
  }

  .sidebar-list {
    display: flex;
    flex-direction: column;
    align-items: center;

    padding-top: 0.5rem;

    height: 100%;
    overflow-y: auto;
  }

  .sidebar-list::-webkit-scrollbar {
    width: 0;
  }

  .sidebar-button {
    background: none;
    border: none;
    border-radius: 0.5rem;
    width: 95%;
    padding: 1rem 0.5rem;

    color: var(--text-color);
    text-align: left;
    font-weight: 600;
    font-size: 0.9rem;

    transition: background 0.2s;
  }

  .sidebar-button:hover {
    background: var(--background-hover);
  }

  .posts-list-container {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    width: 48rem;
  }

  .posts-list-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;

    padding-top: 2rem;
  }

  .posts-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;

    padding-bottom: 1rem;
  }
</style>
