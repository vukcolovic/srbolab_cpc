import mutations from './mutations.js';
import getters from './getters.js';
import actions from './actions.js';

export default {
  state() {
    return {
      token: "",
      firstName: "",
      lastName: "",
      examinationPlaceId: 0,
      activities: [],
    };
  },
  mutations,
  getters,
  actions
};