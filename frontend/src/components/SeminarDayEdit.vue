<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <button class="btn btn-primary" @click="backToSeminar()"><i class="fa fa-arrow-left"></i>Vrati se na seminar
        </button>
      </div>
    </div>
    <div class="row">
      <div class="col-sm-4">
        <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
          <label style="font-size: 1.2em">{{ seminarDay.seminar_theme }} - Dan {{ seminarDay.number }}</label>
          <text-input
              v-model.trim="seminarDay.name"
              :readonly="readonly"
              :required=false
              label="Naziv"
              name="name"
              type="text">
          </text-input>

          <text-input
              v-model="seminarDay.date"
              :readonly="readonly"
              :required=true
              label="Datum"
              name="date"
              type="date">
          </text-input>
          <input class="btn btn-primary m-2" type="submit" value="Snimi">

        </form-tag>

      </div>

      <div class="col-sm-5">
        <h6>Evidencija prisustva</h6>
        <ul>
          <li v-for="pres in seminarDay.presence" :key="pres.client_id" style="list-style-type: none;">
            <input id="verified" v-model="pres.presence" :hidden="readonly" type="checkbox"/>
            {{ pres.client.person.first_name }} {{ pres.client.person.last_name }}
          </li>
        </ul>

      </div>

      <div class="col-sm-3">
        <label :style=styleLabel>Dokumenta: </label>
        <input id="fileId" ref="file" type="file" @change="uploadFile()"/>
        <ul>
          <li v-for="(doc, index) in seminarDay.documents" :key="index" style="list-style-type: none;">
            <label for="index">&nbsp; {{ doc.name }}</label>
            <button class="iconBtn" title="Obriši" @click.prevent="removeFile(index)">
              <i class="fa fa-remove"></i>
            </button>

            <button class="iconBtn" title="Preuzmi" @click.prevent="downloadFile(index)">
              <i class="fa fa-download"></i>
            </button>
          </li>
        </ul>
      </div>
    </div>

    <hr>
    <div class="row">
      <div class="col-sm-1">
        <button class="btn btn-secondary text-white" @click="printMuster()">Prozivka</button>
      </div>
      <div class="col-sm-2">
        <button class="btn btn-secondary text-white" @click="printTeacherEvidence()">Dnevnik predavača</button>
      </div>
    </div>
    <hr>
    <div>
      <h5>Časovi</h5>
      <table class="styled-table" style="width: 80%">
        <thead>
        <tr class="bg-primary text-white">
          <td style="width: 5%;">R.B.</td>
          <td style="width: 55%;">Tema</td>
          <td style="width: 40%;">Predavač</td>
        </tr>
        </thead>
        <tbody>
        <tr v-for="cls in seminarDay.classes" :key="cls.ID">
          <td class="p-1">{{ cls.number}}</td>
          <td>
            <input type="text" size="40" v-model="cls.name">
          </td>
          <td>
          <v-select
              v-model="cls.teacher"
              :options="users"
              :style="styleInput"
              label="full_name"
              placeholder="Traži">
          </v-select>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from "@/router";
import {fileMixin} from "@/mixins/fileMixin";
import {useToast} from "vue-toastification";
import {apiMixin} from "@/mixins/apiMixin";
import {dateMixin} from "@/mixins/dateMixin";
import {styleMixin} from "@/mixins/styleMixin";
import vSelect from "vue-select";

export default {
  name: 'SeminarDayEdit',
  mixins: [fileMixin, apiMixin, dateMixin, styleMixin],
  components: {vSelect, FormTag, TextInput},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      seminarDay: {
        ID: 0,
        date: null,
        number: 0,
        name: "",
        seminar_id: 0,
        seminar: null,
        seminar_theme: "",
        presence: [],
        documents: [],
        classes: [],
      },
      seminarDayId: 0,
    }
  },
  methods: {
    downloadFile(i) {
      const arr = this.seminarDay.documents[i].content.split(',')
      var sampleArr = this.base64ToArrayBuffer(arr[1]);
      const blob = new Blob([sampleArr])

      const link = document.createElement('a')
      link.href = URL.createObjectURL(blob)
      link.download = this.seminarDay.documents[i].name
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
        this.seminarDay.documents.push({content: fileString, name: file.name});
      }
      reader.readAsDataURL(file);
    },
    removeFile(i) {
      this.seminarDay.documents.splice(i, 1);
    },
    async printTeacherEvidence() {
      await axios.get('/print/seminar/teacher-evidence/' + this.seminarDayId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var fileContent = JSON.parse(response.data.Data);
        var sampleArr = this.base64ToArrayBuffer(fileContent);
        const blob = new Blob([sampleArr], {type: 'application/pdf'});

        var iframe = document.createElement('iframe');
        iframe.src = URL.createObjectURL(blob);
        document.body.appendChild(iframe);

        URL.revokeObjectURL(iframe.src);
        iframe.contentWindow.print();
        iframe.setAttribute("hidden", "hidden");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async printMuster() {
      await axios.get('/print/seminar/muster/' + this.seminarDayId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        var fileContent = JSON.parse(response.data.Data);
        var sampleArr = this.base64ToArrayBuffer(fileContent);
        const blob = new Blob([sampleArr], {type: 'application/pdf'});

        var iframe = document.createElement('iframe');
        iframe.src = URL.createObjectURL(blob);
        document.body.appendChild(iframe);

        URL.revokeObjectURL(iframe.src);
        iframe.contentWindow.print();
        iframe.setAttribute("hidden", "hidden");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    backToSeminar() {
      router.push({name: 'SeminarEdit', query: {id: this.seminarDay.seminar_id, action: "update"}});
    },
    async getSeminarDayById() {
      axios.get('/seminar-days/id/' + this.seminarDayId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarDay = JSON.parse(response.data.Data);
        this.seminarDay.date = this.getDateInMMDDYYYYFormat(this.seminarDay.date);
        this.seminarDay.seminar_theme = this.getSeminarFullType(this.seminarDay.seminar.seminar_theme.base_seminar_type, this.seminarDay.seminar.seminar_theme);
        if (this.seminarDay.documents == null) {
          this.seminarDay.documents = [];
        }
        this.seminarDay.classes.forEach(cls => {
          if (cls.teacher) {
            cls.teacher.full_name = cls.teacher.person.first_name + " "+ cls.teacher.person.last_name;
          }
        });
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async submitHandler() {
      this.seminarDay.date = this.getBackendFormat(this.seminarDay.date);
      await axios.post('/seminar-days/update', JSON.stringify(this.seminarDay)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarDay.date = this.getDateInMMDDYYYYFormat(this.seminarDay.date);
        this.toast.info("Uspešno ažuriran seminar dan.");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    this.getAllUsers();
    this.seminarDayId = this.$route.query.id;
    await this.getSeminarDayById();
  }
}
</script>