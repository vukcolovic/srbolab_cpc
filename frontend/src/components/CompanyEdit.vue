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
              v-model.trim="company.name"
              label="Naziv"
              type="text"
              name="name"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="company.pib"
              label="PIB"
              type="text"
              name="pib"
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
  name: 'CompanyEdit',
  mixins: [fileMixin, commonMixin],
  components: {FormTag, TextInput},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      company: {name: "", pib: ""},
      action: "",
    }
  },
  methods: {
    async getCompanyById() {
      axios.get('/companies/id/' + this.companyId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.company = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/companies/id");
      });
    },
    async submitHandler() {
      if (this.companyId != undefined) {
        await this.updateCompany();
      } else {
        await this.createCompany();
      }
    },
    async createCompany() {
      await axios.post('/companies/create', JSON.stringify(this.company)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreirana firma.");
        router.push("/companies");
      }, (error) => {
        this.errorToast(error, "/companies/create");
      });
    },
    async updateCompany() {
      await axios.post('/companies/update', JSON.stringify(this.company)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažurirana firma.");
        router.push("/companies");
      }, (error) => {
        this.errorToast(error, "/companies/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  mounted() {
    if (this.$route.query.id !== '') {
      this.companyId = this.$route.query.id;
      this.getCompanyById();
    }
    this.action = this.$route.query.action;
  }
}
</script>