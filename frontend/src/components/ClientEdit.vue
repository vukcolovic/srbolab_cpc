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
              v-model.trim="client.person.first_name"
              label="Ime"
              type="text"
              name="name"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.person.middle_name"
              label="Ime jednog roditelja"
              type="text"
              name="middleName"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.person.last_name"
              label="Prezime"
              type="text"
              name="lastName"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.person.phone_number"
              label="Broj telefona"
              type="text"
              name="phone_number"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.person.email"
              label="Email"
              type="text"
              name="email"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.jmbg"
              label="JMBG"
              type="text"
              name="jmbg"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.drive_licence"
              label="Broj vozačke"
              type="text"
              name="drive_licence"
              :required=false
              :readonly="readonly">
          </text-input>

          <div class="my-1">
            <label for="verified">Klijent je verifikovan:</label>
            <input id="verified" type="checkbox" :hidden="readonly" v-model="client.verified" />
          </div>

        </div>

        <div class="col-sm-5">
          <text-input
              v-model.trim="client.address.place"
              label="Mesto"
              type="text"
              name="place"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.address.street"
              label="Ulica"
              type="text"
              name="street"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.address.house_number"
              label="Broj"
              type="text"
              name="house_number"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.place_birth"
              label="Mesto rođenja"
              type="text"
              name="place_birth"
              :required=false
              :readonly="readonly">
          </text-input>

          <text-input
              v-model.trim="client.country_birth"
              label="Država rođenja"
              type="text"
              name="country_birth"
              :required=false
              :readonly="readonly">
          </text-input>

          <label>Dokumenta: </label>

          <ul>
            <li v-for="(doc, index) in client.documents" :key="index" style="list-style-type: none;">
              <label for="index">&nbsp; {{ doc.name }}</label>
              <button class="iconBtn" title="Obriši" @click.prevent="removeFile(index)">
                <i class="fa fa-remove"></i>
              </button>

              <button class="iconBtn" title="Preuzmi" @click.prevent="downloadFile(index)">
                <i class="fa fa-download"></i>
              </button>
            </li>
          </ul>

          <input id="fileId" type="file" ref="file" @change="uploadFile()"/>

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
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from "@/router";
import {fileMixin} from "@/mixins/fileMixin";
import {useToast} from "vue-toastification";

export default {
  name: 'ClientEdit',
  mixins: [fileMixin],
  components: {FormTag, TextInput},
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
      client: {
        person: {first_name: "", last_name: "", email: "", phone_number: ""},
        jmbg: "",
        address: {place: "", street: "", house_number: ""},
        drive_licence: "",
        place_birth: "",
        country_birth: "",
        documents: [],
        verified: true
      },
      action: "",
    }
  },
  methods: {
    downloadFile(i) {
      const arr = this.client.documents[i].content.split(',')
      var sampleArr = this.base64ToArrayBuffer(arr[1]);
      const blob = new Blob([sampleArr])

      const link = document.createElement('a')
      link.href = URL.createObjectURL(blob)
      link.download = this.client.documents[i].name
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
        this.client.documents.push({content: fileString, name: file.name});
      }
      reader.readAsDataURL(file);
    },
    removeFile(i) {
      this.client.documents.splice(i, 1);
    },
    async getClientById() {
      axios.get('/clients/id/' + this.clientId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.client = JSON.parse(response.data.Data);
        if (this.client.documents == null) {
          this.client.documents = [];
        }
      }, (error) => {
        this.toast.error(error);
      });
    },
    async submitHandler() {
      if (this.clientId != undefined) {
        await this.updateClient();
      } else {
        await this.createClient();
      }
    },
    async createClient() {
      await axios.post('/clients/create', JSON.stringify(this.client)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.toast.info("Uspešno kreiran klijent.");
        router.push("/clients");
      }, (error) => {
        this.toast.error(error);
      });
    },
    async updateClient() {
      await axios.post('/clients/update', JSON.stringify(this.client)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data.ErrorMessage);
          return;
        }
        this.toast.info("Uspešno ažuriran klijent.");
        router.push("/clients");
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
      this.clientId = this.$route.query.id;
      this.getClientById();
    }
    this.action = this.$route.query.action;
  }
}
</script>