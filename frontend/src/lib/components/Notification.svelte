<script>
  import { navigate } from "svelte-routing";
  import { user, profileData } from "../user/user";
  import { getData, postData } from "../util/postAndGet";
  import { BackendHost } from "../../main";
  export let notif;
  export let joinGroup;
  export let deleteNotificationFromDB;

  const followNotification = (id, bool) => {
    if (bool) {
      let data = {
        NotificationID: id,
        senderID: notif.body.senderID,
        recieverID: notif.body.recieverID,
      };

      postData(`${BackendHost}/notifications`, data);
    }

    let indexToSplice = 0;
    let tempArray = [...$user.allNotificationArray];

    for (let key in tempArray) {
      if (tempArray[key].body.notificationID == id) {
        indexToSplice = Number(key);
        break;
      }
    }

    tempArray.splice(indexToSplice, 1);
    $user.allNotificationArray = [...tempArray];
    deleteNotificationFromDB(id);
  };

  const groupNotifyRemoveRequest = (id, bool, userID) => {
    let indexToSplice = 0;
    if (bool) {
      joinGroup(notif.body.groupID, userID);
    }
    let tempArray = [...$user.allNotificationArray];

    for (let key in tempArray) {
      if (tempArray[key].body.notificationID == id) {
        indexToSplice = Number(key);
        break;
      }
    }

    tempArray.splice(indexToSplice, 1);
    $user.allNotificationArray = [...tempArray];
    deleteNotificationFromDB(id);
  };

  function getGroupData(groupID) {
    getData(`${BackendHost}/group?groupID=${groupID}`).then((data) => {
      $user.currentActiveGroupMembersDetails = data;
    });
  }

  const navToGroup = (groupID, groupName) => {
    getGroupData(groupID);
    setTimeout(() => {
      $user.currentActiveGroupID = groupID;
      $user.currentActiveGroupName = groupName;
      $user.recipientName = groupName;
      navigate("/group", { replace: true });
    }, 200);
  };
</script>

<div class="notification">
  {#if notif.type == "followRequest"}
    <div class="notification-text">
      <b>{notif.body.username}</b> wants to follow you.
    </div>
    <div class="notification-buttons-container">
      <button
        class="notification-button"
        on:click={() => followNotification(notif.body.notificationID, false)}
        >Refuse</button
      >
      <button
        class="notification-button"
        on:click={() => followNotification(notif.body.notificationID, true)}
        >Accept</button
      >
    </div>
  {:else if notif.type == "groupInvite"}
    <div class="notification-text">
      You were invited to join <b>{notif.body.groupName}</b>.
    </div>
    <div class="notification-buttons-container">
      <button
        class="notification-button"
        on:click={() =>
          groupNotifyRemoveRequest(
            notif.body.notificationID,
            false,
            $user.currentUserID,
          )}>Refuse</button
      >
      <button
        class="notification-button"
        on:click={() =>
          groupNotifyRemoveRequest(
            notif.body.notificationID,
            true,
            $user.currentUserID,
          )}>Accept</button
      >
    </div>
  {:else if notif.type == "groupRequest"}
    <div class="notification-text">
      <b>{notif.body.username}</b> requested to join
      <b>{notif.body.groupName}</b>.
    </div>
    <div class="notification-buttons-container">
      <button
        class="notification-button"
        on:click={() =>
          groupNotifyRemoveRequest(
            notif.body.notificationID,
            false,
            parseInt(notif.body.joinerID),
          )}>Refuse</button
      >
      <button
        class="notification-button"
        on:click={() =>
          groupNotifyRemoveRequest(
            notif.body.notificationID,
            true,
            parseInt(notif.body.joinerID),
          )}>Accept</button
      >
    </div>
  {:else if notif.type == "eventNotify"}
    <div class="notification-text">
      <b>{notif.body.groupName}</b> is hosting an event.
    </div>
    <div class="notification-buttons-container">
      <button
        class="notification-button"
        on:click={() => deleteNotificationFromDB(notif.body.notificationID)}
        on:click={() => navToGroup(notif.body.groupID, notif.body.groupName)}
        >Go to Group</button
      >
    </div>
  {/if}
</div>

<style>
  .notification {
    display: flex;
    flex-direction: column;
    align-items: start;

    padding: 1rem;
    width: 34rem;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 5px 10px;
  }

  .notification-buttons-container {
    width: 100%;
    padding-top: 0.5rem;

    display: flex;
    justify-content: end;
    gap: 0.5rem;
  }

  .notification-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 2rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1rem;
    font-weight: lighter;
  }

  .notification-button:hover {
    background: var(--primary-hover);
  }
</style>
