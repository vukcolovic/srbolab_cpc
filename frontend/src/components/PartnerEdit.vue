<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-3">
          <h4 v-if="action === 'add'" class="mt-2">Dodavanje</h4>
          <h4 v-if="action === 'update'" class="mt-2">Partner</h4>
        </div>
      </div>

      <div class="row">

        <text-input
            v-model="partner.name"
            :readonly="readonly"
            :required=true
            label="Naziv"
            name="name"
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
  name: 'PartnerEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextInput, FormTag},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      partner: {
        name: "",
      },
      action: "view",
      partnerId: "",
    }
  },
  methods: {
    async getPartnerById() {
      axios.get('/partners/id/' + this.partnerId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.partner = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/partners/id");
      });
    },
    async submitHandler() {
      if (this.partnerId !== '') {
        await this.updatePartner();
      } else {
        await this.createPartner();
      }
    },
    async createPartner() {
      await axios.post('/partners/create', JSON.stringify(this.partner)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.partner = JSON.parse(response.data.Data);
        this.toast.info("Uspešno kreiran partner!");
        router.push("/partners");
      }, (error) => {
        this.errorToast(error, "/partners/create");
      });
    },
    async updatePartner() {
      await axios.post('/partners/update', JSON.stringify(this.partner)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran partner!");
      }, (error) => {
        this.errorToast(error, "/partners/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    if (this.$route.query.id !== '') {
      this.partnerId = this.$route.query.id;
      await this.getPartnerById();
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