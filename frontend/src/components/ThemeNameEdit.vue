<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-3">
          <h4 v-if="action === 'add'" class="mt-2">Dodavanje</h4>
          <h4 v-if="action === 'update'" class="mt-2">Naziv teme</h4>
        </div>
      </div>

      <div class="row">
            <label :style="styleLabel" class="mb-1 mt-1">Tema seminara</label>
            <v-select
                v-model="seminarDayThemeName.seminar_theme"
                :disabled=readonly
                :options="seminarThemes"
                :style="styleInput"
                label="name"
                placeholder="Traži">
            </v-select>

        <text-input
            v-model.number="seminarDayThemeName.day_number"
            :readonly="readonly"
            :required=true
            label="Redni broj dana"
            name="day_number"
            type="number">
        </text-input>

          <text-input
              v-model="seminarDayThemeName.theme_name"
              :readonly="readonly"
              :required=true
              label="Naziv teme"
              name="theme_name"
              type="text">
          </text-input>

          <div class="col-sm-12">
            <input v-if="this.action === 'add'" class="btn btn-primary m-2" type="submit" value="Snimi">
            <input v-if="this.action === 'update'" class="btn btn-primary m-2" type="submit" value="Snimi">
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
import {commonMixin} from "@/mixins/commonMixin";
import {fileMixin} from "@/mixins/fileMixin";

export default {
  name: 'ThemeNameEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextInput, FormTag, vSelect},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      seminarDayThemeName: {
        seminar_theme: null,
        day_number: 0,
        theme_name: "",
      },
      action: "view",
      seminarThemeNameId: "",
    }
  },
  methods: {
    async getSeminarDayThemeNameById() {
      axios.get('/theme-names/id/' + this.seminarThemeNameId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarDayThemeName = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/theme-names/id");
      });
    },
    async submitHandler() {
      if (this.seminarThemeNameId !== '') {
        await this.updateSeminarDayThemeName();
      } else {
        await this.createSeminarDayThemeName();
      }
    },
    async createSeminarDayThemeName() {
      await axios.post('/theme-names/create', JSON.stringify(this.seminarDayThemeName)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarDayThemeName = JSON.parse(response.data.Data);
        this.toast.info("Uspešno kreiran naziv teme!");
        router.push("/theme-names");
      }, (error) => {
        this.errorToast(error, "/theme-names/create");
      });
    },
    async updateSeminarDayThemeName() {
      await axios.post('/theme-names/update', JSON.stringify(this.seminarDayThemeName)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran naziv teme!");
      }, (error) => {
        this.errorToast(error, "/theme-names/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    this.getAllSeminarThemes();
    if (this.$route.query.id !== '') {
      this.seminarThemeNameId = this.$route.query.id;
      await this.getSeminarDayThemeNameById();
    }
    this.action = this.$route.query.action;
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