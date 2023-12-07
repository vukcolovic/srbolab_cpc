<template>
  <div class="container">
    <div class="row mt-2">
      <div class="btn-group">
        <button class="iconBtn" title="Dodaj" @click="$router.push({name: 'SeminarEdit', query: {id: '', action: 'add' }})">
          <i class="fa fa-user-plus"></i>
        </button>
        <button class="iconBtn" title="Pregledaj" :disabled="table.selectedSeminar == null" @click="$router.push({name: 'SeminarEdit', query: {id: table.selectedSeminar.ID, action: 'view' }})">
          <i class="fa fa-user"></i>
        </button>
        <button class="iconBtn" title="Izmeni" :disabled="table.selectedSeminar == null" @click="$router.push({name: 'SeminarEdit', query: {id: table.selectedSeminar.ID, action: 'update' }})">
          <i class="fa fa-user-md">
          </i></button>
        <label class="m-1" style="font-size: 1.2em; font-style: italic">Seminari</label>

        <button class="iconBtn ms-auto" title="Filter" type="button" data-bs-toggle="collapse" data-bs-target="#filter" aria-expanded="false" aria-controls="filter">
          <i class="fa fa-filter" aria-hidden="true">
          </i>
        </button>
        <button class="iconBtn" title="Traži" type="button" @click="doSearch(0, this.pageSize)">
          <i class="fa fa-search">
          </i>
        </button>
        <button class="iconBtn" title="Polaznici" type="button" @click="printExcelClientsReport()">
          <i class="fa fa-file-excel-o">
          </i>
        </button>
        <button class="iconBtn" title="Predavači" type="button" @click="printExcelTeachersReport()">
          <i class="fa fa-file-excel-o">
          </i>
        </button>
      </div>
    </div>

    <div class="collapse multi-collapse border" style="font-size: 0.7em" id="filter">
      <div class="row">
        <div class="col-sm-2">
          <label for="first_name" style="margin-right: 5px">Datum početka od:</label>
          <Datepicker
              v-model="filter.date_from"
              :style="dateStyleInput"
              inputFormat="dd.MM.yyyy"
              placeholder="dd.MM.yyyy"
              typeable="true"
          />
        </div>
        <div class="col-sm-2">
          <label for="first_name" style="margin-right: 5px">Datum početka do:</label>
          <Datepicker
              v-model="filter.date_to"
              :style="dateStyleInput"
              inputFormat="dd.MM.yyyy"
              placeholder="dd.MM.yyyy"
              typeable="true"
          />
        </div>
        <div class="col-sm-2 my-1">
          <label for="verified" style="margin-right: 5px">Ispitno mesto:</label>
          <v-select
              v-model="filter.location_id"
              :options="locations"
              :style="styleInputSmall"
              :reduce="opt => opt.ID"
              label="address_place"
              placeholder="Traži">
          </v-select>
        </div>
      </div>
    </div>
    <div class="row mt-2">
      <vue-table-lite
          ref="localTable"
          @row-clicked="selectSeminar"
          @dblclick="doubleClick"
          :total= "table.totalCount"
          :page-size=this.pageSize
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
import {dateMixin} from "@/mixins/dateMixin";
import {apiMixin} from "@/mixins/apiMixin";
import {commonMixin} from "@/mixins/commonMixin";
import router from "@/router";
import vSelect from "vue-select";
import Datepicker from "vue3-datepicker";
import {styleMixin} from "@/mixins/styleMixin";
import {fileMixin} from "@/mixins/fileMixin";

export default {
  name: 'SeminarsList',
  mixins: [dateMixin, apiMixin, commonMixin, styleMixin, fileMixin],
  components: {Datepicker, vSelect, VueTableLite },
  data() {
    return {
      pageSize: 50,
      filter: {location_id: null, date_from: null, date_to: null},
    }
  },
  setup() {
    // Table config
    const table = reactive({
      selectedSeminar: null,
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
          label: 'Broj seminara',
          field: 'serial_number_by_location',
          width: '10%',
        },
        {
          label: 'Početak(MM-DD-YYYY)',
          field: 'start_date',
          width: '10%',
        },
        {
          label: 'Lokacija',
          field: 'location_address',
          width: '10%',
        },
        {
          label: 'Učionica',
          field: 'class_room_name',
          width: '10%',
        },
        {
          label: 'Vrsta',
          field: 'type',
          width: '10%',
        },
        {
          label: 'Status',
          field: 'status',
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

    const selectSeminar= (rowData) => {
      // clear all
      Array.from(document.getElementsByClassName('is-rows-el')).map((el) => {
        el.style.backgroundColor = 'white';
      });
      //style checked row
      if (document.getElementsByClassName('row_id_' + rowData.ID)[0]) {
        document.getElementsByClassName('row_id_' + rowData.ID)[0].style.backgroundColor = '#E8E8E8';
      }
      table.selectedSeminar = rowData;
    }

    const doubleClick = () => {
      router.push("/seminar?action=update&id=" + table.selectedSeminar.ID);
    }

    const toast = useToast();
    return {
      toast,
      table,
      selectSeminar,
      doubleClick,
    };
  },
  methods: {
    printExcelClientsReport() {
      axios.post('/excel/seminars-report/clients', JSON.stringify(this.filter))
          .then(response => {
            var fileContent = JSON.parse(response.data.Data);
            var sampleArr = this.base64ToArrayBuffer(fileContent);
            const blob = new Blob([sampleArr], { type: 'application/xlsx' })

            const link = document.createElement('a')
            link.href = URL.createObjectURL(blob)
            link.download = "Seminari-statistika-polaznika.xlsx"
            link.click()
            URL.revokeObjectURL(link.href)
            //FIXME add notie
          }).catch(console.error)
    },
    printExcelTeachersReport() {
      axios.post('/excel/seminars-report/teachers', JSON.stringify(this.filter))
          .then(response => {
            var fileContent = JSON.parse(response.data.Data);
            var sampleArr = this.base64ToArrayBuffer(fileContent);
            const blob = new Blob([sampleArr], { type: 'application/xlsx' })

            const link = document.createElement('a')
            link.href = URL.createObjectURL(blob)
            link.download = "Seminari-statistika-predavaca.xlsx"
            link.click()
            URL.revokeObjectURL(link.href)
            //FIXME add notie
          }).catch(console.error)
    },
    async doSearch(offset, limit, order, sort) {
      console.log(order, sort)
      this.isLoading = true;
      await axios.post('/seminars/list?skip=' + offset + '&take=' + limit, JSON.stringify(this.filter)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.table.rows = JSON.parse(response.data.Data);
        this.table.rows.forEach(s => {
          s.location_address = s.class_room.location.address.place;
          s.class_room_name = s.class_room.name;
          s.type = this.getSeminarFullType(s.seminar_theme.base_seminar_type, s.seminar_theme);
          s.status = s.seminar_status.name;
          s.start_date = this.formatDateWithPoints(s.start_date);
        });
      }, (error) => {
        this.errorToast(error, "/seminars/list");
      });

      this.isLoading = false;
    },
    async countSeminars() {
      await axios.get('/seminars/count').then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.table.totalCount = response.data.Data;
      }, (error) => {
        this.errorToast(error, "/seminars/count");
      });
    }
  },
  async created() {
    await this.countSeminars();
    this.getAllLocations();
    await this.doSearch(0, this.pageSize, 'id', 'asc');
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
