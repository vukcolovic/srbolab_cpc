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
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-4">
          <div v-if="seminar.seminar_status && action !== 'add'">
            Status seminara: {{seminar.seminar_status.name}}
          </div>
          <label :style="styleLabel" class="mb-1 mt-1">Tip seminara</label>
          <v-select
              v-model="selectedBaseSeminarType"
              :disabled=readonly
              :options="seminarBaseTypes"
              :style="styleInput"
              label="name"
              placeholder="Traži"
              @option:selected="onSeminarTypeChange">
          </v-select>

          <div v-if="seminarThemesByType.length > 0">
            <label :style="styleLabel" class="mb-1 mt-1">Tema seminara</label>
            <v-select
                v-model="seminar.seminar_theme"
                :disabled=readonly
                :options="seminarThemesByType"
                :style="styleInput"
                label="name"
                placeholder="Traži">
            </v-select>
          </div>

          <label :style="styleLabel" class="mb-1 mt-1">Lokacija</label>
          <v-select
              @option:selected="onLocationChange"
              v-model="location"
              :disabled=readonly
              :options="locations"
              :style="styleInput"
              label="address_place"
              placeholder="Traži">
          </v-select>

          <label :style="styleLabel" class="mb-1 mt-1">Učionica</label>
          <v-select
              v-model="seminar.class_room"
              :disabled=readonly
              :options="classRoomsByLocationId"
              :style="styleInput"
              label="name"
              placeholder="Traži">
          </v-select>

          <text-input
              v-model="seminar.start_date"
              :readonly="readonly"
              :required=true
              label="Početak seminara (MM/DD/YYYY)"
              name="start"
              type="date">
          </text-input>
        </div>

        <div class="col-sm-3">
        </div>

        <div class="col-sm-5" v-if="action !== 'add'">
          <div>
            <h5>Spisak polaznika</h5>
            <table class="styled-table">
              <thead>
              <tr class="bg-primary text-white">
                <td style="width: 20%;">Plaćeno</td>
                <td style="width: 40%;">Firma</td>
                <td style="width: 40%;">Ime i prezime</td>
              </tr>
              </thead>
              <tbody>
              <tr v-for="trainee in seminar.trainees" :key="trainee.client_id">
                <td :class="[trainee.payed ? 'bg-success' : 'bg-danger']">{{ trainee.payed  ? 'DA' : 'NE' }}</td>
                <td class="p-1">{{trainee.client.company.name}}</td>
                <td class="p-1">{{trainee.client.person.first_name}} {{trainee.client.person.last_name}}</td>
              </tr>
              </tbody>
            </table>
          </div>

            <div class="shell mt-3">
              <div class="bar" :style="{ width: percentFilled  + '%' }">
              </div>
              Broj polaznika: <span>{{seminar.trainees.length}} / {{seminar.class_room.max_students}}</span>
            </div>

        </div>
        <div class="row"></div>
        <div class="col-sm-5">
          <input v-if="this.action === 'add'" class="btn btn-primary m-2" type="submit" value="Snimi">
          <input v-if="this.action === 'update'" class="btn btn-primary m-2" type="submit" value="Snimi">
          <input v-if="this.seminar && this.seminar.seminar_status && (this.seminar.seminar_status.ID === SEMINAR_STATUSES.OPENED || this.seminar.seminar_status.ID === SEMINAR_STATUSES.FILLED)" class="btn btn-primary m-2" @click.prevent="startSeminar()" value="Startuj seminar">
          <input v-if="this.seminar && this.seminar.seminar_status && this.seminar.seminar_status.ID === SEMINAR_STATUSES.IN_PROGRESS" class="btn btn-primary m-2" @click.prevent="finishSeminar()" value="Završi seminar">
        </div>
      </div>
    </form-tag>

    <div v-if="this.seminar && this.seminar.seminar_status && (this.seminar.seminar_status.ID === SEMINAR_STATUSES.IN_PROGRESS || this.seminar.seminar_status.ID === SEMINAR_STATUSES.CLOSED)">
      <hr>
      <div class="row">
        <div class="col-sm-2">
          <button class="btn btn-info text-white" @click="printStudentList()">Spisak polaznika</button>
        </div>
        <div class="col-sm-2">
          <button class="btn btn-info text-white" @click="printConfirmationStatement()">Izjava o pristanku</button>
        </div>
        <div class="col-sm-1">
          <button class="btn btn-info text-white">Prijava</button>
        </div>
        <div class="col-sm-1">
          <button class="btn btn-info text-white">Potvrda</button>
        </div>
        <div class="col-sm-1">
          <button class="btn btn-info text-white">Izjava</button>
        </div>
      </div>
      <hr>
      <h4>Dani seminara</h4>
      <div v-for="day in seminar.days" :key="day.number" @click="openSeminarDayEdit(day.ID)"
           class="border border-info bg-light d-inline-flex rounded m-2" style="width: 10%; height: 120px">
        <div class="m-1">
          <h6>Dan: {{ day.number }}</h6>
          <p style="font-size: 0.8em">{{getDateInMMDDYYYYFormat(day.date)}}</p>
          <hr>
          <p style="font-size: 0.7em; overflow: hidden">Tema: {{ day.name }}</p>

        </div>
      </div>
    </div>
  </div>
