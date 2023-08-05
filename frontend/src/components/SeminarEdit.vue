<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-3">
        <h4 v-if="action === 'add'" class="mt-2">Dodavanje</h4>
        <h4 v-if="action === 'update'" class="mt-2">Seminar</h4>
      </div>
      <div class="col-sm-7">
        <div class="shell mt-1">
          <div :style="{ width: percentFilled  + '%' }" class="bar">
          </div>
          Polaznika: <span>{{ seminar.trainees.length }} / {{ seminar.class_room.max_students }}</span>
        </div>
      </div>
      <hr>
    </div>

      <div class="row">
        <div class="col-sm-4">
          <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
          <div v-if="seminar.seminar_status && action !== 'add'">
            Status seminara: {{ seminar.seminar_status.name }}
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
              v-model="location"
              :disabled=readonly
              :options="locations"
              :style="styleInput"
              label="address_place"
              placeholder="Traži"
              @option:selected="onLocationChange">
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
    </form-tag>
        </div>

        <div v-if="action !== 'add'" class="col-sm-8" style="font-size: 0.8em">
          <div>
            <h5>Spisak polaznika</h5>
            <div class="mb-1">
              <input placeholder="JMBG" style="max-width: 80px; font-size: 0.8em" type="text" id="jmbg" name="jmbg" v-model="filter.jmbg" /> 
              <input placeholder="FIRMA" style="max-width: 80px; font-size: 0.8em; margin-right: 5px" type="text" id="company" name="company" v-model="filter.company" />
            </div>

            <table class="styled-table">
              <thead>
              <tr class="bg-primary text-white">
                <td style="width: 35%;">Ime i prezime</td>
                <td style="width: 20%;">JMBG</td>
                <td style="width: 35%;">Firma</td>
                <td style="width: 10%;">Plaćeno</td>
              </tr>
              </thead>
              <tbody>
              <tr v-for="trainee in filteredClients" :key="trainee.client_id">
                <td class="p-1">{{ trainee.client.person.first_name }} {{ trainee.client.person.last_name }}</td>
                <td class="p-1">{{ trainee.client.jmbg }}</td>
                <td class="p-1">{{ trainee.client.company.name }}</td>
                <td :class="[trainee.payed ? 'bg-success' : 'bg-danger']">{{ trainee.payed ? 'DA' : 'NE' }}</td>
              </tr>
              </tbody>
            </table>
          </div>

        </div>

        <div class="row"></div>
        <div class="col-sm-5">
          <input v-if="this.action === 'add'" class="btn btn-primary m-2" type="submit" value="Snimi">
          <input v-if="this.action === 'update'" class="btn btn-primary m-2" type="submit" value="Snimi">
          <input
              v-if="this.seminar && this.seminar.seminar_status && (this.seminar.seminar_status.ID === SEMINAR_STATUSES.OPENED || this.seminar.seminar_status.ID === SEMINAR_STATUSES.FILLED)"
              class="btn btn-primary m-2" value="Startuj seminar" @click.prevent="startSeminar()">
          <input
              v-if="this.seminar && this.seminar.seminar_status && this.seminar.seminar_status.ID === SEMINAR_STATUSES.IN_PROGRESS"
              class="btn btn-primary m-2" value="Završi seminar" @click.prevent="finishSeminar()">
        </div>
      </div>

    <div
        v-if="this.seminar && this.seminar.seminar_status && (this.seminar.seminar_status.ID === SEMINAR_STATUSES.IN_PROGRESS || this.seminar.seminar_status.ID === SEMINAR_STATUSES.CLOSED)">
      <hr>
      <div class="row">
        <div class="col-sm-2">
          <button class="btn btn-secondary text-white" @click="printStudentList()">Spisak polaznika</button>
        </div>
        <div class="col-sm-2">
          <button class="btn btn-secondary text-white" @click="printConfirmationStatement()">Izjava o pristanku</button>
        </div>
        <div class="col-sm-1">
          <button class="btn btn-secondary text-white">Prijava</button>
        </div>
        <div class="col-sm-1">
          <button class="btn btn-secondary text-white" @click="printConfirmations()">Potvrda</button>
        </div>
        <div class="col-sm-3">
          <button class="btn btn-secondary text-white" @click="printStatementOfReceving()">Izjava o preuzimanju</button>
        </div>
      </div>
      <hr>
      <div class="col-sm-3">
        <label :style=styleLabel>Dokumenta: </label>
        <input id="fileId" ref="file" type="file" @change="uploadFile()"/>
        <ul>
          <li v-for="(doc, index) in seminar.documents" :key="index" style="list-style-type: none;">
            <label for="index">&nbsp; {{ doc.name }}</label>
            <button class="iconBtn" title="Obriši" @click.prevent="removeFile(index)">
              <i class="fa fa-remove"></i>
            </button>

            <button class="iconBtn" title="Preuzmi" @click.prevent="downloadFile(index)">
              <i class="fa fa-download"></i>
            </button>
          </li>
        </ul>
      </div>
      <hr>
      <h4>Dani seminara</h4>
      <div v-for="day in seminar.days" :key="day.number" class="border border-info bg-light d-inline-flex rounded m-2"
           style="width: 10%; height: 120px" @click="openSeminarDayEdit(day.ID)">
        <div class="m-1">
          <h6>Dan: {{ day.number }}</h6>
          <p style="font-size: 0.8em">{{ getDateInMMDDYYYYFormat(day.date) }}</p>
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
        return;
      }
      return (this.seminar.trainees.length / this.seminar.class_room.max_students) * 100;
    },
    filteredClients() {
      return this.seminar.trainees.filter((obj) => {
        var companyName =  obj.client.company ? obj.client.company.name : "";
        return companyName.includes(this.filter.company) && obj.client.jmbg.includes(this.filter.jmbg);
      });
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
        documents: [],
        days: []
      },
      action: "view",
      seminarId: "",
      selectedBaseSeminarType: null,
      location: null,
      seminarThemesByType: [],
      classRoomsByLocationId: [],
      filter: {jmbg: "", company: ""}
    }
  },
  methods: {
    downloadFile(i) {
      const arr = this.seminar.documents[i].content.split(',')
      var sampleArr = this.base64ToArrayBuffer(arr[1]);
      const blob = new Blob([sampleArr])

      const link = document.createElement('a')
      link.href = URL.createObjectURL(blob)
      link.download = this.seminar.documents[i].name
      link.click()
      URL.revokeObjectURL(link.href)
    },
    uploadFile() {
      const file = this.$refs.file.files[0];
      if (file == null) {
        return;
      }
      const reader = new FileReader()
      reader.onloadend = () => {
        const fileString = reader.result;
        this.seminar.documents.push({content: fileString, name: file.name});
      }
      reader.readAsDataURL(file);
    },
    removeFile(i) {
      this.seminar.documents.splice(i, 1);
    },
    async printStudentList() {
      await axios.get('/print/seminar/student-list/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var fileContent = JSON.parse(response.data.Data);
        var sampleArr = this.base64ToArrayBuffer(fileContent);
        const blob = new Blob([sampleArr], {type: 'application/pdf'});

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
        const blob = new Blob([sampleArr], {type: 'application/pdf'});

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
    async printConfirmations() {
      await axios.get('/print/seminar/confirmation/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var fileContent = JSON.parse(response.data.Data);
        var sampleArr = this.base64ToArrayBuffer(fileContent);
        const blob = new Blob([sampleArr], {type: 'application/pdf'});

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
    async printStatementOfReceving() {
      await axios.get('/print/seminar/confirmation-receive/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var fileContent = JSON.parse(response.data.Data);
        var sampleArr = this.base64ToArrayBuffer(fileContent);
        const blob = new Blob([sampleArr], {type: 'application/pdf'});

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
      await axios.get('/class-rooms/location/' + locationId).then((response) => {
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
        if (this.seminar.documents == null) {
          this.seminar.documents = [];
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
  width: 500px;
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