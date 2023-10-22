<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 v-if="action === 'add'" class="mt-1">Dodavanje</h3>
        <h3 v-if="action === 'view'" class="mt-1">Pregled</h3>
        <h3 v-if="action === 'update'" class="mt-1">Ažuriranje</h3>
      </div>
    </div>
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-4">
          <div class="row">
            <!--            <div :class="[clientId ? 'col-sm-12' : 'col-sm-8']">-->
            <div class="col-sm-12">
              <text-input
                  v-model.trim="client.jmbg"
                  :readonly="readonly"
                  :required=true
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="JMBG"
                  min="13"
                  max="13"
                  name="jmbg"
                  type="text"
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
                  :readonly="readonly"
                  :required=true
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Ime"
                  name="name"
                  type="text">
              </text-input>
            </div>

            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.person.middle_name"
                  :readonly="readonly"
                  :required=false
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Ime jednog roditelja"
                  name="middleName"
                  type="text">
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.person.last_name"
              :readonly="readonly"
              :required=true
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="Prezime"
              name="lastName"
              type="text">
          </text-input>

          <div class="row">
            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.person.phone_number"
                  :readonly="readonly"
                  :required=false
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Broj telefona"
                  name="phone_number"
                  type="text">
              </text-input>
            </div>

            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.person.email"
                  :readonly="readonly"
                  :required=false
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Email"
                  name="email"
                  type="text">
              </text-input>
            </div>
          </div>

          <div class="row">
            <div class="col-sm-6">
              <text-input
                  v-model.trim="client.cpc_number"
                  :readonly="readonly"
                  :min=6
                  :max=6
                  :required="cpcRequired"
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Broj CPC kartice"
                  name="cpc_number"
                  type="text">
              </text-input>
            </div>

            <div class="col-sm-6">
              <label :style="styleLabelSmall" class="mb-2 mt-2">CPC datum isticanja</label>
              <Datepicker
                  v-model="client.cpc_date"
                  :disabled="readonly"
                  :style="dateStyleInput"
                  inputFormat="dd.MM.yyyy"
                  placeholder="dd.MM.yyyy"
                  typeable="true"
              />

            </div>
          </div>

          <text-input
              v-model.number="client.initial_completed_seminars"
              :readonly="readonly"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="Redni broj seminara na koji se prijavljuje"
              name="initial_completed_seminars"
              type="number">
          </text-input>

          <text-input
              v-model.trim="client.drive_licence"
              :readonly="readonly"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="Broj vozačke"
              name="drive_licence"
              type="text">
          </text-input>

          <div class="form-check form-switch">
            <input id="flexSwitchCheckDefault" v-model="showNote" class="form-check-input" role="switch"
                   type="checkbox">
            <label :style="styleInputSmall" class="form-check-label" for="flexSwitchCheckDefault">Napomena</label>
          </div>
          <div v-if="showNote">
            <text-area-input
                v-model.trim="client.comment"
                :readonly="readonly"
                :styleInput=styleInputSmall
                :styleLabel=styleLabelSmall
                label="Napomena:"
                name="comment"
                rows="2"
                type="text">
            </text-area-input>
          </div>

          <div class="row">
          <div class="col-sm-6">
            <div class="my-1">
              <label :style=styleLabelSmall for="verified">Klijent je verifikovan:&nbsp;&nbsp;</label>
              <input id="verified" v-model="client.verified" :hidden="readonly" type="checkbox"/>
            </div>
            <div class="my-1">
              <label :style=styleLabelSmall for="wait_seminar">Klijent čeka seminar:&nbsp;&nbsp;</label>
              <input id="wait_seminar" v-model="client.wait_seminar" :hidden="readonly" type="checkbox"/>
            </div>
          </div>
          <div class="col-sm-6">
            <div class="my-1">
              <label :style=styleLabelSmall for="c_licence">C Kategorija:&nbsp;&nbsp;</label>
              <input id="c_licence" v-model="client.c_licence" :hidden="readonly" type="checkbox"/>
            </div>
            <div class="my-1">
              <label :style=styleLabelSmall for="d_licence">D Kategorija:&nbsp;&nbsp;</label>
              <input id="d_licence" v-model="client.d_licence" :hidden="readonly" type="checkbox"/>
            </div>
          </div>
          </div>
        </div>

        <div class="col-sm-4">
          <text-input
              v-model.trim="client.address.place"
              :readonly="readonly"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="Mesto"
              name="place"
              type="text">
          </text-input>

          <div class="row">
            <div class="col-sm-9">
              <text-input
                  v-model.trim="client.address.street"
                  :readonly="readonly"
                  :required=false
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Ulica"
                  name="street"
                  type="text">
              </text-input>
            </div>
            <div class="col-sm-3">
              <text-input
                  v-model.trim="client.address.house_number"
                  :readonly="readonly"
                  :required=false
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Broj"
                  name="house_number"
                  type="text">
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.place_birth"
              :readonly="readonly"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="Mesto rođenja"
              name="place_birth"
              type="text">
          </text-input>

          <text-input
              v-model.trim="client.country_birth"
              :readonly="readonly"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="Država rođenja"
              name="country_birth"
              type="text">
          </text-input>

          <div class="row">
            <div class="my-1 col-sm-4">
              <label :style=styleLabelSmall for="resident">Državljanin:</label>
              <input id="resident" v-model="client.resident" :hidden="readonly" type="checkbox"/>
            </div>
            <div class="my-1 col-sm-8">
              <text-input
                  v-model.trim="client.second_citizenship"
                  :readonly="readonly"
                  :required=false
                  :styleInput=styleInputSmall
                  :styleLabel=styleLabelSmall
                  label="Drugo državljanstvo"
                  name="second_citizenship"
                  type="text">
              </text-input>
            </div>
          </div>

          <text-input
              v-model.trim="client.educational_profile"
              :readonly="readonly"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="Obrazovni profil"
              name="educational_profile"
              type="text">
          </text-input>

          <text-input
              v-model.trim="client.company_pib"
              :required=false
              :styleInput=styleInputSmall
              :styleLabel=styleLabelSmall
              label="PIB Firme"
              name="company_pib"
              type="text">
          </text-input>

          <label :style="styleLabelSmall" class="mb-1 mt-1">Firma</label>
          <v-select
              v-model="client.company"
              :disabled=readonly
              :options="companies"
              :style="styleInputSmall"
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

          <input id="fileId" ref="file" type="file" @change="uploadFile()"/>

        </div>
        <div class="col-sm-4" style="font-size: 0.7em">
          <div v-if="finishedSeminars.length > 0">
            <h6>Odslušani seminari</h6>
            <ul>
              <li v-for="seminarClient in finishedSeminars" :key="seminarClient.ID" style="list-style-type: none">
                {{ seminarClient.seminar.ID }}: {{ seminarClient.seminar.type }}
                {{ formatDateWithPoints(seminarClient.seminar.start_date) }} <span
                  :style="[seminarClient.pass ? {'color':'green'} : {'color':'red'} ]">{{
                  seminarClient.passedText
                }}</span>
              </li>
            </ul>
          </div>
          <div v-if="inProgressSeminars.length > 0">
            <h6>Aktuelni seminari</h6>
            <ul>
              <li v-for="seminarClient in inProgressSeminars" :key="seminarClient.ID" style="list-style-type: none;">
                {{ seminarClient.seminar.ID }}: {{ seminarClient.seminar.type }}
                {{ formatDateWithPoints(seminarClient.seminar.start_date) }}
              </li>
            </ul>
          </div>
          <div v-if="waitingSeminars.length > 0">
            <h6>Prijavljeni seminari</h6>
            <ul>
              <li v-for="seminarClient in waitingSeminars" :key="seminarClient.ID" style="list-style-type: none;">
                <button class="iconBtn" title="Obriši" @click.prevent="removeSeminar(seminarClient)">
                  <i class="fa fa-remove"></i>
                </button>
                <span v-if="seminarClient.payed" class="bg-success">
                  Plaćeno {{ formatDateWithPoints(seminarClient.pay_date) }}
                </span>
                <span v-if="!seminarClient.payed" class="bg-warning">
                  Nije plaćeno
                </span>
                | ID: {{seminarClient.seminar.ID }}:
                {{ seminarClient.seminar.type }} {{ formatDateWithPoints(seminarClient.seminar.start_date) }}
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
              placeholder="Traži"
              @option:selected="onOpenedSeminarChange">
          </v-select>

          <div v-if="selectedOpenSeminar">
            <div class="row my-1">
              <div class="col-sm-2">
                <div class="my-1">
                  <label :style=styleLabelSmall for="payed">Plaćeno:&nbsp;</label>
                  <input id="payed" v-model="selectedOpenSeminar.payed" :hidden="readonly" type="checkbox"/>
                </div>
              </div>
              <div class="col-sm-6" v-if="selectedOpenSeminar.payed">
                <text-input
                    v-model.trim="selectedOpenSeminar.payed_by"
                    :required=false
                    :styleInput=styleInputSmall
                    :styleLabel=styleLabelSmall
                    label="Platio"
                    name="payed_by"
                    type="text">
                </text-input>
              </div>
              <div class="col-sm-4" v-if="selectedOpenSeminar.payed">
                <label :style="styleLabelSmall" class="mb-1 mt-1">Datum plaćanja</label>
                <Datepicker
                    v-model="selectedOpenSeminar.pay_date"
                    :disabled="readonly"
                    :style="dateStyleInput"
                    inputFormat="dd.MM.yyyy"
                    placeholder="dd.MM.yyyy"
                    typeable="true"
                />
              </div>
            </div>
          </div>

        </div>
        <div>
          <input v-if="this.action === 'add'" class="btn btn-primary m-2" type="submit" value="Snimi">
          <input v-if="this.action === 'update'" class="btn btn-primary m-2" type="submit" value="Snimi">
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
import Datepicker from 'vue3-datepicker';

