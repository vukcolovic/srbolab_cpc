<template>
  <div class="container">
    <div class="row mt-2">
      <div class="btn-group">
      <button class="iconBtn" title="Dodaj" @click="$router.push({name: 'ClientEdit', query: {id: '', action: 'add' }})">
        <i class="fa fa-user-plus"></i>
      </button>
      <button class="iconBtn" title="Pregledaj" :disabled="table.selectedClient == null" @click="$router.push({name: 'ClientEdit', query: {id: table.selectedClient.ID, action: 'view' }})">
        <i class="fa fa-user"></i>
      </button>
      <button class="iconBtn" title="Izmeni" :disabled="!table.selectedClient" @click="$router.push({name: 'ClientEdit', query: {id: table.selectedClient.ID, action: 'update' }})">
        <i class="fa fa-user-md">
        </i></button>
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
      <vue-table-lite
          ref="localTable"
          @row-clicked="selectClient"
          :total= "table.totalCount"
          :columns="table.columns"
          :messages="table.messages"
          :rows="table.rows"
          @do-search="doSearch"
          :rowClasses=table.rowClasess
          :is-loading="table.isLoading"
      ></vue-table-lite>
    </div>
  </div>
</template>

<script>
import VueTableLite from "vue3-table-lite";
import axios from "axios";
import {reactive} from "vue";
import {useToast} from "vue-toastification";

export default {
  name: 'ClientsList',
  components: { VueTableLite },
  data() {
    return {
      filter: {verified: true, wait_seminar: true, jmbg: "", first_name: "", last_name: ""}
    }
  },
  setup() {
    // Table config
    const table = reactive({
      selectedClient: null,
      isLoading: false,
      isReSearch: false,
      rowClasess: (row) => { return ['is-rows-el', 'row_id_' + row.ID]},
      // filterObject: {},
      columns: [
        {
          label: 'ID',
          field: 'ID',
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
          label: 'Broj telefona',
          field: 'phone_number',
          width: '10%',
        },
        {
          label: 'Email',
          field: 'email',
          width: '10%',
        },
        {
          label: 'JMBG',
          field: 'jmbg',
          width: '10%',
        },
        {
          label: 'Mesto',
          field: 'place',
          width: '10%',
        },
        {
          label: 'Verifikovan',
          field: 'verified_text',
          width: '10%',
        },
        {
          label: 'Čeka seminar',
          field: 'waiting_seminar_text',
          width: '10%',
        }
      ],
      rows: [],
      totalCount: 0,
      messages: {
        pagingInfo: "Prikaz {0} - {1} od {2}",
        pageSizeChangeLabel: "Broj redova:",
        gotoPageLabel: "Idi na stranu:",
        noDataAvailable: "Nema podataka",
      },
    });

    const selectClient= (rowData) => {
      // clear all
      Array.from(document.getElementsByClassName('is-rows-el')).map((el) => {
        el.style.backgroundColor = 'white';
      });
      //style checked row
      if (document.getElementsByClassName('row_id_' + rowData.ID)[0]) {
        document.getElementsByClassName('row_id_' + rowData.ID)[0].style.backgroundColor = '#E8E8E8';
      }
      table.selectedClient = rowData;
    }

    const toast = useToast();
    return {
      toast,
      table,
      selectClient,
    };
  },
  methods: {
    async doSearch(offset, limit, order, sort) {
      console.log(order, sort)
      this.isLoading = true;
      await axios.post('/clients/list?skip=' + offset + '&take=' + limit, JSON.stringify(this.filter)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.table.rows = JSON.parse(response.data.Data);
        this.table.rows.forEach(vs => {
          vs.first_name = vs.person.first_name;
          vs.last_name = vs.person.last_name;
          vs.email = vs.person.email;
          vs.place = vs.address.place;
          vs.phone_number = vs.person.phone_number;
          vs.verified_text = vs.verified ? "Da" : "Ne";
          vs.waiting_seminar_text = vs.wait_seminar ? "Da" : "Ne";
        });
      }, (error) => {
        this.toast.error(error.message);
      });

      this.isLoading = false;
    },
    async countClients() {
        await axios.get('/clients/count').then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            this.toast.error(response.data != null ? response.data.ErrorMessage : "");
            return;
          }
            this.table.totalCount = response.data.Data;
        }, (error) => {
          this.toast.error(error.message);
        });
    }
  },
  async created() {
    await this.countClients();
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
