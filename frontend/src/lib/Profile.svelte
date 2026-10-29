<script>
  import { navigate } from "svelte-routing";
  import Navbar from "./components/Navbar.svelte";
  import Post from "./components/Post.svelte";
  import Contacts from "./components/Contacts.svelte";
  import Fa from "svelte-fa";
  import { faUser } from "@fortawesome/free-solid-svg-icons";
  import { user, profileData } from "./user/user";
  import { getData, postData } from "./util/postAndGet";
  import { BackendHost } from "../main";
  import { handleFileInput } from "./util/image";
  export let toggleDarkMode;
  export let navToNotifications;
  export let getMessageHistory;
  export let removeFromChatListNotificationArray;
  let followBool = false;
  $user.recipientID = null;
  let imageData = {
    base64String: "",
    name: "",
  };

  let fromFeed = false;

  async function handleFileInputChange(event, id) {
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

  const changeProfilePic = () => {
    postData(
      `${BackendHost}/profile?userID=${$user.currentUserID}&changeImage=true`,
      imageData,
    );
    getUserData($user.currentUserID);
  };

  function getUserData(id) {
    getData(`${BackendHost}/profile?userID=${id}`).then((data) => {
      $profileData.publicBool = data.userStruct.isProfilePublic;
      $profileData.followedList = data.Followed;
      $profileData.followerList = data.Followers;
      $profileData.userData = data.userStruct;
      $profileData.profilePosts = data.Posts;
      navigate(`/loading/profile?userID=${id}`, { replace: false });
    });
  }

  const checkFollowing = (userID) => {
    for (let key in $profileData.followerList) {
      if ($profileData.followerList[key].userID == userID) {
        followBool = true;
        break;
      }
    }
  };

  const toggleFollow = (id) => {
    let shouldFollow;

    let followData = {
      userID: $user.currentUserID,
      name: $user.currentUsername,
    };

    for (let key in $profileData.followerList) {
      if ($profileData.followerList[key].userID == $user.currentUserID) {
        shouldFollow = false;
        break;
      } else {
        shouldFollow = true;
        $profileData.followerList = [...$profileData.followerList, followData];
        break;
      }
    }

    if ($profileData.followerList !== null) {
      if (!$profileData.followerList.length) {
        $profileData.followerList = null;
      }
    }

    if (!$profileData.followerList) {
      shouldFollow = true;
      $profileData.followerList = [followData];
    }

    if (!shouldFollow) {
      followBool = false;
      for (let key in $profileData.followerList) {
        if ($profileData.followerList[key].userID == $user.currentUserID) {
          $profileData.followerList = $profileData.followerList.filter(
            (elem) => {
              return elem !== $profileData.followerList[key];
            },
          );
        }
      }
    }

    if ($profileData.publicBool) {
      if (shouldFollow) {
        followBool = true;
      }
    }

    let data = {
      publicFollowBool: shouldFollow,
    };

    postData(`${BackendHost}/profile?userID=${id}`, data);
  };

  const togglePrivacy = (id) => {
    let data = {
      publicFollowBool: ($profileData.publicBool = !$profileData.publicBool),
    };

    postData(`${BackendHost}/profile?userID=${id}`, data);
  };

  checkFollowing($user.currentUserID);
</script>

<Navbar {toggleDarkMode} {navToNotifications} />
<main class="profile-page-container">
  <div class="profile-container">
    {#if $profileData.publicBool || $profileData.userData.userID == $user.currentUserID || followBool}
      <div class="profile-content-container">
        <div class="profile-infos">
          <div class="profile-info-container">
            <div class="profile-info-content">
              <div class="profile-avatar-container">
                <img
                  class="profile-image"
                  src={`${BackendHost}/profile/image?userID=${
                    $profileData.userData.userID
                  }&timestamp=${Date.now()}`}
                  alt="NO AVATAR"
                />
              </div>
              <div class="profile-user-name">
                {$profileData.userData.firstName}
                {$profileData.userData.lastName}
              </div>
              <div class="profile-user-nickname">
                {$profileData.userData.nickname}
              </div>
              <div class="profile-user-email">
                {$profileData.userData.email}
              </div>
              <div class="profile-user-dob">
                {$profileData.userData.dateOfBirth}
              </div>
              {#if $profileData.userData.userID == $user.currentUserID}
                <button
                  class="profile-user-privacy"
                  on:click={() => togglePrivacy($profileData.userData.userID)}
                >
                  {#if $profileData.publicBool}
                    Public
                  {:else}
                    Private
                  {/if}
                </button>
                <input
                  type="file"
                  accept="image/*"
                  id="avatarIcon"
                  class="avatar-input"
                  on:change={(e) =>
                    handleFileInputChange(e, $profileData.userData.userID)}
                />
                <button
                  class="profile-user-privacy"
                  on:click={() => changeProfilePic()}>Change</button
                >
              {:else}
                <button
                  class="profile-user-privacy"
                  on:click={() => toggleFollow($profileData.userData.userID)}
                >
                  {#if !followBool}
                    Follow
                  {:else}
                    Unfollow
                  {/if}
                </button>
              {/if}
            </div>
          </div>
          <div class="profile-info-about-me-container">
            <div class="profile-info-about-me-content">
              <div class="profile-info-about-me-title">About Me</div>
              <div class="profile-info-about-me-text">
                {$profileData.userData.aboutMe}
              </div>
            </div>
          </div>
        </div>
        <div class="profile-contacts-container">
          <div class="profile-follow-container">
            <div class="profile-follow-title-container">
              {#if $profileData.followedList != null}
                Following ({$profileData.followedList.length})
              {:else}
                Following (0)
              {/if}
            </div>
            <div class="profile-follow-list">
              {#if $profileData.followedList != null}
                {#each $profileData.followedList as _, index (index)}
                  <button
                    id="user-followed-{$profileData.followedList[index].userID}"
                    class="follow-user"
                    on:click={() =>
                      getUserData($profileData.followedList[index].userID)}
                  >
                    <div class="follow-user-avatar">
                      <img
                        class="profile-image"
                        src={`${BackendHost}/profile/image?userID=${
                          $profileData.followedList[index].userID
                        }&timestamp=${Date.now()}`}
                        alt="NO AVATAR"
                      />
                    </div>
                    <div class="follow-user-name">
                      {$profileData.followedList[index].name}
                    </div>
                  </button>
                {/each}
              {/if}
            </div>
          </div>
          <div class="profile-follow-container">
            <div class="profile-follow-title-container">
              {#if $profileData.followerList != null}
                Followers ({$profileData.followerList.length})
              {:else}
                Followers (0)
              {/if}
            </div>
            <div class="profile-follow-list">
              {#if $profileData.followerList != null}
                {#each $profileData.followerList as _, index (index)}
                  <button
                    id="user-follower-{$profileData.followerList[index].userID}"
                    class="follow-user"
                    on:click={() =>
                      getUserData($profileData.followerList[index].userID)}
                  >
                    <div class="follow-user-avatar">
                      <img
                        class="profile-image"
                        src={`${BackendHost}/profile/image?userID=${
                          $profileData.followerList[index].userID
                        }&timestamp=${Date.now()}`}
                        alt="NO AVATAR"
                      />
                    </div>
                    <div class="follow-user-name">
                      {$profileData.followerList[index].name}
                    </div>
                  </button>
                {/each}
              {/if}
            </div>
          </div>
        </div>
        <div class="spacer" />
        <div class="profile-activity-container">
          {#if $profileData.profilePosts != undefined || $profileData.profilePosts != null}
            {#each $profileData.profilePosts as post, index (index)}
              <Post
                data={$profileData.profilePosts}
                currentIndex={index}
                {fromFeed}
              />
            {/each}
          {/if}
        </div>
      </div>
    {:else}
      <div class="profile-content-container">
        <div class="profile-infos">
          <div class="profile-info-container">
            <div class="profile-info-content">
              <div class="profile-avatar-container">
                <Fa icon={faUser} />
              </div>
              <div class="profile-user-name">
                {$profileData.userData.firstName}
                {$profileData.userData.lastName}
              </div>
              <div class="profile-user-nickname">
                {$profileData.userData.nickname}
              </div>
              <div class="profile-user-email">
                {$profileData.userData.email}
              </div>
              <div class="profile-user-dob">
                {$profileData.userData.dateOfBirth}
              </div>
              {#if $profileData.userData.userID == $user.currentUserID}
                <button
                  class="profile-user-privacy"
                  on:click={() => togglePrivacy($profileData.userData.userID)}
                >
                  {#if $profileData.publicBool}
                    Public
                  {:else}
                    Private
                  {/if}
                </button>
              {:else}
                <button
                  class="profile-user-privacy"
                  on:click={() => toggleFollow($profileData.userData.userID)}
                >
                  {#if !followBool}
                    Follow
                  {:else}
                    Unfollow
                  {/if}
                </button>
              {/if}
            </div>
          </div>
          <div class="profile-info-about-me-container">
            <div class="profile-info-about-me-content">
              <div class="profile-info-about-me-title">About Me</div>
              <div class="profile-info-about-me-text">
                {$profileData.userData.aboutMe}
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </div>
  <Contacts {getMessageHistory} {removeFromChatListNotificationArray} />
</main>

<style>
  .profile-page-container {
    position: absolute;
    top: 3rem;
    width: 100%;

    display: flex;
  }

  .profile-container {
    display: flex;

    width: 70%;
    height: fit-content;
    padding-left: 2rem;
  }

  .profile-content-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    width: 100%;
  }

  .profile-infos {
    display: flex;
    gap: 1rem;
    width: 100%;
  }

  .profile-info-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;

    height: 26rem;
    width: 50%;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 0px 15px;
  }

  .profile-info-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    height: 100%;
  }

  .profile-avatar-container {
    display: flex;
    align-items: center;
    justify-content: center;

    background: var(--accent);
    width: 10rem;
    height: 10rem;

    font-size: 5rem;

    border-radius: 50%;
  }

  .profile-image {
    width: inherit;
    height: inherit;
    border-radius: inherit;
    object-fit: cover;
  }

  .avatar-input {
    position: relative;
    padding-top: 0.5rem;
    left: 2rem;
  }

  .profile-user-name {
    padding-top: 1rem;
    font-weight: bold;
    font-size: 2rem;
  }

  .profile-info-about-me-container {
    width: 50%;
    height: fit-content;
    text-overflow: ellipsis;
    overflow: hidden;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 0px 15px;
  }

  .profile-info-about-me-content {
    padding: 1rem;
  }

  .profile-info-about-me-title {
    font-weight: bold;
    font-size: 1.2rem;
  }

  .profile-contacts-container {
    display: flex;
    gap: 1rem;

    width: 100%;
    height: 24rem;
  }

  .profile-follow-container {
    width: 50%;
    height: 100%;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 0px 15px;
  }

  .profile-follow-title-container {
    display: flex;
    justify-content: center;

    padding: 0.5rem;
    font-weight: bold;
    font-size: 1.2rem;
  }

  .profile-follow-list {
    padding: 0.5rem;
    width: 100%;
    height: 80%;

    display: grid;
    grid-template-columns: repeat(6, 1fr);

    overflow-y: auto;
  }

  .profile-follow-list::-webkit-scrollbar {
    width: 0;
  }

  .follow-user {
    background: none;
    border: none;

    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;

    height: 5.5rem;
    min-width: 5.4rem;
    width: 5.4rem;

    border-radius: 1rem;
    color: var(--text-color);
    transition: background 0.1s;

    padding: 0.5rem;
  }

  .follow-user:hover {
    background: var(--secondary-background-hover);
    cursor: pointer;
  }

  .follow-user-avatar {
    background: var(--accent);
    display: flex;
    align-items: center;
    justify-content: center;

    height: 3rem;
    width: 3rem;

    border-radius: 50%;
  }

  .follow-user-name {
    font-weight: bold;
  }

  .profile-user-privacy {
    border: none;
    background: var(--primary);
    font-weight: light;
    font-size: 1.1rem;
    padding: 0.2rem 0.5rem;
    border-radius: 0.5rem;
    transition: background 0.2s;
  }

  .profile-user-privacy:hover {
    background: var(--primary-hover);
  }

  .profile-activity-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;

    padding-bottom: 1rem;
  }
</style>
