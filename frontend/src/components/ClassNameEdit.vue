<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-3">
          <h4 v-if="action === 'add'" class="mt-2">Dodavanje</h4>
          <h4 v-if="action === 'update'" class="mt-2">Naziv časa</h4>
        </div>
      </div>

      <div class="row">
            <label :style="styleLabel" class="mb-1 mt-1">Tema seminara</label>
            <v-select
                v-model="seminarClassName.seminar_theme"
                :disabled=readonly
                :options="seminarThemes"
                :style="styleInput"
                label="name"
                placeholder="Traži">
            </v-select>

        <text-input
            v-model.number="seminarClassName.day_number"
            :readonly="readonly"
            :required=true
            label="Redni broj dana"
            name="day_number"
            type="number">
        </text-input>

        <text-input
            v-model.number="seminarClassName.class_number"
            :readonly="readonly"
            :required=true
            label="Redni broj časa"
            name="class_number"
            type="number">
        </text-input>

          <text-input
              v-model="seminarClassName.class_name"
              :readonly="readonly"
              :required=true
              label="Naziv časa"
              name="class_name"
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
  name: 'ClassNameEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextInput, FormTag, vSelect},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      seminarClassName: {
        seminar_theme: null,
        day_number: 0,
        class_number: 0,
        class_name: "",
      },
      action: "view",
      seminarClassNameId: "",
    }
  },
  methods: {
    async getSeminarClassNameById() {
      axios.get('/class-names/id/' + this.seminarClassNameId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarClassName = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/class-names/id");
      });
    },
    async submitHandler() {
      if (this.seminarClassNameId !== '') {
        await this.updateSeminarClassName();
      } else {
        await this.createSeminarClassName();
      }
    },
    async createSeminarClassName() {
      await axios.post('/class-names/create', JSON.stringify(this.seminarClassName)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarClassName = JSON.parse(response.data.Data);
        this.toast.info("Uspešno kreiran naziv časa!");
        router.push("/class-names");
      }, (error) => {
        this.errorToast(error, "/class-names/create");
      });
    },
    async updateSeminarClassName() {
      await axios.post('/class-names/update', JSON.stringify(this.seminarClassName)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran naziv časa!");
      }, (error) => {
        this.errorToast(error, "/class-names/update");
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
      this.seminarClassNameId = this.$route.query.id;
      await this.getSeminarClassNameById();
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