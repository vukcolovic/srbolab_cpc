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
            <label class="mb-1 mt-1" :style="styleLabel">Lokacija</label>
            <v-select
                :disabled=readonly
                :style="styleInput"
                v-model="seminar.location"
                :options="locations"
                label="place"
                placeholder="Traži">
            </v-select>

            <label class="mb-1 mt-1" :style="styleLabel">Status</label>
            <v-select
                :disabled=readonly
                :style="styleInput"
                v-model="seminar.location"
                :options="locations"
                label="place"
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
import axios from "axios";
import router from "@/router";
import {apiMixin} from "@/mixins/apiMixin";
import {styleMixin} from "@/mixins/styleMixin";
import {useToast} from "vue-toastification";

export default {
  name: 'SeminarEdit',
  mixins: [apiMixin, styleMixin],
  components: {FormTag},
  computed: {
      readonly() {
        if (this.action === 'view') {
          return true;
        }
        return false;
    },
    },
  data() {
    return {
      seminar: {start: null, location: null, seminar_type: null, seminar_status: null},
      action: "view",
      seminarId: "",
    }
  },
  methods: {
    async getSeminarById() {
        axios.get('/seminars/id/' + this.seminarId).then((response) => {
          if (response.data === null || response.data.Status === 'error') {
            this.toast.error(response.data.ErrorMessage);
            return;
          }
          this.seminar = JSON.parse(response.data.Data);
        }, (error) => {
          this.toast.error(error);
        });
    },
    async submitHandler() {
      if (this.seminarId !== '') {
        await this.updateSeminar();
      } else {
        await this.createSeminar();
      }
    },
    async createSeminar() {
      await axios.post('/seminars/create', JSON.stringify(this.seminar)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.toast.info("Uspešno kreiran seminar!");
        router.push("/seminars");
      }, (error) => {
        this.toast.error(error);
      });
    },
    async updateSeminar() {
      await axios.post('/seminars/update', JSON.stringify(this.seminar)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.toast.info("Uspešno ažuriran seminar!");
        router.push("/users");
      }, (error) => {
        this.toast.error(error);
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  mounted() {
    if (this.$route.query.id !== '') {
      this.seminarId = this.$route.query.id;
      this.getSeminarById();
    }
    this.action = this.$route.query.action;
    this.getAllLocations();
  }
}
</script>