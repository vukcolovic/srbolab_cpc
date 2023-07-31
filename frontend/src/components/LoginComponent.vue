<template>
  <div class="container">
    <img src="../assets/srbolab_logo.png" href="">
      <div class="col-sm-5 mx-auto mt-5">
        <hr>
        <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
          <text-input
            v-model.trim="email"
            label="Email"
            type="email"
            style="color: #007bff; font-size: 1.4em"
            name="email"
            :required=true>
          </text-input>
          <br>
          <text-input
              v-model.trim="password"
              label="Šifra"
              type="password"
              style="color: #007bff; font-size: 1.4em"
              name="password"
              :required=true>
          </text-input>
          <input type="submit" class="btn btn-primary mt-3" value="Prijava">
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

export default {
  name: 'LoginComponent',
  components: {FormTag, TextInput},
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
        this.toast.error(error);
      });
      await router.push("/");
    },
    },
  setup() {
    const toast = useToast();
    return {toast}
  },
  }
</script>
