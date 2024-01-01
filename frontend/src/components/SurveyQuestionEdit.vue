<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 v-if="action === 'add'" class="mt-2">Dodavanje</h3>
        <hr>
      </div>
      <div class="col-sm-8">

          <text-area-input
              v-model="question.content"
              :readonly="readonly"
              :required=true
              label="Pitanje"
              name="question"
              type="text">
          </text-area-input>
      </div>
    </div>
      <div>
        <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Snimi">
      </div>
    </form-tag>
    </div>
</template>

<script>
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import {apiMixin} from "@/mixins/apiMixin";
import {styleMixin} from "@/mixins/styleMixin";
import {useToast} from "vue-toastification";
import TextAreaInput from "@/components/forms/TextAreaInput.vue";
import {dateMixin} from "@/mixins/dateMixin";
import {commonMixin} from "@/mixins/commonMixin";
import {fileMixin} from "@/mixins/fileMixin";
import router from "@/router";

export default {
  name: 'SurveyQuestionEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextAreaInput, FormTag},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      question: {
        content: "",
      },
      action: "view",
      questionId: "",
    }
  },
  methods: {
    async getQuestionById() {
      axios.get('/survey-questions/id/' + this.questionId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.question = JSON.parse(response.data.Data);
      }, (error) => {
        this.errorToast(error, "/survey-questions/id");
      });
    },
    async submitHandler() {
        await this.createQuestion();
      router.push("/survey-questions");
    },
    async createQuestion() {
      await axios.post('/survey-questions/create', JSON.stringify(this.question)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreirano pitanje!");
      }, (error) => {
        this.errorToast(error, "/survey-questions/create");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
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