export default {
  name: 'ClientEdit',
  mixins: [fileMixin, styleMixin, apiMixin, commonMixin],
  components: {vSelect, FormTag, TextInput, TextAreaInput, Datepicker},
  computed: {
    readonly() {
      return this.action === 'view';
    },
    cpcRequired() {
      return this.selectedOpenSeminar && (this.selectedOpenSeminar.seminar_theme.base_seminar_type.code === "CYCLE" || this.selectedOpenSeminar.seminar_theme.base_seminar_type.code === "ADDITIONAL");
    }
  },
  data() {
    return {
      client: {
        person: {first_name: "", last_name: "", email: "", phone_number: ""},
        jmbg: "",
        address: {place: "", street: "", house_number: ""},
        // drive_licence: "",
        place_birth: "",
        country_birth: "Србија",
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
        wait_seminar: false,
        seminars: [],
        c_licence: true,
        d_licence: false,
      },
      showNote: false,
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
    onOpenedSeminarChange() {
      if(this.selectedOpenSeminar) {
        if (this.client && this.client.company && this.client.company.ID > 0) {
          this.selectedOpenSeminar.payed_by = this.client.company.name;
        } else {
          this.selectedOpenSeminar.payed_by = this.client.person.first_name + " " + this.client.person.last_name;
        }

      }
    },
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
      var day = this.client.jmbg.substring(0, 2);
      let dayInt = parseInt(day);
      if (!Number.isInteger(dayInt) || dayInt > 31 || dayInt < 1) {
        return "Jmbg nije validan, pogrešan dan rođenja";
      }
      var month = this.client.jmbg.substring(2, 4);
      let monthInt = parseInt(month);
      if (!Number.isInteger(monthInt) || monthInt > 12 || monthInt < 1) {
        return "Jmbg nije validan, pogrešan mesec rođenja";
      }
      var year = this.client.jmbg.substring(4, 6);
      let yearInt = parseInt(year);
      if (!Number.isInteger(yearInt)) {
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
        var location = this.selectedLocation ? this.selectedLocation.address.place : "";
        var seminarID = this.selectedOpenSeminar ? this.selectedOpenSeminar.ID : "";
        var rand = "";
        rand += Math.random();
        router.push("/client?action=update&id=" + foundClient.ID.toString() + "&location=" + location + "&seminar_id=" + seminarID + "&rand=" + rand);
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
        this.client.cpc_date = this.getFullDate(this.client.cpc_date);
        this.client.seminars.forEach(s => {
          s.seminar.type = this.getSeminarFullType(s.seminar.seminar_theme.base_seminar_type, s.seminar.seminar_theme);
          s.passedText = s.pass ? "Položio" : "Nije položio";
        });
        if (this.client.seminars) {
          this.finishedSeminars = this.client.seminars.filter(s => s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.CLOSED);
          this.inProgressSeminars = this.client.seminars.filter(s => s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.IN_PROGRESS);
          this.waitingSeminars = this.client.seminars.filter(s => (s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.OPENED || s.seminar.seminar_status.ID === this.SEMINAR_STATUSES.FILLED));

          this.openedSeminars = this.openedSeminars.filter((el) => !this.waitingSeminars.find(rm => (rm.seminar.ID === el.ID)));
        }
        if (this.client.documents == null) {
          this.client.documents = [];
        }
        if (this.client.seminars == null) {
          this.client.seminars = [];
        }
        if (this.client.company) {
          this.client.company.name_pib = this.client.company.name + "-" + this.client.company.pib;
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
      if (this.client.cpc_number.length > 0 && this.client.cpc_number.length != 6) {
        return "Broj cpc kartice mora imati 6 cifara!";
      }
      // if (!this.client.educational_profile && (!this.client.cpc_number && this.isDateEmpty(this.client.cpc_date))) {
      //   return "Obrazovni profil ili podaci o cpc kartici moraju biti popunjeni!";
      // }

      if (this.selectedOpenSeminar && (this.selectedOpenSeminar.seminar_theme.base_seminar_type.code === "CYCLE" || this.selectedOpenSeminar.seminar_theme.base_seminar_type.code === "ADDITIONAL") && !this.client.cpc_number) {
        return "Za dodatnu i periodičnu obuku broj cpc kartice mora biti popunjen.";
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
        var payDate = null;
        if (this.selectedOpenSeminar.pay_date) {
          payDate = this.selectedOpenSeminar.pay_date;
        }
        this.client.seminars.push({"client_id": this.client.ID, "seminar_id": this.selectedOpenSeminar.ID, payed: this.selectedOpenSeminar.payed, payed_by: this.selectedOpenSeminar.payed_by, pay_date: payDate});
        this.client.wait_seminar = false;
      }
      if (this.clientId) {
        await this.updateClient();
      } else {
        await this.createClient();
      }
    },
    async createClient() {
      await axios.post('/clients/create', JSON.stringify(this.client)).then((response) => {
        if (response.data == null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          this.client.cpc_date = this.getFullDate(this.client.cpc_date);
          return;
        }
        this.toast.info("Uspešno kreiran klijent.");
        var location = this.selectedLocation ? this.selectedLocation.address.place : "";
        var seminarID = this.selectedOpenSeminar ? this.selectedOpenSeminar.ID : "";
        var rand = "";
        rand += Math.random();
        router.push("/client?action=add&id=&location=" + location + "&seminar_id=" + seminarID + "&rand=" + rand);
      }, (error) => {
        this.errorToast(error, "/clients/create");
      });
    },
    async updateClient() {
      await axios.post('/clients/update', JSON.stringify(this.client)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran klijent.");
        var location = this.selectedLocation ? this.selectedLocation.address.place : "";
        var seminarID = this.selectedOpenSeminar ? this.selectedOpenSeminar.ID : "";
        var rand = "";
        rand += Math.random();
        router.push("/client?action=add&id=&location=" + location + "&seminar_id=" + seminarID + "&rand=" + rand);
      }, (error) => {
        this.errorToast(error, "/clients/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  }
  ,
  async mounted() {
    await this.getAllLocations();
    await this.getSeminarsByStatusCode("OPENED").then(result => this.openedSeminars = result);
    this.filteredAndOpenedSeminars = this.openedSeminars;
    await this.getAllCompanies();
    if (this.$route.query.id && this.$route.query.id !== '') {
      this.clientId = this.$route.query.id;
      await this.getClientById();
    }
    this.action = this.$route.query.action;
    if (this.$route.query.location && this.$route.query.location !== '') {
      this.selectedLocation = this.locations.find(s => s.address.place === this.$route.query.location);
      this.filteredAndOpenedSeminars = this.openedSeminars.filter(s => s.class_room.location.ID === this.selectedLocation.ID);
    }
    if (this.selectedLocation && this.$route.query.seminar_id && this.$route.query.seminar_id !== '') {
      this.selectedOpenSeminar = this.filteredAndOpenedSeminars.find(s => s.ID == this.$route.query.seminar_id);
    }
  }
}
</script>

<style scoped>
</style>