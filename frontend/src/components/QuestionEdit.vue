<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <h3 v-if="action === 'add'" class="mt-2">Dodavanje</h3>
        <h3 v-if="action === 'update'" class="mt-2">Ažuriranje</h3>
        <hr>
      </div>
      <div class="col-sm-8">
            <label :style="styleLabel" class="mb-1 mt-1">Tema seminara</label>
            <v-select
                v-model="question.seminar_theme"
                :disabled=readonly
                :options="seminarThemes"
                :style="styleInput"
                label="name"
                placeholder="Traži">
            </v-select>

          <text-area-input
              v-model="question.content"
              :readonly="readonly"
              :required=true
              label="Pitanje"
              name="start"
              type="text">
          </text-area-input>

        <hr>
        <h6>Ponuđeni odgovori:</h6>
        <div class="row" v-for="answer in question.answers" :key="answer.ID">
          <div class="col-sm-1">
            {{answer.letter}})
          </div>
          <div class="col-sm-10">
            <text-area-input
                v-model="answer.content"
                :readonly="readonly"
                :required=true
                label=""
                name="start"
                type="text">
            </text-area-input>
          </div>
          <div class="col-sm-1">
            <label>Da li je tačan?</label>
            <input id="correct" type="checkbox" :hidden="readonly" v-model="answer.correct" />
          </div>
        </div>
      </div>
    </div>
      <div>
        <input type="submit" v-if="this.action === 'add'" class="btn btn-primary m-2" value="Snimi">
        <input type="submit" v-if="this.action === 'update'" class="btn btn-primary m-2" value="Snimi">
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
import {dateMixin} from "@/mixins/dateMixin";
import {commonMixin} from "@/mixins/commonMixin";
import {fileMixin} from "@/mixins/fileMixin";
import router from "@/router";

export default {
  name: 'QuestionEdit',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextAreaInput, FormTag, vSelect},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      question: {
        content: "",
        seminar_theme: null,
        answers: [],
      },
      action: "view",
      questionId: "",
    }
  },
  methods: {
    async getQuestionById() {
      axios.get('/questions/id/' + this.questionId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.question = JSON.parse(response.data.Data);
        if (this.question.answers == null) {
          this.question.answers = [];
        }
      }, (error) => {
        this.errorToast(error, "/questions/id");
      });
    },
    async submitHandler() {
      var correctAnswers = 0;
      for (let i = 0; i < this.question.answers.length; i++) {
        if (this.question.answers[i].correct) {
          correctAnswers++;
        }
      }
      if (correctAnswers > 1) {
        this.toast.warning("Broj tačnih odgovora ne može biti veći od 1");
        return;
      }
      if (correctAnswers === 0) {
        this.toast.warning("Broj tačnih odgovora ne može biti manju od 1");
        return;
      }
      if (this.questionId) {
        await this.updateQuestion();
      } else {
        await this.createQuestion();
      }
      router.push("/questions");
    },
    async createQuestion() {
      await axios.post('/questions/create', JSON.stringify(this.question)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno kreirano pitanje!");
      }, (error) => {
        this.errorToast(error, "/questions/create");
      });
    },
    async updateQuestion() {
      await axios.post('/questions/update', JSON.stringify(this.question)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažurirano pitanje!");
      }, (error) => {
        this.errorToast(error, "/questions/update");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    if (this.$route.query.id !== '') {
      this.questionId = this.$route.query.id;
      await this.getQuestionById();
    }
    this.action = this.$route.query.action;
    if (this.action === "add") {
      this.question.answers = [{letter: "a"}, {letter: "b"}, {letter: "c"}, {letter: "d"}];
    }
    await this.getAllSeminarThemes();
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