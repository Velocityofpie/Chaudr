<template>
  <div
    class="modal fade"
    id="exampleModal"
    tabindex="-1"
    aria-labelledby="exampleModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="exampleModalLabel">Create New Channel</h5>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body">
          <form>
            <h1>Create room</h1>
            <div class="form-group">
              <label for='createRoomUsername'>Username for Room:</label>
              <input
                type="text"
                id='createRoomUsername'
                class="form-control"
                placeholder="Username..."
              />
              <label for='createRoomName'>Room name:</label>
              <input
                type="text"
                id='createRoomName'
                class="form-control"
                placeholder="Channel name..."
              />
            </div>
          </form>
          <button type="button" class="btn btn-save my-2" data-bs-dismiss="modal" @click="onSubmit">
            Create Room
          </button>
          <form>
            <h1>Join a room</h1>
            <div class="form-group">
              <label for='joinRoomUsername'>Username for Room:</label>
              <input
                type="text"
                id="joinRoomUsername"
                class="form-control"
                placeholder="Username..."
              />
              <label for='joinRoomName'>Room name:</label>
              <input
                type="text"
                id="joinRoomName"
                class="form-control"
                placeholder="Room name..."
              />
              <label for='joinRoomId'>Room Id:</label>
              <input
                type="number"
                id="joinRoomId"
                class="form-control"
                placeholder="Room Id..."
              />
              <label for='joinRoomPassphrase'>Room Passphrase:</label>
              <input
                type="text"
                id="joinRoomPassphrase"
                class="form-control"
                placeholder="Passphrase..."
              />
              <label for='joinRoomCode'>Room Code:</label>
              <input
                type="number"
                id="joinRoomCode"
                class="form-control"
                placeholder="Room code..."
              />
            </div>
          </form>
          <button type="button" class="btn btn-save my-2" data-bs-dismiss="modal" @click="joinRoomModal">
            Join Room
          </button>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { inject } from '@vue/runtime-core';

export default {
  setup() {
    const store = inject('store');
    const { createChannel, toggleNavbar, joinRoom } = store();

    const onSubmit = function () {
      const roomName = document.getElementById('createRoomName').value;
      const username = document.getElementById('createRoomUsername').value;
      createChannel(username, roomName);
      toggleNavbar('channels');
    };

    const joinRoomModal = function () {
      const roomCode = document.getElementById('joinRoomCode').value;
      const roomId = document.getElementById('joinRoomId').value;
      const roomName = document.getElementById('joinRoomName').value;
      const roomPassphrase = document.getElementById('joinRoomPassphrase').value;
      const username = document.getElementById('joinRoomUsername').value;
      joinRoom(username, roomId, roomCode, roomName, roomPassphrase);
      toggleNavbar('channels');
    };

    return {
      onSubmit,
      joinRoomModal
    };
  },
};
</script>

<style scoped>
.btn-save {
  background: rgba(106, 0, 255, 1);
  color: var(--white1);
}

.btn-save:hover {
  filter: brightness(1.2);
}
</style>
