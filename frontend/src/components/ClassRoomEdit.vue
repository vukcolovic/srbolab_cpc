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
              v-model.trim="classRoom.name"
              label="Naziv"
              type="text"
              name="name"
              :required=true
              :readonly="readonly">
          </text-input>

          <label :style="styleLabel" class="mb-1 mt-1">Lokacija</label>
          <v-select
              v-model="classRoom.location"
              :disabled=readonly
              :options="locations"
              :style="styleInput"
              label="address_place"
              placeholder="Traži">
          </v-select>

          <text-input
              v-model.number="classRoom.max_students"
              label="Maksimalan broj mesta"
              type="number"
              name="max_students"
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
import vSelect from "vue-select";
import {styleMixin} from "@/mixins/styleMixin";
import {apiMixin} from "@/mixins/apiMixin";

export default {
  name: 'ClassRoomEdit',
  mixins: [fileMixin, styleMixin, apiMixin],
  components: {vSelect, FormTag, TextInput},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      classRoom : {location:null, max_students: 0, name: ""},
      action: "",
      classRoomId: 0,
    }
  },
  methods: {
    async getClassRoomById() {
      axios.get('/class-rooms/id/' + this.classRoomId.toString()).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.classRoom = JSON.parse(response.data.Data);
        this.classRoom.location.address_place = this.classRoom.location.address.place;
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async submitHandler() {
      if (this.classRoomId) {
        await this.updateClassRoom();
      } else {
        await this.createClassRoom();
      }
    },
    async createClassRoom() {
      await axios.post('/class-rooms/create', JSON.stringify(this.classRoom)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreirana učionica.");
        router.push("/class-rooms");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async updateClassRoom() {
      await axios.post('/class-rooms/update', JSON.stringify(this.classRoom)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažurirana učionica.");
        router.push("/class-rooms");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  mounted() {
    this.getAllLocations();
    if (this.$route.query.id !== '') {
      this.classRoomId = this.$route.query.id;
      this.getClassRoomById();
    }
    this.action = this.$route.query.action;
  }
}
</script>