</template>

<script>
import FormTag from "@/components/forms/FormTag";
import vSelect from "vue-select";
import axios from "axios";
import router from "@/router";
import {apiMixin} from "@/mixins/apiMixin";
import {styleMixin} from "@/mixins/styleMixin";
import {useToast} from "vue-toastification";
import TextInput from "@/components/forms/TextInput.vue";
import {dateMixin} from "@/mixins/dateMixin";
import {commonMixin} from "@/mixins/commonMixin";
import {fileMixin} from "@/mixins/fileMixin";

export default {
  name: 'SeminarEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextInput, FormTag, vSelect},
  computed: {
    readonly() {
      return this.action === 'view';
    },
    percentFilled() {
      if (!this.seminarId) {
        return ;
      }
      return (this.seminar.trainees.length/this.seminar.class_room.max_students) *100;
    }
  },
  data() {
    return {
      seminar: {
        start_date: null,
        class_room: {},
        seminar_theme: null,
        seminar_status: null,
        trainees: [],
        days: []
      },
      action: "view",
      seminarId: "",
      selectedBaseSeminarType: null,
      location: null,
      seminarThemesByType: [],
      classRoomsByLocationId: [],
    }
  },
  methods: {
    async printStudentList() {
      await axios.get('/print/seminar/student-list/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var fileContent = JSON.parse(response.data.Data);
        var sampleArr = this.base64ToArrayBuffer(fileContent);
        const blob = new Blob([sampleArr], { type: 'application/pdf' });

        var iframe = document.createElement('iframe');
        iframe.src = URL.createObjectURL(blob);
        document.body.appendChild(iframe);

        URL.revokeObjectURL(iframe.src);
        iframe.contentWindow.print();
        iframe.setAttribute("hidden", "hidden");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async printConfirmationStatement() {
      await axios.get('/print/seminar/confirmation-statement/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var fileContent = JSON.parse(response.data.Data);
        var sampleArr = this.base64ToArrayBuffer(fileContent);
        const blob = new Blob([sampleArr], { type: 'application/pdf' });

        var iframe = document.createElement('iframe');
        iframe.src = URL.createObjectURL(blob);
        document.body.appendChild(iframe);

        URL.revokeObjectURL(iframe.src);
        iframe.contentWindow.print();
        iframe.setAttribute("hidden", "hidden");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    openSeminarDayEdit(dayId) {
     router.push({name: 'SeminarDayEdit', query: {id: dayId}});
    },
    async onSeminarTypeChange() {
      this.seminar.seminar_theme = null;
      this.seminarThemesByType = [];
      await this.getAllSeminarThemesByTypeId(this.selectedBaseSeminarType.ID);
    },
    async getAllSeminarThemesByTypeId(seminarBaseTypeId) {
      await axios.get('/seminar-types/themes/seminar-type/' + seminarBaseTypeId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarThemesByType = JSON.parse(response.data.Data);
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    onLocationChange() {
      this.classRoomsByLocationId = [];
      this.seminar.class_room = null;
      this.getAllClassRoomsLocationId(this.location.ID);
    },
    async getAllClassRoomsLocationId(locationId) {
      await axios.get('/class-rooms/location/' +locationId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.classRoomsByLocationId = JSON.parse(response.data.Data);
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async getSeminarById() {
      axios.get('/seminars/id/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminar = JSON.parse(response.data.Data);
        this.seminar.days.sort((a, b) => {
          return a.number - b.number;
        });
        this.location = this.seminar.class_room.location;
        this.location.address_place = this.seminar.class_room.location.address.place;
        this.seminar.start_date = this.getDateInMMDDYYYYFormat(this.seminar.start_date);
        this.selectedBaseSeminarType = this.seminar.seminar_theme.base_seminar_type;
        this.getAllSeminarThemesByTypeId(this.seminar.seminar_theme.base_seminar_type.ID);
        if (this.seminar.trainees == null) {
          this.seminar.trainees = [];
        }
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async submitHandler() {
      if (this.seminarId !== '') {
        await this.updateSeminar();
      } else {
        await this.createSeminar();
      }
      router.push("/seminars");
    },
    async createSeminar() {
      this.seminar.start_date = this.getBackendFormat(this.seminar.start_date);
      this.seminar.seminar_status = this.seminarStatuses.find(status => status.ID === this.SEMINAR_STATUSES.OPENED);
      await axios.post('/seminars/create', JSON.stringify(this.seminar)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreiran seminar!");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async updateSeminar() {
      this.seminar.start_date = this.getBackendFormat(this.seminar.start_date);
      await axios.post('/seminars/update', JSON.stringify(this.seminar)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran seminar!");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async createSeminarDays() {
      await axios.get('/seminar-days/create-all/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreirani dani seminara!");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    startSeminar() {
      this.seminar.seminar_status = this.seminarStatuses.find(ss => ss.ID == this.SEMINAR_STATUSES.IN_PROGRESS);
      this.updateSeminar();
      this.createSeminarDays();
    },
    finishSeminar() {
      this.seminar.seminar_status = this.seminarStatuses.find(ss => ss.ID == this.SEMINAR_STATUSES.CLOSED);
      this.updateSeminar();
    }
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    if (this.$route.query.id !== '') {
      this.seminarId = this.$route.query.id;
      await this.getSeminarById();
    }
    this.action = this.$route.query.action;
    await this.getAllLocations();
    await this.getAllBaseSeminarTypes();
    await this.getAllSeminarStatuses();
  }
}
</script>

<style scoped>
.shell {
  height: 20px;
  width: 250px;
  border: 1px solid #aaa;
  border-radius: 13px;

}

.bar {
  background: linear-gradient(to right, #007bff, #007bff);
  height: 20px;
  width: 15px;
  border-radius: 9px;
  span {
    float: right;
    color: #fff;
    font-size: 0.7em
  }
}

.styled-table {
  border-collapse: collapse;
  margin: 0px 0;
  font-family: sans-serif;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
}

.styled-table thead tr {
  background-color: #009879;
  color: #ffffff;
  text-align: left;
}

.styled-table th,
.styled-table td {
  padding: 1px 1px;
}

.styled-table tbody tr {
  border-bottom: 1px solid #dddddd;
}

.styled-table tbody tr:nth-of-type(even) {
  background-color: #f3f3f3;
}
</style>