<script>
  import { navigate } from "svelte-routing";
  import Navbar from "./components/Navbar.svelte";
  import Contacts from "./components/Contacts.svelte";
  import { getData, postData } from "./util/postAndGet";
  import { user } from "./user/user";
  import Modal from "./components/Modal.svelte";
  import { BackendHost } from "../main";
  export let toggleDarkMode;
  export let navToNotifications;
  export let requestToJoinGroup;
  export let getMessageHistory;
  export let removeFromChatListNotificationArray;
  export let inviteToGroup;

  let openModal = false;
  let newGroupTitle = "";
  let newGroupDescription = "";
  let type = "group";

  function createGroup() {
    let dto = {
      creatorID: $user.currentUserID,
      name: newGroupTitle,
      description: newGroupDescription,
    };
    postData(`${BackendHost}/createGroup`, dto);
  }

  const navToGroup = (groupID, groupName) => {
    $user.currentActiveGroupID = groupID;
    $user.currentActiveGroupName = groupName;
    navigate("/group", { replace: true });
  };

  const createGroupModal = () => {
    openModal = true;
  };

  getGroupsData();
  async function getGroupsData() {
    getData(`${BackendHost}/groups`).then((data) => {
      $user.allGroups = data;
    });
  }
</script>

{#if openModal}
  <Modal
    bind:newGroupTitle
    bind:newGroupDescription
    bind:openModal
    {type}
    {createGroup}
    {inviteToGroup}
  />
{/if}
<Navbar {toggleDarkMode} {navToNotifications} />
<main class="main-page-container">
  <div class="main-left-sidebar">
    <div class="main-left-sidebar-content">
      <div class="groups-page-button">
        <button class="sidebar-button" on:click={createGroupModal}
          >Create a Group</button
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
        {#if $user.allGroups != null}
          {#each $user.allGroups as group, index (index)}
            <div class="group-container">
              <div class="group-info">
                <div class="group-name">{group.name}</div>
                <div class="group-desc">
                  |{group.description}
                </div>
              </div>
              <button
                class="group-button"
                on:click={() =>
                  requestToJoinGroup(group.groupID, $user.currentUserID)}
                >Join</button
              >
            </div>
          {/each}
        {/if}
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

  .group-container {
    display: flex;
    align-items: center;
    justify-content: space-between;

    padding: 1rem;
    width: 34rem;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 5px 10px;
  }

  .group-info {
    display: flex;
    flex-direction: column;
  }

  .group-name {
    font-size: 1.2rem;
    font-weight: bold;
  }

  .group-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 2rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1rem;
    font-weight: lighter;
  }

  .group-button:hover {
    background: var(--primary-hover);
  }
</style>
