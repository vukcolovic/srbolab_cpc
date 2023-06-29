export default {
    setTokenAction({commit}, token) {
        commit('setToken', token);
    },
    setFirstNameAction({commit}, firstName) {
        commit('setFirstName', firstName);
    },
    setLastNameAction({commit}, lastName) {
        commit('setLastName', lastName);
    },
    setExaminationPlaceIdAction({commit}, examinationPlaceId) {
        commit('setExaminationPlaceId', examinationPlaceId);
    },
    setActivities({commit}, activities) {
        commit('activities', activities);
    }
};