<script>
  import { navigate } from "svelte-routing";
  import Fa from "svelte-fa";
  import { faUser, faComment } from "@fortawesome/free-solid-svg-icons";
  import { user, profileData } from "../user/user";
  import { getData } from "../util/postAndGet";
  import { BackendHost } from "../../main";
  export let getMessageHistory;
  export let removeFromChatListNotificationArray;

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
  async function getRecipientFollowingData(id) {
    getData(`${BackendHost}/profile?userID=${id}`).then((data) => {
      $user.recipientFollowed = data.Followed;
      $user.recipientFollowers = data.Followers;
      if ($user.recipientFollowed === null) {
        $user.recipientFollowed = [];
      }
      if ($user.recipientFollowers === null) {
        $user.recipientFollowers = [];
      }
    });
  }

  async function getCurrentUserData(id) {
    getData(`${BackendHost}/profile?userID=${id}`).then((data) => {
      $user.currentUserFollowed = data.Followed;
      $user.currentUserFollowers = data.Followers;
      if ($user.currentUserFollowed === null) {
        $user.currentUserFollowed = [];
      }
      if ($user.currentUserFollowers === null) {
        $user.currentUserFollowers = [];
      }
    });
  }

  const navToChat = (recipientID, recipientName) => {
    $user.recipientID = recipientID;
    $user.recipientName = recipientName;
    removeFromChatListNotificationArray($user.recipientID);

    getCurrentUserData($user.currentUserID);
    getRecipientFollowingData(recipientID);
    getMessageHistory(recipientID);
    navigate("/loading/chat", { replace: true });
  };

  const navToProfile = (id) => {
    getUserData(id);
  };
</script>

<div class="contacts-container">
  <div class="contacts-title">Contacts</div>
  <div class="contacts-list-container">
    <div class="contacts-list">
      {#each $user.currentChatList as userInList, index (index)}
        <div class="contact">
          <div class="contact-details">
            <div class="contact-avatar-container">
              <img
                class="profile-image"
                src={`${BackendHost}/profile/image?userID=${
                  userInList.userID
                }&timestamp=${Date.now()}`}
                alt="profileImage"
              />
            </div>
            <div class="contact-name">
              ID:{userInList.userID}
              {userInList.firstName + " " + userInList.lastName}
            </div>
          </div>
          <div class="contact-other">
            <div class="contact-buttons">
              <button
                class="contact-button"
                on:click={() => navToProfile(userInList.userID)}
              >
                <Fa icon={faUser} />
              </button>
              <button
                class="contact-button"
                on:click={() =>
                  navToChat(
                    userInList.userID,
                    userInList.firstName + " " + userInList.lastName,
                  )}
              >
                <Fa icon={faComment} />
              </button>
            </div>
            {#if userInList.notification}
              <div class="contact-notify" />
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>

<style>
  .profile-image {
    width: inherit;
    height: inherit;
    border-radius: inherit;
    object-fit: cover;
  }

  .contacts-container {
    position: fixed;
    right: 0;
    height: 100%;
    width: 16rem;
  }

  .contacts-title {
    padding: 0.5rem 0 0 0.5rem;
    color: var(--text-color);
    font-weight: bold;
    font-size: 1.2rem;

    border-bottom: 1px solid #5d5d5d;
  }

  .contacts-list-container {
    padding-top: 0.5rem;
    height: 100%;
  }

  .contacts-list {
    display: flex;
    flex-direction: column;
    align-items: center;
    height: inherit;

    overflow-y: auto;
  }

  .contacts-list::-webkit-scrollbar {
    width: 0;
  }

  .contact {
    background: none;
    border: none;
    border-radius: 0.5rem;
    width: 95%;
    padding: 0.5rem 0.5rem;

    display: flex;
    align-items: center;
    justify-content: space-between;

    color: var(--text-color);
    text-align: left;
    font-weight: 400;
    font-size: 0.9rem;

    transition: background 0.2s;
  }

  .contact:hover {
    background: var(--background-hover);
  }

  .contact-details {
    display: flex;
    align-items: center;
  }

  .contact-other {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 30%;
  }

  .contact-avatar-container {
    background: var(--accent);
    display: flex;
    align-items: center;
    justify-content: center;

    height: 1.8rem;
    width: 1.8rem;

    border-radius: 50%;
  }

  .contact-name {
    padding-left: 0.5rem;
  }

  .contact-buttons {
    display: flex;
    justify-content: space-between;
    padding-left: 0.5rem;
    width: 3rem;
  }

  .contact-button {
    border: none;
    background: var(--primary);
    font-weight: light;
    font-size: 1rem;
    padding: 0.2rem 0.2rem;
    border-radius: 0.5rem;
    transition: background 0.2s;
  }

  .contact-button:hover {
    background: var(--primary-hover);
  }

  .contact-notify {
    background: var(--primary);
    width: 0.5rem;
    height: 0.5rem;
    border-radius: 50%;
  }
</style>
