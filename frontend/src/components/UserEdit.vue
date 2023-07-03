<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 v-if="action === 'add'" class="mt-2">Dodavanje</h3>
        <h3 v-if="action === 'view'" class="mt-2">Pregled</h3>
        <h3 v-if="action === 'update'" class="mt-2">Ažuriranje</h3>
        <hr>
      </div>
    </div>
    <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
      <div class="row">
        <div class="col-sm-5">
          <text-input
              v-model.trim="user.person.first_name"
              label="Ime"
              type="text"
              name="name"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="user.person.middle_name"
              label="Ime jednog roditelja"
              type="text"
              name="middleName"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="user.person.last_name"
              label="Prezime"
              type="text"
              name="lastName"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="user.person.phone_number"
              label="Broj telefona"
              type="text"
              name="phone_number"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="user.person.email"
              label="Email"
              type="text"
              name="email"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              :hidden="readonly"
              v-model.trim="user.password"
              label="Šifra"
              type="password"
              name="password"
              :required="true"
              :readonly="readonly">
          </text-input>
        </div>

        <div class="row"></div>

        <div class="col-sm-5">
          <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Snimi">
          <input type="submit" v-if="this.action === 'update'" class="btn btn-primary m-2" value="Snimi">
        </div>
      </div>
    </form-tag>
  </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from "@/router";
import {fileMixin} from "@/mixins/fileMixin";
import {useToast} from "vue-toastification";
// import {commonMixin} from "@/mixins/commonMixin";

export default {
  name: 'UserEdit',
  mixins: [fileMixin],
  components: {FormTag, TextInput},
  computed: {
    readonly() {
      if (this.$route.query.action === 'view') {
        return true;
      }
      return false;
    },
  },
  data() {
    return {
      user: {
        person: {first_name: "", middle_name: "", last_name: "", email: "", phone_number: ""},
        password: "",
      },
      action: "",
    }
  },
  methods: {
    async getUserById() {
      axios.get('/users/id/' + this.userId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.user = JSON.parse(response.data.Data);
      }, (error) => {
        this.toast.error(error);
      });
    },
    async submitHandler() {
      if (this.userId != undefined) {
        await this.updateUser();
      } else {
        await this.createUser();
      }
    },
    async createUser() {
      await axios.post('/users/register', JSON.stringify(this.user)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.toast.info("Uspešno kreiran korisnik.");
        router.push("/users");
      }, (error) => {
        this.toast.error(error);
      });
    },
    async updateUser() {
      await axios.post('/users/update', JSON.stringify(this.user)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.toast.info("Uspešno ažuriran korisnik.");
        router.push("/users");
      }, (error) => {
        this.toast.error(error);
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  mounted() {
    if (this.$route.query.id !== '') {
      this.userId = this.$route.query.id;
      this.getUserById();
    }
    this.action = this.$route.query.action;
  }
}
</script>