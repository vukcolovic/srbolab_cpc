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
        <label class="m-1" style="font-size: 1.2em; font-style: italic">Lista vozača</label>
        <button class="iconBtn ms-auto" title="Filter" type="button" data-bs-toggle="collapse" data-bs-target="#filter" aria-expanded="false" aria-controls="filter">
          <i class="fa fa-filter" aria-hidden="true">
          </i>
        </button>
        <button class="iconBtn" title="Traži" type="button" @click="doSearch(0, this.pageSize)">
          <i class="fa fa-search">
          </i>
        </button>
        <button class="iconBtn" title="Obriši" :disabled="!table.selectedClient" @click="deleteClient">
          <i class="fa fa-trash">
          </i></button>
        <button class="iconBtn" title="Dodaj vozače na seminar" :disabled="!selectedSeminar" @click.prevent="saveToSeminar()">
          <i class="fa fa-floppy-o"></i>
        </button>
        <button v-can="'ADMINISTRATOR'" class="iconBtn" title="Lista vozača" @click.prevent="downloadExcelWithClientsList()">
          <i class="fa fa-file-excel-o"></i>
        </button>
      </div>
    </div>
    <div class="collapse multi-collapse border" style="font-size: 0.7em" id="filter">
      <div class="row">
        <div class="col-sm-2">
          <label for="first_name" style="margin-right: 5px">Ime</label>
          <input style="max-width: 130px" type="text" id="first_name" name="Ime" v-model="filter.first_name" />
        </div>
        <div class="col-sm-3">
          <label for="last_name" style="margin-right: 5px">Prezime</label>
          <input type="text" id="last_name" name="Prezime" v-model="filter.last_name" />
        </div>
        <div class="col-sm-3">
          <label for="jmbg" style="margin-right: 5px">JMBG</label>
          <input type="text" id="jmbg" name="jmbg" v-model="filter.jmbg" />
        </div>
        <div class="col-sm-4">
        </div>
        <div class="col-sm-3">
          <label :style="styleLabelSmall" class="mb-1">Firma</label>
          <v-select
              v-model="filter.company_id"
              :options="companies"
              :style="styleInputSmall"
              :reduce="opt => opt.ID"
              label="name_pib"
              placeholder="Traži">
          </v-select>
        </div>
        <div class="col-sm-2 my-1">
          <label for="verified" style="margin-right: 5px">Verifikovan:</label>
          <v-select
              v-model="filter.verified"
              :options="yesNoOptions"
              :style="styleInputSmall"
              :reduce="opt => opt.value"
              label="label"
              placeholder="Traži">
          </v-select>
        </div>
        <div class="col-sm-2 my-1">
          <label :style=styleLabelSmall for="wait_seminar">Čeka seminar:&nbsp;&nbsp;</label>
          <v-select
              v-model="filter.wait_seminar"
              :options="yesNoOptions"
              :style="styleInputSmall"
              :reduce="opt => opt.value"
              label="label"
              placeholder="Traži">
          </v-select>
        </div>

        <div v-if="table.selectedClientIds.length" class="col-sm-4 my-1">
        <label :style=styleLabelSmall for="seminar">Seminar:&nbsp;&nbsp;</label>
        <v-select
              v-model="selectedSeminar"
              :disabled=readonly
              :options="openedSeminars"
              :style="styleInputSmall"
              label="base_info"
              placeholder="Traži">
          </v-select>
        </div>

      </div>
    </div>

    <div class="row mt-2">
      <vue-table-lite
          ref="localTable"
          :has-checkbox="true"
          @row-clicked="selectClient"
          @dblclick="doubleClick"
          :total= "table.totalCount"
          :page-size=this.pageSize
          :columns="table.columns"
          :messages="table.messages"
          :rows="table.rows"
          @do-search="doSearch"
          :rowClasses=table.rowClasess
          :is-loading="table.isLoading"
          @return-checked-rows="updateCheckedRows"
          @is-finished="tableLoadingFinish"
      ></vue-table-lite>
    </div>
  </div>
</template>

