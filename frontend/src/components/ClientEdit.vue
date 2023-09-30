<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 v-if="action === 'add'" class="mt-1">Dodavanje</h3>
        <h3 v-if="action === 'view'" class="mt-1">Pregled</h3>
        <h3 v-if="action === 'update'" class="mt-1">Ažuriranje</h3>
      </div>
    </div>
    <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
      <div class="row">
        <div class="col-sm-4">
          <div class="row">
<!--            <div :class="[clientId ? 'col-sm-12' : 'col-sm-8']">-->
              <div class="col-sm-12">
          <text-input
              v-model.trim="client.jmbg"
              label="JMBG"
              type="text"
              name="jmbg"
              :required=true
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              @focusout="onJmbgFocusOut">
          </text-input>
            </div>
<!--            <div class="col-sm-4 mt-4" v-if="!clientId">-->
<!--              <button class="iconBtn" title="Nađi" @click.prevent="getClientByJMBG()">-->
<!--                <i class="fa fa-search"></i>-->
<!--              </button>-->
<!--            </div>-->
          </div>

          <div class="row">
            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.person.first_name"
                  label="Ime"
                  type="text"
                  name="name"
                  :required=true
                  :readonly="readonly"
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall>
              </text-input>
            </div>

            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.person.middle_name"
                  label="Ime jednog roditelja"
                  type="text"
                  name="middleName"
                  :required=true
                  :readonly="readonly"
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.person.last_name"
              label="Prezime"
              type="text"
              name="lastName"
              :required=true
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <div class="row">
            <div class="col-sm-6">
          <text-input
              v-model.trim="client.person.phone_number"
              label="Broj telefona"
              type="text"
              name="phone_number"
              :required=false
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>
            </div>

            <div class="col-sm-6">
          <text-input
              v-model.trim="client.person.email"
              label="Email"
              type="text"
              name="email"
              :required=false
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
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
                  :readonly="readonly"
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall>
              </text-input>
            </div>

            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.cpc_date"
                  label="CPC datum izdavanja"
                  type="date"
                  name="cpc_date"
                  :required=false
                  :readonly="readonly"
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.number="client.initial_completed_seminars"
              label="Broj prethodno odlušanih kurseva"
              type="number"
              name="initial_completed_seminars"
              :required=false
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <text-input
              v-model.trim="client.drive_licence"
              label="Broj vozačke"
              type="text"
              name="drive_licence"
              :required=true
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <text-area-input
              v-model.trim="client.comment"
              label="Napomena:"
              type="text"
              rows="2"
              name="comment"
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-area-input>

          <div class="my-1">
            <label :style=styleLabelSmall for="verified">Klijent je verifikovan:&nbsp;&nbsp;</label>
            <input id="verified" type="checkbox" :hidden="readonly" v-model="client.verified" />
          </div>
          <div class="my-1">
            <label :style=styleLabelSmall for="wait_seminar">Klijent čeka seminar:&nbsp;&nbsp;</label>
            <input id="wait_seminar" type="checkbox" :hidden="readonly" v-model="client.wait_seminar" />
          </div>
        </div>

        <div class="col-sm-4">
          <text-input
              v-model.trim="client.address.place"
              label="Mesto"
              type="text"
              name="place"
              :required=true
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <div class="row">
            <div class="col-sm-9">
              <text-input
                  v-model.trim="client.address.street"
                  label="Ulica"
                  type="text"
                  name="street"
                  :required=true
                  :readonly="readonly"
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall>
              </text-input>
            </div>
            <div class="col-sm-3">
              <text-input
                v-model.trim="client.address.house_number"
                label="Broj"
                type="text"
                name="house_number"
                :required=true
                :readonly="readonly"
                :styleInput=styleInputSmall
                :styleLabel=styleLabelSmall>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.place_birth"
              label="Mesto rođenja"
              type="text"
              name="place_birth"
              :required=true
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <text-input
              v-model.trim="client.country_birth"
              label="Država rođenja"
              type="text"
              name="country_birth"
              :required=true
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <div class="row">
            <div class="my-1 col-sm-4">
              <label :style=styleLabelSmall for="resident">Državljanin:</label>
              <input id="resident" type="checkbox" :hidden="readonly" v-model="client.resident" />
            </div>
            <div class="my-1 col-sm-8">
              <text-input
                  v-model.trim="client.second_citizenship"
                  label="Drugo državljanstvo"
                  type="text"
                  name="second_citizenship"
                  :required=false
                  :readonly="readonly"
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall>
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.educational_profile"
              label="Obrazovni profil"
              type="text"
              name="educational_profile"
              :required=false
              :readonly="readonly"
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <text-input
              v-model.trim="client.company_pib"
              label="PIB Firme"
              type="text"
              name="company_pib"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall>
          </text-input>

          <label :style="styleLabelSmall" class="mb-1 mt-1">Firma</label>
          <v-select
              v-model="client.company"
              :disabled=readonly
              :style="styleInputSmall"
              :options="companies"
              label="name_pib"
              placeholder="Traži">
          </v-select>

          <label :style=styleLabelSmall>Dokumenta: </label>
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
        <div class="col-sm-4" style="font-size: 0.7em">
          <div v-if="finishedSeminars.length > 0">
            <h6>Odslušani seminari</h6>
            <ul>
              <li style="list-style-type: none" v-for="seminarClient in finishedSeminars" :key="seminarClient.ID">
                {{seminarClient.seminar.ID}}: {{seminarClient.seminar.type}} {{getDateInMMDDYYYYFormat(seminarClient.seminar.start_date)}} <span :style="[seminarClient.pass ? {'color':'green'} : {'color':'red'} ]">{{seminarClient.passedText}}</span>
              </li>
            </ul>
          </div>
          <div v-if="inProgressSeminars.length > 0">
            <h6>Aktuelni seminari</h6>
            <ul>
              <li style="list-style-type: none;" v-for="seminarClient in inProgressSeminars" :key="seminarClient.ID">
                {{seminarClient.seminar.ID}}: {{seminarClient.seminar.type}} {{getDateInMMDDYYYYFormat(seminarClient.seminar.start_date)}}
              </li>
            </ul>
          </div>
          <div v-if="waitingSeminars.length > 0">
            <h6>Prijavljeni seminari</h6>
            <ul>
              <li style="list-style-type: none;" v-for="seminarClient in waitingSeminars" :key="seminarClient.ID">
                <button class="iconBtn" title="Obriši" @click.prevent="removeSeminar(seminarClient)">
                  <i class="fa fa-remove"></i>
                </button>
                <span v-if="seminarClient.payed" class="bg-success">Plaćeno</span><span v-if="!seminarClient.payed" class="bg-warning">Nije plaćeno</span>{{seminarClient.seminar.ID}}: {{seminarClient.seminar.type}} {{getDateInMMDDYYYYFormat(seminarClient.seminar.start_date)}}
              </li>
            </ul>
          </div>
          <label :style="styleLabelSmall" class="mb-1 mt-1">Lokacije</label>
          <v-select
              v-model="selectedLocation"
              :disabled=readonly
              :options="locations"
              :style="styleInputSmall"
              label="address_place"
              placeholder="Traži"
              @option:selected="onLocationChange">
          </v-select>

          <label :style="styleLabelSmall" class="mb-1 mt-1">Otvoreni seminari</label>
          <v-select
              v-model="selectedOpenSeminar"
              :disabled=readonly
              :options="filteredAndOpenedSeminars"
              :style="styleInputSmall"
              label="base_info"
              placeholder="Traži">
          </v-select>
        </div>
        <div>
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
import TextAreaInput from "@/components/forms/TextAreaInput";
import axios from "axios";
import router from "@/router";
import {fileMixin} from "@/mixins/fileMixin";
import {useToast} from "vue-toastification";
import {styleMixin} from "@/mixins/styleMixin";
import vSelect from "vue-select";
import {apiMixin} from "@/mixins/apiMixin";
import {commonMixin} from "@/mixins/commonMixin";

