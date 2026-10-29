<script>
  import { navigate } from "svelte-routing";
  import Navbar from "./components/Navbar.svelte";
  import Contacts from "./components/Contacts.svelte";
  import Post from "./components/Post.svelte";
  import Event from "./components/Event.svelte";
  import Modal from "./components/Modal.svelte";
  import Fa from "svelte-fa";
  import { faUser } from "@fortawesome/free-solid-svg-icons";
  import { handleFileInput } from "./util/image";
  import { getData, postData } from "./util/postAndGet";
  import { GetTime } from "./util/getTime";
  import { BackendHost } from "../main";
  import { user } from "./user/user";
  export let toggleDarkMode;
  export let navToNotifications;
  export let getMessageHistory;
  export let updateActiveGroupInBackend;
  export let sendNotificationOfEvent;
  export let removeFromChatListNotificationArray;
  export let inviteToGroup;
  export let createGroup;
  let postContent = "";
  let eventDetails = {
    title: "",
    description: "",
    date: "",
    going: 0,
    notGoing: 0,
  };
  let newGroupTitle = "";
  let newGroupDescription = "";
  let openModal = false;
  let type = "invite";

  let groupFeedData;

  let fromFeed = true;

  let imageData = {
    base64String: "",
    name: "",
  };

  var allMessagesInCurrentGroupChatContainer = [];

  async function handleFileInputChange(event) {
    const result = await handleFileInput(event);
    if (result.dataUrl && result.imageInfo) {
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

  export async function getGroupMessageHistory() {
    await getData(
      "http://localhost:8080/groupMessage?groupID=" +
        $user.currentActiveGroupID,
    ).then((value) => {
      if (value === null) {
        allMessagesInCurrentGroupChatContainer = [];
      } else {
        allMessagesInCurrentGroupChatContainer = value;
      }

      $user.currentGroupChatDetails = allMessagesInCurrentGroupChatContainer;
    });
  }

  const navToGroupChat = () => {
    for (let key in $user.userJoinedGroups) {
      if ($user.userJoinedGroups[key].groupID === $user.currentActiveGroupID) {
        $user.currentActiveGroupName = $user.userJoinedGroups[key].name;
      }
    }
    updateActiveGroupInBackend($user.currentActiveGroupID);
    getGroupMessageHistory();
    navigate("/loading/groupchat", { replace: true });
  };

  const navToMember = () => {};

  const createNewGroupPost = (id) => {
    if (postContent == "") {
      alert("Can't post empty post");
      return;
    }

    let groupPostData = {
      userID: id,
      groupID: $user.currentActiveGroupID,
      image: imageData,
      content: postContent,
      created: GetTime(),
      fullname: $user.currentUsername,
      status: "",
      numberOfCommentsOnPost: 0,
      almostPrivateUsers: null,
      commentsInfo: [],
    };

    postData(`${BackendHost}/group?entryType=post`, groupPostData).then(
      (data) => {
        groupPostData["postID"] = data.entryID;
        groupFeedData = [groupPostData, ...groupFeedData];
      },
    );
  };

  const createNewGroupEvent = (id) => {
    if (eventDetails.title == "") {
      return alert("Can't create an event without a title!");
    }

    if (eventDetails.description == "") {
      return alert("Can't create an event without a description!");
    }

    if (eventDetails.date == "") {
      return alert(
        "Can't create an event without selecting a date for the event!",
      );
    }

    let data = {
      userID: id,
      groupID: $user.currentActiveGroupID,
      title: eventDetails.title,
      content: eventDetails.description,
      created: GetTime(),
      fullname: $user.currentUsername,
      notgoing: 0,
      going: 0,
      date: eventDetails.date,
    };

    postData(`${BackendHost}/group?entryType=event`, data).then((resp) => {
      data["eventID"] = resp.entryID;
    });

    $user.currentActiveGroupMembersDetails.events = [
      ...$user.currentActiveGroupMembersDetails.events,
      data,
    ];

    sendNotificationOfEvent($user.currentActiveGroupID);
  };

  const openInviteModal = () => {
    openModal = true;
  };

  const populateGroupFeed = () => {
    groupFeedData = $user.currentActiveGroupMembersDetails.groupPosts.reverse();
  };

  populateGroupFeed();
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
      <div class="sidebar-title">
        <div class="sidebar-title-name">{$user.currentActiveGroupName}</div>
        <button class="sidebar-title-chat" on:click={() => navToGroupChat()}
          >Chat
        </button>
      </div>
      <div class="sidebar-list">
        {#each $user.currentActiveGroupMembersDetails.user as groupMember, index (index)}
          <button class="sidebar-button" on:click={navToMember}>
            <div class="group-user-img">
              <Fa icon={faUser} />
            </div>
            <div class="group-user-name">
              {groupMember.firstName}
              {groupMember.lastName}
            </div>
          </button>
        {/each}
        <div class="sidebar-invite-container">
          <button class="sidebar-invite-button" on:click={openInviteModal}
            >Invite to Group</button
          >
        </div>
      </div>
    </div>
  </div>
  <div class="posts-list-container">
    <div class="posts-list-content">
      <div class="group-desc-container">
        <div class="group-desc">
          {$user.currentActiveGroupMembersDetails.description}
        </div>
      </div>
      <div class="create-post-container">
        <div class="create-event-title">Create Event</div>
        <div class="input-container">
          <input
            type="text"
            id="createEventTitle"
            autocomplete="off"
            required
            class="input"
            placeholder="Event title..."
            bind:value={eventDetails.title}
          />
        </div>
        <div class="input-container">
          <input
            type="text"
            id="createEventDesc"
            autocomplete="off"
            required
            class="input"
            placeholder="Event description..."
            bind:value={eventDetails.description}
          />
        </div>
        <div class="input-container">
          <input
            type="date"
            id="createEventDate"
            autocomplete="off"
            required
            class="input"
            bind:value={eventDetails.date}
          />
        </div>
        <div class="event-button-container">
          <button
            class="create-post-form-button"
            on:click={() => createNewGroupEvent($user.currentUserID)}
            >Create</button
          >
        </div>
      </div>
      <div class="events-list">
        {#if $user.currentActiveGroupMembersDetails.events != undefined || $user.currentActiveGroupMembersDetails.events != null}
          {#each $user.currentActiveGroupMembersDetails.events as _, index (index)}
            <Event
              data={$user.currentActiveGroupMembersDetails.events}
              currentIndex={index}
            />
          {/each}
        {/if}
      </div>
      <div class="spacer" />
      <div class="create-post-container">
        <div class="create-post-header-text">Create a Post</div>
        <div class="input-container">
          <input
            type="text"
            id="createGroupPostText"
            autocomplete="off"
            required
            class="input"
            placeholder="Create a post..."
            bind:value={postContent}
          />
        </div>
        <div class="post-button-container">
          <input
            type="file"
            accept="image/png, image/jpg, image/jpeg"
            id="addPostImage"
            class="input"
            on:change={(e) => handleFileInputChange(e)}
          />
          <button
            class="create-post-form-button"
            on:click={() => createNewGroupPost($user.currentUserID)}
            >Post</button
          >
        </div>
      </div>
      <div class="posts-list">
        {#if groupFeedData != undefined || groupFeedData != null}
          {#each groupFeedData as _, index (index)}
            <Post data={groupFeedData} currentIndex={index} {fromFeed} />
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
  }

  .sidebar-title {
    display: flex;
    justify-content: space-between;
    align-items: center;

    padding: 0 0 0 0.5rem;
    color: var(--text-color);

    border-bottom: 1px solid #5d5d5d;
  }

  .sidebar-title-name {
    font-weight: bold;
    font-size: 1.2rem;
  }

  .sidebar-title-chat {
    background: var(--primary);
    border: none;
    padding: 0.3rem 0.2rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1rem;
    font-weight: lighter;
  }

  .sidebar-title-chat:hover {
    background: var(--primary-hover);
  }

  .sidebar-list {
    display: flex;
    flex-direction: column;
    align-items: center;

    height: 100%;
    overflow-y: auto;

    padding-top: 0.5rem;
  }

  .sidebar-list::-webkit-scrollbar {
    width: 0;
  }

  .sidebar-button {
    display: flex;
    align-items: center;

    background: none;
    border: none;
    border-radius: 0.5rem;
    width: 95%;
    padding: 0.5rem 0.5rem;

    color: var(--text-color);
    text-align: left;
    font-size: 0.9rem;

    transition: background 0.2s;
  }

  .sidebar-button:hover {
    background: var(--background-hover);
  }

  .sidebar-invite-container {
    padding: 0.5rem;
  }

  .sidebar-invite-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 3rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1.2rem;
    font-weight: lighter;

    cursor: pointer;
  }

  .sidebar-invite-button:hover {
    background: var(--primary-hover);
  }

  .group-user-img {
    background: var(--accent);
    display: flex;
    align-items: center;
    justify-content: center;

    height: 1.8rem;
    width: 1.8rem;

    border-radius: 50%;
  }

  .group-user-name {
    padding-left: 0.5rem;
    font-weight: bold;
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

  .create-post-container {
    display: flex;
    flex-direction: column;
    align-items: center;

    padding: 1rem;
    width: 34rem;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 5px 10px;
  }

  .create-post-header-text {
    padding: 0 0 0.5rem 0.5rem;
    font-weight: bold;
    font-size: 1.2rem;
    width: 100%;
  }

  .post-button-container {
    width: 100%;
    display: flex;
    justify-content: space-between;
    padding: 1rem 0rem 0rem 0rem;
  }

  .event-button-container {
    width: 100%;
    display: flex;
    justify-content: end;
    padding: 1rem 0rem 0rem 0rem;
  }

  .create-post-form-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 4rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1rem;
    font-weight: lighter;
  }

  .create-post-form-button:hover {
    background: var(--primary-hover);
  }

  .posts-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;

    padding-bottom: 1rem;
  }

  .events-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .create-event-title {
    width: 100%;
    padding-left: 1rem;

    font-weight: bold;
    font-size: 1.2rem;
  }

  .group-desc {
    font-weight: bold;
    padding: 1rem;
  }
</style>
