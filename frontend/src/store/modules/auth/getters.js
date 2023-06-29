export default {
  token(state) {
    return state.token;
  },
  isAuthenticated(state) {
    return !!state.token;
  },
  firstName(state) {
    return state.firstName;
  },
  lastName(state) {
    return state.lastName;
  },
  examinationPlaceId(state) {
    return state.examinationPlaceId;
  },
  activities(state) {
    return state.activities;
  },
};