import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";
import Lobby from "./views/Lobby.vue";

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: "/",
            name: "home",
            component: Home
        },
        {
            path: "/lobby",
            redirect: "/"
        },
        {
            path: "/lobby/:lobbyId",
            name: "lobby",
            component: Lobby,
        },
    ]
})