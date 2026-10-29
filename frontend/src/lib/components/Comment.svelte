<script>
  import { navigate } from "svelte-routing";
  import { profileData } from "../user/user";
  import { getData } from "../util/postAndGet";
  import { BackendHost } from "../../main";
  export let data;
  const commentTypeURL = window.location.href;
  const commentType = /[\w]+$/.exec(commentTypeURL);

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

  const navToProfile = (id) => {
    getUserData(id);
  };
</script>

<div class="comment-container">
  <div class="comment-header">
    <div class="comment-user-avatar">
      <img
        class="profile-image"
        src={`${BackendHost}/profile/image?userID=${
          data.userID
        }&timestamp=${Date.now()}`}
        alt="profileImage"
      />
    </div>
    <div class="comment-name-container">
      <button
        class="comment-user-name"
        on:click={() => navToProfile(data.userID)}>{data.fullname}</button
      >
      <div class="comment-timedate">{data.created}</div>
    </div>
  </div>
  <div class="comment-text">
    {data.content}
  </div>
  {#if data.image != undefined && data.image.name != "default_image.png" && data.image.name != ""}
    <div class="comment-img-container">
      <img
        class="comment-img"
        alt="comment-img"
        src={`${BackendHost}/commentImage?commentID=${data.commentID}&postID=${data.postID}&commentType=${commentType[0]}`}
      />
    </div>
  {/if}
</div>

<style>
  .profile-image {
    width: inherit;
    height: inherit;
    border-radius: inherit;
    object-fit: cover;
  }

  .comment-container {
    display: flex;
    flex-direction: column;
    width: 100%;
    padding: 0.2rem 0.5rem;

    background: var(--secondary-background-hover);
    border-radius: 1rem;
  }

  .comment-header {
    display: flex;
    align-items: center;

    background: none;
    border: none;
    border-radius: 1rem;
    color: var(--text-color);
    transition: background 0.1s;

    padding: 0.3rem;
  }

  .comment-user-avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--accent);
    width: 1.8rem;
    height: 1.8rem;

    border-radius: 50%;
  }

  .comment-name-container {
    display: flex;
    flex-direction: column;
    padding-left: 0.5rem;
  }

  .comment-user-name {
    text-align: start;
    font-weight: light;
    font-size: 1rem;
    color: var(--text-color);
    background: none;
    outline: none;
    border: none;
    padding: 0;
    margin: 0;
  }

  .comment-user-name:hover {
    cursor: pointer;
  }

  .comment-timedate {
    color: grey;
    font-size: 0.8rem;
  }

  .comment-text {
    padding: 0.3rem;
  }

  .comment-img-container {
    padding-top: 1rem;
    text-align: center;
    max-height: 16rem;
    max-width: 100%;
  }

  .comment-img {
    max-height: inherit;
  }
</style>
