<script>
  import { navigate, Link } from "svelte-routing";
  import Fa from "svelte-fa";
  import {
    faUser,
    faMoon,
    faMessage,
    faRightFromBracket,
  } from "@fortawesome/free-solid-svg-icons";
  import { user, profileData } from "../user/user";
  import { getData, postData } from "../util/postAndGet";
  import { BackendHost } from "../../main";
  export let toggleDarkMode;
  export let navToNotifications;

  async function getUserData(id) {
    getData(`${BackendHost}/profile?userID=${id}`).then((data) => {
      $profileData.publicBool = data.userStruct.isProfilePublic;
      $profileData.followedList = data.Followed;
      $profileData.followerList = data.Followers;
      $profileData.userData = data.userStruct;
      $profileData.profilePosts = data.Posts;
      navigate(`/loading/profile?userID=${id}`, { replace: false });
    });
  }

  async function navToLogin() {
    $user = {
      currentUserFollowers: [],
      currentUserFollowed: [],
      recipientFollowers: [],
      recipientFollowed: [],

      currentUserID: undefined,
      currentUsername: "",
      chatListOrderArray: [],
      chatListNotificationArray: [],
      userListFromEvent: [],
      isLoggedIn: false,
      currentChatList: [],
      currentChatDetails: [],
      currentGroupChatDetails: [],
      recipientID: null,
      recipientName: null,
      chatOpened: false,
      chatJustOpened: false,
      allGroups: [],
      JoinedGroupIDArray: [],
      currentActiveGroupID: null,
      currentActiveGroupName: undefined,
      userJoinedGroups: [],
      currentActiveGroupMembersDetails: {
        user: [],
        groupPosts: [],
        followers: [],
        events: [],
        description: undefined,
      },
      allNotificationArray: [],
    };
    await postData(`${BackendHost}/logout`).then((data) => {
      if (data) {
        navigate("/login", { replace: true });
      }
    });
  }

  const navToProfile = () => {
    getUserData($user.currentUserID);
  };

  const navToMain = () => {
    navigate("/loading/feed", { replace: true });
  };
</script>

<header class="navbar">
  <div class="navbar-content">
    <div class="navbar-home">
      <button class="navbar-home-button" on:click={navToMain}>
        Social Network
      </button>
    </div>
    <div class="navbar-buttons">
      <button class="toggle-dark-mode" on:click={toggleDarkMode}>
        <Fa icon={faMoon} />
      </button>
      <button class="navbar-profile" on:click={navToProfile}>
        <div class="avatar-container">
          <img
            class="profile-image"
            src={`${BackendHost}/profile/image?userID=${
              $user.currentUserID
            }&timestamp=${Date.now()}`}
            alt="profileImage"
          />
        </div>
        <div class="name-container">{$user.currentUsername}</div>
      </button>
      <button class="navbar-notifications" on:click={navToNotifications}>
        <Fa icon={faMessage} />
      </button>
      <button class="navbar-logout" on:click={navToLogin}>
        <Fa icon={faRightFromBracket} />
      </button>
    </div>
  </div>
</header>

<style>
  .profile-image {
    width: inherit;
    height: inherit;
    border-radius: inherit;
    object-fit: cover;
  }
  .navbar {
    display: flex;

    position: fixed;

    height: 3rem;
    width: 100%;
    background: var(--secondary-background);
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 0px 15px;

    z-index: 1;
  }

  .navbar-content {
    width: 100%;

    display: flex;
    align-items: center;
    justify-content: space-between;

    padding: 1rem;
  }

  .navbar-buttons {
    display: flex;
  }

  .navbar-home-button {
    background: none;
    outline: none;
    border: none;
    color: var(--text-color);
    font-weight: bold;
    font-size: 1.6rem;
  }

  .navbar-home-button:hover {
    cursor: pointer;
  }

  .navbar-profile {
    display: flex;
    align-items: center;

    background: none;
    border: none;
    border-radius: 1rem;
    color: var(--text-color);
    transition: background 0.1s;

    padding: 0.3rem;
  }
  .navbar-profile:hover {
    background: var(--secondary-background-hover);
    cursor: pointer;
  }

  .avatar-container {
    background: var(--accent);
    display: flex;
    align-items: center;
    justify-content: center;

    height: 1.8rem;
    width: 1.8rem;

    border-radius: 50%;
  }

  .name-container {
    padding-left: 0.5rem;
    font-weight: bold;
  }

  .navbar-notifications {
    background: none;
    border: none;
    border-radius: 1rem;
    color: var(--text-color);
    transition: background 0.1s;

    padding: 0.3rem 0.6rem;
  }
  .navbar-notifications:hover {
    background: var(--secondary-background-hover);
    cursor: pointer;
  }

  .navbar-logout {
    background: none;
    border: none;
    border-radius: 1rem;
    color: var(--text-color);
    font-weight: bold;
    transition: background 0.1s;

    padding: 0.3rem 0.6rem;
  }
  .navbar-logout:hover {
    background: var(--secondary-background-hover);
    cursor: pointer;
  }
</style>
