import { createRouter, createWebHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import PlaceDetailView from "@/views/PlaceDetailView.vue";
import LoginView from "@/views/LoginView.vue";
import RegisterView from "@/views/RegisterView.vue";
import VisitHistoryView from "@/views/VisitHistoryView.vue";
import WishlistView from "@/views/WishlistView.vue";
import AddPlaceView from "@/views/AddPlaceView.vue";
import { getToken } from "@/api/client";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/",
      name: "home",
      component: HomeView
    },
    {
      path: "/places/:id",
      name: "place-detail",
      component: PlaceDetailView
    },
    { path: "/login",
      name: "login",
      component: LoginView
    },
    {
      path: "/register",
      name: "register",
      component: RegisterView
    },
    {
      path: "/visits",
      name: "visits",
      component: VisitHistoryView,
      meta: { requiresAuth: true },
    },
    {
      path: "/wishlist",
      name: "wishlist",
      component: WishlistView,
      meta: { requiresAuth: true },
    },
    {
      path: "/add-place",
      name: "add-place",
      component: AddPlaceView,
      meta: { requiresAuth: true },
    },
  ],
});

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !getToken()) {
    return {name: "login", query: {redirect: to.fullPath}}
  }
})

export default router;
