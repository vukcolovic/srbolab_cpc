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

export default {
  name: 'SeminarsList',
  mixins: [dateMixin, apiMixin, commonMixin],
  components: { VueTableLite },
  data() {
    return {
      pageSize: 50,
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
    async doSearch(offset, limit, order, sort) {
      console.log(order, sort)
      this.isLoading = true;
      await axios.get('/seminars/list?skip=' + offset + '&take=' + limit).then((response) => {
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
