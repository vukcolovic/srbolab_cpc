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
          <label :style="styleLabel" class="mb-1 mt-1">Tip seminara</label>
          <v-select
              v-model="seminar.base_seminar_type"
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

        <div  class="col-sm-6">
          <div v-if="action !== 'add'">

          </div>
        </div>
        <div class="col-sm-5">
          <input v-if="this.action === 'add'" class="btn btn-primary m-2" type="submit" value="Snimi">
          <input v-if="this.action === 'update'" class="btn btn-primary m-2" type="submit" value="Snimi">
        </div>
      </div>
    </form-tag>

    <div v-if="action !== 'add'">
      <hr>
      <h4>Dani seminara</h4>
      <div v-for="day in seminar.days" :key="day.number"
           class="border border-info bg-light d-inline-flex rounded m-2" style="width: 10%; height: 120px">
        <div class="m-1">
          <h5>Dan: {{ day.number }}</h5>
          <hr>
          <h5 style="font-size: 0.9em">Tema: {{ day.name }}</h5>

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

export default {
  name: 'SeminarEdit',
  mixins: [apiMixin, styleMixin, dateMixin],
  components: {TextInput, FormTag, vSelect},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      seminar: {
        start_date: null,
        class_room: null,
        base_seminar_type: null,
        seminar_theme: null,
        seminar_status: null,
        days: []
      },
      action: "view",
      seminarId: "",
      location: null,
      seminarThemesByType: [],
      classRoomsByLocationId: [],
    }
  },
  methods: {
    onSeminarTypeClear() {
      this.seminar.seminar_theme = null;
      this.seminarThemesByType = [];
    },
    async onSeminarTypeChange() {
      this.seminar.seminar_theme = null;
      this.seminarThemesByType = [];
      await this.getAllSeminarThemesByTypeId(this.seminar.base_seminar_type.ID);
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
      await axios.get('/locations/class-rooms/location/' +locationId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.classRoomsByLocationId = JSON.parse(response.data.Data);
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    // addSeminarDay() {
    //   this.seminarDay.seminar_id = parseInt(this.seminarId);
    //   this.seminarDay.date = new Date();
    //   axios.post('/seminar-days/create', JSON.stringify(this.seminarDay)).then((response) => {
    //     if (response.data === null || response.data.Status === 'error') {
    //       this.toast.error(response.data != null ? response.data.ErrorMessage : "");
    //       return;
    //     }
    //     this.toast.info("Uspešno kreiran dan seminara!");
    //     this.getSeminarById();
    //   }, (error) => {
    //     this.toast.error(error.message);
    //   });
    // },
    async getSeminarById() {
      axios.get('/seminars/id/' + this.seminarId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminar = JSON.parse(response.data.Data);
        this.location = this.seminar.class_room.location;
        this.location.address_place = this.seminar.class_room.location.address.place;
        this.seminar.start_date = this.getDateInMMDDYYYYFormat(this.seminar.start_date);
        this.getAllSeminarThemesByTypeId(this.seminar.base_seminar_type.ID);
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async submitHandler() {
      this.seminar.start_date = this.getBackendFormat(this.seminar.start_date);
      if (this.seminarId !== '') {
        await this.updateSeminar();
      } else {
        await this.createSeminar();
      }
    },
    async createSeminar() {
      this.seminar.seminar_status = this.seminarStatuses.find(status => status.code === "PENDING");
      await axios.post('/seminars/create', JSON.stringify(this.seminar)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreiran seminar!");
        router.push("/seminars");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async updateSeminar() {
      await axios.post('/seminars/update', JSON.stringify(this.seminar)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran seminar!");
        router.push("/seminars");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
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