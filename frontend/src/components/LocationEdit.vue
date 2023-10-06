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
              v-model.trim="location.address.place"
              label="Mesto"
              type="text"
              name="place"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="location.code"
              label="Skraćeni kod"
              type="text"
              name="code"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="location.address.street"
              label="Ulica"
              type="text"
              name="street"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="location.address.house_number"
              label="Broj stana"
              type="text"
              name="house_number"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="location.address.post_code"
              label="Poštanski kod"
              type="text"
              name="post_code"
              :required=true
              :readonly="readonly">
          </text-input>

          <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Snimi">
          <input type="submit" v-if="this.action === 'update'" class="btn btn-primary m-2" value="Snimi">
        </div>

      </div>
    </form-tag>
  </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from "@/router";
import {fileMixin} from "@/mixins/fileMixin";
import {useToast} from "vue-toastification";
import {commonMixin} from "@/mixins/commonMixin";

export default {
  name: 'LocationEdit',
  mixins: [fileMixin, commonMixin],
  components: {FormTag, TextInput},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      location: {ID: 0, code: "", address: {place:"", street: "", house_number: "", post_code: ""}},
      action: "",
    }
  },
  methods: {
    async getLocationById() {
      axios.get('/locations/id/' + this.locationId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.location = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/locations/id");
      });
    },
    async submitHandler() {
      if (this.locationId != undefined) {
        await this.updateLocation();
      } else {
        await this.createLocation();
      }
    },
    async createLocation() {
      await axios.post('/locations/create', JSON.stringify(this.location)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreirana lokacija.");
        router.push("/locations");
      }, (error) => {
        this.errorToast(error, "/locations/create");
      });
    },
    async updateLocation() {
      await axios.post('/locations/update', JSON.stringify(this.location)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažurirana lokacija.");
        router.push("/locations");
      }, (error) => {
        this.errorToast(error, "/locations/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  mounted() {
    if (this.$route.query.id !== '') {
      this.locationId = this.$route.query.id;
      this.getLocationById();
    }
    this.action = this.$route.query.action;
  }
}
</script>