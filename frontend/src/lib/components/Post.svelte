<script>
  import { navigate } from "svelte-routing";
  import Comment from "./Comment.svelte";
  import { handleFileInput } from "../util/image";
  import { getData, postData } from "../util/postAndGet";
  import { GetTime } from "../util/getTime";
  import { BackendHost } from "../../main";
  import { user, profileData } from "../user/user";
  export let data;
  export let currentIndex;
  export let fromFeed;
  let commentText = "";
  let imageData = {
    base64String: "",
    name: "",
  };
  const locationURL = window.location.href;
  const entryType = /[\w]+$/.exec(locationURL);

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

  function addPostComment(postID) {
    if (commentText == "") {
      return;
    }

    let tempData = {
      postID: postID,
      content: commentText,
      created: GetTime(),
      fullname: $user.currentUsername,
      userID: $user.currentUserID,
      image: imageData,
      commentID: 0,
    };

    commentText = "";

    postData(`${BackendHost}/feed?entryType=comment`, tempData).then((resp) => {
      tempData.commentID = resp.entryID;

      data[currentIndex].commentsInfo = [
        ...data[currentIndex].commentsInfo,
        tempData,
      ];
    });
  }

  function addGroupPostComment(postID) {
    if (commentText == "") {
      return;
    }

    let tempData = {
      commentID: 0,
      groupID: $user.currentActiveGroupID,
      postID: postID,
      content: commentText,
      created: GetTime(),
      fullname: $user.currentUsername,
      userID: $user.currentUserID,
      image: imageData,
    };

    commentText = "";

    postData(`${BackendHost}/group?entryType=comment`, tempData).then(
      (resp) => {
        tempData.commentID = resp.entryID;

        data[currentIndex].commentsInfo = [
          ...data[currentIndex].commentsInfo,
          tempData,
        ];
      },
    );
  }

  const addComment = (postID, url) => {
    if (url === "group") {
      addGroupPostComment(postID);
    } else if (url === "feed") {
      addPostComment(postID);
    } else {
      return;
    }
  };

  const navToProfile = (id) => {
    getUserData(id);
  };
</script>

<div class="post-container">
  <div class="post-header">
    <div class="post-user-img">
      <img
        class="profile-image"
        src={`${BackendHost}/profile/image?userID=${
          data[currentIndex].userID
        }&timestamp=${Date.now()}`}
        alt="profileImg"
      />
    </div>
    <div class="post-name-container">
      <button
        class="post-user-name"
        on:click={() => navToProfile(data[currentIndex].userID)}
        >{data[currentIndex].fullname}</button
      >
      <div class="post-timedate">{data[currentIndex].created}</div>
    </div>
  </div>
  <div class="post-text">
    {data[currentIndex].content}
  </div>
  {#if data[currentIndex].image.name != "default_image.png" && data[currentIndex].image.name != ""}
    <div class="post-img-container">
      <img
        class="post-img"
        alt="post-img"
        src={`${BackendHost}/postImage?postID=${data[currentIndex].postID}&postType=${entryType[0]}`}
      />
    </div>
  {/if}
  <div class="spacer-container">
    <div class="spacer" />
  </div>
  <div class="comment-area-container">
    {#if fromFeed}
      <div class="comment-post-container">
        <div class="input-container">
          <input
            type="text"
            id="createPostText"
            autocomplete="off"
            required
            class="input"
            placeholder="Post a comment..."
            bind:value={commentText}
          />
          <input
            type="file"
            accept="image/png, image/jpg, image/jpeg"
            id="addCommentImage"
            class="input"
            on:change={(e) => handleFileInputChange(e)}
          />
        </div>
        <button
          class="post-comment-button"
          on:click={() => addComment(data[currentIndex].postID, entryType[0])}
          >Comment</button
        >
      </div>
    {/if}
    <div class="comments-list">
      {#if data[currentIndex].commentsInfo != null || data[currentIndex].commentsInfo != undefined}
        {#each data[currentIndex].commentsInfo as comment, index (index)}
          <Comment data={comment} />
        {/each}
      {/if}
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

  .post-container {
    display: flex;
    flex-direction: column;

    padding: 0.3rem 1rem 1rem 1rem;
    width: 34rem;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 5px 10px;
  }

  .post-header {
    display: flex;
    align-items: center;

    background: none;
    border: none;
    border-radius: 1rem;
    color: var(--text-color);
    transition: background 0.1s;

    padding: 0.3rem;
  }

  .post-user-img {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--accent);
    width: 2.5rem;
    height: 2.5rem;

    border-radius: 50%;
  }

  .post-name-container {
    display: flex;
    flex-direction: column;
    padding-left: 0.5rem;
  }

  .post-user-name {
    text-align: start;
    font-weight: bold;
    font-size: 1rem;
    color: var(--text-color);
    background: none;
    outline: none;
    border: none;
    padding: 0;
    margin: 0;
  }

  .post-user-name:hover {
    cursor: pointer;
  }

  .post-timedate {
    color: grey;
    font-size: 0.8rem;
  }

  .post-text {
    padding: 0 0.5rem 0 0.5rem;
  }

  .post-img-container {
    display: flex;
    justify-content: center;

    padding-top: 1rem;
    text-align: center;
    max-height: 20rem;
    max-width: 34rem;
  }

  .post-img {
    max-width: inherit;
    max-height: inherit;
  }

  .spacer-container {
    padding: 1rem 0;
    display: flex;
    justify-content: center;
  }

  .comments-list {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;

    padding: 0 0.5rem;
  }

  .comment-post-container {
    display: flex;
    align-items: center;
  }

  .post-comment-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1rem;
    font-weight: lighter;
  }

  .post-comment-button:hover {
    background: var(--primary-hover);
  }
</style>
