<script>
  import { postData } from "../util/postAndGet";
  import { BackendHost } from "../../main";
  import { user } from "../user/user";
  export let data;
  export let currentIndex;

  const addEventRating = (eventID, ratingType, userID) => {
    let sendData = {
      ratingType: ratingType,
      eventID: eventID,
      userID: userID,
    };
    postData(`${BackendHost}/rating`, sendData).then((resp) => {
      data[currentIndex].going = resp.going;
      data[currentIndex].notgoing = resp.notGoing;
    });
  };
</script>

<div class="event-container">
  <div class="event-header">{data[currentIndex].title}</div>
  <div class="event-text">
    {data[currentIndex].content}
  </div>
  <div class="event-date">{data[currentIndex].date}</div>
  <div class="event-buttons-container">
    <button
      class="event-button"
      on:click={() =>
        addEventRating(
          data[currentIndex].eventID,
          "notGoing",
          $user.currentUserID,
        )}>{data[currentIndex].notgoing} Not Going</button
    >
    <button
      class="event-button"
      on:click={() =>
        addEventRating(
          data[currentIndex].eventID,
          "going",
          $user.currentUserID,
        )}>{data[currentIndex].going} Going</button
    >
  </div>
</div>

<style>
  .event-container {
    display: flex;
    flex-direction: column;

    padding: 0.3rem 1rem 1rem 1rem;
    width: 34rem;

    background: var(--secondary-background);
    border-radius: 1rem;
    box-shadow: hsl(0, 0%, 0%, 30%) 0px 5px 10px;
  }

  .event-header {
    padding: 0.5rem;
    font-weight: bold;
    font-size: 1.2rem;
  }

  .event-text {
    padding: 0 0.5rem 0 0.5rem;
  }

  .event-date {
    color: grey;
    padding: 0 0.5rem 0 0.5rem;
  }

  .event-buttons-container {
    display: flex;
    justify-content: end;
    gap: 0.5rem;

    padding-top: 0.5rem;
  }

  .event-button {
    background: var(--primary);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    transition: background 0.2s;

    font-size: 1rem;
    font-weight: lighter;
  }

  .event-button:hover {
    background: var(--primary-hover);
  }
</style>
