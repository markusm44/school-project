<script>
  import { navigate } from "svelte-routing";
  import Navbar from "./components/Navbar.svelte";
  import Contacts from "./components/Contacts.svelte";
  import Post from "./components/Post.svelte";
  import { handleFileInput } from "./util/image";
  import { user, feedData } from "./user/user";
  import { BackendHost } from "../main";
  import { getData, postData } from "./util/postAndGet";
  import { GetTime } from "./util/getTime";
  export let toggleDarkMode;
  export let navToNotifications;
  export let getMessageHistory;
  export let removeFromChatListNotificationArray;
  export let getMainData;
  export let feedPostData;
  export let navToGroup;

  let postContent = "";
  let selected;
  let users = [];
  let imageData = {
    base64String: "",
    name: "",
  };

  let fromFeed = true;

  $user.recipientID = null;

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

  const navToGroupList = () => {
    navigate("/groups", { replace: true });
  };

  const createNewPost = (id) => {
    if (postContent == "") {
      alert("Can't post empty post");
      return;
    }

    if (selected == "almostPrivate" && users.length == 0) {
      alert("No users selected for almost private post");
      return;
    }

    let tempData = {
      userID: id,
      content: postContent,
      image: imageData,
      created: GetTime(),
      fullname: $user.currentUsername,
      status: selected,
      almostPrivateUsers: users,
      commentsInfo: [],
    };

    postContent = "";

    if (feedPostData.length === 0){
      feedPostData=[]
    }

    postData(`${BackendHost}/feed?entryType=post`, tempData).then((resp) => {
      tempData["postID"] = resp.entryID;
      feedPostData = [tempData, ...feedPostData];
    });
  };

  getMainData();
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
      <div class="create-post-container">
        <div class="create-post-header">
          <div class="create-post-header-text">Create a Post</div>
          <div class="privacy-list">
            {#if selected == "almostPrivate"}
              <select multiple bind:value={users}>
                {#each $user.currentChatList as userList}
                  <option value={parseInt(userList.userID)}
                    >{userList.firstName} {userList.lastName}</option
                  >
                {/each}
              </select>
            {/if}
            <select bind:value={selected}>
              <option value="public" selected>Public</option>
              <option value="private">Private</option>
              <option value="almostPrivate">Almost Private</option>
            </select>
          </div>
        </div>
        <div class="input-container">
          <input
            type="text"
            id="createPostText"
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
            class="form-button"
            on:click={() => createNewPost($user.currentUserID)}>Post</button
          >
        </div>
      </div>
      <div class="posts-list">
        {#if feedPostData != undefined || feedPostData != null}
          {#each feedPostData as _, index (index)}
            <Post data={feedPostData} currentIndex={index} {fromFeed} />
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

  .create-post-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    width: 100%;
  }

  .create-post-header-text {
    padding: 0.5rem;
    font-weight: bold;
    font-size: 1.2rem;
    width: 50%;
  }

  .privacy-list {
    width: 100%;
    position: relative;
    text-align: right;
  }

  .post-button-container {
    width: 100%;
    display: flex;
    justify-content: space-between;
    padding: 1rem 0rem 0rem 0rem;
  }

  .posts-list {
    display: flex;
    flex-direction: column;
    gap: 2rem;

    padding-bottom: 1rem;
  }
</style>
