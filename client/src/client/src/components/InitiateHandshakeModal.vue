<template>
  <div
    class='modal fade'
    id='handshakeModal'
    tabindex='-1'
    aria-labelledby='exampleModalLabel'
    aria-hidden='true'
  >
    <div class='modal-dialog'>
      <div class='modal-content'>
        <div class='modal-header'>
          <h5 class='modal-title' id='exampleModalLabel'>Have another user join</h5>
          <button
            type='button'
            class='btn-close'
            data-bs-dismiss='modal'
            aria-label='Close'
          ></button>
        </div>
        <div class='modal-body'>
          <template v-if='handshakeData'>
            <h4>Code: {{ handshakeData.code }}</h4>
            <h4>Room Id: {{ handshakeData.roomId }}</h4>
            <h4>Passphrase: {{ handshakeData.passphrase }}</h4>
          </template>
          <template v-else>
            <p>UH OH! Something happened in the client that borked it. Sorry</p>
          </template>
        </div>
        <div class='modal-footer'>
          <button type='button' class='btn btn-secondary' data-bs-dismiss='modal'>Close</button>
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
    const { createChannel, toggleNavbar, handshakeData } = store();

    return {
      handshakeData,
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
