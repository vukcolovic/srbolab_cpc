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
<!--          <div class="row">-->
<!--          <div class="col-sm-5">-->
            <text-input
                v-model.trim="client.first_name"
                label="Ime"
                type="text"
                name="name"
                :required= true
                :readonly="readonly">
            </text-input>

<!--            <text-input-->
<!--                v-model.trim="user.last_name"-->
<!--                label="Prezime"-->
<!--                type="text"-->
<!--                name="lastName"-->
<!--                :required= true-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                v-model.trim="user.adress"-->
<!--                label="Adresa"-->
<!--                type="text"-->
<!--                name="address"-->
<!--                :required= false-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                v-model.trim="user.phone_number"-->
<!--                label="Broj telefona"-->
<!--                type="text"-->
<!--                name="phone_number"-->
<!--                :required= false-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                v-model.trim="user.jmbg"-->
<!--                label="JMBG"-->
<!--                type="text"-->
<!--                name="jmbg"-->
<!--                :required= false-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <label class="mb-1 mt-1" :style="styleLabel">Ispitno mesto</label>-->
<!--            <v-select-->
<!--                :disabled=readonly-->
<!--                :style="styleInput"-->
<!--                v-model="user.examination_place_id"-->
<!--                :reduce="(option) => option.id"-->
<!--                :options="examinationPlaces"-->
<!--                placeholder="Traži"-->
<!--                label="name">-->
<!--            </v-select>-->
<!--          </div>-->

<!--          <div class="col-sm-5">-->
<!--            <text-input-->
<!--                v-model.trim="user.email"-->
<!--                label="Email"-->
<!--                type="text"-->
<!--                name="email"-->
<!--                :required= true-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                v-model.trim="user.current_password"-->
<!--                label="Trenutna sifra"-->
<!--                type="password"-->
<!--                name="currentPassword"-->
<!--                :hidden="!showCurentPassword"-->
<!--                :required= false-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                :hidden="readonly"-->
<!--                v-model.trim="user.password"-->
<!--                :label="passwordLabel"-->
<!--                type="password"-->
<!--                name="password"-->
<!--                :required="passwordRequired"-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                v-model.trim="user.started_work"-->
<!--                label="Pocetak rada"-->
<!--                type="date"-->
<!--                name="date"-->
<!--                :required=false-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                v-model.trim="user.contract_number"-->
<!--                label="Broj ugovora"-->
<!--                type="text"-->
<!--                name="contract_number"-->
<!--                :required=false-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->

<!--            <text-input-->
<!--                v-model.trim="user.contract_type"-->
<!--                label="Tip ugovora"-->
<!--                type="text"-->
<!--                name="contract_type"-->
<!--                :required=false-->
<!--                :readonly="readonly">-->
<!--            </text-input>-->
<!--          </div>-->
<!--          <hr>-->
<!--            <h5 class="mt-2">Role</h5>-->
<!--            <div class="col-sm-10 mb-2">-->
<!--              <li v-for="(role, index) in roles" :key="role.code" style="list-style-type: none;">-->
<!--                <input :id="index" :value="role" :disabled="readonly" type="checkbox" v-model="user.roles" />-->
<!--                <label for="index">&nbsp; {{role.code}} ({{role.description}})</label>-->
<!--              </li>-->
<!--            </div>-->
<!--            <hr>-->
<!--            <div class="col-sm-5">-->
<!--          <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Registracija">-->
<!--          <input type="submit" v-if="this.action === 'update'" class="btn btn-primary m-2" value="Azuriranje">-->
<!--            </div>-->
<!--            </div>-->
        </form-tag>
      </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from "@/router";
// import {commonMixin} from "@/mixins/commonMixin";

export default {
  name: 'ClientEdit',
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
      client: {first_name: "", last_name: "", email: "", phone_number: "", jmbg: ""},
      action: "",
    }
  },
  methods: {
    async getClientById() {
        axios.get('/clients/id/' + this.clientId).then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            // notie.alert({
            //   type: 'error',
            //   text: response.data.ErrorMessage,
            //   position: 'bottom',
            // })
            return;
          }
          this.user = JSON.parse(response.data.Data);
        }, (error) => {
          // notie.alert({
          //   type: 'error',
          //   text: "Greska: " + error,
          //   position: 'bottom',
          // })
          alert(error);
        });
    },
    async submitHandler() {
      if (this.clientId !== '') {
        await this.updateClient();
      } else {
        await this.createClient();
      }
    },
    async createClient() {
      await axios.post('/clients/register', JSON.stringify(this.client)).then((response) => {
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
        // notie.alert({
        //   type: 'error',
        //   text: "Greska: " + error,
        //   position: 'bottom',
        // })
        alert(error);
      });
    },
    async updateClient() {
      await axios.post('/clients/update', JSON.stringify(this.client)).then((response) => {
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
        router.push("/clients");
      }, (error) => {
        // notie.alert({
        //   type: 'error',
        //   text: "Greska: " + error,
        //   position: 'bottom',
        // })
        alert(error);
      });
    },
  },
  mounted() {
    if (this.$route.query.id !== '') {
      this.getClientById();
    }
    this.action = this.$route.query.action;
  }
}
</script>