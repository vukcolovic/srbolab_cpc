<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 v-if="action === 'add'" class="mt-2">Dodavanje</h3>
        <hr>
      </div>
      <div class="col-sm-8">

          <text-input
              v-model="survey.name"
              :readonly="readonly"
              :required=true
              label="Naziv"
              name="name"
              type="text">
          </text-input>

          <label :style=styleLabelSmall for="active">Anketa je aktivna:&nbsp;&nbsp;</label>
          <input id="verified" v-model="survey.active" :disabled="readonly" type="checkbox"/>
      </div>

      <hr>
          <h6>Pitanja</h6>
          <div v-for="(question, index) in questions" :key="question.ID" class="row">
            <div class="col-sm-1 d-flex aligns-items-center justify-content-center">
              <input :id="index" :value="question" :disabled="readonly" type="checkbox" v-model="survey.questions" />
            </div>
            <div class="col-sm-11">
              <text-area-input
                  v-model="question.content"
                  :readonly=true
              ></text-area-input>
            </div>
          </div>
    </div>
      <div>
        <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Snimi">
      </div>
    </form-tag>
    </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import TextAreaInput from "@/components/forms/TextAreaInput.vue";
import {apiMixin} from "@/mixins/apiMixin";
import {styleMixin} from "@/mixins/styleMixin";
import {useToast} from "vue-toastification";
import {dateMixin} from "@/mixins/dateMixin";
import {commonMixin} from "@/mixins/commonMixin";
import {fileMixin} from "@/mixins/fileMixin";
import router from "@/router";

export default {
  name: 'SurveyEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {FormTag, TextInput, TextAreaInput},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      survey: {
        name: "",
        questions: [],
        active: false,
      },
      action: "view",
      surveyId: "",
      questions: [],
    }
  },
  methods: {
    async getSurveyById() {
      axios.get('/surveys/id/' + this.surveyId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.survey = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/surveys/id");
      });
    },
    async submitHandler() {
        await this.createSurvey();
      router.push("/surveys");
    },
    async createSurvey() {
      await axios.post('/surveys/create', JSON.stringify(this.survey)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreirana anketa!");
      }, (error) => {
        this.errorToast(error, "/surveys/create");
      });
    },
    async getSurveyQuestons() {
      this.isLoading = true;
      await axios.get('/survey-questions/list').then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.questions = JSON.parse(response.data.Data);
        console.log(this.questions);
      }, (error) => {
        this.errorToast(error, "/surveys/list");
      });

      this.isLoading = false;
    }
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    if (this.$route.query.id !== '') {
      this.surveyId = this.$route.query.id;
      await this.getSurveyById();
    }
    this.action = this.$route.query.action;
    this.getSurveyQuestons();
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