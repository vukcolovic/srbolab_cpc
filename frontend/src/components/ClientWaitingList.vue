<template>
  <div class="container">
    <div class="row mt-2">
      <div class="btn-group">
      <button class="iconBtn" title="Dodaj" @click="$router.push({name: 'ClientEdit', query: {id: '', action: 'add' }})">
        <i class="fa fa-user-plus"></i>
      </button>
<!--      <button class="iconBtn" title="Pregledaj" :disabled="selectedClient == null" @click="$router.push({name: 'ClientEdit', query: {id: selectedClient.ID, action: 'view' }})">-->
<!--        <i class="fa fa-user"></i>-->
<!--      </button>-->
      <button class="iconBtn" title="Izmeni" :disabled="!selectedClient" @click="$router.push({name: 'ClientEdit', query: {id: selectedClient.ID, action: 'update' }})">
        <i class="fa fa-pencil">
        </i></button>
        <label class="m-1" style="font-size: 1.2em; font-style: italic">Čekaonica</label>
        <button class="iconBtn ms-auto" title="Filter" type="button" data-bs-toggle="collapse" data-bs-target="#filter" aria-expanded="false" aria-controls="filter">
          <i class="fa fa-filter" aria-hidden="true">
          </i>
        </button>
        <button class="iconBtn" title="Traži" type="button" @click="doSearch(0, 10)">
          <i class="fa fa-search">
          </i>
        </button>
      </div>
    </div>
    <div class="collapse multi-collapse border" style="font-size: 0.7em" id="filter">
      <div class="row">
          <div class="col-sm-3">
            <label for="first_name" style="margin-right: 5px">Ime</label>
            <input type="text" id="first_name" name="Ime" v-model="filter.first_name" />
          </div>
        <div class="col-sm-3">
            <label for="last_name" style="margin-right: 5px">Prezime</label>
            <input type="text" id="last_name" name="Prezime" v-model="filter.last_name" />
          </div>
        <div class="col-sm-3">
          <label for="jmbg" style="margin-right: 5px">JMBG</label>
          <input type="text" id="jmbg" name="jmbg" v-model="filter.jmbg" />
        </div>
      </div>
    </div>

    <div class="row mt-2">
      <table class="table">
        <thead>
          <tr class="bg-primary text-white">
            <td style="width: 4%;">ID</td>
            <td style="width: 10%;">Ime i Prezime</td>
            <td style="width: 10%;">JMBG</td>
            <td style="width: 10%;">Verifikovan</td>
            <td style="width: 10%;">Firma</td>
            <td style="width: 10%;">Plaćeno</td>
            <td style="width: 30%;">Seminar</td>
            <td style="width: 5%;"></td>
          </tr>
        </thead>
        <tbody>
          <tr @click="selectRow(row, i)" :id="'id' + i" v-for="(row, i) in rows" :key="row.ID">
            <td>{{row.ID}}</td>
            <td>{{row.first_name}} {{row.last_name}}</td>
            <td>{{row.jmbg}}</td>
            <td>{{row.verified_text}}</td>
            <td>TODO</td>
            <td><v-select
                v-model="row.payed"
                :options="yesNoOptions"
                :disabled="!row.verified"
                :style="styleInputSmall"
                label="label"
                :reduce="opt => opt.value"
                placeholder="Traži">
            </v-select></td>
            <td>
              <v-select
                  v-model="row.selectedSeminar"
                  :options="openedSeminars"
                  :style="styleInputSmall"
                  :disabled="!row.verified"
                  label="base_info"
                  placeholder="Traži">
              </v-select>
            </td>
            <td><button @click.prevent="saveClient(row)" class="btn btn-primary">Snimi</button></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import vSelect from "vue-select";
import {styleMixin} from "@/mixins/styleMixin";
import {apiMixin} from "@/mixins/apiMixin";
import {commonMixin} from "@/mixins/commonMixin";
import router from "@/router";
import {useToast} from "vue-toastification";
// import {useToast} from "vue-toastification";

export default {
  name: 'ClientWaitingList',
  components: {vSelect},
  mixins:[styleMixin, apiMixin, commonMixin],
  data() {
    return {
      rows: [],
      totalCount: 0,
      selectedClient: null,
      filter: {verified: "", wait_seminar: "", jmbg: "", first_name: "", last_name: "", waiting_room: true},
      openedSeminars: [],
    }
  },
  methods: {
    selectRow(item, i) {
      this.selectedClient = item;
      var rows = document.getElementsByTagName("table")[0].rows;
      Array.from(rows).forEach(row => row.style.backgroundColor = "white");

      let selectedRow = document.getElementById("id" + i);
      selectedRow.style.backgroundColor = '#E8E8E8';
    },
    async saveClient(row) {
      if (row.selectedSeminar) {
        const payed = row.payed=== "true" ? true : false;
        row.seminars.push({"client_id": row.ID, "seminar_id": row.selectedSeminar.ID, "payed": payed});
        row.wait_seminar = false;
      }
      delete row.selectedSeminar;
      delete row.payed;
      await axios.post('/clients/update', JSON.stringify(row)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran klijent.");
        router.push("/clients");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async doSearch(offset, limit, order, sort) {
      console.log(order, sort)
      this.isLoading = true;
      await axios.post('/clients/list?skip=' + offset + '&take=' + limit, JSON.stringify(this.filter)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.rows = JSON.parse(response.data.Data);
        this.rows.forEach(vs => {
          vs.first_name = vs.person.first_name;
          vs.last_name = vs.person.last_name;
          vs.verified_text = vs.verified ? "Da" : "Ne";
          vs.waiting_seminar_text = vs.wait_seminar ? "Da" : "Ne";
          if (vs.seminars == null) {
            vs.seminars = [];
          }
        });
      }, (error) => {
        this.toast.error(error);
      });

      this.isLoading = false;
    },
    async countClients() {
        await axios.get('/clients/count').then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            this.toast.error(response.data != null ? response.data.ErrorMessage : "");
            return;
          }
            this.totalCount = response.data.Data;
        }, (error) => {
          this.toast.error(error.message);
        });
    }
  },
  async created() {
    await this.countClients();
    await this.doSearch(0, 10, 'id', 'asc');
    await this.getSeminarsByStatusCode("OPENED").then(result => this.openedSeminars = result);
  }, setup() {
    const toast = useToast();
    return {toast}
  }
  }
</script>

<style scoped>
table {
  table-layout: fixed;
}
</style>