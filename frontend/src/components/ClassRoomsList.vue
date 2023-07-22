<template>
  <div class="container">
    <div class="row mt-2">
      <div class="btn-group">
        <button class="iconBtn" title="Dodaj" @click="$router.push({name: 'ClassRoomEdit', query: {id: '', action: 'add' }})">
          <i class="fa fa-user-plus"></i>
        </button>
        <button class="iconBtn" title="Pregledaj" :disabled="table.selectedClassRoom == null" @click="$router.push({name: 'ClassRoomEdit', query: {id: table.selectedClassRoom.ID, action: 'view' }})">
          <i class="fa fa-user"></i>
        </button>
        <button class="iconBtn" title="Izmeni" :disabled="table.selectedClassRoom == null" @click="$router.push({name: 'ClassRoomEdit', query: {id: table.selectedClassRoom.ID, action: 'update' }})">
          <i class="fa fa-user-md">
          </i></button>
      </div>
    </div>
    <div class="row mt-2">
      <vue-table-lite
          ref="localTable"
          @row-clicked="selectClassRoom"
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
import {dateMixin} from "@/mixins/dateMixin";

export default {
  name: 'ClassRoomsList',
  mixins: [dateMixin],
  components: { VueTableLite },
  setup() {
    // Table config
    const table = reactive({
      selectedClassRoom: null,
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
          label: 'Naziv',
          field: 'name',
          width: '20%',
        },
        {
          label: 'Broj mesta',
          field: 'max_students',
          width: '20%',
        }
        ,
        {
          label: 'Mesto',
          field: 'location_place',
          width: '20%',
        }
        ,
        {
          label: 'Ulica i broj',
          field: 'location_street_and_number',
          width: '20%',
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

    const selectClassRoom= (rowData) => {
      // clear all
      Array.from(document.getElementsByClassName('is-rows-el')).map((el) => {
        el.style.backgroundColor = 'white';
      });
      //style checked row
      if (document.getElementsByClassName('row_id_' + rowData.ID)[0]) {
        document.getElementsByClassName('row_id_' + rowData.ID)[0].style.backgroundColor = '#E8E8E8';
      }
      table.selectedClassRoom = rowData;
    }

    const toast = useToast();
    return {
      toast,
      table,
      selectClassRoom,
    };
  },
  methods: {
    async doSearch(offset, limit, order, sort) {
      console.log(order, sort, offset, limit)
      this.isLoading = true;
      await axios.get('/class-rooms/list').then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.table.rows = JSON.parse(response.data.Data);
        this.table.rows.forEach(s => {
          s.location_place = s.location.address.place;
          s.location_street_and_number = s.location.address.street + " " + s.location.address.house_number;
        });
        this.table.totalCount = this.table.rows.length;
      }, (error) => {
        this.toast.error(error.message);
      });

      this.isLoading = false;
    },
  },
  async created() {
    await this.doSearch(0, 0, 'id', 'asc');
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
