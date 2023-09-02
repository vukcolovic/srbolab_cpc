<template>
  <div class="container">
    <img href="" src="../assets/srbolab_logo.png">
    <div class="col-sm-5 mx-auto mt-5">
      <hr>
      <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
        <text-input
            v-model.trim="email"
            :required=true
            label="Email"
            name="email"
            style="color: #007bff; font-size: 1.4em"
            type="email">
        </text-input>
        <br>
        <text-input
            v-model.trim="password"
            :required=true
            label="Šifra"
            name="password"
            style="color: #007bff; font-size: 1.4em"
            type="password">
        </text-input>
        <input class="btn btn-primary mt-3" type="submit" value="Prijava">
      </form-tag>
    </div>
  </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from "@/router";
import {useToast} from "vue-toastification";
import {apiMixin} from "@/mixins/apiMixin";
import {commonMixin} from "@/mixins/commonMixin";

export default {
  name: 'LoginComponent',
  components: {FormTag, TextInput},
  mixins: [apiMixin, commonMixin],
  data() {
    return {
      email: "",
      password: "",
      errorMsg: "",
    }
  },
  methods: {
    async submitHandler() {
      const payload = {
        email: this.email,
        password: this.password,
      }
      await axios.post('/login', JSON.stringify(payload)).then((response) => {
        if (response.data == null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        const loginData = JSON.parse(response.data.Data);

        this.$store.dispatch('setTokenAction', loginData.token);
        this.$store.dispatch('setLastNameAction', loginData.last_name);
        this.$store.dispatch('setFirstNameAction', loginData.first_name);
        localStorage.setItem("roles", JSON.stringify(loginData.roles))
      }, (error) => {
        this.errorToast(error, "/login");
      });
      await router.push("/home");
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  mounted() {
    this.isCorporateIp();
  }
}
</script>