export default {
  name: 'ClientEdit',
  mixins: [fileMixin, styleMixin, apiMixin, commonMixin],
  components: {vSelect, FormTag, TextInput, TextAreaInput},
  computed: {
    readonly() {
      return this.action === 'view';
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
        country_birth: "Srbija",
        documents: [],
        company: null,
        company_pib: "",
        comment: "",
        resident: true,
        second_citizenship: "",
        cpc_number: "",
        cpc_date: null,
        educational_profile: "",
        verified: true,
        initial_completed_seminars: 0,
        wait_seminar: true,
        seminars: []
      },
      finishedSeminars: [],
      inProgressSeminars: [],
      waitingSeminars: [],
      selectedOpenSeminar: null,
      openedSeminars: [],
      filteredAndOpenedSeminars: [],
      action: "",
      clientId: "",
      selectedLocation: null,
    }
  },
  methods: {
    onJmbgFocusOut() {
      const errMsg = this.jmbgValidation();
      if (errMsg) {
        this.toast.warning(errMsg);
        return;
      }
      this.getClientByJMBG();
    },
    jmbgValidation() {
      if (this.client.jmbg.length != 13) {
        return "Jmbg mora imati 13 cifara!";
      }
      var day = this.client.jmbg.substring(0,2);
      let dayInt = parseInt(day);
      if (!Number.isInteger(dayInt) || dayInt > 31 || dayInt < 1) {
        return "Jmbg nije validan, pogrešan dan rođenja";
      }
      var month = this.client.jmbg.substring(2,4);
      let monthInt = parseInt(month);
      if (!Number.isInteger(monthInt) || monthInt > 12 || monthInt < 1) {
        return "Jmbg nije validan, pogrešan mesec rođenja";
      }
      var year = this.client.jmbg.substring(4,6);
      let yearInt = parseInt(year);
      if (!Number.isInteger(yearInt)  ) {
        return "Jmbg nije validan, pogrešna godina rođenja";
      }
    },
    onLocationChange() {
      this.filteredAndOpenedSeminars = this.openedSeminars;
      this.filteredAndOpenedSeminars = this.openedSeminars.filter(s => s.class_room.location.ID === this.selectedLocation.ID);
    },
    getClientByJMBG() {
      axios.get('/clients/jmbg/' + this.client.jmbg).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          // this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var foundClient = JSON.parse(response.data.Data);
        router.push("/client?action=update&id=" + foundClient.ID.toString());
      }, (error) => {
        this.errorToast(error, "/clients/jmbg");
      });
    },
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
    removeSeminar(seminar) {
      const index = this.client.seminars.indexOf(seminar);
      if (index > -1) { // only splice array when item is found
        this.client.seminars.splice(index, 1); // 2nd parameter means remove one item only
      }
      const idx = this.waitingSeminars.indexOf(seminar);
      if (idx > -1) { // only splice array when item is found
        this.waitingSeminars.splice(idx, 1); // 2nd parameter means remove one item only
      }
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
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.client = JSON.parse(response.data.Data);
        this.client.cpc_date = this.getDateInMMDDYYYYFormat(this.client.cpc_date);
        this.client.seminars.forEach(s => {
          s.seminar.type = this.getSeminarFullType(s.seminar.seminar_theme.base_seminar_type, s.seminar.seminar_theme);
          s.passedText = s.pass ? "Položio" : "Nije položio";
        });
        if (this.client.seminars) {
          this.finishedSeminars = this.client.seminars.filter(s => s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.CLOSED);
          this.inProgressSeminars = this.client.seminars.filter(s => s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.IN_PROGRESS);
          this.waitingSeminars = this.client.seminars.filter(s => (s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.OPENED || s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.FILLED));

          this.openedSeminars = this.openedSeminars.filter( ( el ) => !this.waitingSeminars.find(rm => (rm.seminar.ID === el.ID)));
        }
        if (this.client.documents == null) {
          this.client.documents = [];
        }
        if (this.client.seminars == null) {
          this.client.seminars = [];
        }
      }, (error) => {
        this.errorToast(error, "/clients/id");
      });
    },
    validateClient() {
      if (!this.client.jmbg) {
        return "Polje jmbg mora biti popunjeno!";
      }
      if (this.client.jmbg.length != 13) {
        return "Jmbg mora imati 13 cifara!";
      }
      if (!this.client.verified && this.selectedOpenSeminar) {
        return "Ne možete prijaviti kurs za klijenta koji nije verifikovan!";
      }
      if (!this.client.educational_profile && (!this.client.cpc_number && this.isDateEmpty(this.client.cpc_date))) {
        return "Obrazovni profil ili podaci o cpc kartici moraju biti popunjeni!";
      }

      return "";
    },
    async submitHandler() {
      const errMsg = this.validateClient();
      if (errMsg) {
        this.toast.warning(errMsg);
        return;
      }
      if (this.selectedOpenSeminar) {
        this.client.seminars.push({"client_id": this.client.ID, "seminar_id": this.selectedOpenSeminar.ID});
        this.client.wait_seminar =false;
      }
      if (this.clientId) {
        await this.updateClient();
      } else {
        await this.createClient();
      }
    },
    async createClient() {
      this.client.cpc_date = this.getBackendFormat(this.client.cpc_date);
      await axios.post('/clients/create', JSON.stringify(this.client)).then((response) => {
        if (response.data == null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          this.client.cpc_date = this.getDateInMMDDYYYYFormat(this.client.cpc_date);
          return;
        }
        this.toast.info("Uspešno kreiran klijent.");
        router.push("/clients");
      }, (error) => {
        this.errorToast(error, "/clients/create");
      });
    },
    async updateClient() {
      this.client.cpc_date = this.getBackendFormat(this.client.cpc_date);
      await axios.post('/clients/update', JSON.stringify(this.client)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.client.cpc_date = this.getDateInMMDDYYYYFormat(this.client.cpc_date);
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran klijent.");
        router.push("/clients");
      }, (error) => {
        this.errorToast(error, "/clients/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    await this.getAllLocations();
    await this.getSeminarsByStatusCode("OPENED").then(result => this.openedSeminars = result);
    this.filteredAndOpenedSeminars = this.openedSeminars;
    await this.getAllCompanies();
    if (this.$route.query.id !== '') {
      this.clientId = this.$route.query.id;
      this.getClientById();
    }
    this.action = this.$route.query.action;
  }
}
</script>

<style scoped>
</style>