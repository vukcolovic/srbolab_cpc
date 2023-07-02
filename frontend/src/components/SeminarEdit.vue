<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 v-if="action === 'add'" class="mt-2">Registracija</h3>
        <h3 v-if="action === 'view'" class="mt-2">Pregled</h3>
        <h3 v-if="action === 'update'" class="mt-2">Azuriranje</h3>
        <hr>
      </div>
    </div>
        <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
          <div class="row">
          <div class="col-sm-5">
            <text-input
                v-model.trim="user.first_name"
                label="Ime"
                type="text"
                name="name"
                :required= true
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.last_name"
                label="Prezime"
                type="text"
                name="lastName"
                :required= true
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.adress"
                label="Adresa"
                type="text"
                name="address"
                :required= false
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.phone_number"
                label="Broj telefona"
                type="text"
                name="phone_number"
                :required= false
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.jmbg"
                label="JMBG"
                type="text"
                name="jmbg"
                :required= false
                :readonly="readonly">
            </text-input>

            <label class="mb-1 mt-1" :style="styleLabel">Ispitno mesto</label>
            <v-select
                :disabled=readonly
                :style="styleInput"
                v-model="user.examination_place_id"
                :reduce="(option) => option.id"
                :options="examinationPlaces"
                placeholder="Traži"
                label="name">
            </v-select>
          </div>

          <div class="col-sm-5">
            <text-input
                v-model.trim="user.email"
                label="Email"
                type="text"
                name="email"
                :required= true
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.current_password"
                label="Trenutna sifra"
                type="password"
                name="currentPassword"
                :hidden="!showCurentPassword"
                :required= false
                :readonly="readonly">
            </text-input>

            <text-input
                :hidden="readonly"
                v-model.trim="user.password"
                :label="passwordLabel"
                type="password"
                name="password"
                :required="passwordRequired"
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.started_work"
                label="Pocetak rada"
                type="date"
                name="date"
                :required=false
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.contract_number"
                label="Broj ugovora"
                type="text"
                name="contract_number"
                :required=false
                :readonly="readonly">
            </text-input>

            <text-input
                v-model.trim="user.contract_type"
                label="Tip ugovora"
                type="text"
                name="contract_type"
                :required=false
                :readonly="readonly">
            </text-input>
          </div>
          <hr>
            <h5 class="mt-2">Role</h5>
            <div class="col-sm-10 mb-2">
              <li v-for="(role, index) in roles" :key="role.code" style="list-style-type: none;">
                <input :id="index" :value="role" :disabled="readonly" type="checkbox" v-model="user.roles" />
                <label for="index">&nbsp; {{role.code}} ({{role.description}})</label>
              </li>
            </div>
            <hr>
            <div class="col-sm-5">
          <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Registracija">
          <input type="submit" v-if="this.action === 'update'" class="btn btn-primary m-2" value="Azuriranje">
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
import {commonMixin} from "@/mixins/commonMixin";

export default {
  name: 'SeminarEdit',
  props: {
    action: {
      default: 'add',
      type: String
    },
    userId: {
      default: '',
      type: String
    }
  },
  mixins: [commonMixin],
  components: {FormTag, TextInput},
  computed: {
    readonly() {
      if (this.action === 'view') {
        return true;
      }
      return false;
    },
    passwordRequired() {
      if (this.action === 'add') {
        return true;
      }
      return false;
    },
    showCurentPassword() {
      if (this.action === 'update') {
        return true;
      }
      return false;
    },
    passwordLabel() {
      if (this.action === 'update') {
        return "Nova sifra";
      }
      return "Sifra";
    }
  },
  data() {
    return {
      user: {first_name: "", last_name: "", email: "", password: "", current_password: "", phone_number: "", contract_number: "", contract_type: "", jmbg: "", adress: "", started_work: "", roles: [], examination_place_id: 0},
      roles: [],
      styleInput: "font-size: 16px",
      styleLabel: "font-size: 16px",
    }
  },
  methods: {
    getDate(date) {
      return date.split('T')[0];
    },
    async getUserById() {
        axios.get('/users/id/' + this.userId).then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            // notie.alert({
            //   type: 'error',
            //   text: response.data.ErrorMessage,
            //   position: 'bottom',
            // })
            return;
          }
          this.user = JSON.parse(response.data.Data);
          this.user.started_work = this.getDate(this.user.started_work);
        }, (error) => {
          // notie.alert({
          //   type: 'error',
          //   text: "Greska: " + error,
          //   position: 'bottom',
          // })
        });
    },
    async submitHandler() {
      if (this.userId !== '') {
        await this.updateUser();
      } else {
        await this.createUser();
      }
    },
    async createUser() {
      await axios.post('/users/register', JSON.stringify(this.user)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          // notie.alert({
          //   type: 'error',
          //   text: response.data.ErrorMessage,
          //   position: 'bottom',
          // })
          return;
        }
        // notie.alert({
        //   type: 'success',
        //   text: 'Uspešno kreiran korisnik!',
        //   position: 'bottom',
        //   time: 2,
        // })
        router.push("/users");
      }, (error) => {
        notie.alert({
          type: 'error',
          text: "Greska: " + error,
          position: 'bottom',
        })
      });
    },
    async updateUser() {
      await axios.post('/users/update', JSON.stringify(this.user)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          // notie.alert({
          //   type: 'error',
          //   text: response.data.ErrorMessage,
          //   position: 'bottom',
          // })
          return;
        }
        // notie.alert({
        //   type: 'success',
        //   text: 'Uspešno ažuriran korisnik!',
        //   position: 'bottom',
        //   time: 2,
        // })
        router.push("/users");
      }, (error) => {
        // notie.alert({
        //   type: 'error',
        //   text: "Greska: " + error,
        //   position: 'bottom',
        // })
      });
    },
    async getAllRoles() {
      await axios.get('/role/list?skip=0&take=1000').then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          // notie.alert({
          //   type: 'error',
          //   text: response.data.ErrorMessage,
          //   position: 'bottom',
          // })
          return;
        }
        this.roles = JSON.parse(response.data.Data);
      }, (error) => {
        // notie.alert({
        //   type: 'error',
        //   text: "Greska: " + error,
        //   position: 'bottom',
        // })
      });
    },
  },
  mounted() {
    this.getAllRoles();
    this.getAllExaminationPlaces();
    if (this.userId !== '') {
      this.getUserById();
    }
  }
}
</script>