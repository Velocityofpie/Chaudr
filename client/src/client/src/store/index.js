import { reactive, toRefs } from 'vue';
import axios from 'axios';
import AES from 'crypto-js/aes';
import CryptoJS from 'crypto-js';

const state = reactive({
  error: null,
  passphrase: null,
  connection: null,
  handshakeData: {},
  channels: [],
  currentUser: {},
  currentChannel: {},
  channelUsers: [],
  channelMessages: [],
  navbar: {
    channels: true,
    users: true,
  },
});

const instance = axios.create({
  baseURL: 'http://' + document.location.host,
});

export default function() {
  const setUser = (userId) => {
    try {
      state.passphrase = userId;
      getChannels();
    } catch (err) {
      state.error = err;
    }
  };

  const toggleNavbar = (nav) => {
    state.navbar = {
      channels: true,
      users: true,
      [nav]: !state.navbar[nav],
    };
  };

  const setChannelUsers = async (channelId) => {
    try {
      //const data = await instance.get('/users?channelId=' + channelId).then((res) => res.data);
      state.channelUsers = [
        'dodowater',
        'someguy',
      ];
    } catch (err) {
      state.error = err;
    }
  };

  const addMessage = (msgData) => {
    console.log(`inside addMessage`);
    if (msgData.type == 'message') {
      const bytes = AES.decrypt(msgData.message, state.currentChannel.passphrase);
      const originalText = bytes.toString(CryptoJS.enc.Utf8);
      console.log(`decrypted ${msgData.message} to ${originalText}`);
      msgData.message = originalText;
    }
    state.channelMessages = [...state.channelMessages, msgData];
  };

  const getChannels = () => {
    try {
      const jsonRooms = localStorage.getItem('rooms');
      if (jsonRooms) {
        console.log(`loaded from localstorage ${jsonRooms}`);
        state.channels = JSON.parse(jsonRooms);
      } else {
        console.log(`could not load from localstorage`);
      }
      // state.channels = [
      //   {
      //     "name": "room 1",
      //     "username": "freeguy",
      //     "id": 1121
      //   },
      //   {
      //     "name": "room 2",
      //     "username": "freeguy",
      //     "id": 2324
      //   },
      // ];
      if (!state.currentChannel?.id && state.channels.length > 0) {
        console.log('setting channel');
        setChannel(state.channels[0].id);
      }
    } catch (err) {
      state.error = err;
    }
  };

  const sendMessage = async (username, channelId, message) => {
    try {
      // connection.send(JSON.stringify({
      //   "type": "message",
      //   "username": username,
      //   "message": message
      // }));
      console.log(`trying to send message`);
      // encrypt first
      var encrypted = AES.encrypt(message, state.currentChannel.passphrase);
      console.log(`encrypted message: ${encrypted.toString()}`);
      console.log(`encrypted obj: ${encrypted}`);
      state.connection.send(JSON.stringify({
        'type': 'message',
        // 'username': state.currentChannel.username,
        'username': username,
        'message': encrypted.toString(),
      }));
    } catch (err) {
      state.error = err;
      console.log(`err while sending message: ${err}`);
    }
  };

  const setChannel = (channelId) => {
    try {
      const idx = state.channels.findIndex((channel) => channel.id == channelId);
      console.log(`idx is ${idx}`);
      state.currentChannel = state.channels[idx];
      state.channelMessages = [];
      state.currentUser = state.currentChannel.username;
      if (state.connection) {
        console.log(`closing existing connection`);
        state.connection.close(1000);
      }

      state.connection = new WebSocket('ws://' + document.location.host + '/room/connect?roomId=' + channelId + '&username=' + state.currentUser);
      state.connection.onopen = function(message) {
        console.log(`connected successfully`);
      };
      state.connection.onerror = function(err) {
        console.log(`websocket err: ${err.data}`);
      };
      state.connection.onmessage = function(message) {
        console.log(`received data in socket: ${message.data}`);
        const data = JSON.parse(message.data);

        addMessage(data);
      };
      //await setChannelUsers(state.currentChannel.id);
    } catch (err) {
      state.error = err;
    }
  };

  const genPassphrase = () => {
    const length = 15;
    var result = '';
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    for (var i = 0; i < length; i++) {
      result += characters.charAt(Math.floor(Math.random() *
        charactersLength));
    }
    return result;
  };

  const initiateHandshake = async () => {
    try {
      console.log(`initiating handshake`);
      const data = await instance
        .put('/room/handshake', {
          username: state.currentUser,
          roomId: state.currentChannel.id,
        })
        .then((res) => {
          console.log(res.data);
          console.log(`success initiating handshake: ${res.data}`);
          state.handshakeData = {
            code: res.data.code,
            roomId: state.currentChannel.id,
            passphrase: state.currentChannel.passphrase,
          };
        }).catch(err => {
          console.log(`failed initiating handshake: ${err}`);
        });

      localStorage.setItem('rooms', JSON.stringify(state.channels));
      console.log(`saved to localstorage: ${JSON.stringify(state.channels)}`);

      await getChannels(userId);
      await setChannel(data.id);
    } catch (err) {
      state.error = err;
    }
  };

  const joinRoom = async (username, roomId, code, name, passphrase) => {
    try {
      console.group(`joinRoom`);
      console.log(`username: ${username}`);
      console.log(`roomId: ${roomId}`);
      console.log(`code: ${code}`);
      console.log(`name: ${name}`);
      console.log(`passphrase: ${passphrase}`);
      const data = await instance
        .put('/room/user', {
          username: username,
          roomId: parseInt(roomId, 10),
          code: parseInt(code, 10),
        })
        .then((res) => {
          console.log(res.data);
          state.channels = [...state.channels, {
            'name': name,
            'username': username,
            'id': roomId,
            'passphrase': passphrase,
          }];
        }).catch(err => {
          console.log(`err join room: ${err}`);
        });

      localStorage.setItem('rooms', JSON.stringify(state.channels));
      console.log(`saved to localstorage: ${JSON.stringify(state.channels)}`);

      getChannels();
    } catch (err) {
      state.error = err;
    }
  };

  const createChannel = async (username, name) => {
    try {
      const data = await instance
        .put('/room', {
          username: username,
        })
        .then((res) => {
          console.log(res.data);
          state.channels = [...state.channels, {
            'name': name,
            'username': username,
            'id': res.data.roomId,
            'passphrase': genPassphrase(),
          }];
        });

      localStorage.setItem('rooms', JSON.stringify(state.channels));
      console.log(`saved to localstorage: ${JSON.stringify(state.channels)}`);

      await getChannels(userId);
      await setChannel(data.id);
    } catch (err) {
      state.error = err;
    }
  };

  return {
    // States
    ...toRefs(state),

    // Actions
    getChannels,
    setChannel,
    setChannelUsers,
    setUser,
    sendMessage,
    createChannel,
    toggleNavbar,
    addMessage,
    initiateHandshake,
    joinRoom,
  };
}
