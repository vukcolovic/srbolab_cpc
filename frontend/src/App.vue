<template>
  <template v-if="!true">
    <LoginComponent></LoginComponent>
  </template>
  <template v-else>
    <div class="d-flex flex-column vh-100">
      <div class="sticky-top">
        <HeaderComponent></HeaderComponent>
      </div>
      <div class="d-flex align-items-stretch" style="height: 90%">
        <div class="col-sm-12 overflow-scroll">
          <router-view :key="$route.fullPath"></router-view>
        </div>
      </div>
      <div>
      </div>
    </div>
  </template>
</template>

<script>
import HeaderComponent from "./components/HeaderComponent.vue";
import LoginComponent from "@/components/LoginComponent";

export default {
  name: 'App',
  components: {
    // SideBar,
    LoginComponent,
    HeaderComponent,
  },
  computed : {
    // isLoggedIn : function(){
    //   return this.$store.getters.isAuthenticated;
    // }
  },
  beforeCreate() {
    // this.$store.commit('initialiseStore');
  }
}
</script>

<style>
@import './assets/css/main.css';
</style>