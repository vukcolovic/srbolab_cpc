<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-11 mx-auto">
          <h3 v-if="action === 'add'" class="mt-2">Dodavanje</h3>
          <h3 v-if="action === 'update'" class="mt-2">Ažuriranje</h3>
          <hr>
        </div>
        <div class="col-sm-12">
          <text-input
              v-model="test.name"
              :readonly="readonly"
              :required=true
              label="Naziv testa"
              name="name"
              type="text">
          </text-input>

          <label :style="styleLabel" class="mb-1 mt-1">Tema seminara</label>
          <v-select
              v-model="test.seminar_theme"
              :disabled=readonly
              :options="seminarThemes"
              :style="styleInput"
              label="name"
              placeholder="Traži"
              @option:selected="onSeminarThemeChange">
          </v-select>

          <hr>
          <h6>Pitanja</h6>
          <div v-for="(question, index) in questions" :key="question.ID" class="row">
            <div class="col-sm-1 d-flex aligns-items-center justify-content-center">
              <input :id="index" :value="question" :disabled="readonly" type="checkbox" v-model="test.questions" />
            </div>
            <div class="col-sm-11">
              <text-area-input
                  v-model="question.content"
                  :readonly=true
              ></text-area-input>
            </div>
          </div>
        </div>
      </div>
      <div>
        <input v-if="this.action === 'add'" class="btn btn-primary m-2" type="submit" value="Snimi">
        <input v-if="this.action === 'update'" class="btn btn-primary m-2" type="submit" value="Snimi">
      </div>
    </form-tag>
  </div>
</template>

<script>
import FormTag from "@/components/forms/FormTag";
import vSelect from "vue-select";
import axios from "axios";
import {apiMixin} from "@/mixins/apiMixin";
import {styleMixin} from "@/mixins/styleMixin";
import {useToast} from "vue-toastification";
import TextAreaInput from "@/components/forms/TextAreaInput.vue";
import TextInput from "@/components/forms/TextInput.vue";
import {dateMixin} from "@/mixins/dateMixin";
import {commonMixin} from "@/mixins/commonMixin";
import {fileMixin} from "@/mixins/fileMixin";
import router from "@/router";

export default {
  name: 'TestEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextAreaInput, TextInput, FormTag, vSelect},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      test: {
        name: "",
        seminar_theme: null,
        questions: [],
      },
      questions: [],
      action: "view",
      testId: "",
    }
  },
  methods: {
    async onSeminarThemeChange() {
      this.test.questions = [];
      if (!this.test.seminar_theme) {
        this.questions = [];
        return;
      }
      await this.getQuestionsBySeminarTheme();
    },
    async getQuestionsBySeminarTheme() {
      await axios.get('/questions/list/seminar-theme/' + this.test.seminar_theme.ID).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.questions = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/questions/list/seminar-theme");
      });
    },
    async getTestById() {
      await axios.get('/tests/id/' + this.testId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.test = JSON.parse(response.data.Data);
        if (this.test.questions == null) {
          this.test.questions = [];
        }
      }, (error) => {
        this.errorToast(error, "/tests/id");
      });
    },
    async submitHandler() {
      if (this.testId) {
        await this.updateTest();
      } else {
        await this.createTest();
      }
    },
    async createTest() {
      await axios.post('/tests/create', JSON.stringify(this.test)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreiran test!");
        router.push("/tests");
      }, (error) => {
        this.errorToast(error, "/tests/create");
      });
    },
    async updateTest() {
      await axios.post('/tests/update', JSON.stringify(this.test)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran test!");
        router.push("/tests");
      }, (error) => {
        this.errorToast(error, "/tests/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    await this.getAllSeminarThemes();
    if (this.$route.query.id !== '') {
      this.testId = this.$route.query.id;
      await this.getTestById();
      await this.getQuestionsBySeminarTheme();
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