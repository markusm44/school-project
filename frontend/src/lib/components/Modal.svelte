<script>
  import { user } from "../user/user";
  export let newGroupTitle;
  export let newGroupDescription;
  export let openModal;
  export let type;
  export let createGroup;
  export let inviteToGroup;

  let inviteList = [];

  const closeModal = () => {
    openModal = false;
  };

  const confirmModal = () => {
    if (type == "group") {
      createGroup();
    } else if (type == "invite") {
      for (let i = 0; i < inviteList.length; i++) {
        inviteToGroup($user.currentActiveGroupID, inviteList[i]);
      }
    }
    openModal = false;
  };
</script>

<div class="modal-container">
  <div class="modal">
    {#if type == "group"}
      <div class="modal-title">Create a Group</div>
      <div class="input-container">
        <input
          type="text"
          id="createGroupText"
          autocomplete="off"
          class="input"
          placeholder="Title..."
          bind:value={newGroupTitle}
        />
      </div>
      <div class="input-container">
        <input
          type="text"
          id="createGroupDescription"
          autocomplete="off"
          class="input"
          placeholder="Description..."
          bind:value={newGroupDescription}
        />
      </div>
    {:else if type == "invite"}
      <div class="modal-title">Invite followers to Group</div>
      <select multiple bind:value={inviteList}>
        {#each $user.currentActiveGroupMembersDetails.followers as userList}
          <option value={parseInt(userList.userID)}>{userList.name}</option>
        {/each}
      </select>
    {/if}
    <div class="modal-buttons">
      <button class="modal-button" on:click={closeModal}> Cancel </button>
      <button class="modal-button" on:click={confirmModal}>Confirm</button>
    </div>
  </div>
</div>

<style>
  .modal-container {
    position: absolute;
    width: 100vw;
    height: 100vh;
    z-index: 2;
  }

  .modal {
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

  .modal-buttons {
    display: flex;
    justify-content: space-around;
    align-items: center;

    width: 90%;
  }

  .modal-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 3rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1.2rem;
    font-weight: lighter;

    cursor: pointer;
  }

  .modal-button:hover {
    background: var(--primary-hover);
  }
</style>
