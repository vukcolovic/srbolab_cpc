<template>
  <template v-if="!isLoggedIn && isCorporated">
    <LoginComponent></LoginComponent>
  </template>
  <template v-else>
    <div class="d-flex flex-column vh-100">
      <div class="sticky-top">
        <HeaderComponent v-if="isCorporated"></HeaderComponent>
      </div>
      <div class="d-flex align-items-stretch" :style="{height:isCorporated ? '90%' : '100%'}">
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
    isLoggedIn : function(){
      return this.$store.getters.isAuthenticated;
    },
    isCorporated : function(){
      return this.$store.getters.isCorporate;
    }
  },
  beforeCreate() {
    this.$store.commit('initialiseStore');
  }
}
</script>

<style>
@import './assets/css/main.css';
</style>