import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue')
    },
    {
      path: '/rooms',
      children: [
        {
          path: '',
          name: 'roomlist',
          component: () => import('../views/RoomListView.vue')
        },
        {
          path: '/rooms/new',
          name: 'new-room',
          component: () => import('../views/NewRoomView.vue')
        },
        {
          path: '/rooms/:roomId',
          children: [
            {
              path: '',
              name: 'room',
              component: () => import('@/views/RoomView.vue'),
            },
            {
              path:'/rooms/:roomId/beer/:beerId',
              name: 'beer',
              component: () => import('@/views/BeerView.vue'),
            }
          ]
        }
      ]
    }
  ],
})

const checkToken = async () => {
  const user = localStorage.getItem('token')
  if (!user) return;

  try {
    const result = await fetch(`${import.meta.env.VITE_API_URL}/api/verifyToken`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      }
    );
    if (result.status.toString().includes('4')) {
      localStorage.removeItem('token');

      router.replace('/login');
    }
  } catch (error) {
    console.error(error);
  }
};

checkToken();

router.beforeEach(async (to, from) => {
  const isAuthenticated = localStorage.getItem('token')
  const publicPages = ['home', 'login', 'register'];
  if (
    !isAuthenticated &&
    !publicPages.includes(to.name)
  ) {
    return { name: 'login' }
  }
})

export default router


