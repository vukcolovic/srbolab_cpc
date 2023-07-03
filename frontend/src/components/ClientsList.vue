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
      <button class="iconBtn" title="Izmeni" :disabled="table.selectedClient == null" @click="$router.push({name: 'ClientEdit', query: {id: table.selectedClient.ID, action: 'update' }})">
        <i class="fa fa-user-md">
        </i></button>
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
      await axios.get('/clients/list?skip=' + offset + '&take=' + limit).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.table.rows = JSON.parse(response.data.Data);
        this.table.rows.forEach(vs => {
          vs.first_name = vs.person.first_name;
          vs.last_name = vs.person.last_name;
          vs.email = vs.person.email;
          vs.phone_number = vs.person.phone_number;
        });
      }, (error) => {
        this.toast.error(error);
      });

      this.isLoading = false;
    },
    async countClients() {
        await axios.get('/clients/count').then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            this.toast.error(response.data.ErrorMessage);
            return;
          }
            this.table.totalCount = response.data.Data;
        }, (error) => {
          this.toast.error(error);
          alert(error);
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
