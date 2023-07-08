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
        <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
          <div class="row">
          <div class="col-sm-5">
            <text-input
                v-model="seminar.start_date"
                label="Početak seminara (MM/DD/YYYY)"
                type="date"
                name="start"
                :required=true
                :readonly="readonly">
            </text-input>

            <label class="mb-1 mt-1" :style="styleLabel">Lokacija</label>
            <v-select
                :disabled=readonly
                :style="styleInput"
                v-model="seminar.location"
                :options="locations"
                label="address_place"
                placeholder="Traži">
            </v-select>

            <label class="mb-1 mt-1" :style="styleLabel">Tip seminara</label>
            <v-select
                :disabled=readonly
                :style="styleInput"
                v-model="seminar.seminar_type"
                :options="seminarTypes"
                label="name"
                placeholder="Traži">
            </v-select>
          </div>

          <div class="col-sm-5">

          </div>
            <div class="col-sm-5">
          <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Snimi">
          <input type="submit" v-if="this.action === 'update'" class="btn btn-primary m-2" value="Snimi">
            </div>
            </div>
        </form-tag>
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
      seminar: {start_date: null, location: null, seminar_type: null, seminar_status: null},
      action: "view",
      seminarId: "",
    }
  },
  methods: {
    async getSeminarById() {
        axios.get('/seminars/id/' + this.seminarId).then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            this.toast.error(response.data != null ? response.data.ErrorMessage : "");
            return;
          }
          this.seminar = JSON.parse(response.data.Data);
          this.seminar.location.address_place = this.seminar.location.address.place;
          this.seminar.start_date = this.getDateInMMDDYYYYFormat(this.seminar.start_date);
          console.log(this.seminar);
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
    await this.getAllSeminarTypes();
    await this.getAllSeminarStatuses();
  }
}
</script>