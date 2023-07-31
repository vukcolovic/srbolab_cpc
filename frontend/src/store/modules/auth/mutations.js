export default {
  setToken(state, newToken) {
    state.token = newToken;
  },
  setFirstName(state, firstName) {
    state.firstName = firstName;
  },
  setLastName(state, lastName) {
    state.lastName = lastName;
  },
  setRoles(state, roles) {
    state.roles = roles;
  },
  initialiseStore(state) {
    // Check if the token exists
    if(localStorage.getItem('token')) {
      this.replaceState(
          Object.assign(state, JSON.parse(localStorage.getItem('token')))
      );
    }
  }
};