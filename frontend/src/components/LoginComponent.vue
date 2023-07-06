<template>
  <div class="container">
    <div class="row">
      <div class="col-5 mx-auto">
        <h1 class="mt-5">Prijava</h1>
        <hr>
        <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
          <text-input
              v-model.trim="email"
            label="Email"
            type="email"
            name="email"
            :required=true>
          </text-input>

          <text-input
              v-model.trim="password"
              label="Šifra"
              type="password"
              name="password"
              :required=true>
          </text-input>
          <hr>
          <input type="submit" class="btn btn-primary m-2" value="Prijava">
        </form-tag>
      </div>
    </div>
  </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from './../router/index.js';

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

      await axios.post('/users/login', JSON.stringify(payload)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          // console.log(error)//fixme notie
        }
        const loginData = JSON.parse(response.data.Data);

        this.$store.dispatch('setTokenAction', loginData.token);
        this.$store.dispatch('setLastNameAction', loginData.last_name);
        this.$store.dispatch('setFirstNameAction', loginData.first_name);
        this.$store.dispatch('setExaminationPlaceIdAction', loginData.examination_place_id);
        // this.$store.dispatch('setActivities', loginData.activities);
        localStorage.setItem("activities", JSON.stringify(loginData.activities))
        }, (error) => {
        this.toast.error(error.message);
      });
      await router.push("/");
    },
    }
  }
</script>