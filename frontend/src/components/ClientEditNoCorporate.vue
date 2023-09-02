<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 class="mt-1">Prijava za slušanje seminara</h3>
        <p>Nakon što se prijavite vrlo brzo će vas kontaktirati osoblje firme Srbolab d.o.o., kako bi dogovorili termin i mesto slušanja seminara.</p>
      </div>
    </div>
    <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
      <div class="row">
        <div class="col-sm-6">
          <div class="row">
          <text-input
              v-model.trim="client.jmbg"
              label="JMBG"
              type="text"
              name="jmbg"
              :required=true
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>
            </div>

          <div class="row">
            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.person.first_name"
                  label="Ime"
                  type="text"
                  name="name"
                  :required=true
                  :styleInput=styleInput
                  :styleLabel=styleLabel>
              </text-input>
            </div>

            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.person.middle_name"
                  label="Ime jednog roditelja"
                  type="text"
                  name="middleName"
                  :required=true
                  :styleInput=styleInput
                  :styleLabel=styleLabel>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.person.last_name"
              label="Prezime"
              type="text"
              name="lastName"
              :required=true
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <div class="row">
            <div class="col-sm-6">
          <text-input
              v-model.trim="client.person.phone_number"
              label="Broj telefona"
              type="text"
              name="phone_number"
              :required=true
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>
            </div>

            <div class="col-sm-6">
          <text-input
              v-model.trim="client.person.email"
              label="Email"
              type="text"
              name="email"
              :required=false
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>
            </div>
          </div>

          <div class="row">
            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.cpc_number"
                  label="Broj CPC kartice"
                  type="text"
                  name="cpc_number"
                  :required=false
                  :styleInput=styleInput
                  :styleLabel=styleLabel>
              </text-input>
            </div>

            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.cpc_date"
                  label="CPC datum izdavanja"
                  type="date"
                  name="cpc_date"
                  :required=false
                  :styleInput=styleInput
                  :styleLabel=styleLabel>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.number="client.initial_completed_seminars"
              label="Broj prethodno odlušanih kurseva"
              type="number"
              name="initial_completed_seminars"
              :required=false
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <text-input
              v-model.trim="client.drive_licence"
              label="Broj vozačke"
              type="text"
              name="drive_licence"
              :required=true
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <text-area-input
              v-model.trim="client.comment"
              label="Napomena:"
              type="text"
              rows="2"
              name="comment"
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-area-input>

        </div>

        <div class="col-sm-6">
          <text-input
              v-model.trim="client.address.place"
              label="Mesto"
              type="text"
              name="place"
              :required=true
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <div class="row">
            <div class="col-sm-9">
              <text-input
                  v-model.trim="client.address.street"
                  label="Ulica"
                  type="text"
                  name="street"
                  :required=true
                  :styleInput=styleInput
                  :styleLabel=styleLabel>
              </text-input>
            </div>
            <div class="col-sm-3">
              <text-input
                v-model.trim="client.address.house_number"
                label="Broj"
                type="text"
                name="house_number"
                :required=true
                :styleInput=styleInput
                :styleLabel=styleLabel>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.place_birth"
              label="Mesto rođenja"
              type="text"
              name="place_birth"
              :required=true
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <text-input
              v-model.trim="client.country_birth"
              label="Država rođenja"
              type="text"
              name="country_birth"
              :required=true
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <div class="row">
            <div class="my-1 col-sm-4">
              <label :style=styleLabel for="resident">Državljanin:</label>
              <input id="resident" type="checkbox" v-model="client.resident" />
            </div>
            <div class="my-1 col-sm-8">
              <text-input
                  v-model.trim="client.second_citizenship"
                  label="Drugo državljanstvo"
                  type="text"
                  name="second_citizenship"
                  :required=false
                  :styleInput=styleInput
                  :styleLabel=styleLabel>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.educational_profile"
              label="Obrazovni profil"
              type="text"
              name="educational_profile"
              :required=false
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <text-input
              v-model.trim="client.company_pib"
              label="PIB Firme"
              type="text"
              name="company_pib"
              :required=false
              :styleInput=styleInput
              :styleLabel=styleLabel>
          </text-input>

          <label :style=styleLabel>Dokumenta: </label>
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
        <div class="row "></div>
        <div class="col-sm-3">
          <input type="submit" class="btn btn-primary m-2" value="Snimi">
        </div>
      </div>
    </form-tag>
  </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import TextAreaInput from "@/components/forms/TextAreaInput";
import axios from "axios";
import {fileMixin} from "@/mixins/fileMixin";
import {useToast} from "vue-toastification";
import {styleMixin} from "@/mixins/styleMixin";
import {apiMixin} from "@/mixins/apiMixin";
import {commonMixin} from "@/mixins/commonMixin";

export default {
  name: 'ClientEditNoCorporate',
  mixins: [fileMixin, styleMixin, apiMixin, commonMixin],
  components: {FormTag, TextInput, TextAreaInput},
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
        company_pib: "",
        comment: "",
        resident: true,
        second_citizenship: "",
        cpc_number: "",
        cpc_date: null,
        educational_profile: "",
        verified: false,
        initial_completed_seminars: 0,
        wait_seminar: false,
        seminars: []
      },
      finishedSeminars: [],
      inProgressSeminars: [],
      waitingSeminars: [],
      selectedOpenSeminar: null,
      openedSeminars: [],
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
    validateClient() {
      if (!this.client.jmbg) {
        return "Polje jmbg mora biti popunjeno!";
      }
      if (this.client.jmbg.length != 13) {
        return "Jmbg mora imati 13 cifara!";
      }
      if (!this.client.educational_profile && !(this.client.cpc_number && this.isDateEmpty(this.client.cpc_date))) {
        return "Obrazovni profil ili podaci o cpc kartici moraju biti popunjeni";
      }

      return "";
    },
    async submitHandler() {
      const errMsg = this.validateClient();
      if (errMsg) {
        this.toast.warning(errMsg);
        return;
      }
        await this.createClient();
    },
    async createClient() {
      this.client.cpc_date = this.getBackendFormat(this.client.cpc_date);
      await axios.post('/clients/create-not-verified', JSON.stringify(this.client)).then((response) => {
        if (response.data == null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno snimljeno.");
      }, (error) => {
        this.errorToast(error, "/clients/create-not-verified");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    this.isCorporateIp();
  },
  create() {

  }
}
</script>

<style scoped>
</style>