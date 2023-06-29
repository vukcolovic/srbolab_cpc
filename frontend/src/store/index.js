import { createStore } from 'vuex';
import authModule from './modules/auth/index.js';

const store = createStore({
  modules: {
    auth: authModule
  },
});

store.subscribe((mutation, state) => {
  // Store the state object as a JSON string
  localStorage.setItem('token', JSON.stringify(state));
});

export default store;