export default {
  token(state) {
    return state.token;
  },
  isAuthenticated(state) {
    return !!state.token;
  },
  isCorporate(state) {
    return !!state.isCorporate;
  },
  firstName(state) {
    return state.firstName;
  },
  lastName(state) {
    return state.lastName;
  },
  roles(state) {
    return state.roles;
  },
};