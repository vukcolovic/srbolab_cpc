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
    setRoless({commit}, roles) {
        commit('roles', roles);
    }
};