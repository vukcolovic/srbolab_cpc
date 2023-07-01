<template>
  <nav class="navbar navbar-dark bg-dark navbar-expand-lg">
    <div class="container-fluid">
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse d-flex">
        <a class="navbar-brand" href="/"><img src="../assets/left_logo.png" href=""></a>
        <ul class="navbar-nav">
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown"
               aria-expanded="false">
              Vozači
            </a>
            <ul class="dropdown-menu" aria-labelledby="navBarDropDown">
              <li>
                <router-link class="dropdown-item" to="/clients">Vozači</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/clients-notverified">Vozači za verifikaciju</router-link>
              </li>
            </ul>
          </li>
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown"
               aria-expanded="false">
              Seminari
            </a>
            <ul class="dropdown-menu" aria-labelledby="navBarDropDown">
              <li>
                <router-link class="dropdown-item" to="/clients">Vozači</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/clients-notverified">Vozači za verifikaciju</router-link>
              </li>
            </ul>
          </li>
          <li v-can="'AdminPanelMeni'" class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown"
               aria-expanded="false">
              Admin panel
            </a>
            <ul class="dropdown-menu" aria-labelledby="navBarDropDown">
              <li>
                <router-link v-can="'PregledKorisnika'" class="dropdown-item" to="/users">Korisnici</router-link>
              </li>
              <li>
                <router-link v-can="'UpravljanjeRolama'" class="dropdown-item" to="/roles">Role i aktivnosti</router-link>
              </li>
            </ul>
          </li>
          <li v-can="'SifarniciMeni'" class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="navBarDropDown" role="button" data-bs-toggle="dropdown"
               aria-expanded="false">
              Šifarnici
            </a>
            <ul class="dropdown-menu" aria-labelledby="navBarDropDown">
              <li>
                <router-link class="dropdown-item" to="/irregularity-levels">Nivoi nepravilnosti</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/examination-types">Vrste ispitivanja</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/fuels">Goriva</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/vehicle-types">Vrsta Vozila</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/clients">Klijenti</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/examination-places">Ispitna Mesta</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/vehicle-brands">Marke Vozila</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/notes">Napomene</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/conclusions">Zaključci</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/examination-prices">Cene</router-link>
              </li>
            </ul>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/support">Korisnička podrška</router-link>
          </li>
        </ul>
        <ul class="ms-auto">
          <label style="color: white; font-size: 0.7em">{{firstName}} {{lastName}}</label>
          <button class="iconBtn" title="Odjava" @click="logout"><i class="fa fa-sign-out"></i></button>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>

export default {
  name: 'HeaderComponenet',
  data() {
    return {}
  },
  computed : {
    firstName: function () {
      // return this.$store.getters.firstName;
      return "Vuk";
    },
    lastName: function () {
      // return this.$store.getters.lastName;
      return "Colovic";
    }
  },
  methods: {
    logout() {
      this.$store.dispatch('setTokenAction', "");
      localStorage.removeItem('activities');
    }
  },
}
</script>