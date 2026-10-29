<script>
  import { BackendHost, WSBackendHost } from "./main";
  import { user, feedData } from "./lib/user/user";
  import { getData } from "./lib/util/postAndGet";
  import { Router, Route, navigate } from "svelte-routing";
  import Login from "./lib/Login.svelte";
  import Register from "./lib/Register.svelte";
  import Main from "./lib/Main.svelte";
  import Profile from "./lib/Profile.svelte";
  import Group from "./lib/Group.svelte";
  import Groups from "./lib/Groups.svelte";
  import Notifications from "./lib/Notifications.svelte";
  import Chat from "./lib/Chat.svelte";
  import GroupChat from "./lib/GroupChat.svelte";
  import Loading from "./lib/Loading.svelte";
  export let conn;
  let feedPostData=[];

  const url = "/";

  const toggleDarkMode = () => {
    if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
      document.body.classList.toggle("light-theme");
    } else {
      document.body.classList.toggle("dark-theme");
    }
  };

  const navToLogin = () => {
    navigate("/login", { replace: true });
  };

  async function getMainData() {
    getData(`${BackendHost}/feed`).then((data) => {
      $feedData = data;
      feedPostData = $feedData.allPosts.reverse();
      $user.userJoinedGroups = data.userJoinedGroups;
    });
  }

  export function addToChatListNotificationArray(interactingUserID) {
    const index = $user.chatListNotificationArray.indexOf(interactingUserID);
    if (index > -1) {
      $user.chatListNotificationArray.splice(index, 1);
    }
    $user.chatListNotificationArray.push(interactingUserID);
    updateChatUserList();
  }

  export function removeFromChatListNotificationArray(interactingUserID) {
    const index = $user.chatListNotificationArray.indexOf(interactingUserID);
    if (index > -1) {
      $user.chatListNotificationArray.splice(index, 1);
    }
    updateChatUserList();
  }

  export function updateChatListOrderArray(interactingUserID) {
    if ($user.chatListOrderArray[0] === interactingUserID) {
      return;
    }
    const index = $user.chatListOrderArray.indexOf(interactingUserID);
    if (index > -1) {
      $user.chatListOrderArray.splice(index, 1);
    }
    $user.chatListOrderArray.unshift(interactingUserID);
    updateChatUserList();
    sendEvent("chat_list_order_array", $user.chatListOrderArray);
  }

  export function updateChatUserList() {
    let chatList = [];
    let onlineClientList = [];
    let offlineClientList = [];
    let failed;
    $user.userListFromEvent.sort(function (a, b) {
      if (a.firstname < b.firstname) {
        return -1;
      }
      if (b.firstname < a.firstname) {
        return 1;
      }
      return 0;
    });
    if ($user.chatListOrderArray.length !== 0) {
      for (let val of $user.chatListOrderArray) {
        for (let userFromList of $user.userListFromEvent) {
          if (
            userFromList.userID === val &&
            userFromList.userID !== $user.currentUserID
          ) {
            chatList.push(userFromList);
            break;
          }
        }
      }
    }

    for (let userFromList of $user.userListFromEvent) {
      for (let val of $user.chatListOrderArray) {
        if (userFromList.userID === val) {
          failed = true;
          break;
        }
      }
      if (userFromList.userID === $user.currentUserID) {
        failed = true;
      }
      if (!failed) {
        if (userFromList.onlineStatus) {
          onlineClientList.push(userFromList);
        } else {
          offlineClientList.push(userFromList);
        }
      }
      failed = false;
    }

    chatList = chatList.concat(onlineClientList.concat(offlineClientList));

    for (let userFromList of chatList) {
      userFromList.notification = false;
      for (let val of $user.chatListNotificationArray) {
        if (
          userFromList.userID !== $user.recipientID &&
          userFromList.userID === val
        ) {
          userFromList.notification = true;
          break;
        }
      }
    }
    $user.currentChatList = chatList;
  }

  let messagesFromEnd = 0;

  class SendMessageEvent {
    constructor(message, sender, recipient) {
      this.message = message;
      this.senderID = sender;
      this.recipientID = recipient;
      this.date = getMessageSentTime();
    }
  }

  class SendGroupMessageEvent {
    constructor(groupID, message, senderID, senderName, recipient) {
      this.groupID = groupID;
      this.message = message;
      this.senderID = senderID;
      this.senderName = senderName;
      this.date = getMessageSentTime();
    }
  }
  var allMessagesInCurrentChatContainer = [];

  export async function getMessageHistory(receiverID) {
    let messages;
    let addToStart = false;
    await getData(
      "http://localhost:8080/message?senderID=" +
        $user.currentUserID +
        "&receiverID=" +
        receiverID,
    ).then((value) => {
      if (value === null) {
        allMessagesInCurrentChatContainer = [];
      } else {
        allMessagesInCurrentChatContainer = value;
      }

      $user.currentChatDetails = allMessagesInCurrentChatContainer;
    });
  }

  export function sendChatMessage(e, inputText) {
    e.preventDefault();
    if (inputText != "" && inputText != null) {
      let outgoingEvent = new SendMessageEvent(
        inputText,
        $user.currentUserID,
        $user.recipientID,
      );
      updateChatListOrderArray($user.recipientID);
      sendEvent("send_message", outgoingEvent);
    }
  }

  export function sendGroupChatMessage(e, inputText, senderName) {
    e.preventDefault();
    if (inputText != "" && inputText != null) {
      let outgoingEvent = new SendGroupMessageEvent(
        $user.currentActiveGroupID,
        inputText,
        $user.currentUserID,
        senderName,
      );
      sendEvent("send_group_message", outgoingEvent);
    }
  }

  function getMessageSentTime() {
    const currentDate = new Date();

    const datetime = new Intl.DateTimeFormat("en-GB", {
      dateStyle: "short",
      timeStyle: "short",
    });
    return datetime.format(currentDate);
  }

  export function startWS() {
    if (window["WebSocket"]) {
      console.log("Browser supports websockets");
      conn = new WebSocket(WSBackendHost + "/ws");
      console.log("Attempting Connection...");
      conn.addEventListener("open", () => {
        console.log("WS Successfully Connected");
        fetchJoinedGroupIDArray($user.currentUserID);
      });
      conn.onclose = () => {
        console.log("Connection closed");
      };
      // Add a listener to the onmessage event
      conn.onmessage = function (evt) {
        console.log(evt);
        // parse websocket message as JSON
        const eventData = JSON.parse(evt.data);
        // Assign JSON data to new Event Object
        const event = Object.assign(new Event(), eventData);
        // Let router manage message
        routeEvent(event);
      };
    } else {
      var item = document.createElement("div");
      item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    }
  }

  export function closeWs(conn) {
    conn.close();
  }

  export function routeEvent(event) {
    if (event.type === undefined) {
      alert("no 'type' field in event");
    }
    switch (event.type) {
      case "new_notification":
        // Route Notification
        routeNotification(event.payload);
        break;
      case "send_joined_groupID_array":
        $user.JoinedGroupIDArray = event.payload;
        break;
      case "new_group_message":
        if (event.payload.groupID === $user.currentActiveGroupID) {
          let tempMessage = {
            content: event.payload.message,
            date: event.payload.date,
            GroupID: event.payload.groupID,
            senderID: event.payload.senderID,
            senderName: event.payload.senderName,
            messageID: null,
          };
          $user.currentGroupChatDetails = [
            ...$user.currentGroupChatDetails,
            tempMessage,
          ];
          setTimeout(() => {
            document.querySelector(".chat-text-area").scrollTop =
              document.querySelector(".chat-text-area").scrollHeight;
          }, 100);
        }
        break;
      case "new_message":
        addToChatListNotificationArray(event.payload.senderID);
        updateChatListOrderArray(event.payload.senderID);
        if (event.payload.senderID === $user.recipientID) {
          let tempMessage = {
            content: event.payload.message,
            date: event.payload.date,
            receiverID: event.payload.recipientID,
            senderID: event.payload.senderID,
            messageID: null,
          };
          $user.currentChatDetails = [...$user.currentChatDetails, tempMessage];
        }
        setTimeout(() => {
          document.querySelector(".chat-text-area").scrollTop =
            document.querySelector(".chat-text-area").scrollHeight;
        }, 100);
        break;
      case "clients_online_status":
        if (event.payload === null) {
          $user.userListFromEvent = [];
          return;
        }
        for (let val of event.payload) {
          val.notification = false;
        }
        $user.userListFromEvent = event.payload;
        updateChatUserList();
        break;
      default:
        alert(event.type + " unsupported message type");
        break;
    }
  }

  export function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload);
    conn.send(JSON.stringify(event));
  }

  export class Event {
    constructor(type, payload) {
      this.type = type;
      this.payload = payload;
    }
  }

  export function updateActiveGroupInBackend(currentActiveGroup) {
    sendEvent("update_active_group", currentActiveGroup);
  }

  export function updateActiveChatInBackend(currentActiveChat) {
    sendEvent("update_active_chat", currentActiveChat);
  }

  export function fetchJoinedGroupIDArray(userID) {
    sendEvent("get_joined_groupID_array", userID);
  }

  export function routeNotification(payload) {
    if (payload.subtype === undefined) {
      alert("no 'subtype' field in notification event");
    }
    switch (payload.subtype) {
      case "request_to_join_group":
        formatRequestToJoinGroup(payload.body);
        break;
      case "follow_request":
        formatFollowRequest(payload.body);
        break;
      case "invitation_to_join_group":
        formatGroupInvitation(payload.body);
        break;
      case "group_event_notification":
        formatGroupEventNotification(payload.body);
        break;
    }
  }

  export async function fetchAllNotifications() {
    $user.allNotificationArray = [];
    getData(`${BackendHost}/notifications?userID=${$user.currentUserID}`).then(
      (data) => {
        data.forEach((notification) => {
          routeNotification(notification);
        });
        navigate("/notifications", { replace: true });
      },
    );
  }

  export function deleteNotificationFromDB(notificationID) {
    sendEvent("delete_notification", notificationID);
  }

  export function formatRequestToJoinGroup(body) {
    let arr = body.content.split("|");
    let username = arr[0];
    let groupName = arr[1];
    let groupID = arr[2];
    let joinerID = arr[3];
    let requestToJoinGroup = {
      type: "groupRequest",
      body: {
        notificationID: body.notificationID,
        username: username,
        groupName: groupName,
        groupID: groupID,
        joinerID: joinerID,
      },
    };
    $user.allNotificationArray.push(requestToJoinGroup);
  }

  export function formatFollowRequest(body) {
    let arr = body.content.split("|");
    let followerName = arr[0];
    let followRequest = {
      type: "followRequest",
      body: {
        notificationID: body.notificationID,
        groupName: "",
        username: followerName,
        senderID: body.senderID,
        recieverID: body.recieverID,
      },
    };
    $user.allNotificationArray.push(followRequest);
  }

  export function formatGroupInvitation(body) {
    let arr = body.content.split("|");
    let groupName = arr[0];
    let groupID = arr[1];
    let groupInvitation = {
      type: "groupInvite",
      body: {
        notificationID: body.notificationID,
        groupName: groupName,
        username: "",
        groupID: groupID,
      },
    };
    $user.allNotificationArray.push(groupInvitation);
  }

  export function formatGroupEventNotification(body) {
    let arr = body.content.split("|");
    let groupName = arr[0];
    let groupID = arr[1];
    let groupEventNotification = {
      type: "eventNotify",
      body: {
        notificationID: body.notificationID,
        groupName: groupName,
        username: "",
        groupID: groupID,
      },
    };
    $user.allNotificationArray.push(groupEventNotification);
  }

  async function login(usernameAndPasswordObj, isRegister) {
    fetch(`${BackendHost}/login`, {
      method: "POST",
      mode: "cors",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(usernameAndPasswordObj),
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw "unauthorized";
        }
      })
      .then(async (resp) => {
        if (resp.isLoggedIn) {
          console.log("correct password, starting WS");
          startWS();
          $user.currentUserID = resp.userID;
          let currentNickname = resp.nickname;
          $user.currentUsername = resp.firstName + " " + resp.lastName;
          let msg = resp.msg;
          console.log("msg: ", msg);
          conn.onopen = () => {
            if (isRegister) {
              sendEvent("update_all_users_from_DB");
              sendEvent("get_joined_groupID_array", $user.currentUserID);
            }
          };
          setTimeout(() => {
            navigate("/feed", { replace: false });
          }, 250);
          return;
        } else {
          alert(resp.msg);
        }
      })
      .catch((e) => {
        console.log(e);
      });
    return false;
  }

  async function cookieLogin(resp) {
    console.log("active session, starting WS");
    startWS();
    $user.currentUserID = resp.userID;
    let currentNickname = resp.nickname;
    $user.currentUsername = resp.firstName + " " + resp.lastName;
    let msg = resp.msg;
    console.log("msg: ", msg);
    conn.onopen = () => {
      sendEvent("update_all_users_from_DB");
      sendEvent("get_joined_groupID_array", $user.currentUserID);
    };
    setTimeout(() => {
      navigate("/feed", { replace: false });
    }, 250);
    return;
  }

  const navToNotifications = () => {
    fetchAllNotifications();
  };

  export function joinGroup(groupID, userID) {
    sendEvent("join_group", groupID + "|" + userID);
  }
  export function leaveGroup(groupID, userID) {
    sendEvent("leave_group", groupID + "|" + userID);
  }

  export function updateActiveGroup(groupID, groupName) {
    $user.currentActiveGroupID = groupID;
    $user.currentActiveGroupName = groupName;
  }

  export function inviteToGroup(groupID, userID) {
    sendEvent("invite_to_group", groupID + "|" + userID);
  }

  export function requestToJoinGroup(groupID, userID) {
    sendEvent("request_to_join_group", groupID + "|" + userID);
  }

  export function sendNotificationOfEvent(groupID) {
    sendEvent("notify_of_group_event", groupID);
  }

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

