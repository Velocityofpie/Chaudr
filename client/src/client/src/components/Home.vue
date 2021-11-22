<template>
  <div class="jumbotron vh-100 d-flex flex-column">
    <Navbar />
    <div class="main row h-100 m-0">
      <Channels />
      <ChatScreen />
      <Users />
    </div>
  </div>
</template>

<script>
import Navbar from './Navbar.vue';
import Channels from './Channels.vue';
import ChatScreen from './ChatScreen.vue';
import Users from './Users.vue';

export default {
  components: {
    Channels,
    ChatScreen,
    Users,
    Navbar,
  },
};
</script>