<script>
import VueTableLite from "vue3-table-lite";
import axios from "axios";
import {reactive} from "vue";
import {useToast} from "vue-toastification";
import {styleMixin} from "@/mixins/styleMixin";
import vSelect from "vue-select";
import {commonMixin} from "@/mixins/commonMixin";
import router from "@/router";
import {apiMixin} from "@/mixins/apiMixin";
import { fileMixin } from "@/mixins/fileMixin";

export default {
  name: 'ClientsList',
  mixins: [styleMixin, commonMixin, apiMixin, fileMixin],
  components: {vSelect, VueTableLite },
  data() {
    return {
      filter: {verified: "", wait_seminar: "", jmbg: "", first_name: "", last_name: "", company_id: null},
      pageSize: 50,
      openedSeminars: [],
      selectedSeminar: null,
    }
  },
  setup() {
    // Table config
    const table = reactive({
      selectedClient: null,
      selectedClientIds: [],
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
          label: 'Firma',
          field: 'company_name',
          width: '15%',
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
        },
        {
          label: '',
          field: '',
          width: '2%',
          display: function (row) {
            return (
                '<button data-id="' +
                row.ID +
                '" class="is-rows-el name-btn"><i class="fa fa-external-link" aria-hidden="true"></i></button>'
            );
          },
        },
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

    const tableLoadingFinish = (elements) => {
      table.isLoading = false;
      Array.prototype.forEach.call(elements, function (element) {
        if (element.classList.contains("name-btn")) {
          element.addEventListener("click", function () {
            const routeData = router.resolve({name: 'ClientEdit', query: {id: this.dataset.id, action: 'update'}});
            window.open(routeData.href, '_blank');
          });
        }
      });
    };

    const doubleClick = () => {
      router.push("/client?action=update&id=" + table.selectedClient.ID);
    }

     /**
     * Row checked event
     */
     const updateCheckedRows = (rowsKey) => {
      table.selectedClientIds = rowsKey;
      console.log(table.selectedClientIds);
    };

    const toast = useToast();
    return {
      toast,
      table,
      selectClient,
      doubleClick,
      updateCheckedRows,
      tableLoadingFinish,
    };
  },
  methods: {
    async downloadExcelWithClientsList() {
      axios.get('/excel/clients')
          .then(response => {
            var fileContent = JSON.parse(response.data.Data);
            var sampleArr = this.base64ToArrayBuffer(fileContent);
            const blob = new Blob([sampleArr], { type: 'application/xlsx' })

            const link = document.createElement('a')
            link.href = URL.createObjectURL(blob)
            link.download = "vozaci.xlsx"
            link.click()
            URL.revokeObjectURL(link.href)
            //FIXME add notie
          }).catch(console.error)
    },
    async saveToSeminar() {
      const response = confirm("Da li ste sigurni da želite da dodate vozače na seminar?");
      if (response) {
        var bulkRequest = {seminar_id: this.selectedSeminar.ID, client_ids: this.table.selectedClientIds}
        await axios.post('/client-seminar/insert-bulk', JSON.stringify(bulkRequest)).then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            this.toast.error(response.data != null ? response.data.ErrorMessage : "");
            return;
          }
          location.reload();
          this.toast.info("Uspešno dodati vozači na seminar.");
        }, (error) => {
          this.errorToast(error, "/clients/delete");
        });
      }  
    },
    async deleteClient() {
      const response = confirm("Da li ste sigurni da želite da obrišete vozača?");
      if (response) {
        await axios.get('/clients/delete/' + this.table.selectedClient.ID).then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            this.toast.error(response.data != null ? response.data.ErrorMessage : "");
            return;
          }
          location.reload();
        }, (error) => {
          this.errorToast(error, "/clients/delete");
        });
      }
    },
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
          vs.company_name = vs.company.name;
          vs.verified_text = vs.verified ? "Da" : "Ne";
          vs.waiting_seminar_text = vs.wait_seminar ? "Da" : "Ne";
        });
      }, (error) => {
        this.errorToast(error, "/clients/list");
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
        this.errorToast(error, "/clients/count");
      });
    }
  },
  async created() {
    await this.countClients();
    await this.getAllCompanies();
    await this.doSearch(0, this.pageSize, 'id', 'asc');
    await this.getSeminarsByStatusCode("OPENED").then(result => this.openedSeminars = result);
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
