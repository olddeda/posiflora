import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom'
import { TelegramIntegrationPage } from './pages/TelegramIntegrationPage'
import './i18n'
import './index.css'

const router = createBrowserRouter(
  [
    {
      path: '/shops/:shopId/growth/telegram',
      element: <TelegramIntegrationPage />,
    },
    {
      path: '*',
      element: (
        <Navigate
          to="/shops/1/growth/telegram"
          replace
        />
      ),
    },
  ],
  {
    future: {
      v7_relativeSplatPath: true,
    },
  },
)

ReactDOM.createRoot(document.getElementById('root')!).render(
  <RouterProvider
    router={router}
    future={{ v7_startTransition: true }}
  />,
)
