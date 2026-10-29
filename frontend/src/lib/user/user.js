import { writable } from "svelte/store";

export const profileData = writable({
    publicBool: false,
    followedList: [],
    followerList: [],
    profilePosts: [],
    userData: {}
})

export const user = writable({
    currentUserID: undefined,
    currentUsername: "",

    currentUserFollowers: [],
    currentUserFollowed: [],

    recipientFollowers: [],
    recipientFollowed: [],

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

    allNotificationArray: []
})

export const feedData = writable({
    allEvents: [],
    allPosts: [],
    userJoinedGroups: []
})