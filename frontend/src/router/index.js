import { createRouter, createWebHistory } from 'vue-router'
import BodyComponent from './../components/BodyComponent.vue'
import LoginComponent from './../components/LoginComponent.vue'
import ClientsList from '../components/ClientsList.vue'
import ClientEdit from './../components/ClientEdit.vue'
import UsersList from '../components/UsersList.vue'
import UserEdit from './../components/UserEdit.vue'
import SupportComponent from './../components/SupportComponent.vue'
import SeminarsList from '../components/SeminarsList.vue'
import SeminarEdit from './../components/SeminarEdit.vue'
import CompaniesList from '../components/CompaniesList.vue'
import CompanyEdit from './../components/CompanyEdit.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: BodyComponent,
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginComponent,
  },
  {
    path: '/support',
    name: 'Support',
    component: SupportComponent,
  },
  {
    path: '/clients',
    name: 'Clients',
    component: ClientsList,
  },
  {
    path: '/client',
    name: 'ClientEdit',
    component: ClientEdit,
    props: true,
  },
  {
    path: '/seminars',
    name: 'Seminars',
    component: SeminarsList,
  },
  {
    path: '/seminar',
    name: 'SeminarEdit',
    component: SeminarEdit,
    props: true,
  },
  {
    path: '/users',
    name: 'Users',
    component: UsersList,
  },
  {
    path: '/user',
    name: 'UserEdit',
    component: UserEdit,
    props: true,
  },
  {
    path: '/companies',
    name: 'Companies',
    component: CompaniesList,
  },
  {
    path: '/company',
    name: 'CompanyEdit',
    component: CompanyEdit,
    props: true,
  },
]

const router = createRouter({history: createWebHistory(), routes})
export default router