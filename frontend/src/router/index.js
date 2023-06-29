import { createRouter, createWebHistory } from 'vue-router'
import BodyComponent from './../components/BodyComponent.vue'
import LoginComponent from './../components/LoginComponent.vue'
// import UsersList from '../components/UsersList.vue'
// import UserEdit from './../components/UserEdit.vue'
import SupportComponent from './../components/SupportComponent.vue'

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
  // {
  //   path: '/users',
  //   name: 'Users',
  //   component: UsersList,
  // },
  // {
  //   path: '/user',
  //   name: 'UserEdit',
  //   component: UserEdit,
  //   props: true,
  // },
]

const router = createRouter({history: createWebHistory(), routes})
export default router