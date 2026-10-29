<script>
  import Navbar from "./components/Navbar.svelte";
  import Contacts from "./components/Contacts.svelte";
  import Fa from "svelte-fa";
  import { faFaceSmile } from "@fortawesome/free-solid-svg-icons";
  import { user } from "./user/user";
  import { BackendHost } from "../main";
  import EmojiPicker from "svelte-emoji-picker";
  export let toggleDarkMode;
  export let navToNotifications;
  export let sendChatMessage;
  export let getMessageHistory;
  export let removeFromChatListNotificationArray;
  let inputText = "";
  let mutualFollowers;

  areMutualFollowers();

  function areMutualFollowers() {
    for (const follower of $user.currentUserFollowers) {
      if (follower.userID === $user.recipientID) {
        //check if current user is followed by recipient
        for (const recipientsFollower of $user.recipientFollowers) {
          if (recipientsFollower.userID === $user.currentUserID) {
            //check if recipient is followed by current user
            mutualFollowers = true;
            return;
          }
        }
      }
    }
    mutualFollowers = false;
  }

  const resetInput = (text) => {
    let sentChat = {
      content: text,
      senderID: $user.currentUserID,
      receiverID: $user.recipientID,
      date: getMessageSentTime(),
      messageID: null,
    };
    $user.currentChatDetails = [...$user.currentChatDetails, sentChat];
    inputText = "";
  };

  function getMessageSentTime() {
    const currentDate = new Date();

    const datetime = new Intl.DateTimeFormat("en-GB", {
      dateStyle: "short",
      timeStyle: "short",
    });
    return datetime.format(currentDate);
  }

  let pickerOpen = false;
  const togglePicker = () => {
    pickerOpen = !pickerOpen;
  };
</script>

<Navbar {toggleDarkMode} {navToNotifications} />
<main class="main-page-container">
  <div class="chat-container">
    <div class="chat-header">
      <div class="chat-user-avatar">
        <img
          class="profile-image"
          src={`${BackendHost}/profile/image?userID=${
            $user.recipientID
          }&timestamp=${Date.now()}`}
          alt="profileImage"
        />
      </div>
      <div class="chat-user-name">{$user.recipientName}</div>
    </div>
    {#if mutualFollowers}
      <div class="chat-content">
        <div class="chat-text-area">
          {#each $user.currentChatDetails as message, index (index)}
            {#if message.senderID == $user.currentUserID}
              <div class="chat-message-me-container">
                <div class="chat-message-name">{$user.currentUsername}</div>
                <div class="chat-message-date">{message.date}</div>
                <div class="chat-message-me">{message.content}</div>
              </div>
            {:else if message.receiverID == $user.currentUserID}
              <div class="chat-message-other-container">
                <div class="chat-message-name">{$user.recipientName}</div>
                <div class="chat-message-date">{message.date}</div>
                <div class="chat-message-other">{message.content}</div>
              </div>
            {/if}
          {/each}
        </div>
      </div>
      <div class="chat-footer">
        <div class="chat-footer-inputs">
          <input
            type="text"
            placeholder="Send a message.."
            autocomplete="off"
            class="chat-footer-textarea"
            bind:value={inputText}
          />
          <div class="chat-emoji-container">
            <button class="chat-emoji" on:click={() => togglePicker()}>
              <Fa icon={faFaceSmile} />
            </button>
            {#if pickerOpen}
              <EmojiPicker bind:value={inputText} />
            {/if}
          </div>
        </div>
        <button
          class="chat-send-button"
          on:click={(event) => sendChatMessage(event, inputText)}
          on:click={() => resetInput(inputText)}>Send</button
        >
      </div>
    {:else}
      <div class="chat-user-name">
        To chat you need to be mutual followers with {$user.recipientName}
      </div>
    {/if}
  </div>
  <Contacts {getMessageHistory} {removeFromChatListNotificationArray} />
</main>

<style>
  .profile-image {
    width: inherit;
    height: inherit;
    border-radius: inherit;
    object-fit: cover;
  }

  .main-page-container {
    position: absolute;
    top: 3rem;
    width: 100%;

    display: flex;
  }

  .chat-container {
    padding: 2rem 5rem;
    height: 100%;
    width: 60rem;

    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .chat-header {
    display: flex;
    align-items: center;

    background: none;
    border: none;
    border-radius: 1rem;
    color: var(--text-color);
    transition: background 0.1s;

    padding: 0.5rem;
  }

  .chat-user-avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--accent);
    width: 2.5rem;
    height: 2.5rem;

    border-radius: 50%;
  }

  .chat-user-name {
    padding-left: 0.5rem;
    font-weight: bold;
  }

  .chat-content {
    width: 100%;
    height: 40rem;
    background: var(--background-hover);
    border-radius: 0.5rem;
  }

  .chat-text-area {
    display: flex;
    flex-direction: column;
    height: 100%;
    overflow-y: auto;
  }

  .chat-text-area::-webkit-scrollbar {
    background: none;
  }

  .chat-text-area::-webkit-scrollbar-thumb {
    background: var(--primary);
    border-radius: 0.5rem;
  }

  .chat-text-area::-webkit-scrollbar-thumb:hover {
    background: var(--primary-hover);
  }

  .chat-message-me-container {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    position: relative;
    padding-right: 0.5rem;
    margin-bottom: 0.5rem;
    margin-left: 2rem;
  }

  .chat-message-other-container {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    position: relative;
    padding-left: 0.5rem;
    margin-bottom: 0.5rem;
    margin-right: 2rem;
  }

  .chat-message-name {
    font-size: 0.7rem;
    color: var(--text-color);
  }

  .chat-message-date {
    font-size: 0.7rem;
    color: #bebebe;
  }

  .chat-message-me {
    color: var(--text-color);
    background: var(--primary-hover);
    padding: 0.5rem;
    font-size: 0.9rem;
    border-radius: 0.5rem;
  }

  .chat-message-other {
    color: var(--text-color);
    background: var(--accent);
    padding: 0.5rem;
    font-size: 0.9rem;
    border-radius: 0.5rem;
  }

  .chat-footer {
    height: 2rem;
    width: 100%;

    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .chat-footer-inputs {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .chat-footer-textarea {
    color: var(--text-color);
    background: var(--background-hover);
    border: none;
    border-radius: 0.5rem;
    padding: 0.8rem 1rem;
    width: 34rem;
  }

  .chat-footer-textarea:focus {
    outline: none;
  }

  .chat-emoji {
    border: none;
    border-radius: 0.5rem;
    background: var(--background-hover);
    color: var(--accent);

    padding: 0.5rem 0.6rem;
    font-size: 1.3rem;
  }

  .chat-emoji:hover {
    background: #1c2f2f;
  }

  .chat-send-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 2rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1rem;
    font-weight: lighter;
  }

  .chat-send-button:hover {
    background: var(--primary-hover);
  } 
</style>