<svelte:window on:load={navToLogin} />

<Router {url}>
  <Route
    path="/login"
    component={Login}
    {toggleDarkMode}
    {login}
    {cookieLogin}
  />
  <Route path="/register" component={Register} {toggleDarkMode} {login} />
  <Route
    path="/feed"
    component={Main}
    {toggleDarkMode}
    {navToNotifications}
    {getMessageHistory}
    {updateActiveGroup}
    {removeFromChatListNotificationArray}
    {getMainData}
    {feedPostData}
    {navToGroup}
  />
  <Route
    path="/profile"
    component={Profile}
    {toggleDarkMode}
    {navToNotifications}
    {getMessageHistory}
    {removeFromChatListNotificationArray}
  />
  <Route
    path="/group"
    component={Group}
    {toggleDarkMode}
    {navToNotifications}
    {getMessageHistory}
    {updateActiveGroupInBackend}
    {sendNotificationOfEvent}
    {removeFromChatListNotificationArray}
    {inviteToGroup}
  />
  <Route
    path="/groups"
    component={Groups}
    {toggleDarkMode}
    {navToNotifications}
    {updateActiveGroup}
    {requestToJoinGroup}
    {getMessageHistory}
    {inviteToGroup}
  />
  <Route
    path="/notifications"
    component={Notifications}
    {toggleDarkMode}
    {navToNotifications}
    {joinGroup}
    {deleteNotificationFromDB}
    {getMessageHistory}
    {navToGroup}
  />
  <Route
    path="/chat"
    component={Chat}
    {toggleDarkMode}
    {navToNotifications}
    {sendChatMessage}
    {getMessageHistory}
    {removeFromChatListNotificationArray}
  />
  <Route
    path="/groupchat"
    component={GroupChat}
    {toggleDarkMode}
    {navToNotifications}
    {sendChatMessage}
    {getMessageHistory}
    {sendGroupChatMessage}
    {removeFromChatListNotificationArray}
  />
  <Route
    path="/loading/:path"
    component={Loading}
    {toggleDarkMode}
    {navToNotifications}
  />
</Router>
