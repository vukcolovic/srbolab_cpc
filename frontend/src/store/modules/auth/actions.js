export default {
    setTokenAction({commit}, token) {
        commit('setToken', token);
    },
    setIsCorporateAction({commit}, isCorporate) {
        commit('setIsCorporate', isCorporate);
    },
    setFirstNameAction({commit}, firstName) {
        commit('setFirstName', firstName);
    },
    setLastNameAction({commit}, lastName) {
        commit('setLastName', lastName);
    },
    setRoles({commit}, roles) {
        commit('roles', roles);
    }
};