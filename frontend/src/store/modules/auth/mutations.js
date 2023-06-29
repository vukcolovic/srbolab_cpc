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
  setExaminationPlaceId(state, examinationPlaceId) {
    state.examinationPlaceId = examinationPlaceId;
  },
  setActivities(state, activities) {
    state.activities = activities;
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