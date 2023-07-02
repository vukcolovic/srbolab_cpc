<template>
  <div class="container">
    <div class="row m-3">
      <div class="btn-group">
      <button v-can="'UpravljanjeKorisnicima'" class="iconBtn" title="Dodaj" @click="$router.push({name: 'UserEdit', params: {userId: '', action: 'add' }})">
        <i class="fa fa-user-plus"></i>
      </button>
      <button class="iconBtn" title="Pregledaj" :disabled="selectedUser == null" @click="$router.push({name: 'UserEdit', params: {userId: selectedUser.id, action: 'view' }})">
        <i class="fa fa-user"></i>
      </button>
      <button v-can="'UpravljanjeKorisnicima'" class="iconBtn" title="Izmeni" :disabled="selectedUser == null" @click="$router.push({name: 'UserEdit', params: {userId: selectedUser.id, action: 'update' }})">
        <i class="fa fa-user-md">
        </i></button>
      <button v-can="'UpravljanjeKorisnicima'" class="iconBtn" title="Obrisi" @click="deleteUser" :disabled="selectedUser == null">
        <i class="fa fa-user-times"></i>
      </button>
      </div>
    </div>
    <div class="row m-3">
      <vue-table-lite
          @row-clicked="selectUser"
          :total= "totalCount"
          :columns="columns"
          :messages="messages"
          :rows="rows"
          @do-search="doSearch"
          :is-loading="isLoading"
      ></vue-table-lite>
    </div>
  </div>
</template>

<script>
import VueTableLite from "vue3-table-lite";
import axios from "axios";

export default {
  name: 'SeminarsList',
  components: { VueTableLite },
  data() {
    return {
      columns: [
        {
          label: 'ID',
          field: 'id',
          width: '3%',
          isKey: true,
        },
        {
          label: 'Ime',
          field: 'first_name',
          width: '10%',
        },
        {
          label: 'Prezime',
          field: 'last_name',
          width: '10%',
        },
        {
          label: 'Email',
          field: 'email',
          width: '10%',
        },
        {
          label: 'Broj telefona',
          field: 'phone_number',
          width: '10%',
        },
        {
          label: 'JMBG',
          field: 'jmbg',
          width: '10%',
        }
      ],
      messages: {
        pagingInfo: "Prikaz {0} - {1} od {2}",
        pageSizeChangeLabel: "Broj redova:",
        gotoPageLabel: "Idi na stranu:",
        noDataAvailable: "Nema podataka",
      },
    rows: [],
      selectedUser: null,
      isLoading: false,
      totalCount: 0
    }
  },
  methods: {
    selectUser(rowData) {
      this.selectedUser = rowData;
    },
    async doSearch(offset, limit, order, sort) {
      console.log(order, sort)
      this.isLoading = true;
      await axios.get('/users/list?skip=' + offset + '&take=' + limit).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          // notie.alert({
          //   type: 'error',
          //   text: response.data.ErrorMessage,
          //   position: 'bottom',
          // })
          return;
        }
        this.rows = JSON.parse(response.data.Data);
      }, (error) => {
        // notie.alert({
        //   type: 'error',
        //   text: "Greska: " + error,
        //   position: 'bottom',
        // })
      });

      this.isLoading = false;
    },
    async deleteUser() {
      if (confirm("Da li ste sigurni da zelite da obrisete korisnika?")) {
        await axios.get('/users/delete/' + this.selectedUser.id).then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            // notie.alert({
            //   type: 'error',
            //   text: response.data.ErrorMessage,
            //   position: 'bottom',
            // })
          }
        }, (error) => {
          // notie.alert({
          //   type: 'error',
          //   text: "Greska: " + error,
          //   position: 'bottom',
          // })
        });
      }

      location.reload();
    },
    async countUsers() {
        await axios.get('/users/count').then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            // notie.alert({
            //   type: 'error',
            //   text: response.data.ErrorMessage,
            //   position: 'bottom',
            // })
            return;
          }
            this.totalCount = response.data.Data;
        }, (error) => {
          // notie.alert({
          //   type: 'error',
          //   text: "Greska: " + error,
          //   position: 'bottom',
          // })
        });
    }
  },
  async created() {
    await this.countUsers();
    await this.doSearch(0, 10, 'id', 'asc');
  }
  }
</script>

<style scoped>
::v-deep(.vtl-table .vtl-thead .vtl-thead-th) {
  font-size: 12px;
}
::v-deep(.vtl-table td),
::v-deep(.vtl-table tr) {
  font-size: 12px;
  padding: 5px;
}
::v-deep(.vtl-paging-info) {
  font-size: 12px;
  padding: 5px;
}
::v-deep(.vtl-paging-count-label),
::v-deep(.vtl-paging-page-label),
::v-deep(.vtl-paging-count-dropdown),
::v-deep(.vtl-paging-page-dropdown){
  font-size: 12px;
  padding: 5px;
}
::v-deep(.vtl-paging-pagination-page-link) {
  font-size: 12px;
  padding: 5px;
}
</style>
